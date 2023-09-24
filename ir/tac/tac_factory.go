package tac

type TACFactory struct {
	LabelCount int
	TempCount  int
	Block      TACBlock
	heapPtr    *HeapPtr
	stackPtr   *StackPtr
}

func NewTACFactory() *TACFactory {
	hp := &HeapPtr{}
	sp := &StackPtr{}
	return &TACFactory{0, 0, make(TACBlock, 0), hp, sp}
}

func (f *TACFactory) AppendToBlock(stmt TACStmtI) {
	f.Block = append(f.Block, stmt)
}

func (f *TACFactory) NewLabel() *Label {
	f.LabelCount++
	label := &Label{
		ID: f.LabelCount,
	}
	return label
}

func (f *TACFactory) NewTemp() *Temp {
	f.TempCount++
	temp := &Temp{
		ID: f.TempCount,
	}
	return temp
}

func (f *TACFactory) NewMethodDcl() *MethodDcl {
	return &MethodDcl{
		Block: make(TACBlock, 0),
	}
}

func (f *TACFactory) NewMethodCall(name string) *MethodCall {
	return &MethodCall{
		Name: name,
	}
}

func (f *TACFactory) NewCompoundAssignment() *CompoundAssignment {
	return &CompoundAssignment{}
}

func (f *TACFactory) NewSimpleAssignment() *SimpleAssignment {
	return &SimpleAssignment{}
}

func (f *TACFactory) NewBoolExpression() *BoolExpression {
	return &BoolExpression{}
}

func (f *TACFactory) NewConditionalJump() *ConditionalJump {
	return &ConditionalJump{}
}

func (f *TACFactory) NewUnconditionalJump() *UnconditionalJump {
	return &UnconditionalJump{}
}

func (f *TACFactory) NewHeapPtr() *HeapPtr {
	return f.heapPtr
}

func (f *TACFactory) NewStackPtr() *StackPtr {
	return f.stackPtr
}

func (f *TACFactory) NewLiteral() *Literal {
	return &Literal{}
}

func (f *TACFactory) NewHeapIndexed() *HeapIndexedValue {
	return &HeapIndexedValue{}
}

func (f *TACFactory) NewStackIndexed() *StackIndexedValue {
	return &StackIndexedValue{}
}

func (f *TACFactory) NewPrint() *Print {
	return &Print{}
}

func (f *TACFactory) NewComment() *Comment {
	return &Comment{}
}
