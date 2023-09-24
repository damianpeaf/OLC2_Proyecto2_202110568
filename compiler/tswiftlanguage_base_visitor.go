// Code generated from compiler\TSwiftLanguage.g4 by ANTLR 4.13.0. DO NOT EDIT.

package compiler // TSwiftLanguage
import "github.com/antlr4-go/antlr/v4"

type BaseTSwiftLanguageVisitor struct {
	*antlr.BaseParseTreeVisitor
}

func (v *BaseTSwiftLanguageVisitor) VisitProgram(ctx *ProgramContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTSwiftLanguageVisitor) VisitDelimiter(ctx *DelimiterContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTSwiftLanguageVisitor) VisitStmt(ctx *StmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTSwiftLanguageVisitor) VisitTypeValueDecl(ctx *TypeValueDeclContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTSwiftLanguageVisitor) VisitValueDecl(ctx *ValueDeclContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTSwiftLanguageVisitor) VisitTypeDecl(ctx *TypeDeclContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTSwiftLanguageVisitor) VisitVectorItemList(ctx *VectorItemListContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTSwiftLanguageVisitor) VisitVectorItem(ctx *VectorItemContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTSwiftLanguageVisitor) VisitVectorProp(ctx *VectorPropContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTSwiftLanguageVisitor) VisitVectorFunc(ctx *VectorFuncContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTSwiftLanguageVisitor) VisitRepeating(ctx *RepeatingContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTSwiftLanguageVisitor) VisitVar_type(ctx *Var_typeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTSwiftLanguageVisitor) VisitType(ctx *TypeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTSwiftLanguageVisitor) VisitVector_type(ctx *Vector_typeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTSwiftLanguageVisitor) VisitMatrix_type(ctx *Matrix_typeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTSwiftLanguageVisitor) VisitAux_matrix_type(ctx *Aux_matrix_typeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTSwiftLanguageVisitor) VisitDirectAssign(ctx *DirectAssignContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTSwiftLanguageVisitor) VisitArithmeticAssign(ctx *ArithmeticAssignContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTSwiftLanguageVisitor) VisitVectorAssign(ctx *VectorAssignContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTSwiftLanguageVisitor) VisitIdPattern(ctx *IdPatternContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTSwiftLanguageVisitor) VisitIntLiteral(ctx *IntLiteralContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTSwiftLanguageVisitor) VisitFloatLiteral(ctx *FloatLiteralContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTSwiftLanguageVisitor) VisitStringLiteral(ctx *StringLiteralContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTSwiftLanguageVisitor) VisitBoolLiteral(ctx *BoolLiteralContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTSwiftLanguageVisitor) VisitNilLiteral(ctx *NilLiteralContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTSwiftLanguageVisitor) VisitLiteralExp(ctx *LiteralExpContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTSwiftLanguageVisitor) VisitIdExp(ctx *IdExpContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTSwiftLanguageVisitor) VisitStructVectorExp(ctx *StructVectorExpContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTSwiftLanguageVisitor) VisitRepeatingExp(ctx *RepeatingExpContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTSwiftLanguageVisitor) VisitParenExp(ctx *ParenExpContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTSwiftLanguageVisitor) VisitVectorPropExp(ctx *VectorPropExpContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTSwiftLanguageVisitor) VisitVectorItemExp(ctx *VectorItemExpContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTSwiftLanguageVisitor) VisitFuncCallExp(ctx *FuncCallExpContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTSwiftLanguageVisitor) VisitVectorFuncExp(ctx *VectorFuncExpContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTSwiftLanguageVisitor) VisitVectorExp(ctx *VectorExpContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTSwiftLanguageVisitor) VisitUnaryExp(ctx *UnaryExpContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTSwiftLanguageVisitor) VisitBinaryExp(ctx *BinaryExpContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTSwiftLanguageVisitor) VisitIfStmt(ctx *IfStmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTSwiftLanguageVisitor) VisitIfChain(ctx *IfChainContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTSwiftLanguageVisitor) VisitElseStmt(ctx *ElseStmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTSwiftLanguageVisitor) VisitSwitchStmt(ctx *SwitchStmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTSwiftLanguageVisitor) VisitSwitchCase(ctx *SwitchCaseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTSwiftLanguageVisitor) VisitDefaultCase(ctx *DefaultCaseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTSwiftLanguageVisitor) VisitWhileStmt(ctx *WhileStmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTSwiftLanguageVisitor) VisitForStmt(ctx *ForStmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTSwiftLanguageVisitor) VisitNumericRange(ctx *NumericRangeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTSwiftLanguageVisitor) VisitGuardStmt(ctx *GuardStmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTSwiftLanguageVisitor) VisitReturnStmt(ctx *ReturnStmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTSwiftLanguageVisitor) VisitBreakStmt(ctx *BreakStmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTSwiftLanguageVisitor) VisitContinueStmt(ctx *ContinueStmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTSwiftLanguageVisitor) VisitFuncCall(ctx *FuncCallContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTSwiftLanguageVisitor) VisitArgList(ctx *ArgListContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTSwiftLanguageVisitor) VisitFuncArg(ctx *FuncArgContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTSwiftLanguageVisitor) VisitFuncDecl(ctx *FuncDeclContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTSwiftLanguageVisitor) VisitParamList(ctx *ParamListContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTSwiftLanguageVisitor) VisitFuncParam(ctx *FuncParamContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTSwiftLanguageVisitor) VisitStructDecl(ctx *StructDeclContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTSwiftLanguageVisitor) VisitStructAttr(ctx *StructAttrContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTSwiftLanguageVisitor) VisitStructFunc(ctx *StructFuncContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseTSwiftLanguageVisitor) VisitStructVector(ctx *StructVectorContext) interface{} {
	return v.VisitChildren(ctx)
}
