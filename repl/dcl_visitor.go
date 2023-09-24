package repl

import (
	"github.com/damianpeaf/OLC2_Proyecto2_202110568/compiler"
	"github.com/damianpeaf/OLC2_Proyecto2_202110568/value"
	"log"

	"github.com/antlr4-go/antlr/v4"
)

type DclVisitor struct {
	compiler.BaseTSwiftLanguageVisitor
	ScopeTrace  *ScopeTrace
	ErrorTable  *ErrorTable
	StructNames []string
}

func NewDclVisitor(errorTable *ErrorTable) *DclVisitor {
	return &DclVisitor{
		ScopeTrace:  NewScopeTrace(),
		ErrorTable:  errorTable,
		StructNames: []string{},
	}
}

func (v *DclVisitor) Visit(tree antlr.ParseTree) interface{} {

	switch val := tree.(type) {
	case *antlr.ErrorNodeImpl:
		log.Fatal(val.GetText())
		return nil
	default:
		return tree.Accept(v)
	}

}

func (v *DclVisitor) VisitProgram(ctx *compiler.ProgramContext) interface{} {

	for _, stmt := range ctx.AllStmt() {
		v.Visit(stmt)
	}

	return nil
}

func (v *DclVisitor) VisitStmt(ctx *compiler.StmtContext) interface{} {

	if ctx.Func_dcl() != nil {
		v.Visit(ctx.Func_dcl())
	} else if ctx.Strct_dcl() != nil {
		v.Visit(ctx.Strct_dcl())
	}

	return nil
}

func (v *DclVisitor) VisitFuncDecl(ctx *compiler.FuncDeclContext) interface{} {

	if v.ScopeTrace.CurrentScope != v.ScopeTrace.GlobalScope {
		v.ErrorTable.NewSemanticError(ctx.GetStart(), "Las funciones solo pueden ser declaradas en el scope global")
	}

	funcName := ctx.ID().GetText()

	params := make([]*Param, 0)

	if ctx.Param_list() != nil {
		params = v.Visit(ctx.Param_list()).([]*Param)
	}

	if len(params) > 0 {

		baseParamType := params[0].ParamType()

		for _, param := range params {
			if param.ParamType() != baseParamType {
				v.ErrorTable.NewSemanticError(param.Token, "Todos los parametros de la funcion deben ser del mismo tipo")
				return nil
			}
		}
	}

	returnType := value.IVOR_NIL
	var returnTypeToken antlr.Token = nil

	if ctx.Type_() != nil {
		returnType = ctx.Type_().GetText()
		returnTypeToken = ctx.Type_().GetStart()
	}

	body := ctx.AllStmt()

	function := &Function{ // pointer ?
		Name:            funcName,
		Param:           params,
		ReturnType:      returnType,
		Body:            body,
		DeclScope:       v.ScopeTrace.CurrentScope,
		ReturnTypeToken: returnTypeToken,
		Token:           ctx.GetStart(),
	}

	ok, msg := v.ScopeTrace.AddFunction(funcName, function)

	if !ok {
		v.ErrorTable.NewSemanticError(ctx.GetStart(), msg)
	}

	return nil
}

func (v *DclVisitor) VisitParamList(ctx *compiler.ParamListContext) interface{} {

	params := make([]*Param, 0)

	for _, param := range ctx.AllFunc_param() {
		params = append(params, v.Visit(param).(*Param))
	}

	return params
}

func (v *DclVisitor) VisitFuncParam(ctx *compiler.FuncParamContext) interface{} {

	externName := ""
	innerName := ""

	// at least ID(0) is defined
	// only 1 ID defined
	if ctx.ID(1) == nil {
		// innerName : type
		// _ : type
		innerName = ctx.ID(0).GetText()
	} else {
		// externName innerName : type
		externName = ctx.ID(0).GetText()
		innerName = ctx.ID(1).GetText()
	}

	passByReference := false

	if ctx.INOUT_KW() != nil {
		passByReference = true
	}

	paramType := ctx.Type_().GetText()

	return &Param{
		ExternName:      externName,
		InnerName:       innerName,
		PassByReference: passByReference,
		Type:            paramType,
		Token:           ctx.GetStart(),
	}

}

func (v *DclVisitor) VisitStructDecl(ctx *compiler.StructDeclContext) interface{} {
	v.StructNames = append(v.StructNames, ctx.ID().GetText())
	return nil
}
