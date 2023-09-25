package tac

import "strconv"

type Utility struct {
	factory *TACFactory
}

func NewUtility(factory *TACFactory) *Utility {
	return &Utility{factory: factory}
}

func (u *Utility) SaveValOnStack(val SimpleValue) int {

	indexLiteral := u.factory.NewStackPtr().SetCast(CAST_INT)                       // (int) P
	stackIndexed := u.factory.NewStackIndexed().SetIndex(indexLiteral)              // stack[P]
	assign := u.factory.NewSimpleAssignment().SetAssignee(stackIndexed).SetVal(val) // stack[P] = val

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
	indexLiteral := u.factory.NewHeapPtr().SetCast(CAST_INT)                        // H(int)
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
