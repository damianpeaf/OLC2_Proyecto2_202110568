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

	u.IncreaseStackPtr()

	return u.factory.StackCurr
}

func (u *Utility) IncreaseStackPtr() {
	stackPtr := u.factory.NewStackPtr()
	assign := u.factory.NewCompoundAssignment().SetAssignee(stackPtr).SetLeft(stackPtr).SetRight(u.factory.NewLiteral().SetValue("1")).SetOperator(PLUS) // P = P + 1

	// u.factory.AppendToBlock(u.factory.NewComment().SetComment("Increasing stack pointer"))
	u.factory.AppendToBlock(assign)
	u.factory.StackCurr++
}

func (u *Utility) ArithmeticOperation(left SimpleValue, right SimpleValue, operator string) *Temp {
	temp := u.factory.NewTemp()
	assign := u.factory.NewCompoundAssignment().SetAssignee(temp).SetLeft(left).SetRight(right).SetOperator(operator) // temp = left operator right

	u.factory.AppendToBlock(u.factory.NewComment().SetComment("Arithmetic operation"))
	u.factory.AppendToBlock(assign)

	return temp
}

func (u *Utility) NilValue() *Literal {
	return u.factory.NewLiteral().SetValue("9999999827968.00")
}

func (u *Utility) SaveValOnHeap(val SimpleValue) int {
	indexLiteral := u.factory.NewHeapPtr()                                          // H
	stackIndexed := u.factory.NewHeapIndexed().SetIndex(indexLiteral)               // heap[H(int)]
	assign := u.factory.NewSimpleAssignment().SetAssignee(stackIndexed).SetVal(val) // heap[H(int)] = val

	// u.factory.AppendToBlock(u.factory.NewComment().SetComment("Saving value on heap"))
	u.factory.AppendToBlock(assign)

	u.IncreaseHeapPtr()

	return u.factory.HeapCurr
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
	params := u.factory.GetBuiltinParams("__concat")

	firstAddress := params[0]
	secondAddress := params[1]
	resultAddress := params[2]

	// assign the address
	assignS1 := u.factory.NewSimpleAssignment().SetAssignee(firstAddress).SetVal(s1)
	u.factory.AppendToBlock(assignS1)

	assignS2 := u.factory.NewSimpleAssignment().SetAssignee(secondAddress).SetVal(s2)
	u.factory.AppendToBlock(assignS2)

	// call builtin
	call := u.factory.NewMethodCall("__concat")
	u.factory.AppendToBlock(call)

	return resultAddress
}
