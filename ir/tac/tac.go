package tac

import (
	"strconv"
)

// tac stands for three-address code

type TACStmtI interface {
	String() string
}

type TACBlock = []TACStmtI

// some other tac stmts

// ** Label
type Label struct {
	ID int
}

func (l *Label) String() string {
	return l.LabelName() + ":"
}

func (l *Label) LabelName() string {
	return "L" + strconv.Itoa(l.ID)
}

// ** Print

type PrintModeType string

const (
	PRINT_DIGIT      PrintModeType = "%d"
	PRINT_FLOAT      PrintModeType = "%f"
	PRINT_CHAR       PrintModeType = "%c"
	PRINT_SCIENTIFIC PrintModeType = "%e"
)

type Print struct {
	Val  SimpleValue
	Mode PrintModeType
	Cast string
}

func (p *Print) String() string {
	castStr := ""
	if p.Cast != "" {
		castStr = "(" + p.Cast + ")"
	}

	return "printf(\"" + string(p.Mode) + "\", " + castStr + " " + p.Val.String() + ");"
}

// builder utils
func (p *Print) SetVal(val SimpleValue) *Print {
	p.Val = val
	return p
}

func (p *Print) SetMode(mode PrintModeType) *Print {
	p.Mode = mode
	return p
}

func (p *Print) SetCast(cast string) *Print {
	p.Cast = cast
	return p
}

// ** Comment
type Comment struct {
	Comment string
}

func (c *Comment) String() string {
	return "// " + c.Comment
}

// builder utils
func (c *Comment) SetComment(comment string) *Comment {
	c.Comment = comment
	return c
}
