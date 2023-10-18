package visitor

import (
	"github.com/damianpeaf/OLC2_Proyecto2_202110568/compiler"
	"github.com/damianpeaf/OLC2_Proyecto2_202110568/ir/abstract"
	"github.com/damianpeaf/OLC2_Proyecto2_202110568/ir/tac"
	"github.com/damianpeaf/OLC2_Proyecto2_202110568/ir/value"
)

// * FUNC CALL

func (v *IrVisitor) VisitFuncCall(ctx *compiler.FuncCallContext) interface{} {

	// find if its a func or constructor of a struct

	canditateName := ctx.Id_pattern().GetText()
	funcObj := v.ScopeTrace.GetFunction(canditateName)
	// structObj, msg2 := v.ScopeTrace.GlobalScope.GetStruct(canditateName)

	args := make([]*abstract.Argument, 0)
	if ctx.Arg_list() != nil {
		args = v.Visit(ctx.Arg_list()).([]*abstract.Argument)
	}

	// struct has priority over func
	// if structObj != nil {
	// 	if IsArgValidForStruct(args) {
	// 		return NewObjectValue(v, canditateName, ctx.Id_pattern().GetStart(), args, false)
	// 	} else {
	// 		v.ErrorTable.NewSemanticError(ctx.GetStart(), "Si bien "+canditateName+" es un struct, no se puede llamar a su constructor con los argumentos especificados. Ni tampoco es una funcion.")
	// 		return value.DefaultNilValue
	// 	}
	// }

	if funcObj == nil {
		return v.GetNilVW()
	}

	if funcObj.Type == abstract.BUILTIN_FUNCTION {
		return v.BuiltinHandler(funcObj.Name, args)
	}

	return v.GetNilVW()
}

func (v *IrVisitor) VisitArgList(ctx *compiler.ArgListContext) interface{} {

	args := make([]*abstract.Argument, 0)

	for _, arg := range ctx.AllFunc_arg() {
		args = append(args, v.Visit(arg).(*abstract.Argument))
	}

	return args

}

func (v *IrVisitor) VisitFuncArg(ctx *compiler.FuncArgContext) interface{} {

	argName := ""
	passByReference := false

	var argValue *value.ValueWrapper = v.GetNilVW()
	var argVariableRef *abstract.IVOR = nil

	if ctx.Id_pattern() != nil {
		// Because is a reference to a variable, the treatment is a bit different
		argName = ctx.Id_pattern().GetText()
		argVariableRef = v.ScopeTrace.GetVariable(argName)

		if argVariableRef != nil {
			temp := v.Factory.NewTemp()
			stackAddress := argVariableRef.GetStackStmt(v.Factory)
			assign := v.Factory.NewSimpleAssignment().SetAssignee(temp).SetVal(stackAddress)
			v.Factory.AppendToBlock(assign)
			argValue = &value.ValueWrapper{
				Val:      temp,
				Metadata: argVariableRef.Type,
				// ? address
			}
		}
	} else {
		argValue = v.Visit(ctx.Expr()).(*value.ValueWrapper)
	}

	if ctx.ID() != nil {
		argName = ctx.ID().GetText()
	}

	if ctx.ANPERSAND() != nil {
		passByReference = true
	}

	return &abstract.Argument{
		Name:            argName,
		Wrapper:         argValue,
		PassByReference: passByReference,
		//TODO:  VariableRefAddress:     ,
	}

}

// * FUNC DCL

func (v *IrVisitor) VisitFuncParam(ctx *compiler.FuncParamContext) interface{} {

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

	paramType := v.Visit(ctx.Type_()).(string)

	return &abstract.Param{
		ExternName:      externName,
		InnerName:       innerName,
		PassByReference: passByReference,
		Type:            paramType,
	}

}

func (v *IrVisitor) VisitParamList(ctx *compiler.ParamListContext) interface{} {

	params := make([]*abstract.Param, 0)

	for _, param := range ctx.AllFunc_param() {
		params = append(params, v.Visit(param).(*abstract.Param))
	}

	return params
}

func (v *IrVisitor) VisitFuncDecl(ctx *compiler.FuncDeclContext) interface{} {

	if v.ScopeTrace.CurrentScope != v.ScopeTrace.GlobalScope {
		// aready declared by dcl_visitor
		return nil
	}

	funcName := ctx.ID().GetText()

	params := make([]*abstract.Param, 0)

	if ctx.Param_list() != nil {
		params = v.Visit(ctx.Param_list()).([]*abstract.Param)
	}

	returnType := abstract.IVOR_NIL

	if ctx.Type_() != nil {
		returnType = v.Visit(ctx.Type_()).(string)
	}

	frameVisitor := NewFrameVisitor(true, len(params))
	staticScopeTrace := frameVisitor.VisitStmts(ctx.AllStmt())
	prevScope := v.ScopeTrace.CurrentScope

	// link scopes
	v.ScopeTrace.GlobalScope.AddChild(staticScopeTrace.GlobalScope)
	v.ScopeTrace.CurrentScope = staticScopeTrace.GlobalScope
	returnTemp := v.Factory.NewTemp()

	function := &abstract.Function{ // pointer ?
		Name:       funcName,
		Params:     params,
		ReturnType: returnType,
		ScopeTrace: staticScopeTrace,
		Type:       abstract.USER_DEFINED_FUNCTION,
		ReturnTemp: returnTemp,
	}

	v.ScopeTrace.NewFunction(funcName, function)

	// Al the content will be added to the function block
	funcBlock := make(tac.TACBlock, 0)
	prevBlock := v.Factory.MainBlock
	v.Factory.MainBlock = &funcBlock

	returnLabel := v.Factory.NewLabel()

	v.Transfer.Return = &TransferReturn{
		ReturnTemp:  returnTemp,
		ReturnType:  returnType,
		ReturnLabel: returnLabel,
	}

	// add params to scope
	for i, param := range params {
		paramVar := &abstract.IVOR{
			Name:          param.InnerName,
			Type:          param.Type,
			Address:       i,
			FrameRelative: true,
			Offset:        1, // just skip the header
		}
		staticScopeTrace.GlobalScope.DirectVariable(paramVar)
	}

	for _, stmt := range ctx.AllStmt() {
		v.Visit(stmt)
	}

	// add return label
	v.Factory.AppendToBlock(returnLabel)
	// TODO: assign return value to return temp

	// now we have to add the func tac obj to outer block
	tacFunc := v.Factory.NewMethodDcl(funcBlock).SetName(funcName)
	v.Factory.OutBlock = append(v.Factory.OutBlock, tacFunc)

	v.Factory.MainBlock = prevBlock
	v.ScopeTrace.CurrentScope = prevScope
	return nil
}
