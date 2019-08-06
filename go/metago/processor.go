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
	"strconv"
	"strings"

	"github.com/go-toolsmith/astcopy"
	"golang.org/x/tools/go/ast/astutil"
	"golang.org/x/tools/go/packages"
)

const metagoPackagePath = "github.com/vvakame/til/go/metago"

var valueTypeString string
var fieldTypeString string
var valueOfTypeString string
var fieldNameMethodName string
var fieldValueMethodName string
var fieldStructTagGetMethodName string

const valueFieldMethodName = "Fields"

func init() {
	{
		var v Value
		valueTypeString = reflect.ValueOf(&v).Elem().Type().Name()
	}
	{
		var f Field
		fieldTypeString = reflect.ValueOf(&f).Elem().Type().Name()
	}
	{
		s := runtime.FuncForPC(reflect.ValueOf(ValueOf).Pointer()).Name()
		valueOfTypeString = strings.TrimPrefix(s, metagoPackagePath+".")
	}
	{
		var f Field
		method, ok := reflect.TypeOf(&f).Elem().MethodByName("Name")
		if !ok {
			panic("metago.Field#Name method is missing")
		}
		fieldNameMethodName = method.Name
	}
	{
		var f Field
		method, ok := reflect.TypeOf(&f).Elem().MethodByName("Value")
		if !ok {
			panic("metago.Field#Value method is missing")
		}
		fieldValueMethodName = method.Name
	}
	{
		var f Field
		method, ok := reflect.TypeOf(&f).Elem().MethodByName("StructTagGet")
		if !ok {
			panic("metago.Field#StructTagGet method is missing")
		}
		fieldStructTagGetMethodName = method.Name
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
			p := &metaProcessor{
				currentPkg:        pkg,
				currentFile:       file,
				removeNodes:       make(map[ast.Node]bool),
				replaceNodes:      make(map[ast.Node]ast.Node),
				valueMapping:      make(map[*ast.Object]ast.Expr),
				fieldMapping:      make(map[*ast.Object]ast.Expr),
				fieldBlockMapping: make(map[*ast.BlockStmt]*ast.Object),
			}

			if !astutil.UsesImport(file, metagoPackagePath) {
				continue
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
	currentPkg         *packages.Package
	currentFile        *ast.File
	currentTargetField *ast.Object

	removeNodes  map[ast.Node]bool
	replaceNodes map[ast.Node]ast.Node

	// mv → obj への変換用
	valueMapping map[*ast.Object]ast.Expr
	// mf → obj.X への変換用 X はわからんので obj 部分を持つ
	fieldMapping map[*ast.Object]ast.Expr
	// あるBlockStmt中でどの*ast.Identに紐づくか
	// { ... } は struct { ID int64 } の `ID` に紐づく！
	// *ast.IdentのObjをたぐるとフィールドのTypeとかTagも掘れる
	fieldBlockMapping map[*ast.BlockStmt]*ast.Object

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
	if n := p.replaceNodes[current]; n != nil {
		cursor.Replace(n)
		return false
	}

	switch node := current.(type) {
	case *ast.AssignStmt:
		if p.checkMetagoValueOfAssignStmt(cursor, node) {
			return true
		}

	case *ast.Ident:
		if p.checkReplaceTargetIdent(cursor, node) {
			return false
		}

	case *ast.RangeStmt:
		if p.checkMetagoFieldRange(cursor, node) {
			return false
		}

	case *ast.CallExpr:
		if p.checkInlineTemplateCallExpr(cursor, node) {
			return false
		}
		if p.checkUseMetagoFieldValue(cursor, node) {
			return false
		}
		if p.checkUseMetagoFieldName(cursor, node) {
			return false
		}
		if p.checkUseMetagoStructTagGet(cursor, node) {
			return false
		}

	case *ast.IfStmt:
		if p.checkIfStmtInInitWithTypeAssert(cursor, node) {
			return false
		}
		if p.checkIfStmtInCondWithTypeAssert(cursor, node) {
			// TypeAssertExprの置き換えやらで子要素を歩く必要があるのでtrue返す
			return true
		}

	case *ast.TypeSwitchStmt:
		if p.checkTypeSwitchStmt(cursor, node) {
			return false
		}

	case *ast.FuncDecl:
		if p.isInlineTemplateFuncDecl(cursor, node) {
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

func (p *metaProcessor) relateToMetagoPackage(ident *ast.Ident) bool {
	// [mv] ← metago packageですか？ とか
	// mv.[Value]() ← metago packageですか？ とか
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

func (p *metaProcessor) isMetagoField(ident *ast.Ident) bool {
	// var ident metago.Field ← true
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
	} else if typeName.Name() != fieldTypeString {
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

func (p *metaProcessor) isCallMetagoFieldValue(node *ast.CallExpr) bool {
	// mf.Value() 系かどうか
	selectorExpr, ok := node.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}

	objIdent, ok := selectorExpr.X.(*ast.Ident)
	if !ok {
		return false
	}

	target := p.fieldMapping[objIdent.Obj]
	if target == nil {
		return false
	}

	if selectorExpr.Sel.Name != fieldValueMethodName {
		return false
	}

	return true
}

func (p *metaProcessor) isCallMetagoFieldName(node *ast.CallExpr) bool {
	// mf.Name() 系かどうか
	selectorExpr, ok := node.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}

	objIdent, ok := selectorExpr.X.(*ast.Ident)
	if !ok {
		return false
	}

	target := p.fieldMapping[objIdent.Obj]
	if target == nil {
		return false
	}

	if selectorExpr.Sel.Name != fieldNameMethodName {
		return false
	}

	return true
}

func (p *metaProcessor) isAssignable(expr1 ast.Expr, expr2 ast.Expr) bool {
	if p.isSameType(expr1, expr2) {
		return true
	}

	// TODO ものすごく素晴らしくする
	// time.Time が json.Marshaler に assign できることがわかると素敵なコードが書けるぞ！
	// …といっても p.currentPkg.TypesInfo.Uses とかは astcopy との相性が悪くて死だ！
	// https://golang.org/pkg/go/types/#AssignableTo
	// https://golang.org/pkg/go/types/#Info.TypeOf
	// ↑この辺ちゃう？ってtenntennさんが言ってた

	return false
}

func (p *metaProcessor) isSameType(expr1 ast.Expr, expr2 ast.Expr) bool {
	{
		ident1, ok1 := expr1.(*ast.Ident)
		ident2, ok2 := expr2.(*ast.Ident)
		if ok1 && ok2 {
			if ident1.Obj == nil && ident2.Obj == nil {
				type1 := types.Universe.Lookup(ident1.Name)
				type2 := types.Universe.Lookup(ident2.Name)
				if type1 != nil && type2 != nil {
					return type1 == type2
				} else if type1 == nil && type2 == nil {
					// TODO package の ident の場合の比較が甘い…！
					// import hoge "hoge" と import h "hoge" で hoge と h 比較した時にtrueにならない
					return ident1.Name == ident2.Name
				}
				return false
			} else if ident1.Obj == ident2.Obj {
				return true
			}
			return false
		} else if ok1 {
			return false
		} else if ok2 {
			return false
		}
	}
	{
		sel1, ok1 := expr1.(*ast.SelectorExpr)
		sel2, ok2 := expr2.(*ast.SelectorExpr)
		if ok1 && ok2 {
			if p.isSameType(sel1.X, sel2.X) && p.isSameType(sel1.Sel, sel2.Sel) {
				return true
			}
			return false

		} else if ok1 {
			return false
		} else if ok2 {
			return false
		}
	}

	panic("unreachable")
}

func (p *metaProcessor) isInlineTemplateFuncDecl(cursor *astutil.Cursor, node *ast.FuncDecl) bool {
	// TODO func (obj *Foo) Template(mv metago.Value) 的なメソッドは除外する

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
		p.Errorf(arg1, "argument must be ident")
		return nil
	}

	return targetIdent
}

func (p *metaProcessor) checkReplaceTargetIdent(cursor *astutil.Cursor, node *ast.Ident) bool {
	// mv系の単純な置き換え
	if target := p.valueMapping[node.Obj]; target != nil {
		cursor.Replace(astcopy.Node(target))
		return true
	}
	// mf系の単純な置き換え
	if target := p.fieldMapping[node.Obj]; target != nil {
		field := p.currentTargetField
		cursor.Replace(&ast.SelectorExpr{
			X: target,
			Sel: &ast.Ident{
				Name: field.Name,
				Obj:  field,
			},
		})
		return true
	}

	return false
}

// checkMetagoValueOfAssignStmt is capture `mv := metago.ValueOf(foo)` format assignment.
// it marks up convert rule about `mv` to `foo` and remove these assignment.
func (p *metaProcessor) checkMetagoValueOfAssignStmt(cursor *astutil.Cursor, stmt *ast.AssignStmt) bool {
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
		p.valueMapping[ident.Obj] = &ast.Ident{
			Name: targetIdent.Name,
			Obj:  targetIdent.Obj,
		}

		p.removeNodes[lhs] = true
		p.removeNodes[rhs] = true
	}

	return found
}

func (p *metaProcessor) checkMetagoFieldRange(cursor *astutil.Cursor, node *ast.RangeStmt) bool {
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
	cursor.Delete()

	// RangeStmtのBody部分をフィールドの数だけコピペする
	var targetObjDef *ast.Object
	switch targetNode := target.(type) {
	case *ast.Ident:
		// TODO 今出てくるパターンを決め打ちでサポートしてるだけなのでいい感じにする
		// ↓は foo *Foo 的なやつの *Foo から Foo の定義を取ろうとしている
		targetObjDef = targetNode.Obj.Decl.(*ast.Field).Type.(*ast.StarExpr).X.(*ast.Ident).Obj
	default:
		panic("unknown type")
	}
	defFields := targetObjDef.Decl.(*ast.TypeSpec).Type.(*ast.StructType).Fields
	for _, field := range defFields.List {
		for _, name := range field.Names {
			bk := p.currentTargetField

			bodyStmt := astcopy.BlockStmt(node.Body)
			p.fieldBlockMapping[bodyStmt] = name.Obj
			p.currentTargetField = name.Obj
			astutil.Apply(
				bodyStmt,
				p.ApplyPre,
				p.ApplyPost,
			)
			cursor.InsertBefore(bodyStmt)

			p.currentTargetField = bk
		}
	}

	return true
}

func (p *metaProcessor) checkIfStmtInInitWithTypeAssert(cursor *astutil.Cursor, node *ast.IfStmt) bool {
	// if v, ok := mf.Value().(time.Time); ok { ... } 系
	// condがokの参照か !ok か ok == false 以外の場合怒る
	if node.Init == nil {
		return false
	}

	assignStmt, ok := node.Init.(*ast.AssignStmt)
	if !ok {
		return false
	}

	// mf.Value().(time.Time) 的なやつかチェック
	if len(assignStmt.Rhs) != 1 {
		return false
	}
	typeAssertExpr, ok := assignStmt.Rhs[0].(*ast.TypeAssertExpr)
	if !ok {
		return false
	}
	callExpr, ok := typeAssertExpr.X.(*ast.CallExpr)
	if !ok {
		return false
	}
	if !p.isCallMetagoFieldValue(callExpr) {
		return false
	}

	// v, ok := 的なやつかチェック
	if len(assignStmt.Lhs) != 2 {
		p.Errorf(assignStmt, "lhs assignment should be 2")
		return false
	}

	varIdent, ok := assignStmt.Lhs[0].(*ast.Ident)
	if !ok {
		p.Errorf(assignStmt.Lhs[0], "var assignment should be ident")
		return false
	}
	okIdent, ok := assignStmt.Lhs[1].(*ast.Ident)
	if !ok {
		p.Errorf(assignStmt.Lhs[1], "ok assignment should be ident")
		return false
	}

	// IfStmtでスコープに新しい変数が導入されるので置き換えルールを登録
	p.valueMapping[varIdent.Obj] = &ast.SelectorExpr{
		X: p.fieldMapping[callExpr.Fun.(*ast.SelectorExpr).X.(*ast.Ident).Obj],
		Sel: &ast.Ident{
			Name: p.currentTargetField.Name,
			Obj:  p.currentTargetField,
		},
	}

	_, ok = node.Cond.(*ast.Ident)
	// TODO この辺もうちょっと柔軟性をもたせる 静的にboolに還元できる範囲であれば許容してあげたい
	if !ok {
		p.Errorf(node.Cond, "must be '%s'", okIdent.Name)
		return false
	}

	if p.isAssignable(typeAssertExpr.Type, p.currentTargetField.Decl.(*ast.Field).Type) {
		// Bodyが評価される & if全体を置き換え
		astutil.Apply(
			node.Body,
			p.ApplyPre,
			p.ApplyPost,
		)
		cursor.Replace(node.Body)
	} else if node.Else != nil {
		// Elseが評価される & if全体を置き換え
		astutil.Apply(
			node.Else,
			p.ApplyPre,
			p.ApplyPost,
		)
		cursor.Replace(node.Else)
	}

	return true
}

func (p *metaProcessor) checkIfStmtInCondWithTypeAssert(cursor *astutil.Cursor, node *ast.IfStmt) bool {
	// if mf.Value().(time.Time).IsZero() { ... } 系のハンドリング
	// condで、もしvalueの型とassert先がマッチしてたらBlockStmt残し、それ以外は削除

	var typeAssertExpr *ast.TypeAssertExpr
	ast.Walk(astVisitorFunc(func(node ast.Node) bool {
		found, ok := node.(*ast.TypeAssertExpr)
		if ok {
			typeAssertExpr = found
			return false
		}
		return true
	}), node.Cond)
	if typeAssertExpr == nil {
		return true
	}

	callExpr, ok := typeAssertExpr.X.(*ast.CallExpr)
	if !ok {
		return false
	}
	if !p.isCallMetagoFieldValue(callExpr) {
		return false
	}

	targetField := p.currentTargetField
	if targetField == nil {
		p.Errorf(callExpr, "invalid context. not in metago.Field range statement")
		return false
	}

	if !p.isSameType(typeAssertExpr.Type, targetField.Decl.(*ast.Field).Type) {
		cursor.Delete()
		return false
	}

	p.replaceNodes[typeAssertExpr] = &ast.SelectorExpr{
		X: p.fieldMapping[callExpr.Fun.(*ast.SelectorExpr).X.(*ast.Ident).Obj],
		Sel: &ast.Ident{
			Name: p.currentTargetField.Name,
			Obj:  p.currentTargetField,
		},
	}

	return true
}

func (p *metaProcessor) checkTypeSwitchStmt(cursor *astutil.Cursor, node *ast.TypeSwitchStmt) bool {
	assignStmt, ok := node.Assign.(*ast.AssignStmt)
	if !ok {
		return false
	}

	// mf.Value().(type) 的なやつかチェック
	if len(assignStmt.Rhs) != 1 {
		return false
	}
	typeAssertExpr, ok := assignStmt.Rhs[0].(*ast.TypeAssertExpr)
	if !ok {
		return false
	}
	if typeAssertExpr.Type != nil {
		p.Errorf(typeAssertExpr, "unknown type assert expr. must be use foo.(type)")
		return false
	}
	callExpr, ok := typeAssertExpr.X.(*ast.CallExpr)
	if !ok {
		return false
	}
	if !p.isCallMetagoFieldValue(callExpr) {
		return false
	}

	// 新しい変数が導入される場合そのマッピングルールを記録
	if len(assignStmt.Lhs) != 1 {
		p.Errorf(assignStmt, "var assignment must required")
		return false
	}
	varIdent, ok := assignStmt.Lhs[0].(*ast.Ident)
	if !ok {
		p.Errorf(assignStmt.Lhs[0], "var assignment should be ident")
		return false
	}
	p.valueMapping[varIdent.Obj] = &ast.SelectorExpr{
		X: p.fieldMapping[callExpr.Fun.(*ast.SelectorExpr).X.(*ast.Ident).Obj],
		Sel: &ast.Ident{
			Name: p.currentTargetField.Name,
			Obj:  p.currentTargetField,
		},
	}

	var targetBody []ast.Stmt // nil と len == 0 を区別する
	for _, stmt := range node.Body.List {
		switch stmt := stmt.(type) {
		case *ast.CaseClause:
			if targetBody == nil && len(stmt.List) == 0 {
				// 長さ 0 は default
				targetBody = stmt.Body
			} else {
				for _, typeExpr := range stmt.List {
					if p.isAssignable(typeExpr, p.currentTargetField.Decl.(*ast.Field).Type) {
						targetBody = stmt.Body
					}
				}
			}

		default:
			panic("unreachable")
		}
	}

	newBlock := &ast.BlockStmt{
		List: targetBody,
	}
	astutil.Apply(newBlock, p.ApplyPre, p.ApplyPost)
	cursor.Replace(newBlock)

	return true
}

func (p *metaProcessor) checkInlineTemplateCallExpr(cursor *astutil.Cursor, node *ast.CallExpr) bool {
	// func fooBarTemplate(mv metago.Value, a, b string) bool 的なやつの変換
	// 第一引数が metago.Value だったら対象

	// foo(mv) 形式のみ対応 メソッド類は対応が大変
	funcName, ok := node.Fun.(*ast.Ident)
	if !ok {
		return false
	}

	funcDecl, ok := funcName.Obj.Decl.(*ast.FuncDecl)
	if !ok {
		return false
	}

	// 引数無しは対象外
	if len(funcDecl.Type.Params.List) == 0 {
		return false
	}
	// 引数の最初が metago.Value じゃないものは対象外
	metaValueArg := funcDecl.Type.Params.List[0].Names[0]
	if !p.isMetagoValue(metaValueArg) {
		return false
	}

	// 実引数側の数が0ってことはここまで来たらないだろうけど一応
	if len(node.Args) == 0 {
		return false
	}
	arg, ok := node.Args[0].(*ast.Ident)
	if !ok {
		return false
	}

	funcDecl = astcopy.FuncDecl(funcDecl)

	// 実引数側の mv がマッピングされる先を 仮引数側の mv にも継承させる
	// *ast.Ident#Obj はコピーされないので metaValueArg 取り直さなくても大丈夫
	p.valueMapping[metaValueArg.Obj] = p.valueMapping[arg.Obj]

	// 引数が metago.Value だけならinline展開するやつ
	// goroutineの境界変わったりするとめんどいので即時実行関数で包む
	newCallExpr := &ast.CallExpr{
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

	// 操作したNodeには入っていってくれないので自分で歩く必要がある
	astutil.Apply(
		newCallExpr,
		p.ApplyPre,
		p.ApplyPost,
	)
	cursor.Replace(newCallExpr)

	return true
}

func (p *metaProcessor) checkUseMetagoFieldValue(cursor *astutil.Cursor, node *ast.CallExpr) bool {
	// mf.Value() 系を obj.Foo 的なのに置き換える
	selectorExpr, ok := node.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}

	objIdent, ok := selectorExpr.X.(*ast.Ident)
	if !ok {
		return false
	}

	target := p.fieldMapping[objIdent.Obj]
	if target == nil {
		return false
	}

	if selectorExpr.Sel.Name != fieldValueMethodName {
		return false
	}

	cursor.Replace(&ast.SelectorExpr{
		X: target,
		Sel: &ast.Ident{
			Name: p.currentTargetField.Name,
			Obj:  p.currentTargetField,
		},
	})

	return false
}

func (p *metaProcessor) checkUseMetagoFieldName(cursor *astutil.Cursor, node *ast.CallExpr) bool {
	// mf.Name() 系を "Foo" 的なのに置き換える
	if !p.isCallMetagoFieldName(node) {
		return false
	}

	cursor.Replace(&ast.BasicLit{
		Kind:  token.STRING,
		Value: strconv.Quote(p.currentTargetField.Name),
	})

	return false
}

func (p *metaProcessor) checkUseMetagoStructTagGet(cursor *astutil.Cursor, node *ast.CallExpr) bool {
	// mf.Name() 系を "Foo" 的なのに置き換える
	selectorExpr, ok := node.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}

	objIdent, ok := selectorExpr.X.(*ast.Ident)
	if !ok {
		return false
	}

	target := p.fieldMapping[objIdent.Obj]
	if target == nil {
		return false
	}

	if selectorExpr.Sel.Name != fieldStructTagGetMethodName {
		return false
	}

	if len(node.Args) != 1 {
		p.Errorf(node, "string literal argument must required")
		return false
	}

	basicLit, ok := node.Args[0].(*ast.BasicLit)
	if !ok || basicLit.Kind != token.STRING {
		p.Errorf(node.Args[0], "string literal argument must required")
		return false
	}

	tagName, err := strconv.Unquote(basicLit.Value)
	if err != nil {
		p.Errorf(node.Args[0], "unexpected string literal format. %s: %s", basicLit.Value, err.Error())
		return false
	}

	targetField := p.currentTargetField.Decl.(*ast.Field)
	if targetField.Tag == nil {
		cursor.Replace(&ast.BasicLit{
			Kind:  token.STRING,
			Value: `""`,
		})
		return false
	}
	structTagValue, err := strconv.Unquote(targetField.Tag.Value)
	if err != nil {
		p.Errorf(targetField, "unexpected string literal format. %s: %s", targetField.Tag.Value, err.Error())
		return false
	}

	tagValue := reflect.StructTag(structTagValue).Get(tagName)

	cursor.Replace(&ast.BasicLit{
		Kind:  token.STRING,
		Value: strconv.Quote(tagValue),
	})

	return false
}

func (p *metaProcessor) Debugf(node ast.Node, format string, a ...interface{}) {
	p.nodeErrors = append(p.nodeErrors, &NodeError{
		ErrorLevel: ErrorLevelDebug,
		Fset:       p.currentPkg.Fset,
		Node:       node,
		Message:    fmt.Sprintf(format, a...),
	})
}

func (p *metaProcessor) Noticef(node ast.Node, format string, a ...interface{}) {
	p.nodeErrors = append(p.nodeErrors, &NodeError{
		ErrorLevel: ErrorLevelNotice,
		Fset:       p.currentPkg.Fset,
		Node:       node,
		Message:    fmt.Sprintf(format, a...),
	})
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
