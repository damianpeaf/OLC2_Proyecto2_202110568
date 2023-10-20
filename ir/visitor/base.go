package visitor

import (
	"log"

	"github.com/antlr4-go/antlr/v4"
	"github.com/damianpeaf/OLC2_Proyecto2_202110568/compiler"
	"github.com/damianpeaf/OLC2_Proyecto2_202110568/ir/abstract"
	"github.com/damianpeaf/OLC2_Proyecto2_202110568/ir/value"
)

func (v *IrVisitor) Visit(tree antlr.ParseTree) interface{} {

	switch val := tree.(type) {
	case *antlr.ErrorNodeImpl:
		log.Fatal(val.GetText())
		return nil
	default:
		return tree.Accept(v)
	}

}

func (v *IrVisitor) VisitProgram(ctx *compiler.ProgramContext) interface{} {
	// main scope static analysis
	frame := NewFrameVisitor(false, 0, v.Factory)
	mainScope := frame.VisitStmts(ctx.AllStmt())
	v.ScopeTrace = mainScope
	v.ReserveVariables()

	for _, stmt := range ctx.AllStmt() {
		v.Visit(stmt)
	}

	return nil
}

func (v *IrVisitor) VisitStmt(ctx *compiler.StmtContext) interface{} {

	if ctx.Decl_stmt() != nil {
		v.Visit(ctx.Decl_stmt())
	} else if ctx.Assign_stmt() != nil {
		v.Visit(ctx.Assign_stmt())
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
	} else if ctx.Transfer_stmt() != nil {
		v.Visit(ctx.Transfer_stmt())
	} else if ctx.Func_call() != nil {
		v.Visit(ctx.Func_call())
	} else if ctx.Func_dcl() != nil {
		v.Visit(ctx.Func_dcl())
	} else if ctx.Strct_dcl() != nil {
		v.Visit(ctx.Strct_dcl())
	} else if ctx.Vector_func() != nil {
		v.Visit(ctx.Vector_func())
	} else {
		log.Fatal("Statement not found " + ctx.GetText())
	}

	return nil
}

func (v *IrVisitor) GetNilVW() *value.ValueWrapper {
	return &value.ValueWrapper{
		Val:      v.Utility.NilValue(),
		Metadata: abstract.IVOR_NIL,
	}
}

// will take the current scope trace and dcl all the variables
func (v *IrVisitor) ReserveVariables() {

	v.Utility.Comment("Reserving variables")
	for i := 0; i < v.ScopeTrace.Correlative; i++ {
		v.Utility.SaveValOnStack(v.Factory.NewLiteral().SetValue("0"))
	}
}
