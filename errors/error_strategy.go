package errors

import (
	"github.com/antlr4-go/antlr/v4"
)

type CustomErrorStrategy struct {
	*antlr.DefaultErrorStrategy
}

func NewCustomErrorStrategy() *CustomErrorStrategy {
	return &CustomErrorStrategy{
		DefaultErrorStrategy: antlr.NewDefaultErrorStrategy(),
	}
}

// spanish translation of the error message
func (es *CustomErrorStrategy) ReportInputMisMatch(recognizer antlr.Parser, e *antlr.InputMisMatchException) {
	t1 := recognizer.GetTokenStream().LT(-1)
	msg := "Se recibió " + t1.GetText() + ", se esperaba " + es.GetExpectedTokens(recognizer).String()
	recognizer.NotifyErrorListeners(msg, e.GetOffendingToken(), e)
}

// Nota: que pereza traducir todos los mensajes de error
