package metago

import (
	"bytes"
	"errors"
	"fmt"
	"go/ast"
	"go/constant"
	"go/format"
	"go/token"
	"go/types"
	"reflect"
	"strconv"

	"github.com/vvakame/astcopy"
	"golang.org/x/tools/go/ast/astutil"
	"golang.org/x/tools/go/packages"
)

const metagoBuildTag = "metago"
const metagoPackagePath = "github.com/vvakame/til/go/metago"

const (
	valueOfFuncNameString       = "ValueOf"      // metago.[ValueOf](obj)
	valueTypeString             = "Value"        // metago.[Value]
	valueFieldsMethodName       = "Fields"       // mv.[Fields]()
	fieldNameMethodName         = "Name"         // mf.[Name]()
	fieldValueMethodName        = "Value"        // mf.[Value]()
	fieldStructTagGetMethodName = "StructTagGet" // mf.[StructTagGet]("json")
)

type fieldInfo struct {
	recvExpr   ast.Expr
	recvType   *types.Struct
	fieldIndex int
}

func (f *fieldInfo) Field() *types.Var {
	return f.recvType.Field(f.fieldIndex)
}

func (f *fieldInfo) Tag() string {
	return f.recvType.Tag(f.fieldIndex)
}

type metaProcessor struct {
	cfg       *Config
	typesInfo *types.Info

	currentPkg  *packages.Package
	currentFile *ast.File

	hasMetagoBuildTag     bool
	copyNodeMap           astcopy.CopyNodeMap
	removeNodes           map[ast.Node]bool
	replaceNodes          map[ast.Node]ast.Node
	gotoCounter           int
	requiredContinueLabel []string
	requiredBreakLabel    []string

	// mv → obj への変換用 *ast.Object は mv のもの
	valueMapping map[*ast.Object]ast.Expr
	// mf → obj.X への変換用 *ast.Object は mf のもの
	fieldMapping map[*ast.Object]*fieldInfo

	nodeErrors NodeErrors
}

func (p *metaProcessor) Process(cfg *Config) (*Result, error) {
	p.cfg = cfg
	p.typesInfo = &types.Info{}
	p.copyNodeMap = make(astcopy.CopyNodeMap)
	p.removeNodes = make(map[ast.Node]bool)
	p.replaceNodes = make(map[ast.Node]ast.Node)
	p.valueMapping = make(map[*ast.Object]ast.Expr)
	p.fieldMapping = make(map[*ast.Object]*fieldInfo)
	p.requiredBreakLabel = nil
	p.requiredContinueLabel = nil

	pkgCfg := &packages.Config{
		Mode: packages.NeedName |
			packages.NeedFiles |
			packages.NeedImports |
			packages.NeedDeps |
			packages.NeedTypes |
			packages.NeedSyntax |
			packages.NeedTypesInfo,
		BuildFlags: []string{"-tags", metagoBuildTag},
	}
	pkgs, err := packages.Load(pkgCfg, p.cfg.TargetPackages...)
	if err != nil {
		return nil, err
	}

	var errs []packages.Error
	packages.Visit(pkgs, nil, func(pkg *packages.Package) {
		errs = append(errs, pkg.Errors...)
	})
	if len(errs) != 0 {
		return &Result{CompileErrors: errs}, errors.New("some errors occured")
	}

	{ // merge typesInfo
		mergedPkg := make(map[string]bool)
		var mergePkgTypes func(pkg *packages.Package)
		mergePkgTypes = func(pkg *packages.Package) {
			if mergedPkg[pkg.PkgPath] {
				return
			}
			mergedPkg[pkg.PkgPath] = true
			if pkg.TypesInfo.Types != nil {
				if p.typesInfo.Types == nil {
					p.typesInfo.Types = make(map[ast.Expr]types.TypeAndValue)
				}
				for k, v := range pkg.TypesInfo.Types {
					p.typesInfo.Types[k] = v
				}
			}
			if pkg.TypesInfo.Defs != nil {
				if p.typesInfo.Defs == nil {
					p.typesInfo.Defs = make(map[*ast.Ident]types.Object)
				}
				for k, v := range pkg.TypesInfo.Defs {
					p.typesInfo.Defs[k] = v
				}
			}
			if pkg.TypesInfo.Uses != nil {
				if p.typesInfo.Uses == nil {
					p.typesInfo.Uses = make(map[*ast.Ident]types.Object)
				}
				for k, v := range pkg.TypesInfo.Uses {
					p.typesInfo.Uses[k] = v
				}
			}
			if pkg.TypesInfo.Implicits != nil {
				if p.typesInfo.Implicits == nil {
					p.typesInfo.Implicits = make(map[ast.Node]types.Object)
				}
				for k, v := range pkg.TypesInfo.Implicits {
					p.typesInfo.Implicits[k] = v
				}
			}
			if pkg.TypesInfo.Selections != nil {
				if p.typesInfo.Selections == nil {
					p.typesInfo.Selections = make(map[*ast.SelectorExpr]*types.Selection)
				}
				for k, v := range pkg.TypesInfo.Selections {
					p.typesInfo.Selections[k] = v
				}
			}
			if pkg.TypesInfo.Scopes != nil {
				if p.typesInfo.Scopes == nil {
					p.typesInfo.Scopes = make(map[ast.Node]*types.Scope)
				}
				for k, v := range pkg.TypesInfo.Scopes {
					p.typesInfo.Scopes[k] = v
				}
			}
			// InitOrder はコピーできない

			for _, nextPkg := range pkg.Imports {
				mergePkgTypes(nextPkg)
			}
		}
		for _, pkg := range pkgs {
			mergePkgTypes(pkg)
		}
	}

	result := &Result{}
	for _, pkg := range pkgs {
	file:
		for idx, file := range pkg.Syntax {
			p.hasMetagoBuildTag = false
			p.currentPkg = pkg
			p.currentFile = file
			if len(p.requiredContinueLabel) != 0 {
				panic("unknown state about requiredContinueLabel")
			}
			if len(p.requiredBreakLabel) != 0 {
				panic("unknown state about requiredBreakLabel")
			}

			fileResult := &FileResult{
				Package:  pkg,
				File:     file,
				FilePath: pkg.GoFiles[idx], // TODO これほんとうに安全？
			}
			result.Results = append(result.Results, fileResult)

			useMetagoPackage := astutil.UsesImport(file, metagoPackagePath)

			// . import してたら殺す
			for _, importSpec := range file.Imports {
				importPath, err := strconv.Unquote(importSpec.Path.Value)
				if err != nil {
					return nil, err
				}
				if metagoPackagePath == importPath && importSpec.Name != nil && importSpec.Name.Name == "." {
					p.Errorf(importSpec.Name, "don't use '.' import")
					fileResult.Errors = append(fileResult.Errors, p.nodeErrors...)
					p.nodeErrors = nil
					continue file
				}
			}

			for _, commentGroup := range file.Comments {
				// FileのCommentsは自動では歩かれないので自分でやる
				astutil.Apply(
					commentGroup,
					p.ApplyPre,
					p.ApplyPost,
				)
			}

			astutil.Apply(
				file,
				p.ApplyPre,
				p.ApplyPost,
			)

			// 長さ0のものが残ってると astutil.DeleteNamedImport で panic になる
			var newComments []*ast.CommentGroup
			for _, comments := range file.Comments {
				if len(comments.List) == 0 {
					continue
				}
				newComments = append(newComments, comments)
			}
			file.Comments = newComments

			// clean-up ununsed import
			for _, importSpec := range file.Imports {
				if importSpec.Name != nil && importSpec.Name.Name == "_" {
					continue
				}

				importPath, err := strconv.Unquote(importSpec.Path.Value)
				if err != nil {
					return nil, err
				}
				if !astutil.UsesImport(file, importPath) {
					var importName string
					if importSpec.Name != nil {
						importName = importSpec.Name.Name
					}
					astutil.DeleteNamedImport(pkg.Fset, file, importName, importPath)
				}
			}

			if !p.hasMetagoBuildTag {
				if useMetagoPackage {
					p.Noticef(file, "this file has %s buildtag but doesn't use metago package. ignored", metagoBuildTag)
				}
				fileResult.Errors = append(fileResult.Errors, p.nodeErrors...)
				p.nodeErrors = nil

				continue
			}

			fileResult.Errors = append(fileResult.Errors, p.nodeErrors...)
			p.nodeErrors = nil

			for _, nErr := range fileResult.Errors {
				if nErr.ErrorLevel == ErrorLevelError {
					continue file
				}
			}

			var buf bytes.Buffer
			buf.Write([]byte("// Code generated by metago. DO NOT EDIT.\n\n"))
			buf.Write([]byte(fmt.Sprintf("//+build !%s\n\n", metagoBuildTag)))
			err := format.Node(&buf, pkg.Fset, file)
			if err != nil {
				return nil, err
			}

			fileResult.GeneratedCode = buf.String()
		}
	}

	var nErrs NodeErrors
	for _, fileResult := range result.Results {
		for _, nErr := range fileResult.Errors {
			if nErr.ErrorLevel <= ErrorLevelWarning {
				nErrs = append(nErrs, nErr)
			}
		}
	}
	if len(nErrs) != 0 {
		return result, nErrs
	}

	return result, nil
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
	case *ast.Comment:
		if p.checkMetagoBuildTagComment(cursor, node) {
			return false
		}

	case *ast.AssignStmt:
		if p.checkMetagoValueOfAssignStmt(cursor, node) {
			return true
		}

	case *ast.Ident:
		if p.checkReplaceTargetIdent(cursor, node) {
			return false
		}
		if p.checkUnimportedPackage(cursor, node) {
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
			return false
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

	case *ast.BranchStmt:
		if p.checkInRangeBranchStmt(cursor, node) {
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
	v := p.typesInfo.ObjectOf(ident)
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
	for _, importSpec := range p.currentFile.Imports {
		pkgPath, err := strconv.Unquote(importSpec.Path.Value)
		if err != nil {
			panic(err)
		}
		if pkgPath != metagoPackagePath {
			continue
		}
		if importSpec.Name != nil && importSpec.Name.Name == lhsIdent.Name {
			found = true
			break
		}
		if pkg := p.currentPkg.Imports[pkgPath]; pkg != nil && pkg.Name == lhsIdent.Name {
			found = true
			break
		}
	}
	if !found {
		return false
	}

	if selectorExpr.Sel.Name != valueOfFuncNameString {
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

func (p *metaProcessor) isInlineTemplateFuncDecl(cursor *astutil.Cursor, node *ast.FuncDecl) bool {
	var found bool
outer:
	for _, params := range node.Type.Params.List {
		for _, param := range params.Names {
			if p.isMetagoValue(param) {
				found = true
				break outer
			}
		}
	}
	if !found {
		return false
	}

	if node.Recv != nil {
		p.Errorf(node, "method (function with receiver) can't become inline template")
		return true // 子方向の展開処理されるとコンパイルエラーになりがちなので
	}

	return true
}

func (p *metaProcessor) checkMetagoBuildTagComment(cursor *astutil.Cursor, node *ast.Comment) bool {
	if node.Text != fmt.Sprintf("//+build %s", metagoBuildTag) {
		return false
	}

	p.hasMetagoBuildTag = true
	cursor.Delete()
	return true
}

func (p *metaProcessor) checkReplaceTargetIdent(cursor *astutil.Cursor, node *ast.Ident) bool {
	// mv系の単純な置き換え
	if target := p.valueMapping[node.Obj]; target != nil {
		cursor.Replace(target)
		return true
	}
	// mf系の単純な置き換え
	if fi := p.fieldMapping[node.Obj]; fi != nil {
		cursor.Replace(&ast.SelectorExpr{
			X: fi.recvExpr,
			Sel: &ast.Ident{
				Name: fi.Field().Name(),
			},
		})
		return true
	}

	return false
}

func (p *metaProcessor) checkUnimportedPackage(cursor *astutil.Cursor, node *ast.Ident) bool {
	base := p.copyNodeMap[node]
	if base == nil {
		// コピーされてなければ元のコードにあったものなのでimportは揃ってる
		return false
	}

	// 今見てるIdentがパッケージを指してなかったら気にしない
	obj := p.typesInfo.ObjectOf(base.(*ast.Ident))
	pkgName, ok := obj.(*types.PkgName)
	if !ok {
		return false
	}

	importPkg := pkgName.Imported()
	var importName string
	if importPkg.Name() != node.Name {
		importName = node.Name
	}
	// 現状使ってるか使ってないかを調べずにとりあえず追加する
	astutil.AddNamedImport(p.currentPkg.Fset, p.currentFile, importName, importPkg.Path())

	return true
}

// checkMetagoValueOfAssignStmt is capture `mv := metago.ValueOf(foo)` format assignment.
// it marks up convert rule about `mv` to `foo` and remove these assignment.
func (p *metaProcessor) checkMetagoValueOfAssignStmt(cursor *astutil.Cursor, stmt *ast.AssignStmt) bool {
	// mv := metago.ValueOf(foo) 系を処理する。
	// mv と foo の紐付けを覚える。

	extractMetagoBaseVariable := func(expr ast.Expr) *ast.Ident {
		// metago.ValueOf(foo) から foo を切り出す
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
		targetIdent := extractMetagoBaseVariable(rhs)
		if targetIdent == nil {
			p.Errorf(rhs, "definition of 'metago.Value' variable must be introduced from metago.ValueOf(obj)")
			continue
		}

		found = true
		targetIdent = astcopy.Ident(targetIdent, p.copyNodeMap)
		targetIdent.NamePos = token.NoPos // formatした時に見た目が崩れるのを防ぐ
		p.valueMapping[ident.Obj] = targetIdent

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
	if selectorExpr.Sel.Name != valueFieldsMethodName {
		return false
	}

	// 大本のfor句は全部捨てる必要がある
	cursor.Delete()

	var exprToStructType func(node ast.Expr) *ast.StructType
	var objToStructType func(node *ast.Object) *ast.StructType
	exprToStructType = func(node ast.Expr) *ast.StructType {
		switch node := node.(type) {
		case *ast.Ident:
			return objToStructType(node.Obj)
		case *ast.StarExpr:
			return exprToStructType(node.X)
		case *ast.CallExpr:
			if funcIdent, ok := node.Fun.(*ast.Ident); ok && funcIdent.Name == "new" && funcIdent.Obj == nil {
				// new(Foo) の対応
				return exprToStructType(node.Args[0])
			}
			// TODO 関数呼び出しによる定義の拡充
			return nil
		case *ast.UnaryExpr:
			if node.Op != token.AND {
				panic("unknown op")
			}
			return exprToStructType(node.X)
		case *ast.CompositeLit:
			return exprToStructType(node.Type)
		case *ast.StructType:
			return node
		default:
			panic("unknown type")
		}
	}
	objToStructType = func(node *ast.Object) *ast.StructType {
		switch decl := node.Decl.(type) {
		case *ast.Field:
			return exprToStructType(decl.Type)
		case *ast.AssignStmt:
			if len(decl.Lhs) != 1 || len(decl.Rhs) != 1 {
				p.Errorf(target, "assignment stmt of var definition must be 1 by 1")
				return nil
			}
			return exprToStructType(decl.Rhs[0])
		case *ast.TypeSpec:
			return exprToStructType(decl.Type)
		case *ast.ValueSpec:
			return exprToStructType(decl.Type)
		default:
			panic("unknown type")
		}
	}

	// RangeStmtのBody部分をフィールドの数だけコピペする
	structType := exprToStructType(target)
	if structType == nil {
		p.Errorf(target, "can't extract fields from %s", xIdent.Name)
		return true
	}

	structTypeTypes := p.typesInfo.TypeOf(structType)
	if structTypeTypes == nil {
		return false
	}
	structTypeTypes2, ok := structTypeTypes.(*types.Struct)
	if !ok {
		return false
	}

	for i := 0; i < structTypeTypes2.NumFields(); i++ {
		p.fieldMapping[fieldIdent.Obj] = &fieldInfo{
			recvExpr:   target,
			recvType:   structTypeTypes2,
			fieldIndex: i,
		}

		bodyStmt := astcopy.BlockStmt(node.Body, p.copyNodeMap)
		astutil.Apply(
			bodyStmt,
			p.ApplyPre,
			p.ApplyPost,
		)
		cursor.InsertBefore(bodyStmt)

		for _, labelName := range p.requiredContinueLabel {
			cursor.InsertBefore(&ast.LabeledStmt{
				Label: &ast.Ident{
					Name: labelName,
				},
				Stmt: &ast.EmptyStmt{},
			})
		}
		p.requiredContinueLabel = nil
	}
	delete(p.fieldMapping, fieldIdent.Obj)

	for _, labelName := range p.requiredBreakLabel {
		cursor.InsertAfter(&ast.LabeledStmt{
			Label: &ast.Ident{
				Name: labelName,
			},
			Stmt: &ast.EmptyStmt{},
		})
	}
	p.requiredBreakLabel = nil

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
	fi := p.fieldMapping[callExpr.Fun.(*ast.SelectorExpr).X.(*ast.Ident).Obj]
	if fi == nil {
		p.Errorf(callExpr, "invalid context. not in metago.Field range statement")
		return false
	}
	p.valueMapping[varIdent.Obj] = &ast.SelectorExpr{
		X: fi.recvExpr,
		Sel: &ast.Ident{
			Name: fi.Field().Name(),
		},
	}

	// Cond 部分が静的にboolに評価できるかやってみる
	var condBoolValue bool
	{
		var buf bytes.Buffer
		err := format.Node(&buf, token.NewFileSet(), node.Cond)
		if err != nil {
			panic(err)
		}
		tmpPkg := types.NewPackage("main", "main")
		okVar := types.NewConst(
			token.NoPos,
			tmpPkg,
			okIdent.Name,
			types.Universe.Lookup("true").Type(),
			constant.MakeBool(true),
		)
		tmpPkg.Scope().Insert(okVar)
		ret, err := types.Eval(token.NewFileSet(), tmpPkg, token.NoPos, buf.String())
		if err != nil {
			p.Errorf(node.Cond, "cond can't evaluate statically: %s", err.Error())
			return false
		}
		_ = ret
		retType, ok := ret.Type.(*types.Basic)
		if !ok {
			p.Errorf(node.Cond, "cond is not evaluate to bool")
			return false
		}
		if retType.Kind() != types.UntypedBool {
			p.Errorf(node.Cond, "cond is not evaluate to bool")
			return false
		}
		condBoolValue = constant.BoolVal(ret.Value)
	}

	// fi が存在している == 常にコピーされたコンテキストの中
	targetType := p.typesInfo.TypeOf(p.copyNodeMap[typeAssertExpr.Type].(ast.Expr))

	if condBoolValue == types.AssignableTo(fi.Field().Type(), targetType) {
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
		return false
	}

	callExpr, ok := typeAssertExpr.X.(*ast.CallExpr)
	if !ok {
		return false
	}
	if !p.isCallMetagoFieldValue(callExpr) {
		return false
	}

	fi := p.fieldMapping[callExpr.Fun.(*ast.SelectorExpr).X.(*ast.Ident).Obj]
	if fi == nil {
		p.Errorf(callExpr, "invalid context. not in metago.Field range statement")
		return false
	}

	// fi が存在している == 常にコピーされたコンテキストの中
	targetType := p.typesInfo.TypeOf(p.copyNodeMap[typeAssertExpr.Type].(ast.Expr))

	if !types.AssignableTo(fi.Field().Type(), targetType) {
		cursor.Delete()
		return true // 子をApplyされたくない
	}

	p.replaceNodes[typeAssertExpr] = &ast.SelectorExpr{
		X: fi.recvExpr,
		Sel: &ast.Ident{
			Name: fi.Field().Name(),
		},
	}

	// 呼び出し元でreturn falseするので必要なモノを自分でApplyしてやる必要がある
	astutil.Apply(
		node.Init,
		p.ApplyPre,
		p.ApplyPost,
	)
	astutil.Apply(
		node.Cond,
		p.ApplyPre,
		p.ApplyPost,
	)
	astutil.Apply(
		node.Body,
		p.ApplyPre,
		p.ApplyPost,
	)
	astutil.Apply(
		node.Else,
		p.ApplyPre,
		p.ApplyPost,
	)

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

	fi := p.fieldMapping[callExpr.Fun.(*ast.SelectorExpr).X.(*ast.Ident).Obj]
	if fi == nil {
		p.Errorf(callExpr, "invalid context. not in metago.Field range statement")
		return false
	}

	p.valueMapping[varIdent.Obj] = &ast.SelectorExpr{
		X: fi.recvExpr,
		Sel: &ast.Ident{
			Name: fi.Field().Name(),
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
					// fi が存在している == 常にコピーされたコンテキストの中
					targetType := p.typesInfo.TypeOf(p.copyNodeMap[typeExpr].(ast.Expr))

					if types.AssignableTo(fi.Field().Type(), targetType) {
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

func (p *metaProcessor) checkInRangeBranchStmt(cursor *astutil.Cursor, node *ast.BranchStmt) bool {
	switch node.Tok {
	case token.CONTINUE:
		labelName := fmt.Sprintf("metagoGoto%d", p.gotoCounter)
		p.gotoCounter++
		p.requiredContinueLabel = append(p.requiredContinueLabel, labelName)
		cursor.Replace(&ast.BranchStmt{
			Tok: token.GOTO,
			Label: &ast.Ident{
				Name: labelName,
			},
		})
	case token.BREAK:
		labelName := fmt.Sprintf("metagoGoto%d", p.gotoCounter)
		p.gotoCounter++
		p.requiredBreakLabel = append(p.requiredBreakLabel, labelName)
		cursor.Replace(&ast.BranchStmt{
			Tok: token.GOTO,
			Label: &ast.Ident{
				Name: labelName,
			},
		})
	default:
		return false
	}

	return true
}

func (p *metaProcessor) checkInlineTemplateCallExpr(cursor *astutil.Cursor, node *ast.CallExpr) bool {
	// func fooBarTemplate(mv metago.Value, basic, b string) bool 的なやつの変換
	// 第一引数が metago.Value だったら対象

	// foo(mv) 形式のみ対応 メソッド類は対応が大変
	var funcName *ast.Ident
	switch funcExpr := node.Fun.(type) {
	case *ast.Ident:
		funcName = funcExpr
	case *ast.SelectorExpr:
		funcName = funcExpr.Sel
	case *ast.ArrayType:
		// []byte(foo) とか
		return false
	case *ast.ParenExpr:
		// (*[]byte) とか
		return false
	case *ast.FuncLit:
		if funcExpr.Pos() == token.NoPos {
			// checkInlineTemplateCallExpr で置き換えたヤツ
			return false
		}
		p.Debugf(node, "ignore %T in checkInlineTemplateCallExpr", funcExpr)
		return false
	default:
		p.Debugf(node, "ignore %T in checkInlineTemplateCallExpr", funcExpr)
		return false
	}

	obj := p.typesInfo.ObjectOf(funcName)
	if obj == nil {
		return false
	}

	var findFuncDef func(pkg *packages.Package) (*ast.FuncType, *ast.BlockStmt)
	findFuncDef = func(pkg *packages.Package) (*ast.FuncType, *ast.BlockStmt) {
		for _, file := range pkg.Syntax {
			path, exact := astutil.PathEnclosingInterval(file, obj.Pos(), obj.Pos())
			if exact {
				funcDecl := path[1].(*ast.FuncDecl)
				return funcDecl.Type, funcDecl.Body
			}
		}

		for _, nextPkg := range pkg.Imports {
			funcType, funcBody := findFuncDef(nextPkg)
			if funcType != nil {
				return funcType, funcBody
			}
		}

		return nil, nil
	}

	funcType, funcBody := findFuncDef(p.currentPkg)
	if funcType == nil {
		return false
	}

	// 引数無しは対象外
	if len(funcType.Params.List) == 0 {
		return false
	}
	// 引数の最初が metago.Value じゃないものは対象外
	metaValueArg := funcType.Params.List[0].Names[0]
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

	funcBody = astcopy.BlockStmt(funcBody, p.copyNodeMap)

	// 実引数側の mv がマッピングされる先を 仮引数側の mv にも継承させる
	// *ast.Ident#Obj はコピーされないので metaValueArg 取り直さなくても大丈夫
	p.valueMapping[metaValueArg.Obj] = p.valueMapping[arg.Obj]

	// 引数が metago.Value だけならinline展開するやつ
	// goroutineの境界変わったりするとめんどいので即時実行関数で包む
	newCallExpr := &ast.CallExpr{
		Fun: &ast.FuncLit{
			Type: &ast.FuncType{
				Params: &ast.FieldList{
					List: funcType.Params.List[1:], // 先頭は metago.Valueなので
				},
				Results: &ast.FieldList{
					List: funcType.Results.List,
				},
			},
			Body: funcBody,
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

	fi := p.fieldMapping[objIdent.Obj]
	if fi == nil {
		return false
	}

	if selectorExpr.Sel.Name != fieldValueMethodName {
		return false
	}

	cursor.Replace(&ast.SelectorExpr{
		X: fi.recvExpr,
		Sel: &ast.Ident{
			Name: fi.Field().Name(),
		},
	})

	return false
}

func (p *metaProcessor) checkUseMetagoFieldName(cursor *astutil.Cursor, node *ast.CallExpr) bool {
	// mf.Name() 系を "Foo" 的なのに置き換える
	if !p.isCallMetagoFieldName(node) {
		return false
	}

	fi := p.fieldMapping[node.Fun.(*ast.SelectorExpr).X.(*ast.Ident).Obj]
	if fi == nil {
		return false
	}

	cursor.Replace(&ast.BasicLit{
		Kind:  token.STRING,
		Value: strconv.Quote(fi.Field().Name()),
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

	fi := p.fieldMapping[objIdent.Obj]
	if fi == nil {
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

	targetTag := fi.Tag()
	if targetTag == "" {
		cursor.Replace(&ast.BasicLit{
			Kind:  token.STRING,
			Value: `""`,
		})
		return false
	}

	cursor.Replace(&ast.BasicLit{
		Kind:  token.STRING,
		Value: strconv.Quote(reflect.StructTag(targetTag).Get(tagName)),
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
