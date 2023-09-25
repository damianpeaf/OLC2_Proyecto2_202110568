package tac

const (
	PLUS     string = "+"
	MINUS    string = "-"
	MULTIPLY string = "*"
	DIVIDE   string = "/"
	MOD      string = "%"
)

const (
	EQ  string = "=="
	NEQ string = "!="
	GT  string = ">"
	GTE string = ">="
	LT  string = "<"
	LTE string = "<="
)

type TACAssigment struct {
	Assignee SimpleValue // just for simplicity
}

// ** CompoundAssignment
type CompoundAssignment struct {
	TACAssigment
	Left     SimpleValue
	Right    SimpleValue
	Operator string
}

func (c *CompoundAssignment) String() string {
	return c.Assignee.String() + " = " + c.Left.String() + " " + string(c.Operator) + " " + c.Right.String() + ";"
}

// builder utils
func (c *CompoundAssignment) SetLeft(left SimpleValue) *CompoundAssignment {
	c.Left = left
	return c
}

func (c *CompoundAssignment) SetRight(right SimpleValue) *CompoundAssignment {
	c.Right = right
	return c
}

func (c *CompoundAssignment) SetOperator(operator string) *CompoundAssignment {
	c.Operator = operator
	return c
}

func (c *CompoundAssignment) SetAssignee(assignee SimpleValue) *CompoundAssignment {
	c.Assignee = assignee
	return c
}

// ** SimpleAssignment
type SimpleAssignment struct {
	TACAssigment
	Val SimpleValue
}

func (s *SimpleAssignment) String() string {
	return s.Assignee.String() + " = " + s.Val.String() + ";"
}

// builder utils

func (s *SimpleAssignment) SetVal(val SimpleValue) *SimpleAssignment {
	s.Val = val
	return s
}

func (s *SimpleAssignment) SetAssignee(assignee SimpleValue) *SimpleAssignment {
	s.Assignee = assignee
	return s
}

// ** BoolExpression
type BoolExpression struct {
	Operator string
	Left     SimpleValue
	Right    SimpleValue
}

func (b *BoolExpression) String() string {
	return b.Left.String() + " " + string(b.Operator) + " " + b.Right.String()
}

// builder utils
func (b *BoolExpression) SetLeft(left SimpleValue) *BoolExpression {
	b.Left = left
	return b
}

func (b *BoolExpression) SetRight(right SimpleValue) *BoolExpression {
	b.Right = right
	return b
}
