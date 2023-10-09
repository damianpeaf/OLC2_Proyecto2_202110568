package repl

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/damianpeaf/OLC2_Proyecto2_202110568/compiler"
	"github.com/damianpeaf/OLC2_Proyecto2_202110568/value"

	"github.com/antlr4-go/antlr/v4"
)

type ReplVisitor struct {
	compiler.BaseTSwiftLanguageVisitor
	ScopeTrace  *ScopeTrace
	CallStack   *CallStack
	Console     *Console
	ErrorTable  *ErrorTable
	StructNames []string
}

func NewVisitor(dclVisitor *DclVisitor) *ReplVisitor {
	return &ReplVisitor{
		ScopeTrace:  dclVisitor.ScopeTrace,
		ErrorTable:  dclVisitor.ErrorTable,
		StructNames: dclVisitor.StructNames,
		CallStack:   NewCallStack(),
		Console:     NewConsole(),
	}
}

func (v *ReplVisitor) GetReplContext() *ReplContext {
	return &ReplContext{
		Console:    v.Console,
		ScopeTrace: v.ScopeTrace,
		CallStack:  v.CallStack,
		ErrorTable: v.ErrorTable,
	}
}

func (v *ReplVisitor) ValidType(_type string) bool {
	return v.ScopeTrace.GlobalScope.ValidType(_type)
}

func (v *ReplVisitor) Visit(tree antlr.ParseTree) interface{} {

	switch val := tree.(type) {
	case *antlr.ErrorNodeImpl:
		log.Fatal(val.GetText())
		return nil
	default:
		return tree.Accept(v)
	}

}

func (v *ReplVisitor) VisitProgram(ctx *compiler.ProgramContext) interface{} {

	for _, stmt := range ctx.AllStmt() {
		v.Visit(stmt)
	}

	return nil
}

func (v *ReplVisitor) VisitStmt(ctx *compiler.StmtContext) interface{} {

	if ctx.Decl_stmt() != nil {
		v.Visit(ctx.Decl_stmt())
	} else if ctx.Assign_stmt() != nil {
		v.Visit(ctx.Assign_stmt())
	} else if ctx.If_stmt() != nil {
		v.Visit(ctx.If_stmt())
	} else if ctx.Switch_stmt() != nil {
		v.Visit(ctx.Switch_stmt())
	} else if ctx.While_stmt() != nil {
		v.Visit(ctx.While_stmt())
	} else if ctx.For_stmt() != nil {
		v.Visit(ctx.For_stmt())
	} else if ctx.Guard_stmt() != nil {
		v.Visit(ctx.Guard_stmt())
	} else if ctx.Transfer_stmt() != nil {
		v.Visit(ctx.Transfer_stmt())
	} else if ctx.Func_call() != nil {
		v.Visit(ctx.Func_call())
	} else if ctx.Func_dcl() != nil {
		v.Visit(ctx.Func_dcl())
	} else if ctx.Strct_dcl() != nil {
		v.Visit(ctx.Strct_dcl())
	} else if ctx.Vector_func() != nil {
		v.Visit(ctx.Vector_func())
	} else {
		log.Fatal("Statement not found " + ctx.GetText())
	}

	return nil
}

func isDeclConst(lexval string) bool {
	return lexval == "let"
}

func (v *ReplVisitor) VisitTypeValueDecl(ctx *compiler.TypeValueDeclContext) interface{} {

	isConst := isDeclConst(ctx.Var_type().GetText())
	varName := ctx.ID().GetText()
	varType := v.Visit(ctx.Type_()).(string)
	varValue := v.Visit(ctx.Expr()).(value.IVOR)

	// copy object
	if obj, ok := varValue.(*ObjectValue); ok {
		varValue = obj.Copy()
	}

	if IsVectorType(varValue.Type()) {
		varValue = varValue.Copy()
	}

	variable, msg := v.ScopeTrace.AddVariable(varName, varType, varValue, isConst, false, ctx.GetStart())

	// Variable already exists
	if variable == nil {
		v.ErrorTable.NewSemanticError(ctx.GetStart(), msg)
	}

	return nil
}

func (v *ReplVisitor) VisitValueDecl(ctx *compiler.ValueDeclContext) interface{} {

	isConst := isDeclConst(ctx.Var_type().GetText())
	varName := ctx.ID().GetText()
	varValue := v.Visit(ctx.Expr()).(value.IVOR)
	varType := varValue.Type()

	if varType == "[]" {
		v.ErrorTable.NewSemanticError(ctx.GetStart(), "No se puede inferir el tipo de un vector vacio '"+varName+"'")
		return nil
	}

	// copy object
	if obj, ok := varValue.(*ObjectValue); ok {
		varValue = obj.Copy()
	}

	if IsVectorType(varValue.Type()) {
		varValue = varValue.Copy()
	}

	variable, msg := v.ScopeTrace.AddVariable(varName, varType, varValue, isConst, false, ctx.GetStart())

	// Variable already exists
	if variable == nil {
		v.ErrorTable.NewSemanticError(ctx.GetStart(), msg)
	}
	return nil
}

func (v *ReplVisitor) VisitTypeDecl(ctx *compiler.TypeDeclContext) interface{} {

	isConst := isDeclConst(ctx.Var_type().GetText())
	varName := ctx.ID().GetText()
	varType := v.Visit(ctx.Type_()).(string)

	if isConst {
		v.ErrorTable.NewSemanticError(ctx.GetStart(), "Las constantes Deben tener un valor asignado")
		return nil
	}

	variable, msg := v.ScopeTrace.AddVariable(varName, varType, value.DefaultNilValue, isConst, true, ctx.GetStart())

	if variable == nil {
		v.ErrorTable.NewSemanticError(ctx.GetStart(), msg)
	}

	return nil
}

func (v *ReplVisitor) VisitVectorItemList(ctx *compiler.VectorItemListContext) interface{} {

	var vectorItems []value.IVOR

	if len(ctx.AllExpr()) == 0 {
		return NewVectorValue(vectorItems, "[]", value.IVOR_ANY)
	}

	for _, item := range ctx.AllExpr() {
		vectorItems = append(vectorItems, v.Visit(item).(value.IVOR))
	}

	var itemType = value.IVOR_NIL

	if ctx.Expr(0) != nil {
		itemType = vectorItems[0].Type()

		for _, item := range vectorItems {
			if item.Type() != itemType {
				v.ErrorTable.NewSemanticError(ctx.GetStart(), "Todos los items de la coleccion deben ser del mismo tipo")
				return value.DefaultNilValue
			}
		}
	}

	_type := "[" + itemType + "]"

	if IsVectorType(_type) {
		return NewVectorValue(vectorItems, _type, itemType)
	}

	if IsMatrixType(_type) {

		if !value.IsPrimitiveType(RemoveBrackets(_type)) {
			v.ErrorTable.NewSemanticError(ctx.GetStart(), "Las matrices deben contener unicamente tipos primitivos")
			return value.DefaultNilValue
		}

		return NewMatrixValue(vectorItems, _type, itemType)
	}

	v.ErrorTable.NewSemanticError(ctx.GetStart(), "Tipo "+_type+" no encontrado")
	return value.DefaultNilValue
}

func (v *ReplVisitor) VisitType(ctx *compiler.TypeContext) interface{} {

	// remove white spaces
	_type := ctx.GetText()

	if v.ValidType(_type) {
		return _type
	}

	if IsVectorType(_type) {
		// remove [ ]
		internType := RemoveBrackets(_type)
		if v.ValidType(internType) {
			return _type
		}

		v.ErrorTable.NewSemanticError(ctx.GetStart(), "El tipo "+internType+" no es valido para un vector")
		return value.IVOR_NIL
	}

	if IsMatrixType(_type) {
		// remove [[]]
		internType := RemoveBrackets(_type)
		if value.IsPrimitiveType(internType) {
			return _type
		}

		v.ErrorTable.NewSemanticError(ctx.GetStart(), "Las matrices solo pueden contener tipos primitivos")
		return value.IVOR_NIL
	}

	v.ErrorTable.NewSemanticError(ctx.GetStart(), "Tipo "+ctx.GetText()+" no encontrado")
	return value.IVOR_NIL
}

func (v *ReplVisitor) VisitVector_type(ctx *compiler.Vector_typeContext) interface{} {
	return ctx.GetText()
}

func (v *ReplVisitor) VisitRepeating(ctx *compiler.RepeatingContext) interface{} {

	if ctx.ID(0).GetText() != "repeating" {
		v.ErrorTable.NewSemanticError(ctx.ID(0).GetSymbol(), "La sintaxis de la función espera el parametro 'repeating'")
		return value.DefaultNilValue
	}

	if ctx.ID(1).GetText() != "count" {
		v.ErrorTable.NewSemanticError(ctx.ID(1).GetSymbol(), "La sintaxis de la función espera el parametro 'count'")
		return value.DefaultNilValue
	}

	reapeating_val := v.Visit(ctx.Expr(0)).(value.IVOR)
	count_val := v.Visit(ctx.Expr(1)).(value.IVOR)

	count, validCount := count_val.(*value.IntValue)
	if !validCount {
		v.ErrorTable.NewSemanticError(ctx.Expr(1).GetStart(), "El parametro count debe ser un entero")
		return value.DefaultNilValue
	}

	if ctx.Vector_type() != nil {
		vector_type := ctx.Vector_type().GetText()
		primitive_type := RemoveBrackets(vector_type)

		if primitive_type != reapeating_val.Type() {
			v.ErrorTable.NewSemanticError(ctx.Expr(0).GetStart(), "El tipo del valor repetido debe ser "+primitive_type)
			return value.DefaultNilValue
		}

		var vectorItems []value.IVOR

		for i := 0; i < count.InternalValue; i++ {
			vectorItems = append(vectorItems, reapeating_val.Copy()) // ? indepedent values
		}

		return NewVectorValue(vectorItems, vector_type, primitive_type)

	} else if ctx.Matrix_type() != nil {

		matrix_type := ctx.Matrix_type().GetText()

		if !(IsMatrixType(reapeating_val.Type()) || IsVectorType(reapeating_val.Type())) {
			v.ErrorTable.NewSemanticError(ctx.Expr(0).GetStart(), "Para crear una matriz con valores repetidos, el valor debe ser un vector o una matriz, se obtuvo '"+reapeating_val.Type()+"'")
			return value.DefaultNilValue
		}

		if !value.IsPrimitiveType(RemoveBrackets(matrix_type)) {
			v.ErrorTable.NewSemanticError(ctx.Expr(0).GetStart(), "Las matrices solo pueden contener tipos primitivos")
			return value.DefaultNilValue
		}

		// must be a lower order collection
		if matrix_type != "["+reapeating_val.Type()+"]" {
			v.ErrorTable.NewSemanticError(ctx.Expr(0).GetStart(), "Para conseguir un valor de tipo '"+matrix_type+"' no es posible contruirlo con un valores repetidos del tipo '"+reapeating_val.Type()+"'")
			return value.DefaultNilValue
		}

		var matrixItems []value.IVOR

		for i := 0; i < count.InternalValue; i++ {
			matrixItems = append(matrixItems, reapeating_val.Copy()) // ? indepedent values
		}

		return NewMatrixValue(matrixItems, matrix_type, RemoveBrackets(reapeating_val.Type()))
	}
	return value.DefaultNilValue
}

func (v *ReplVisitor) VisitRepeatingExp(ctx *compiler.RepeatingExpContext) interface{} {
	return v.Visit(ctx.Repeating())
}

func (v *ReplVisitor) VisitVectorItem(ctx *compiler.VectorItemContext) interface{} {

	varName := ctx.Id_pattern().GetText()

	variable := v.ScopeTrace.GetVariable(varName)

	if variable == nil {
		v.ErrorTable.NewSemanticError(ctx.GetStart(), "Variable "+varName+" no encontrada")
		return nil
	}

	if !(IsVectorType(variable.Type) || IsMatrixType(variable.Type)) {
		v.ErrorTable.NewSemanticError(ctx.GetStart(), "La variable "+varName+" no es un vector o una matriz")
		return nil
	}

	structType := value.IVOR_VECTOR

	index := v.Visit(ctx.Expr(0)).(value.IVOR)

	if len(ctx.AllExpr()) != 1 {
		structType = value.IVOR_MATRIX
	}

	indexes := []int{}

	for _, expr := range ctx.AllExpr() {

		val := v.Visit(expr).(value.IVOR)

		if val.Type() != value.IVOR_INT {
			v.ErrorTable.NewSemanticError(ctx.GetStart(), "Los indices de acceso deben ser enteros")
			return nil
		}

		indexes = append(indexes, val.Value().(int))
	}

	if structType == value.IVOR_VECTOR {

		switch vectorValue := variable.Value.(type) {

		case *VectorValue:
			indexValue := index.(*value.IntValue).InternalValue

			if !vectorValue.ValidIndex(indexValue) {
				v.ErrorTable.NewSemanticError(ctx.GetStart(), "El indice "+strconv.Itoa(indexValue)+" esta fuera de rango")
				return nil
			}

			return &VectorItemReference{
				Vector: vectorValue,
				Index:  indexValue,
				Value:  vectorValue.Get(indexValue),
			}
		default:
			v.ErrorTable.NewSemanticError(ctx.GetStart(), "La variable "+varName+" no es un vector")
			return nil
		}

	} else if structType == value.IVOR_MATRIX {

		switch matrixValue := variable.Value.(type) {

		case *MatrixValue:

			if !matrixValue.ValidIndexes(indexes) {
				v.ErrorTable.NewSemanticError(ctx.GetStart(), "El indice "+fmt.Sprint(indexes)+" esta fuera de rango")
				return nil
			}

			return &MatrixItemReference{
				Matrix: matrixValue,
				Index:  indexes,
				Value:  matrixValue.Get(indexes),
			}
		default:
			v.ErrorTable.NewSemanticError(ctx.GetStart(), "La variable "+varName+" no es una matriz")
			return nil
		}

	} else {
		log.Fatal("Invalid struct type")
	}
	return nil
}

func (v *ReplVisitor) VisitDirectAssign(ctx *compiler.DirectAssignContext) interface{} {

	varName := v.Visit(ctx.Id_pattern()).(string)
	varValue := v.Visit(ctx.Expr()).(value.IVOR)

	variable := v.ScopeTrace.GetVariable(varName)

	if variable == nil {
		v.ErrorTable.NewSemanticError(ctx.GetStart(), "Variable "+varName+" no encontrada")
	} else {

		// copy object
		if obj, ok := varValue.(*ObjectValue); ok {
			varValue = obj.Copy()
		}

		if IsVectorType(varValue.Type()) {
			varValue = varValue.Copy()
		}

		canMutate := true

		if v.ScopeTrace.CurrentScope.isStruct {
			canMutate = v.ScopeTrace.IsMutatingEnvironment()
		}

		ok, msg := variable.Assign(varValue, canMutate)

		if !ok {
			v.ErrorTable.NewSemanticError(ctx.GetStart(), msg)
		}
	}

	return nil

}

func (v *ReplVisitor) VisitArithmeticAssign(ctx *compiler.ArithmeticAssignContext) interface{} {
	varName := v.Visit(ctx.Id_pattern()).(string)

	variable := v.ScopeTrace.GetVariable(varName)

	if variable == nil {
		v.ErrorTable.NewSemanticError(ctx.GetStart(), "Variable "+varName+" no encontrada")
	} else {

		leftValue := variable.Value
		rightValue := v.Visit(ctx.Expr()).(value.IVOR)

		op := string(ctx.GetOp().GetText()[0])

		strat, ok := BinaryStrats[op]

		if !ok {
			log.Fatal("Binary operator not found")
		}

		ok, msg, varValue := strat.Validate(leftValue, rightValue)

		if !ok {
			v.ErrorTable.NewSemanticError(ctx.GetStart(), msg)
			return nil
		}

		canMutate := true

		if v.ScopeTrace.CurrentScope.isStruct {
			canMutate = v.ScopeTrace.IsMutatingEnvironment()
		}

		ok, msg = variable.Assign(varValue, canMutate)

		if !ok {
			v.ErrorTable.NewSemanticError(ctx.GetStart(), msg)
		}
	}

	return nil
}

func (v *ReplVisitor) VisitVectorAssign(ctx *compiler.VectorAssignContext) interface{} {

	rightValue := v.Visit(ctx.Expr()).(value.IVOR)

	switch itemRef := v.Visit(ctx.Vector_item()).(type) {
	case *VectorItemReference:

		leftValue := itemRef.Value

		// check type, improve cast -> ¿? idk what i was thinking
		if rightValue.Type() != itemRef.Vector.ItemType {
			v.ErrorTable.NewSemanticError(ctx.GetStart(), "No se puede asignar un valor de tipo "+rightValue.Type()+" a un vector de tipo "+itemRef.Vector.ItemType)
			return nil
		}
		op := string(ctx.GetOp().GetText()[0])

		if op == "=" {
			itemRef.Vector.InternalValue[itemRef.Index] = rightValue
			return nil
		}

		strat, ok := BinaryStrats[op]

		if !ok {
			log.Fatal("Binary operator not found")
		}

		ok, msg, varValue := strat.Validate(leftValue, rightValue)

		if !ok {
			v.ErrorTable.NewSemanticError(ctx.GetStart(), msg)
			return nil
		}

		itemRef.Vector.InternalValue[itemRef.Index] = varValue

		return nil
	case *MatrixItemReference:
		leftValue := itemRef.Value

		// check type, improve cast -> ¿? idk what i was thinking
		if rightValue.Type() != RemoveBrackets(itemRef.Matrix.Type()) {
			v.ErrorTable.NewSemanticError(ctx.GetStart(), "No se puede asignar un valor de tipo "+rightValue.Type()+" a una matriz de tipo "+RemoveBrackets(itemRef.Matrix.Type()))
			return nil
		}

		op := string(ctx.GetOp().GetText()[0])

		if op == "=" {
			itemRef.Matrix.Set(itemRef.Index, rightValue)
			return nil
		}

		strat, ok := BinaryStrats[op]

		if !ok {
			log.Fatal("Binary operator not found")
		}

		ok, msg, varValue := strat.Validate(leftValue, rightValue)

		if !ok {
			v.ErrorTable.NewSemanticError(ctx.GetStart(), msg)
			return nil
		}

		itemRef.Matrix.Set(itemRef.Index, varValue)
		return nil
	}

	return nil
}

func (v *ReplVisitor) VisitIdPattern(ctx *compiler.IdPatternContext) interface{} {
	return ctx.GetText()
}

func (v *ReplVisitor) VisitIntLiteral(ctx *compiler.IntLiteralContext) interface{} {

	intVal, _ := strconv.Atoi(ctx.GetText())

	return &value.IntValue{
		InternalValue: intVal,
	}

}

func (v *ReplVisitor) VisitFloatLiteral(ctx *compiler.FloatLiteralContext) interface{} {

	floatVal, _ := strconv.ParseFloat(ctx.GetText(), 64)

	return &value.FloatValue{
		InternalValue: floatVal,
	}

}

func (v *ReplVisitor) VisitStringLiteral(ctx *compiler.StringLiteralContext) interface{} {

	// remove quotes
	stringVal := ctx.GetText()[1 : len(ctx.GetText())-1]

	// \" \\ \n \r \
	stringVal = strings.ReplaceAll(stringVal, "\\\"", "\"")
	stringVal = strings.ReplaceAll(stringVal, "\\\\", "\\")
	stringVal = strings.ReplaceAll(stringVal, "\\n", "\n")
	stringVal = strings.ReplaceAll(stringVal, "\\r", "\r")

	// Character literal
	if len(stringVal) == 1 {
		return &value.CharacterValue{
			InternalValue: stringVal,
		}
	}

	// String literal
	return &value.StringValue{
		InternalValue: stringVal,
	}

}

func (v *ReplVisitor) VisitBoolLiteral(ctx *compiler.BoolLiteralContext) interface{} {

	boolVal, _ := strconv.ParseBool(ctx.GetText())

	return &value.BoolValue{
		InternalValue: boolVal,
	}

}

func (v *ReplVisitor) VisitNilLiteral(ctx *compiler.NilLiteralContext) interface{} {
	return value.DefaultNilValue
}

func (v *ReplVisitor) VisitLiteralExp(ctx *compiler.LiteralExpContext) interface{} {
	return v.Visit(ctx.Literal())
}

func (v *ReplVisitor) VisitIdExp(ctx *compiler.IdExpContext) interface{} {
	varName := ctx.Id_pattern().GetText()

	variable := v.ScopeTrace.GetVariable(varName)

	if variable == nil {
		v.ErrorTable.NewSemanticError(ctx.GetStart(), "Variable "+varName+" no encontrada")
		return value.DefaultNilValue
	}

	// ? pointer
	return variable.Value
}

func (v *ReplVisitor) VisitParenExp(ctx *compiler.ParenExpContext) interface{} {
	return v.Visit(ctx.Expr())
}

func (v *ReplVisitor) VisitVectorItemExp(ctx *compiler.VectorItemExpContext) interface{} {

	switch itemRef := v.Visit(ctx.Vector_item()).(type) {
	case *VectorItemReference:
		return itemRef.Value
	case *MatrixItemReference:
		return itemRef.Value
	}
	return value.DefaultNilValue
}

func (v *ReplVisitor) VisitFuncCallExp(ctx *compiler.FuncCallExpContext) interface{} {
	return v.Visit(ctx.Func_call())
}

func (v *ReplVisitor) VisitVectorExp(ctx *compiler.VectorExpContext) interface{} {
	return v.Visit(ctx.Vector_expr())
}

func (v *ReplVisitor) VisitUnaryExp(ctx *compiler.UnaryExpContext) interface{} {

	exp := v.Visit(ctx.Expr()).(value.IVOR)

	strat, ok := UnaryStrats[ctx.GetOp().GetText()]

	if !ok {
		log.Fatal("Unary operator not found")
	}

	ok, msg, result := strat.Validate(exp)

	if !ok {
		v.ErrorTable.NewSemanticError(ctx.GetOp(), msg)
		return value.DefaultNilValue
	}

	return result

}

func (v *ReplVisitor) VisitBinaryExp(ctx *compiler.BinaryExpContext) interface{} {

	op := ctx.GetOp().GetText()
	left := v.Visit(ctx.GetLeft()).(value.IVOR)

	earlyCheck, ok := EarlyReturnStrats[op]

	if ok {
		ok, _, result := earlyCheck.Validate(left)

		if ok {
			return result
		}
	}

	right := v.Visit(ctx.GetRight()).(value.IVOR)

	strat, ok := BinaryStrats[op]

	if !ok {
		log.Fatal("Binary operator not found")
	}

	ok, msg, result := strat.Validate(left, right)

	if !ok {
		v.ErrorTable.NewSemanticError(ctx.GetOp(), msg)
		return value.DefaultNilValue
	}

	return result
}

func (v *ReplVisitor) VisitIfStmt(ctx *compiler.IfStmtContext) interface{} {

	runChain := true

	for _, ifStmt := range ctx.AllIf_chain() {

		runChain = !v.Visit(ifStmt).(bool)

		if !runChain {
			break
		}
	}

	if runChain && ctx.Else_stmt() != nil {
		v.Visit(ctx.Else_stmt())
	}

	return nil
}

func (v *ReplVisitor) VisitIfChain(ctx *compiler.IfChainContext) interface{} {

	condition := v.Visit(ctx.Expr()).(value.IVOR)

	if condition.Type() != value.IVOR_BOOL {
		v.ErrorTable.NewSemanticError(ctx.GetStart(), "La condicion del if debe ser un booleano")
		return false

	}

	if condition.(*value.BoolValue).InternalValue {

		// Push scope
		v.ScopeTrace.PushScope("if")

		for _, stmt := range ctx.AllStmt() {
			v.Visit(stmt)
		}

		// Pop scope
		v.ScopeTrace.PopScope()

		return true
	}

	return false
}

func (v *ReplVisitor) VisitElseStmt(ctx *compiler.ElseStmtContext) interface{} {

	// Push scope
	v.ScopeTrace.PushScope("else")

	for _, stmt := range ctx.AllStmt() {
		v.Visit(stmt)
	}

	// Pop scope
	v.ScopeTrace.PopScope()

	return nil
}

func (v *ReplVisitor) VisitSwitchStmt(ctx *compiler.SwitchStmtContext) interface{} {

	mainValue := v.Visit(ctx.Expr()).(value.IVOR)

	v.ScopeTrace.PushScope("switch")

	// Push break switchItem to call stack [breakable]
	switchItem := &CallStackItem{
		ReturnValue: value.DefaultNilValue,
		Type: []string{
			BreakItem,
		},
	}

	v.CallStack.Push(switchItem)

	// handle break statements from call stack
	defer func() {

		v.ScopeTrace.PopScope()       // pop switch scope
		v.CallStack.Clean(switchItem) // clean item if it's still in call stack

		if item, ok := recover().(*CallStackItem); item != nil && ok {

			// Not a switch item, propagate panic
			if item != switchItem {
				panic(item)
			}

			return // break
		}
	}()

	visited := false

	// evaluate cases
	for _, switchCase := range ctx.AllSwitch_case() {

		caseValue := v.GetCaseValue(switchCase)

		// ? use binary strat
		if caseValue.Type() != mainValue.Type() {
			// warning
			continue
		}

		if caseValue.Value() == mainValue.Value() {
			v.Visit(switchCase)
			visited = true
			break // implicit break
		}

	}

	// evaluate default
	if ctx.Default_case() != nil && !visited {
		v.Visit(ctx.Default_case())
	}

	return nil
}

func (v *ReplVisitor) GetCaseValue(tree antlr.ParseTree) value.IVOR {

	switch val := tree.(type) {
	case *compiler.SwitchCaseContext:
		return v.Visit(val.Expr()).(value.IVOR)
	default:
		return nil
	}

}

func (v *ReplVisitor) VisitSwitchCase(ctx *compiler.SwitchCaseContext) interface{} {

	// * all cases inside switch case will share the same scope

	for _, stmt := range ctx.AllStmt() {
		v.Visit(stmt)
	}
	return nil
}

func (v *ReplVisitor) VisitDefaultCase(ctx *compiler.DefaultCaseContext) interface{} {
	for _, stmt := range ctx.AllStmt() {
		v.Visit(stmt)
	}
	return nil
}

func (v *ReplVisitor) VisitWhileStmt(ctx *compiler.WhileStmtContext) interface{} {

	condition := v.Visit(ctx.Expr()).(value.IVOR)
	// Push scope
	whileScope := v.ScopeTrace.PushScope("while")

	// Push whileItem to call stack [breakable, continuable]
	whileItem := &CallStackItem{
		ReturnValue: value.DefaultNilValue,
		Type: []string{
			BreakItem,
			ContinueItem,
		},
	}

	v.CallStack.Push(whileItem)

	v.VisitInnerWhile(ctx, condition, whileScope, whileItem)

	v.ScopeTrace.PopScope()      // pop while scope
	v.CallStack.Clean(whileItem) // clean item if it's still in call stack

	return nil
}

func (v *ReplVisitor) VisitInnerWhile(ctx *compiler.WhileStmtContext, condition value.IVOR, whileScope *BaseScope, whileItem *CallStackItem) {

	// ? use binary strat
	if condition.Type() != value.IVOR_BOOL {
		v.ErrorTable.NewSemanticError(ctx.GetStart(), "La condicion del ciclo debe ser un booleano")
		return
	}

	// reset scope
	whileScope.Reset()

	// handle break and continue statements from call stack
	defer func() {

		if item, ok := recover().(*CallStackItem); item != nil && ok {
			v.ScopeTrace.CurrentScope = whileScope // reset scope to while scope
			// Not a while item, propagate panic
			if item != whileItem {
				panic(item)
			}

			// Continue
			if item.IsAction(ContinueItem) {
				item.ResetAction()                                       // reset action, can be used again
				condition = v.Visit(ctx.Expr()).(value.IVOR)             // update condition
				v.VisitInnerWhile(ctx, condition, whileScope, whileItem) // continue
			} else if item.IsAction(BreakItem) {
				// Break
				return
			}
		}
	}()

	for condition.(*value.BoolValue).InternalValue {

		for _, stmt := range ctx.AllStmt() {
			v.Visit(stmt)
		}

		condition = v.Visit(ctx.Expr()).(value.IVOR)

		if condition.Type() != value.IVOR_BOOL {
			v.ErrorTable.NewSemanticError(ctx.GetStart(), "La condicion del ciclo debe ser un booleano")
			return
		}

		// reset scope
		whileScope.Reset()
	}
}

func (v *ReplVisitor) VisitForStmt(ctx *compiler.ForStmtContext) interface{} {

	varName := ctx.ID().GetText()
	var iterableItem *VectorValue = DefaultEmptyVectorValue

	if ctx.Range_() != nil {
		rangeItem, ok := v.Visit(ctx.Range_()).(*VectorValue)

		if !ok {
			v.ErrorTable.NewSemanticError(ctx.GetStart(), "El valor del rango debe ser un vector")
			return nil
		}

		iterableItem = rangeItem
	}

	if ctx.Expr() != nil {
		iterableValue := v.Visit(ctx.Expr()).(value.IVOR)

		if IsVectorType(iterableValue.Type()) {
			iterableItem = iterableValue.(*VectorValue)
		} else if iterableValue.Type() == value.IVOR_STRING {
			iterableItem = StringToVector(iterableValue.(*value.StringValue))
		} else {
			v.ErrorTable.NewSemanticError(ctx.GetStart(), "El valor del rango debe ser un vector o una cadena")
			return nil
		}
	}

	if iterableItem.Size() == 0 {
		return nil
	}

	// Push scope outer scope
	outerForScope := v.ScopeTrace.PushScope("outer_for")

	// create the associated variable to the iterable
	iterableVariable, msg := outerForScope.AddVariable(varName, iterableItem.ItemType, iterableItem.Current(), true, false, ctx.ID().GetSymbol())

	if iterableVariable == nil {
		v.ErrorTable.NewSemanticError(ctx.GetStart(), msg)
		log.Fatal("This should not happen")
		return nil
	}

	// Push forItem to call stack [breakable, continuable]

	forItem := &CallStackItem{
		ReturnValue: value.DefaultNilValue,
		Type: []string{
			BreakItem,
			ContinueItem,
		},
	}

	v.CallStack.Push(forItem)

	// Push inner for scope
	innerForScope := v.ScopeTrace.PushScope("inner_for")

	v.VisitInnerFor(ctx, outerForScope, innerForScope, forItem, iterableItem, iterableVariable)

	iterableItem.Reset()
	v.ScopeTrace.PopScope()    // pop inner for scope
	v.ScopeTrace.PopScope()    // pop outer for scope
	v.CallStack.Clean(forItem) // ? clean item if it's still in call stack

	return nil
}

func (v *ReplVisitor) VisitInnerFor(ctx *compiler.ForStmtContext, outerForScope *BaseScope, innerForScope *BaseScope, forItem *CallStackItem, iterableItem *VectorValue, iterableVariable *Variable) {

	// handle break and continue statements from call stack
	defer func() {

		// reset scope
		innerForScope.Reset()
		if item, ok := recover().(*CallStackItem); item != nil && ok {

			// Not a for item, propagate panic
			if item != forItem {
				panic(item)
			}

			// Continue
			if item.IsAction(ContinueItem) {
				item.ResetAction()                                                                          // reset action, can be used again
				iterableItem.Next()                                                                         // next item
				v.VisitInnerFor(ctx, outerForScope, innerForScope, forItem, iterableItem, iterableVariable) // continue
			}

			// Break
			if item.IsAction(BreakItem) {
				return
			}

		}
	}()

	// iterableItem.Size()
	for iterableItem.CurrentIndex < iterableItem.Size() {

		// update variable value
		iterableVariable.Value = iterableItem.Current()

		for _, stmt := range ctx.AllStmt() {
			v.Visit(stmt)
		}

		iterableItem.Next()
		innerForScope.Reset()
	}
}

func (v *ReplVisitor) VisitNumericRange(ctx *compiler.NumericRangeContext) interface{} {

	leftExpr := v.Visit(ctx.Expr(0)).(value.IVOR)
	rightExpr := v.Visit(ctx.Expr(1)).(value.IVOR)

	if leftExpr.Type() != value.IVOR_INT || rightExpr.Type() != value.IVOR_INT {
		v.ErrorTable.NewSemanticError(ctx.GetStart(), "Los valores de los rangos deben ser enteros")
		return value.DefaultNilValue
	}

	left := leftExpr.(*value.IntValue).InternalValue
	right := rightExpr.(*value.IntValue).InternalValue

	if left > right {
		v.ErrorTable.NewSemanticError(ctx.GetStart(), "El valor izquierdo del rango debe ser menor o igual al valor derecho")
	}

	var values []value.IVOR

	for i := left; i <= right; i++ {
		values = append(values, &value.IntValue{
			InternalValue: i,
		})
	}

	return &VectorValue{
		InternalValue: values,
		CurrentIndex:  0,
		ItemType:      value.IVOR_INT,
	}
}

func (v *ReplVisitor) VisitGuardStmt(ctx *compiler.GuardStmtContext) interface{} {

	condition := v.Visit(ctx.Expr()).(value.IVOR)

	if condition.Type() != value.IVOR_BOOL {
		v.ErrorTable.NewSemanticError(ctx.GetStart(), "La condicion del guard debe ser un booleano")
	}

	if !condition.(*value.BoolValue).InternalValue {

		// Push scope
		v.ScopeTrace.PushScope("guard")

		for _, stmt := range ctx.AllStmt() {
			v.Visit(stmt)
		}

		// Pop scope
		v.ScopeTrace.PopScope()
	}

	return nil
}

func (v *ReplVisitor) VisitReturnStmt(ctx *compiler.ReturnStmtContext) interface{} {

	exits, item := v.CallStack.IsReturnEnv()

	if !exits {
		v.ErrorTable.NewSemanticError(ctx.GetStart(), "La sentencia return debe estar dentro de una funcion")
		return nil
	}

	item.ReturnValue = value.DefaultNilValue
	item.Action = ReturnItem

	if ctx.Expr() != nil {
		item.ReturnValue = v.Visit(ctx.Expr()).(value.IVOR)
	}

	panic(item)
}

func (v *ReplVisitor) VisitBreakStmt(ctx *compiler.BreakStmtContext) interface{} {

	exits, item := v.CallStack.IsBreakEnv()

	if !exits {
		v.ErrorTable.NewSemanticError(ctx.GetStart(), "La sentencia break debe estar dentro de un ciclo o un switch")
		return nil
	}

	item.Action = BreakItem
	panic(item)
}

func (v *ReplVisitor) VisitContinueStmt(ctx *compiler.ContinueStmtContext) interface{} {

	exits, item := v.CallStack.IsContinueEnv()

	if !exits {
		v.ErrorTable.NewSemanticError(ctx.GetStart(), "La sentencia continue debe estar dentro de un ciclo")
		return nil
	}

	item.Action = ContinueItem
	panic(item)
}

func (v *ReplVisitor) VisitFuncCall(ctx *compiler.FuncCallContext) interface{} {

	// find if its a func or constructor of a struct

	canditateName := v.Visit(ctx.Id_pattern()).(string)
	funcObj, msg1 := v.ScopeTrace.GetFunction(canditateName)
	structObj, msg2 := v.ScopeTrace.GlobalScope.GetStruct(canditateName)

	if funcObj == nil && structObj == nil {
		v.ErrorTable.NewSemanticError(ctx.GetStart(), msg1+msg2)
		return value.DefaultNilValue
	}

	args := make([]*Argument, 0)
	if ctx.Arg_list() != nil {
		args = v.Visit(ctx.Arg_list()).([]*Argument)
	}

	// struct has priority over func
	if structObj != nil {
		if IsArgValidForStruct(args) {
			return NewObjectValue(v, canditateName, ctx.Id_pattern().GetStart(), args, false)
		} else {
			v.ErrorTable.NewSemanticError(ctx.GetStart(), "Si bien "+canditateName+" es un struct, no se puede llamar a su constructor con los argumentos especificados. Ni tampoco es una funcion.")
			return value.DefaultNilValue
		}
	}

	switch funcObj := funcObj.(type) {
	case *BuiltInFunction:
		returnValue, ok, msg := funcObj.Exec(v.GetReplContext(), args)

		if !ok {

			if msg != "" {
				v.ErrorTable.NewSemanticError(ctx.GetStart(), msg)
			}

			return value.DefaultNilValue

		}

		return returnValue

	case *Function:
		funcObj.Exec(v, args, ctx.GetStart())
		return funcObj.ReturnValue

	case *ObjectBuiltInFunction:
		funcObj.Exec(v, args, ctx.GetStart())
		return funcObj.ReturnValue

	default:
		log.Fatal("Function type not found")
	}

	return value.DefaultNilValue
}

func (v *ReplVisitor) VisitArgList(ctx *compiler.ArgListContext) interface{} {

	args := make([]*Argument, 0)

	for _, arg := range ctx.AllFunc_arg() {
		args = append(args, v.Visit(arg).(*Argument))
	}

	return args

}

func (v *ReplVisitor) VisitFuncArg(ctx *compiler.FuncArgContext) interface{} {

	argName := ""
	passByReference := false

	var argValue value.IVOR = value.DefaultNilValue
	var argVariableRef *Variable = nil

	if ctx.Id_pattern() != nil {
		// Because is a reference to a variable, the treatment is a bit different
		argName = ctx.Id_pattern().GetText()
		argVariableRef = v.ScopeTrace.GetVariable(argName)

		if argVariableRef != nil {
			argValue = argVariableRef.Value
		} else {
			v.ErrorTable.NewSemanticError(ctx.GetStart(), "Variable "+argName+" no encontrada")
		}
	} else {
		argValue = v.Visit(ctx.Expr()).(value.IVOR)
	}

	if ctx.ID() != nil {
		argName = ctx.ID().GetText()
	}

	if ctx.ANPERSAND() != nil {
		passByReference = true
	}

	return &Argument{
		Name:            argName,
		Value:           argValue,
		PassByReference: passByReference,
		Token:           ctx.GetStart(),
		VariableRef:     argVariableRef,
	}

}

func (v *ReplVisitor) VisitFuncDecl(ctx *compiler.FuncDeclContext) interface{} {

	if v.ScopeTrace.CurrentScope == v.ScopeTrace.GlobalScope {
		// aready declared by dcl_visitor
		return nil
	}

	if v.ScopeTrace.CurrentScope != v.ScopeTrace.GlobalScope && !v.ScopeTrace.CurrentScope.isStruct {
		v.ErrorTable.NewSemanticError(ctx.GetStart(), "Las funciones solo pueden ser declaradas en el scope global o en un struct")
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
		returnType = v.Visit(ctx.Type_()).(string)
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
		return nil
	}

	return function
}

func (v *ReplVisitor) VisitParamList(ctx *compiler.ParamListContext) interface{} {

	params := make([]*Param, 0)

	for _, param := range ctx.AllFunc_param() {
		params = append(params, v.Visit(param).(*Param))
	}

	return params
}

func (v *ReplVisitor) VisitFuncParam(ctx *compiler.FuncParamContext) interface{} {

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

	return &Param{
		ExternName:      externName,
		InnerName:       innerName,
		PassByReference: passByReference,
		Type:            paramType,
		Token:           ctx.GetStart(),
	}

}

// * Structs

func (v *ReplVisitor) VisitStructDecl(ctx *compiler.StructDeclContext) interface{} {
	if v.ScopeTrace.CurrentScope != v.ScopeTrace.GlobalScope {
		v.ErrorTable.NewSemanticError(ctx.GetStart(), "Los structs solo pueden ser declaradas en el scope global")
		return nil
	}

	structAdded, msg := v.ScopeTrace.GlobalScope.AddStruct(ctx.ID().GetText(), &Struct{
		Name:   ctx.ID().GetText(),
		Fields: ctx.AllStruct_prop(),
		Token:  ctx.GetStart(),
	})

	if !structAdded {
		v.ErrorTable.NewSemanticError(ctx.ID().GetSymbol(), msg)
	}

	return nil
}

func (v *ReplVisitor) VisitStructAttr(ctx *compiler.StructAttrContext) interface{} {

	if ctx.Type_() != nil && ctx.Expr() != nil {
		v.ErrorTable.NewSemanticError(ctx.GetStart(), "Los atributos de un struct deben ser declarados con un tipo o con un valor")
		return nil
	}

	isConst := isDeclConst(ctx.Var_type().GetText())
	varName := ctx.ID().GetText()
	var varValue value.IVOR = value.DefaultUnInitializedValue
	explicitType := ""
	implicitType := ""
	finalType := ""

	// value is defined
	if ctx.Expr() != nil {
		varValue = v.Visit(ctx.Expr()).(value.IVOR)
		implicitType = varValue.Type()
	}

	if ctx.Type_() != nil {
		explicitType = v.Visit(ctx.Type_()).(string)
	}

	// explicit type and implicit type are defined
	if explicitType != "" && implicitType != "" {
		if explicitType != implicitType {
			v.ErrorTable.NewSemanticError(ctx.GetStart(), "El tipo explicito y el tipo implicito no coinciden")
			return nil
		}
	}

	// only explicit type is defined
	if explicitType != "" && implicitType == "" {
		finalType = explicitType
	} else {
		// only implicit type is defined
		finalType = implicitType
	}

	variable, msg := v.ScopeTrace.AddVariable(varName, finalType, varValue, isConst, false, ctx.ID().GetSymbol())

	if variable == nil {
		v.ErrorTable.NewSemanticError(ctx.GetStart(), msg)
	}

	return nil
}

func (v *ReplVisitor) VisitStructFunc(ctx *compiler.StructFuncContext) interface{} {

	funcDcl := v.Visit(ctx.Func_dcl())

	if ctx.MUTATING_KW() != nil {
		structFunc, ok := funcDcl.(*Function)

		if !ok {
			return nil
		}
		structFunc.IsMutating = true
	}

	return nil
}

func (v *ReplVisitor) VisitStructVector(ctx *compiler.StructVectorContext) interface{} {

	_type := ctx.ID().GetText()

	stc, msg := v.ScopeTrace.GlobalScope.GetStruct(_type)

	if stc == nil {
		v.ErrorTable.NewSemanticError(ctx.GetStart(), msg)
		return value.DefaultNilValue
	}

	return NewVectorValue(make([]value.IVOR, 0), "["+_type+"]", _type)
}

func (v *ReplVisitor) VisitStructVectorExp(ctx *compiler.StructVectorExpContext) interface{} {
	return v.Visit(ctx.Struct_vector())
}

func (v *ReplVisitor) VisitVectorFuncExp(ctx *compiler.VectorFuncExpContext) interface{} {
	return v.Visit(ctx.Vector_func())
}

func (v *ReplVisitor) VisitVectorPropExp(ctx *compiler.VectorPropExpContext) interface{} {
	return v.Visit(ctx.Vector_prop())
}

func (v *ReplVisitor) VisitVectorProp(ctx *compiler.VectorPropContext) interface{} {

	var objectCandidate value.IVOR

	switch itemRef := v.Visit(ctx.Vector_item()).(type) {
	case *VectorItemReference:
		objectCandidate = itemRef.Value
	case *MatrixItemReference:
		objectCandidate = itemRef.Value
	}

	obj, ok := objectCandidate.(*ObjectValue)

	if !ok {
		v.ErrorTable.NewSemanticError(ctx.GetStart(), "El item del vector no es un struct")
		return value.DefaultNilValue
	}

	lastScope := v.ScopeTrace.CurrentScope
	v.ScopeTrace.CurrentScope = obj.InternalScope

	defer func() {
		v.ScopeTrace.CurrentScope = lastScope
	}()

	variable := v.ScopeTrace.GetVariable(ctx.Id_pattern().GetText())

	if variable == nil {
		v.ErrorTable.NewSemanticError(ctx.GetStart(), "Propiedad "+ctx.Id_pattern().GetText()+" no encontrada item del vector")
		return value.DefaultNilValue
	}

	return variable.Value
}

func (v *ReplVisitor) VisitVectorFunc(ctx *compiler.VectorFuncContext) interface{} {

	var objectCandidate value.IVOR

	switch itemRef := v.Visit(ctx.Vector_item()).(type) {
	case *VectorItemReference:
		objectCandidate = itemRef.Value
	case *MatrixItemReference:
		objectCandidate = itemRef.Value
	}

	obj, ok := objectCandidate.(*ObjectValue)

	if !ok {
		v.ErrorTable.NewSemanticError(ctx.GetStart(), "El objeto no es un struct")
		return value.DefaultNilValue
	}

	lastScope := v.ScopeTrace.CurrentScope
	v.ScopeTrace.CurrentScope = obj.InternalScope

	defer func() {
		v.ScopeTrace.CurrentScope = lastScope
	}()

	return v.Visit(ctx.Func_call())
}

func (s *ScopeTrace) Print() {

	fmt.Println("Global Scope")
	fmt.Println("============")

	fmt.Println("Variables")
	for k, v := range s.GlobalScope.variables {
		fmt.Println(k, v.Value.Value(), v.Type)
	}

	fmt.Println("Funciones")
	for k, v := range s.GlobalScope.functions {
		fmt.Println(k, v)
	}

	fmt.Println("Child Scopes")
	fmt.Println("============")
	fmt.Println("")

	for _, child := range s.GlobalScope.children {

		fmt.Println(child.name)
		fmt.Println("============")

		fmt.Println("Variables")
		for k, v := range child.variables {
			fmt.Println(k, v.Value.Value())
		}

		fmt.Println("Funciones")
		for k, v := range child.functions {
			fmt.Println(k, v)
		}

		fmt.Println("")
	}

}
