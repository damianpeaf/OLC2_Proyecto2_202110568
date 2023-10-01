package visitor

import (
	"strconv"

	"github.com/damianpeaf/OLC2_Proyecto2_202110568/compiler"
	"github.com/damianpeaf/OLC2_Proyecto2_202110568/ir/abstract"
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
			index := v.Factory.NewLiteral().SetValue(strconv.Itoa(argVariableRef.Address))
			stackValue := v.Factory.NewStackIndexed().SetIndex(index)
			assign := v.Factory.NewSimpleAssignment().SetAssignee(temp).SetVal(stackValue)
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
func (v *IrVisitor) VisitParamList(ctx *compiler.ParamListContext) interface{} {

	params := make([]*abstract.Param, 0)

	for _, param := range ctx.AllFunc_param() {
		params = append(params, v.Visit(param).(*abstract.Param))
	}

	return params
}
