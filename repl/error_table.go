package repl

import "github.com/antlr4-go/antlr/v4"

const (
	LexicalError  = "Error léxico"
	SyntaxError   = "Error sintáctico"
	SemanticError = "Error semántico"
	RuntimeError  = "Error en tiempo de ejecución"
)

type Error struct {
	Line   int
	Column int
	Msg    string
	Type   string
}

type ErrorTable struct {
	Errors []Error
}

func (et *ErrorTable) AddError(line int, column int, msg string, errorType string) {
	et.Errors = append(et.Errors, Error{line, column, msg, errorType})
}

func (et *ErrorTable) NewLexicalError(line int, column int, msg string) {
	et.AddError(line, column, msg, LexicalError)
}

func (et *ErrorTable) NewSyntaxError(line int, column int, msg string) {
	et.AddError(line, column, msg, SyntaxError)
}

func (et *ErrorTable) NewSemanticError(token antlr.Token, msg string) {
	et.AddError(token.GetLine(), token.GetColumn(), msg, SemanticError)
}

func (et *ErrorTable) NewRuntimeError(line int, column int, msg string) {
	et.AddError(line, column, msg, RuntimeError)
}

func NewErrorTable() *ErrorTable {
	return &ErrorTable{}
}
