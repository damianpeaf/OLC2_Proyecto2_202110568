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
	Left      SimpleValue
	Right     SimpleValue
	LeftCast  string
	RightCast string
	Operator  string
}

func (c *CompoundAssignment) String() string {
	lcast := ""
	rcast := ""

	if c.LeftCast != "" {
		lcast = "(" + c.LeftCast + ")"
	}

	if c.RightCast != "" {
		rcast = "(" + c.RightCast + ")"
	}

	return c.Assignee.String() + " = " + lcast + " " + c.Left.String() + " " + string(c.Operator) + " " + rcast + " " + c.Right.String() + ";"
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

func (c *CompoundAssignment) SetLeftCast(cast string) *CompoundAssignment {
	c.LeftCast = cast
	return c
}

func (c *CompoundAssignment) SetRightCast(cast string) *CompoundAssignment {
	c.RightCast = cast
	return c
}

// ** SimpleAssignment
type SimpleAssignment struct {
	TACAssigment
	Val  SimpleValue
	cast string
}

func (s *SimpleAssignment) String() string {
	c := ""
	if s.cast != "" {
		c = "(" + s.cast + ")"
	}

	return s.Assignee.String() + " = " + c + " " + s.Val.String() + ";"
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

func (s *SimpleAssignment) SetCast(cast string) *SimpleAssignment {
	s.cast = cast
	return s
}

// ** BoolExpression
type BoolExpression struct {
	Operator  string
	Left      SimpleValue
	Right     SimpleValue
	leftCast  string
	rightCast string
}

func (b *BoolExpression) String() string {

	lcast := ""
	rcast := ""

	if b.leftCast != "" {
		lcast = "(" + b.leftCast + ")"
	}

	if b.rightCast != "" {
		rcast = "(" + b.rightCast + ")"
	}

	return lcast + " " + b.Left.String() + " " + string(b.Operator) + " " + rcast + " " + b.Right.String()
}

// builder utils
func (b *BoolExpression) SetLeft(left SimpleValue) *BoolExpression {
	b.Left = left
	return b
}

func (b *BoolExpression) SetLeftCast(leftCast string) *BoolExpression {
	b.leftCast = leftCast
	return b
}

func (b *BoolExpression) SetRight(right SimpleValue) *BoolExpression {
	b.Right = right
	return b
}

func (b *BoolExpression) SetRightCast(rightCast string) *BoolExpression {
	b.rightCast = rightCast
	return b
}

func (b *BoolExpression) SetOp(operator string) *BoolExpression {
	b.Operator = operator
	return b
}
