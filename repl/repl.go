package repl

type ReplContext struct {
	// The console is the output of the REPL
	Console *Console
	// The scope is the current scope of the REPL
	ScopeTrace *ScopeTrace
	// The call stack is the stack of breakable, continueable and returnable items
	CallStack *CallStack
	// Error table is the table of errors
	ErrorTable *ErrorTable
}
