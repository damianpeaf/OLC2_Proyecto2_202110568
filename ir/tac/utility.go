package tac

import (
	"strconv"
)

type Utility struct {
	factory *TACFactory
}

func NewUtility(factory *TACFactory) *Utility {
	return &Utility{factory: factory}
}

func (u *Utility) SaveValOnStack(val SimpleValue) int {

	indexLiteral := u.factory.NewStackPtr()                                         //  P
	stackIndexed := u.factory.NewStackIndexed().SetIndex(indexLiteral)              // stack[ (int) P]
	assign := u.factory.NewSimpleAssignment().SetAssignee(stackIndexed).SetVal(val) // stack[ (int) P] = val

	// u.factory.AppendToBlock(u.factory.NewComment().SetComment("Saving value on stack"))
	u.factory.AppendToBlock(assign)

	address := u.factory.StackCurr
	u.IncreaseStackPtr()

	return address
}

func (u *Utility) IncreaseStackPtr() {
	stackPtr := u.factory.NewStackPtr()
	assign := u.factory.NewCompoundAssignment().SetAssignee(stackPtr).SetLeft(stackPtr).SetRight(u.factory.NewLiteral().SetValue("1")).SetOperator(PLUS) // P = P + 1

	// u.factory.AppendToBlock(u.factory.NewComment().SetComment("Increasing stack pointer"))
	u.factory.AppendToBlock(assign)
	u.factory.StackCurr++
}

func (u *Utility) BasicOperation(left SimpleValue, right SimpleValue, operator string) *Temp {

	var endLabel *Label

	// surround division with dynamic check
	if operator == DIVIDE || operator == MOD {
		endLabel = u.factory.NewLabel()
		params := u.factory.GetBuiltinParams("__zero_division")
		param := params[0]
		// assign param to denominator
		assign := u.factory.NewSimpleAssignment().SetAssignee(param).SetVal(right)
		u.factory.AppendToBlock(assign)
		// call builtin
		call := u.factory.NewMethodCall("__zero_division")
		u.factory.AppendToBlock(call)
		// if param is 1, then is an error

		condition := u.factory.NewBoolExpression().SetLeft(param).SetRight(u.factory.NewLiteral().SetValue("1")).SetOp(EQ) // if(t2 == 0)
		conditional := u.factory.NewConditionalJump().SetCondition(condition).SetTarget(endLabel)                          // goto end_print_str
		u.factory.AppendToBlock(conditional)
	}

	temp := u.factory.NewTemp()
	assign := u.factory.NewCompoundAssignment().SetAssignee(temp).SetLeft(left).SetRight(right).SetOperator(operator) // temp = left operator right

	u.factory.AppendToBlock(u.factory.NewComment().SetComment("Arithmetic operation"))
	u.factory.AppendToBlock(assign)

	if operator == DIVIDE || operator == MOD {
		u.factory.AppendToBlock(endLabel)
	}
	return temp
}

func (u *Utility) NilValue() *Literal {
	return u.factory.NewLiteral().SetValue("9999999827968.00")
}

func (u *Utility) BoolOperation(left SimpleValue, right SimpleValue, operator, lcast, rcast string) *Temp {

	temp := u.factory.NewTemp()

	// condition for true
	trueLabel := u.factory.NewLabel()
	endLabel := u.factory.NewLabel()

	condition := u.factory.NewBoolExpression().SetLeft(left).SetRight(right).SetOp(operator).SetLeftCast(lcast).SetRightCast(rcast) // left operator right
	trueJump := u.factory.NewConditionalJump().SetCondition(condition).SetTarget(trueLabel)                                         // goto true
	u.factory.AppendToBlock(trueJump)

	// false
	falseAssign := u.factory.NewSimpleAssignment().SetAssignee(temp).SetVal(u.factory.NewLiteral().SetValue("0")) // temp = 0
	u.factory.AppendToBlock(falseAssign)

	falseJump := u.factory.NewUnconditionalJump().SetTarget(endLabel) // goto end
	u.factory.AppendToBlock(falseJump)

	// true
	u.factory.AppendToBlock(trueLabel)
	trueAssign := u.factory.NewSimpleAssignment().SetAssignee(temp).SetVal(u.factory.NewLiteral().SetValue("1")) // temp = 1
	u.factory.AppendToBlock(trueAssign)

	// end
	u.factory.AppendToBlock(endLabel)

	return temp
}

func (u Utility) AndOperation(left, right SimpleValue) *Temp {

	params := u.factory.GetBuiltinParams("__and")

	// assign the address
	assignS1 := u.factory.NewSimpleAssignment().SetAssignee(params[0]).SetVal(left)
	u.factory.AppendToBlock(assignS1)

	assignS2 := u.factory.NewSimpleAssignment().SetAssignee(params[1]).SetVal(right)
	u.factory.AppendToBlock(assignS2)

	// call builtin
	call := u.factory.NewMethodCall("__and")
	u.factory.AppendToBlock(call)

	return params[2]
}

func (u Utility) OrOperation(left, right SimpleValue) *Temp {
	params := u.factory.GetBuiltinParams("__or")

	// assign the address
	assignS1 := u.factory.NewSimpleAssignment().SetAssignee(params[0]).SetVal(left)
	u.factory.AppendToBlock(assignS1)

	assignS2 := u.factory.NewSimpleAssignment().SetAssignee(params[1]).SetVal(right)
	u.factory.AppendToBlock(assignS2)

	// call builtin
	call := u.factory.NewMethodCall("__or")
	u.factory.AppendToBlock(call)

	return params[2]
}

func (u Utility) NotOperation(left SimpleValue) *Temp {

	params := u.factory.GetBuiltinParams("__not")

	// assign the address
	assignS1 := u.factory.NewSimpleAssignment().SetAssignee(params[0]).SetVal(left)
	u.factory.AppendToBlock(assignS1)

	// call builtin
	call := u.factory.NewMethodCall("__not")
	u.factory.AppendToBlock(call)

	return params[1]
}

func (u *Utility) SaveValOnHeap(val SimpleValue) int {
	indexLiteral := u.factory.NewHeapPtr()                                          // H
	stackIndexed := u.factory.NewHeapIndexed().SetIndex(indexLiteral)               // heap[H(int)]
	assign := u.factory.NewSimpleAssignment().SetAssignee(stackIndexed).SetVal(val) // heap[H(int)] = val

	// u.factory.AppendToBlock(u.factory.NewComment().SetComment("Saving value on heap"))
	u.factory.AppendToBlock(assign)

	address := u.factory.HeapCurr
	u.IncreaseHeapPtr()

	return address
}

func (u *Utility) IncreaseHeapPtr() {
	heapPtr := u.factory.NewHeapPtr()
	assign := u.factory.NewCompoundAssignment().SetAssignee(heapPtr).SetLeft(heapPtr).SetRight(u.factory.NewLiteral().SetValue("1")).SetOperator(PLUS) // P = P + 1

	// u.factory.AppendToBlock(u.factory.NewComment().SetComment("Increasing heap pointer"))
	u.factory.AppendToBlock(assign)
	u.factory.HeapCurr++
}

func (u *Utility) SaveString(stream string) *Temp {
	// temporal for saving the start address of the string
	u.factory.AppendToBlock(u.factory.NewComment().SetComment("---- string literal ----"))
	u.factory.AppendToBlock(u.factory.NewComment().SetComment("start address of the string"))
	startAddressTemporal := u.factory.NewTemp()
	assign := u.factory.NewSimpleAssignment().SetAssignee(startAddressTemporal).SetVal(u.factory.NewHeapPtr())
	u.factory.AppendToBlock(assign)

	// save the string on heap
	u.factory.AppendToBlock(u.factory.NewComment().SetComment("Saving string"))
	for _, char := range stream {
		u.SaveValOnHeap(u.factory.NewLiteral().SetValue(strconv.Itoa(int(char))))
	}
	u.SaveValOnHeap(u.factory.NewLiteral().SetValue("0"))

	return startAddressTemporal
}

// Concats strings on heap and returns the temporal with the start address of the new string
func (u *Utility) ConcatStrings(s1, s2 SimpleValue) *Temp {
	params := u.factory.GetBuiltinParams("__concat_str")

	firstAddress := params[0]
	secondAddress := params[1]
	resultAddress := params[2]

	// assign the address
	assignS1 := u.factory.NewSimpleAssignment().SetAssignee(firstAddress).SetVal(s1)
	u.factory.AppendToBlock(assignS1)

	assignS2 := u.factory.NewSimpleAssignment().SetAssignee(secondAddress).SetVal(s2)
	u.factory.AppendToBlock(assignS2)

	// call builtin
	call := u.factory.NewMethodCall("__concat_str")
	u.factory.AppendToBlock(call)

	return resultAddress
}

func (u *Utility) CompareStrings(s1, s2 SimpleValue, op string) *Temp {

	params := u.factory.GetBuiltinParams("__compare_str")

	// assign the address
	assignS1 := u.factory.NewSimpleAssignment().SetAssignee(params[0]).SetVal(s1)
	u.factory.AppendToBlock(assignS1)

	assignS2 := u.factory.NewSimpleAssignment().SetAssignee(params[1]).SetVal(s2)
	u.factory.AppendToBlock(assignS2)

	// call builtin
	call := u.factory.NewMethodCall("__compare_str")
	u.factory.AppendToBlock(call)

	// code:
	// 0 -> s1 == s2
	// -1 -> s1 < s2
	// 1 -> s1 > s2

	return u.BoolOperation(params[2], u.factory.NewLiteral().SetValue("0"), op, "int", "")
}

func (u *Utility) PrintString(s1 SimpleValue) {
	params := u.factory.GetBuiltinParams("__print_str")

	// assign the address
	assignS1 := u.factory.NewSimpleAssignment().SetAssignee(params[0]).SetVal(s1)
	u.factory.AppendToBlock(assignS1)

	// call builtin
	call := u.factory.NewMethodCall("__print_str")
	u.factory.AppendToBlock(call)
}

func (u *Utility) ConcatCharStrings(s1, s2 SimpleValue, charFirst bool) *Temp {

	u.factory.AppendToBlock(u.factory.NewComment().SetComment("---- converting char to string ----"))
	charHeapAddress := u.factory.NewTemp()
	assign := u.factory.NewSimpleAssignment().SetAssignee(charHeapAddress).SetVal(u.factory.NewHeapPtr())
	u.factory.AppendToBlock(assign)

	// convert char to strings
	if charFirst {
		u.SaveValOnHeap(s1)
	} else {
		u.SaveValOnHeap(s2)
	}
	// null terminated
	u.SaveValOnHeap(u.factory.NewLiteral().SetValue("0"))

	// concat strings
	u.factory.AppendToBlock(u.factory.NewComment().SetComment("---- concatenating strings ----"))
	if charFirst {
		return u.ConcatStrings(charHeapAddress, s2)
	}
	return u.ConcatStrings(s1, charHeapAddress)
}

func (u *Utility) PrintBool(val SimpleValue) {

	params := u.factory.GetBuiltinParams("__print_bool")

	// assign the address
	assignS1 := u.factory.NewSimpleAssignment().SetAssignee(params[0]).SetVal(val)
	u.factory.AppendToBlock(assignS1)

	// call builtin
	call := u.factory.NewMethodCall("__print_bool")
	u.factory.AppendToBlock(call)
}

func (u *Utility) PrintNil() {
	u.factory.AppendToBlock(u.factory.NewPrint().SetMode(PRINT_CHAR).SetVal(u.factory.NewLiteral().SetValue(strconv.Itoa('n'))))
	u.factory.AppendToBlock(u.factory.NewPrint().SetMode(PRINT_CHAR).SetVal(u.factory.NewLiteral().SetValue(strconv.Itoa('i'))))
	u.factory.AppendToBlock(u.factory.NewPrint().SetMode(PRINT_CHAR).SetVal(u.factory.NewLiteral().SetValue(strconv.Itoa('l'))))
}

func (u *Utility) PrintStringStream(stream string) []TACStmtI {
	block := []TACStmtI{}
	for _, char := range stream {
		printStmt := u.factory.NewPrint().SetMode(PRINT_CHAR).SetVal(u.factory.NewLiteral().SetValue(strconv.Itoa(int(char))))
		block = append(block, printStmt)
	}
	return block
}

func (u *Utility) PrintSpace() {
	u.factory.AppendToBlock(u.factory.NewPrint().SetMode(PRINT_CHAR).SetVal(u.factory.NewLiteral().SetValue(strconv.Itoa(' '))))
}

func (u *Utility) PrintNewLine() {
	u.factory.AppendToBlock(u.factory.NewPrint().SetMode(PRINT_CHAR).SetVal(u.factory.NewLiteral().SetValue(strconv.Itoa('\n'))))
}
