package tac

import (
	"fmt"
	"strconv"
)

type TACFactory struct {
	LabelCount         int
	TempCount          int
	MainBlock          TACBlock
	OutBlock           TACBlock // TODO
	HeapCurr           int      // ?
	StackCurr          int      // ?
	Utility            *Utility
	RegisteredBuiltins map[string][]*Temp
}

func NewTACFactory() *TACFactory {
	return &TACFactory{0, 0, make(TACBlock, 0), make(TACBlock, 0), 0, 0, nil, make(map[string][]*Temp)}
}

func (f *TACFactory) AppendToBlock(stmt TACStmtI) {
	f.MainBlock = append(f.MainBlock, stmt)
}

func (f *TACFactory) AppendBlock(block TACBlock) {
	f.MainBlock = append(f.MainBlock, block...)
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
	return &HeapPtr{}
}

func (f *TACFactory) NewStackPtr() *StackPtr {
	return &StackPtr{}
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

func (f *TACFactory) registerBuiltins() {

	if f.RegisteredBuiltins["__concat"] != nil {
		fmt.Println("Builtin __concat already registered")
		f.OutBlock = append(f.OutBlock, f.ConcatBuiltIn())
	}

}

func (s *TACFactory) GetBuiltinParams(name string) []*Temp {
	params := s.RegisteredBuiltins[name]

	if params == nil {
		params = s.reserveParams(name)
		s.RegisteredBuiltins[name] = params
	}

	return params
}

func (f *TACFactory) String() string {

	header := "#include <stdio.h>\n" + "float stack[10000];\n" + "float heap[10000];\n" + "float P;\n" + "float H;\n"

	var temps = ""

	f.registerBuiltins()

	for i := 0; i < f.TempCount; i++ {
		if i == 0 {
			temps = "float "
		}
		temps += "t" + strconv.Itoa(i+1)

		if i != f.TempCount-1 {
			temps += ", "
		} else {
			temps += ";\n"
		}
	}
	header += temps

	main_block := "int main() {\n"
	for _, stmt := range f.MainBlock {
		main_block += "\t" + stmt.String() + "\n"
	}
	main_block += "return 0;\n}\n"

	out_block := ""
	for _, stmt := range f.OutBlock {
		out_block += stmt.String() + "\n"
	}

	return header + out_block + main_block
}
