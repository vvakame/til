package metago

import (
	"errors"
	"fmt"
	"go/ast"
	"go/format"
	"go/token"
	"go/types"
	"os"
	"reflect"
	"runtime"
	"strings"

	"github.com/go-toolsmith/astcopy"
	"golang.org/x/tools/go/ast/astutil"
	"golang.org/x/tools/go/packages"
)

const metagoPackagePath = "github.com/vvakame/til/go/metago"

var valueTypeString string
var valueOfTypeString string

const valueFieldMethodName = "Fields"

func init() {
	{
		var v Value
		valueTypeString = reflect.ValueOf(&v).Elem().Type().Name()
	}
	{
		s := runtime.FuncForPC(reflect.ValueOf(ValueOf).Pointer()).Name()
		valueOfTypeString = strings.TrimPrefix(s, metagoPackagePath+".")
	}
}

func Process(patterns ...string) error {
	cfg := &packages.Config{
		Mode: packages.NeedName |
			packages.NeedFiles |
			packages.NeedImports |
			packages.NeedDeps |
			packages.NeedTypes |
			packages.NeedSyntax |
			packages.NeedTypesInfo,
	}
	pkgs, err := packages.Load(cfg, patterns...)
	if err != nil {
		return err
	}

	if packages.PrintErrors(pkgs) > 0 {
		return errors.New("some errors occured")
	}

	var nErrs NodeErrors

	for _, pkg := range pkgs {
		fmt.Println(pkg.ID, pkg.GoFiles)
		for _, file := range pkg.Syntax {
			if !astutil.UsesImport(file, metagoPackagePath) {
				continue
			}

			p := &metaProcessor{
				currentPkg:   pkg,
				currentFile:  file,
				removeNodes:  make(map[ast.Node]bool),
				valueMapping: make(map[*ast.Object]*ast.Object),
				fieldMapping: make(map[*ast.Object]*ast.Object),
			}

			// TODO . import してたら殺す

			astutil.Apply(
				file,
				p.ApplyPre,
				p.ApplyPost,
			)

			nErrs = append(nErrs, p.nodeErrors...)

			err := format.Node(os.Stdout, pkg.Fset, file)
			if err != nil {
				panic(err)
			}
			fmt.Print("\n\n")
		}
	}

	if len(nErrs) != 0 {
		return nErrs
	}

	return nil
}

type metaProcessor struct {
	currentPkg  *packages.Package
	currentFile *ast.File

	removeNodes map[ast.Node]bool

	// mv → obj への変換用
	valueMapping map[*ast.Object]*ast.Object
	// mf → obj.X への変換用 X はわからんので obj 部分を持つ
	fieldMapping map[*ast.Object]*ast.Object

	nodeErrors NodeErrors
}

func (p *metaProcessor) ApplyPre(cursor *astutil.Cursor) bool {
	current := cursor.Node()
	if current == nil {
		return true
	}

	if p.removeNodes[current] {
		cursor.Delete()
		return false
	}

	switch node := current.(type) {
	case *ast.AssignStmt:
		if p.checkMetagoValueOfAssignStmt(node) {
			// mv := metago.ValueOf(foo) 系
			return true
		}

	case *ast.Ident:
		if target := p.valueMapping[node.Obj]; target != nil {
			cursor.Replace(&ast.Ident{
				Name: target.Name,
				Obj:  target,
			})
			return true
		}

	case *ast.RangeStmt:
		// 特殊対応ポイント
		if p.checkMetagoFieldRange(node) {
			cursor.Delete()
			// TODO BlockStatementに対して繰り返し処理をアレする
			return false
		}

	case *ast.CallExpr:
		if newCallExpr, ok := p.checkCallExprWithMetagoValue(node); ok {
			// 操作したNodeには入っていってくれないので自分で歩く必要がある
			astutil.Apply(
				newCallExpr,
				p.ApplyPre,
				p.ApplyPost,
			)
			cursor.Replace(newCallExpr)
			return false // replace前のASTを見る必要はない
		}

	case *ast.CompositeLit:
	// 特殊対応ポイント

	case *ast.TypeAssertExpr:
		// metago関連だったら…？

	case *ast.FuncDecl:
		if p.isInlineTemplateFuncDecl(node) {
			// 仮引数に metago.Value があったら、展開処理の対象ではないのでskip
			return false
		}
	}

	return true
}

func (p *metaProcessor) ApplyPost(cursor *astutil.Cursor) bool {
	// NOTE Postでは基本的に return false しない panic+recover されて処理がわからなくなるぞ！

	current := cursor.Node()
	if current == nil {
		return true
	}

	if p.removeNodes[current] {
		cursor.Delete()
		return true
	}

	switch node := current.(type) {
	case *ast.AssignStmt:
		if len(node.Lhs) == 0 && len(node.Rhs) == 0 {
			// metago.ValueOf を消した結果空になる場合がある
			cursor.Delete()
			return true
		}
	}

	return true
}

func (p *metaProcessor) isMetagoValue(ident *ast.Ident) bool {
	// var ident metago.Value ← true
	// var ident FooBar ← false
	v := p.currentPkg.TypesInfo.Defs[ident]
	if v == nil {
		return false
	}
	t := v.Type()
	tn, ok := t.(*types.Named)
	if !ok {
		return false
	}
	typeName := tn.Obj()
	typePkg := typeName.Pkg()
	if typePkg == nil {
		return false
	}
	if typePkg.Path() != metagoPackagePath {
		return false
	} else if typeName.Name() != valueTypeString {
		return false
	}

	return true
}

func (p *metaProcessor) isMetagoValueOf(selectorExpr *ast.SelectorExpr) bool {
	// selectorExpr == metago.ValueOf
	lhsIdent, ok := selectorExpr.X.(*ast.Ident)
	if !ok {
		return false
	}
	var found bool
	for _, pkg := range p.currentPkg.Types.Imports() {
		if pkg.Path() != metagoPackagePath {
			continue
		}
		if pkg.Name() == lhsIdent.Name {
			found = true
			break
		}
	}
	if !found {
		return false
	}

	if selectorExpr.Sel.Name != valueOfTypeString {
		return false
	}

	return true
}

func (p *metaProcessor) isInlineTemplateFuncDecl(node *ast.FuncDecl) bool {
	// TODO メソッドは除外する

	for _, params := range node.Type.Params.List {
		for _, param := range params.Names {
			if p.isMetagoValue(param) {
				return true
			}
		}
	}

	return false
}

func (p *metaProcessor) extractMetagoBaseVariable(expr ast.Expr) *ast.Ident {
	callExpr, ok := expr.(*ast.CallExpr)
	if !ok {
		return nil
	}
	selectorExpr, ok := callExpr.Fun.(*ast.SelectorExpr)
	if !ok {
		return nil
	}

	if !p.isMetagoValueOf(selectorExpr) {
		return nil
	}

	arg1 := callExpr.Args[0]
	targetIdent, ok := arg1.(*ast.Ident)
	if !ok {
		// TODO エラーにするべきかも
		return nil
	}

	return targetIdent
}

// checkMetagoValueOfAssignStmt is capture `mv := metago.ValueOf(foo)` format assignment.
// it marks up convert rule about `mv` to `foo` and remove these assignment.
func (p *metaProcessor) checkMetagoValueOfAssignStmt(stmt *ast.AssignStmt) bool {
	// mv := metago.ValueOf(foo) 系を処理する。
	// mv と foo の紐付けを覚える。
	// 該当のassignmentをNode毎削除するようマークする。

	var found bool
	for idx, lhs := range stmt.Lhs {
		ident, ok := lhs.(*ast.Ident)
		if !ok {
			continue
		}
		if !p.isMetagoValue(ident) {
			continue
		}

		rhs := stmt.Rhs[idx]
		targetIdent := p.extractMetagoBaseVariable(rhs)
		if targetIdent == nil {
			// TODO なんらかの警告を出したほうがよさそう
			continue
		}

		found = true
		p.valueMapping[ident.Obj] = targetIdent.Obj

		p.removeNodes[lhs] = true
		p.removeNodes[rhs] = true
	}

	return found
}

func (p *metaProcessor) checkMetagoFieldRange(node *ast.RangeStmt) bool {
	// for _, mf := range mv.Fields() {
	// ↑的なヤツをサポートしていく

	ident, ok := node.Key.(*ast.Ident)
	if !ok {
		return false
	}
	if ident.Name != "_" {
		p.Warningf(ident, "index part assignment should be '_'")
		return false
	}

	// assignStmt → _, mf := range mv.Fields() 部分相当
	assignStmt, ok := ident.Obj.Decl.(*ast.AssignStmt)
	if !ok {
		return false
	}
	if len(assignStmt.Lhs) != 2 {
		// := の左が2未満だとindexとかしか取ってない
		p.Errorf(assignStmt, "value part assignment must required")
		return false
	}
	// _, mf ← の mf 相当の場所の名前取る
	fieldIdent, ok := assignStmt.Lhs[1].(*ast.Ident)
	if !ok {
		return false
	}

	// range 以外は対応しないめんどいから
	unaryExpr, ok := assignStmt.Rhs[0].(*ast.UnaryExpr)
	if !ok {
		return false
	}
	if unaryExpr.Op != token.RANGE {
		return false
	}
	if len(assignStmt.Rhs) != 1 {
		// 右辺が1以外は知らないパターンだ
		return false
	}

	// mv.Fields() 相当の部分を調べていく
	// とりあえずCallExprじゃなかったら知らないパターン
	callExpr, ok := unaryExpr.X.(*ast.CallExpr)
	if !ok {
		return false
	}
	// mv.Fields の部分調べる
	// mv 部分がが変換候補じゃない場合処理対象ではない
	selectorExpr, ok := callExpr.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}
	xIdent, ok := selectorExpr.X.(*ast.Ident)
	if !ok {
		return false
	}
	target, ok := p.valueMapping[xIdent.Obj]
	if !ok {
		return false
	}
	// Fields 部分が Fields 以外だったら処理対象ではない
	if selectorExpr.Sel.Name != valueFieldMethodName {
		return false
	}

	p.fieldMapping[fieldIdent.Obj] = target
	// 大本のfor句は全部捨てる必要がある
	p.removeNodes[node] = true

	return true
}

func (p *metaProcessor) checkCallExprWithMetagoValue(node *ast.CallExpr) (*ast.CallExpr, bool) {
	// func fooBarTemplate(mv metago.Value, a, b string) bool 的なやつの変換
	// 第一引数が metago.Value だったら対象

	// foo(mv) 形式のみ対応 メソッド類は対応が大変
	funcName, ok := node.Fun.(*ast.Ident)
	if !ok {
		return nil, false
	}

	funcDecl, ok := funcName.Obj.Decl.(*ast.FuncDecl)
	if !ok {
		return nil, false
	}

	// 引数無しは対象外
	if len(funcDecl.Type.Params.List) == 0 {
		return nil, false
	}
	// 引数の最初が metago.Value じゃないものは対象外
	metaValueArg := funcDecl.Type.Params.List[0].Names[0]
	if !p.isMetagoValue(metaValueArg) {
		return nil, false
	}

	// 実引数側の数が0ってことはここまで来たらないだろうけど一応
	if len(node.Args) == 0 {
		return nil, false
	}
	arg, ok := node.Args[0].(*ast.Ident)
	if !ok {
		return nil, false
	}

	funcDecl = astcopy.FuncDecl(funcDecl)

	// 実引数側の mv がマッピングされる先を 仮引数側の mv にも継承させる
	// *ast.Ident#Obj はコピーされないので metaValueArg 取り直さなくても大丈夫
	p.valueMapping[metaValueArg.Obj] = p.valueMapping[arg.Obj]

	// 引数が metago.Value だけならinline展開するやつ
	// goroutineの境界変わったりするとめんどいので即時実行関数で包む
	newNode := &ast.CallExpr{
		Fun: &ast.FuncLit{
			Type: &ast.FuncType{
				Params: &ast.FieldList{
					List: funcDecl.Type.Params.List[1:], // 先頭は metago.Valueなので
				},
				Results: &ast.FieldList{
					List: funcDecl.Type.Results.List,
				},
			},
			Body: funcDecl.Body,
		},
		Args: node.Args[1:], // 先頭は metago.Valueなので
	}

	return newNode, true
}

func (p *metaProcessor) Warningf(node ast.Node, format string, a ...interface{}) {
	p.nodeErrors = append(p.nodeErrors, &NodeError{
		ErrorLevel: ErrorLevelWarning,
		Fset:       p.currentPkg.Fset,
		Node:       node,
		Message:    fmt.Sprintf(format, a...),
	})
}

func (p *metaProcessor) Errorf(node ast.Node, format string, a ...interface{}) {
	p.nodeErrors = append(p.nodeErrors, &NodeError{
		ErrorLevel: ErrorLevelError,
		Fset:       p.currentPkg.Fset,
		Node:       node,
		Message:    fmt.Sprintf(format, a...),
	})
}
