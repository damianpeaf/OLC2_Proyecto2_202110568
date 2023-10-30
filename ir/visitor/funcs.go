package visitor

import (
	"fmt"
	"strconv"

	"github.com/damianpeaf/OLC2_Proyecto2_202110568/compiler"
	"github.com/damianpeaf/OLC2_Proyecto2_202110568/ir/abstract"
	"github.com/damianpeaf/OLC2_Proyecto2_202110568/ir/tac"
	"github.com/damianpeaf/OLC2_Proyecto2_202110568/ir/value"
)

// * FUNC CALL

func IsArgValidForStruct(arg []*abstract.Argument) bool {
	for _, a := range arg {

		if a.Name == "" {
			return false
		}
	}
	return true
}

func (v *IrVisitor) VisitFuncCall(ctx *compiler.FuncCallContext) interface{} {

	// find if its a func or constructor of a struct

	canditateName := ctx.Id_pattern().GetText()
	funcObj := v.ScopeTrace.GetFunction(canditateName)
	structObj := v.ScopeTrace.GlobalScope.GetStruct(canditateName)

	args := make([]*abstract.Argument, 0)
	if ctx.Arg_list() != nil {
		args = v.Visit(ctx.Arg_list()).([]*abstract.Argument)
	}

	fmt.Println("func call: ", canditateName)
	fmt.Println("obj: ", funcObj)

	// struct has priority over func
	if structObj != nil {
		if IsArgValidForStruct(args) {
			return v.BuildStruct(canditateName, structObj, args)
		} else {
			return v.GetNilVW()
		}
	}

	if funcObj == nil {
		return v.GetNilVW()
	}

	if funcObj.Type == abstract.BUILTIN_FUNCTION {
		return v.BuiltinHandler(funcObj, args)
	} else if funcObj.Type == abstract.USER_DEFINED_FUNCTION {

		// save the current value of stack pointer
		prevStack := v.Factory.NewTemp()
		v.Factory.AppendToBlock(v.Factory.NewSimpleAssignment().SetAssignee(prevStack).SetVal(v.Factory.NewStackPtr()))

		allocParams := v.Factory.GetBuiltinParams("__alloc_frame")
		size := allocParams[0]
		prevFrame := allocParams[1]

		// set size
		frameSize := len(funcObj.Params) + funcObj.ScopeTrace.Correlative
		v.Factory.AppendToBlock(v.Factory.NewSimpleAssignment().SetAssignee(size).SetVal(v.Factory.NewLiteral().SetValue(strconv.Itoa(frameSize))))

		// set prev frame
		v.Factory.AppendToBlock(v.Factory.NewSimpleAssignment().SetAssignee(prevFrame).SetVal(v.Factory.GetFramePointer()))

		// call alloc frame
		v.Factory.AppendToBlock(v.Factory.NewMethodCall("__alloc_frame"))

		// set frame pointer
		v.Factory.AppendToBlock(v.Factory.NewSimpleAssignment().SetAssignee(v.Factory.GetFramePointer()).SetVal(prevStack))

		// set params
		v.SetParamsOnFrame(funcObj, args)

		// call function
		v.Factory.AppendToBlock(v.Factory.NewMethodCall(funcObj.Name))

		// save the return value on other temp
		returnTemp := v.Factory.NewTemp()
		v.Factory.AppendToBlock(v.Factory.NewSimpleAssignment().SetAssignee(returnTemp).SetVal(funcObj.ReturnTemp))

		/*
			 	summary:
				t1 = P // at this moment P is the next free space on stack (will be occupied by the prev_frame value)
				size = #params + #variables on scope
				alloc_frame(size, t1) // prev_frame | ... | ... | ...
				F = t1
				func()
		*/
		return &value.ValueWrapper{
			Val:      returnTemp,
			Metadata: funcObj.ReturnType,
		}
	}

	return v.GetNilVW()
}

func (v *IrVisitor) SetParamsOnFrame(funcObj *abstract.Function, args []*abstract.Argument) {
	// need to map the params to the args, and save the variable. {name, value}
	argsOk, argsMap := funcObj.ValidateArgs(args)

	if !argsOk {
		return
	}

	// set the scope
	prevScope := v.ScopeTrace.CurrentScope
	v.ScopeTrace.CurrentScope = funcObj.ScopeTrace.GlobalScope
	for _, param := range funcObj.Params {
		argVar := v.ScopeTrace.GetVariable(param.InnerName)
		arg := argsMap[param.InnerName]

		wasPointer := argVar.Pointer
		argVar.Pointer = false // just for the assign

		// assign the value to the variable
		stackAddress := argVar.GetStackStmt(v.Factory)
		assign := v.Factory.NewSimpleAssignment().SetAssignee(stackAddress).SetVal(arg.Wrapper.Val)
		v.Factory.AppendToBlock(assign)

		argVar.Pointer = wasPointer
	}
	v.ScopeTrace.CurrentScope = prevScope
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

	if ctx.ANPERSAND() != nil {
		passByReference = true
	}

	var argValue *value.ValueWrapper = v.GetNilVW()
	var argVariableRef *abstract.IVOR = nil

	if ctx.Id_pattern() != nil {
		// Because is a reference to a variable, the treatment is a bit different
		argName = ctx.Id_pattern().GetText()
		argVariableRef = v.ScopeTrace.GetVariable(argName)

		if argVariableRef != nil {
			temp := v.Factory.NewTemp()
			var stackAddress tac.SimpleValue

			if passByReference {
				stackAddress = argVariableRef.GetStackIndex(v.Factory)
			} else {
				stackAddress = argVariableRef.GetStackStmt(v.Factory)
			}

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

	return &abstract.Argument{
		Name:            argName,
		Wrapper:         argValue,
		PassByReference: passByReference,
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

	frameVisitor := NewFrameVisitor(true, len(params), v.Factory)
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

	prevScope.NewFunction(funcName, function)

	// Al the content will be added to the function block
	funcBlock := make(tac.TACBlock, 0)
	prevBlock := v.Factory.MainBlock
	v.Factory.MainBlock = &funcBlock

	returnLabel := v.Factory.NewLabel()

	prevReturn := v.Transfer.Return
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

		if param.PassByReference {
			paramVar.Pointer = true
		}

		staticScopeTrace.GlobalScope.DirectVariable(paramVar)
	}

	for _, stmt := range ctx.AllStmt() {
		v.Visit(stmt)
	}

	// add return label
	v.Factory.AppendToBlock(returnLabel)

	// return to previous frame
	prevFrameTemp := v.Factory.NewTemp()
	prevFrameStack := v.Factory.NewStackIndexed().SetIndex(v.Factory.GetFramePointer())                                     // stack[fp]
	v.Factory.AppendToBlock(v.Factory.NewSimpleAssignment().SetAssignee(prevFrameTemp).SetVal(prevFrameStack))              // t1 = stack[fp]
	v.Factory.AppendToBlock(v.Factory.NewSimpleAssignment().SetAssignee(v.Factory.GetFramePointer()).SetVal(prevFrameTemp)) // fp = t1

	// life saver
	lifesaver := tac.NewBlockLifesaver(&funcBlock, funcName, v.Factory, returnTemp, staticScopeTrace.Correlative+staticScopeTrace.ParamOffset+1)
	lifesaver.EvalBlock()

	// now we have to add the func tac obj to outer block
	tacFunc := v.Factory.NewMethodDcl(funcBlock).SetName(funcName)
	v.Factory.UserBlock = append(v.Factory.UserBlock, tacFunc)

	v.Transfer.Return = prevReturn
	v.Factory.MainBlock = prevBlock
	v.ScopeTrace.CurrentScope = prevScope
	return nil
}
