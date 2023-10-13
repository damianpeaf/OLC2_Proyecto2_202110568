package visitor

import (
	"log"

	"github.com/antlr4-go/antlr/v4"
	"github.com/damianpeaf/OLC2_Proyecto2_202110568/compiler"
	"github.com/damianpeaf/OLC2_Proyecto2_202110568/ir/abstract"
)

/*
Frame visitor will be in charge of analyzing the "block" of code.
Will return the structure of scope trace.

This includes block statements:
1. if
2. switch
2.1. case
2.2. default
3. while
4. for
5. guard

and the other statements:
declaration
*/
type FrameVisitor struct {
	compiler.BaseTSwiftLanguageVisitor
	ScopeTrace *abstract.ScopeTrace
}

func NewFrameVisitor() *FrameVisitor {
	return &FrameVisitor{
		ScopeTrace: abstract.NewScopeTrace(),
	}
}

func (v *FrameVisitor) VisitStmts(stmts []compiler.IStmtContext) *abstract.ScopeTrace {
	for _, stmt := range stmts {
		v.Visit(stmt)
	}
	return v.ScopeTrace
}

func (v *FrameVisitor) Visit(tree antlr.ParseTree) interface{} {

	switch val := tree.(type) {
	case *antlr.ErrorNodeImpl:
		log.Fatal(val.GetText())
		return nil
	default:
		return tree.Accept(v)
	}

}

func (v *FrameVisitor) VisitStmt(ctx *compiler.StmtContext) interface{} {

	if ctx.Decl_stmt() != nil {
		v.Visit(ctx.Decl_stmt())
	} else if ctx.If_stmt() != nil {
		v.Visit(ctx.If_stmt())
	} else if ctx.Switch_stmt() != nil {
		v.Visit(ctx.Switch_stmt())
	} else if ctx.While_stmt() != nil {
		v.Visit(ctx.While_stmt())
	} else if ctx.For_stmt() != nil {
		v.Visit(ctx.For_stmt())
	} else if ctx.Guard_stmt() != nil {
		v.Visit(ctx.Guard_stmt())
	}

	return nil
}

// * Assignments

func (v *FrameVisitor) VisitTypeValueDecl(ctx *compiler.TypeValueDeclContext) interface{} {

	varName := ctx.ID().GetText()
	varType := ctx.Type_().GetText()

	v.ScopeTrace.NewVariable(varName, varType)

	return nil
}

func (v *FrameVisitor) VisitValueDecl(ctx *compiler.ValueDeclContext) interface{} {
	varName := ctx.ID().GetText()

	v.ScopeTrace.CurrentScope.NewVariable(varName, "")

	return nil
}

func (v *FrameVisitor) VisitTypeDecl(ctx *compiler.TypeDeclContext) interface{} {

	varName := ctx.ID().GetText()
	varType := v.Visit(ctx.Type_()).(string)

	v.ScopeTrace.CurrentScope.NewVariable(varName, varType)
	return nil
}

// * IF
func (v *FrameVisitor) VisitIfStmt(ctx *compiler.IfStmtContext) interface{} {
	for _, child := range ctx.AllIf_chain() {
		v.Visit(child)
	}
	if ctx.Else_stmt() != nil {
		v.Visit(ctx.Else_stmt())
	}
	return nil
}
func (v *FrameVisitor) VisitIfChain(ctx *compiler.IfChainContext) interface{} {
	v.ScopeTrace.PushScope("if")
	for _, child := range ctx.AllStmt() {
		v.Visit(child)
	}
	v.ScopeTrace.PopScope()
	return nil
}

func (v *FrameVisitor) VisitElseStmt(ctx *compiler.ElseStmtContext) interface{} {
	v.ScopeTrace.PushScope("else")
	for _, child := range ctx.AllStmt() {
		v.Visit(child)
	}
	v.ScopeTrace.PopScope()
	return nil
}

// * SWITCH

func (v *FrameVisitor) VisitSwitchStmt(ctx *compiler.SwitchStmtContext) interface{} {
	v.ScopeTrace.PushScope("switch")
	for _, switchCase := range ctx.AllSwitch_case() {
		v.Visit(switchCase)
	}
	if ctx.Default_case() != nil {
		v.Visit(ctx.Default_case())
	}
	v.ScopeTrace.PopScope()
	return nil
}

func (v *FrameVisitor) VisitSwitchCase(ctx *compiler.SwitchCaseContext) interface{} {
	for _, child := range ctx.AllStmt() {
		v.Visit(child)
	}
	return nil
}

func (v *FrameVisitor) VisitDefaultCase(ctx *compiler.DefaultCaseContext) interface{} {
	for _, child := range ctx.AllStmt() {
		v.Visit(child)
	}
	return nil
}

// * WHILE

func (v *FrameVisitor) VisitWhileStmt(ctx *compiler.WhileStmtContext) interface{} {
	v.ScopeTrace.PushScope("while")
	for _, child := range ctx.AllStmt() {
		v.Visit(child)
	}
	v.ScopeTrace.PopScope()
	return nil
}

// * FOR

func (v *FrameVisitor) VisitForStmt(ctx *compiler.ForStmtContext) interface{} {
	v.ScopeTrace.PushScope("outer_for")
	varName := ctx.ID().GetText()
	v.ScopeTrace.CurrentScope.NewVariable(varName, "")

	v.ScopeTrace.PushScope("inner_for")
	for _, child := range ctx.AllStmt() {
		v.Visit(child)
	}
	v.ScopeTrace.PopScope()
	v.ScopeTrace.PopScope()
	return nil
}

// * GUARD

func (v *FrameVisitor) VisitGuardStmt(ctx *compiler.GuardStmtContext) interface{} {
	v.ScopeTrace.PushScope("guard")
	for _, child := range ctx.AllStmt() {
		v.Visit(child)
	}
	v.ScopeTrace.PopScope()
	return nil
}
