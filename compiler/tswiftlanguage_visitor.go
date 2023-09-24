// Code generated from compiler\TSwiftLanguage.g4 by ANTLR 4.13.0. DO NOT EDIT.

package compiler // TSwiftLanguage
import "github.com/antlr4-go/antlr/v4"

// A complete Visitor for a parse tree produced by TSwiftLanguage.
type TSwiftLanguageVisitor interface {
	antlr.ParseTreeVisitor

	// Visit a parse tree produced by TSwiftLanguage#program.
	VisitProgram(ctx *ProgramContext) interface{}

	// Visit a parse tree produced by TSwiftLanguage#delimiter.
	VisitDelimiter(ctx *DelimiterContext) interface{}

	// Visit a parse tree produced by TSwiftLanguage#stmt.
	VisitStmt(ctx *StmtContext) interface{}

	// Visit a parse tree produced by TSwiftLanguage#TypeValueDecl.
	VisitTypeValueDecl(ctx *TypeValueDeclContext) interface{}

	// Visit a parse tree produced by TSwiftLanguage#ValueDecl.
	VisitValueDecl(ctx *ValueDeclContext) interface{}

	// Visit a parse tree produced by TSwiftLanguage#TypeDecl.
	VisitTypeDecl(ctx *TypeDeclContext) interface{}

	// Visit a parse tree produced by TSwiftLanguage#VectorItemList.
	VisitVectorItemList(ctx *VectorItemListContext) interface{}

	// Visit a parse tree produced by TSwiftLanguage#VectorItem.
	VisitVectorItem(ctx *VectorItemContext) interface{}

	// Visit a parse tree produced by TSwiftLanguage#VectorProp.
	VisitVectorProp(ctx *VectorPropContext) interface{}

	// Visit a parse tree produced by TSwiftLanguage#VectorFunc.
	VisitVectorFunc(ctx *VectorFuncContext) interface{}

	// Visit a parse tree produced by TSwiftLanguage#repeating.
	VisitRepeating(ctx *RepeatingContext) interface{}

	// Visit a parse tree produced by TSwiftLanguage#var_type.
	VisitVar_type(ctx *Var_typeContext) interface{}

	// Visit a parse tree produced by TSwiftLanguage#type.
	VisitType(ctx *TypeContext) interface{}

	// Visit a parse tree produced by TSwiftLanguage#vector_type.
	VisitVector_type(ctx *Vector_typeContext) interface{}

	// Visit a parse tree produced by TSwiftLanguage#matrix_type.
	VisitMatrix_type(ctx *Matrix_typeContext) interface{}

	// Visit a parse tree produced by TSwiftLanguage#aux_matrix_type.
	VisitAux_matrix_type(ctx *Aux_matrix_typeContext) interface{}

	// Visit a parse tree produced by TSwiftLanguage#DirectAssign.
	VisitDirectAssign(ctx *DirectAssignContext) interface{}

	// Visit a parse tree produced by TSwiftLanguage#ArithmeticAssign.
	VisitArithmeticAssign(ctx *ArithmeticAssignContext) interface{}

	// Visit a parse tree produced by TSwiftLanguage#VectorAssign.
	VisitVectorAssign(ctx *VectorAssignContext) interface{}

	// Visit a parse tree produced by TSwiftLanguage#IdPattern.
	VisitIdPattern(ctx *IdPatternContext) interface{}

	// Visit a parse tree produced by TSwiftLanguage#IntLiteral.
	VisitIntLiteral(ctx *IntLiteralContext) interface{}

	// Visit a parse tree produced by TSwiftLanguage#FloatLiteral.
	VisitFloatLiteral(ctx *FloatLiteralContext) interface{}

	// Visit a parse tree produced by TSwiftLanguage#StringLiteral.
	VisitStringLiteral(ctx *StringLiteralContext) interface{}

	// Visit a parse tree produced by TSwiftLanguage#BoolLiteral.
	VisitBoolLiteral(ctx *BoolLiteralContext) interface{}

	// Visit a parse tree produced by TSwiftLanguage#NilLiteral.
	VisitNilLiteral(ctx *NilLiteralContext) interface{}

	// Visit a parse tree produced by TSwiftLanguage#LiteralExp.
	VisitLiteralExp(ctx *LiteralExpContext) interface{}

	// Visit a parse tree produced by TSwiftLanguage#IdExp.
	VisitIdExp(ctx *IdExpContext) interface{}

	// Visit a parse tree produced by TSwiftLanguage#StructVectorExp.
	VisitStructVectorExp(ctx *StructVectorExpContext) interface{}

	// Visit a parse tree produced by TSwiftLanguage#RepeatingExp.
	VisitRepeatingExp(ctx *RepeatingExpContext) interface{}

	// Visit a parse tree produced by TSwiftLanguage#ParenExp.
	VisitParenExp(ctx *ParenExpContext) interface{}

	// Visit a parse tree produced by TSwiftLanguage#VectorPropExp.
	VisitVectorPropExp(ctx *VectorPropExpContext) interface{}

	// Visit a parse tree produced by TSwiftLanguage#VectorItemExp.
	VisitVectorItemExp(ctx *VectorItemExpContext) interface{}

	// Visit a parse tree produced by TSwiftLanguage#FuncCallExp.
	VisitFuncCallExp(ctx *FuncCallExpContext) interface{}

	// Visit a parse tree produced by TSwiftLanguage#VectorFuncExp.
	VisitVectorFuncExp(ctx *VectorFuncExpContext) interface{}

	// Visit a parse tree produced by TSwiftLanguage#VectorExp.
	VisitVectorExp(ctx *VectorExpContext) interface{}

	// Visit a parse tree produced by TSwiftLanguage#UnaryExp.
	VisitUnaryExp(ctx *UnaryExpContext) interface{}

	// Visit a parse tree produced by TSwiftLanguage#BinaryExp.
	VisitBinaryExp(ctx *BinaryExpContext) interface{}

	// Visit a parse tree produced by TSwiftLanguage#IfStmt.
	VisitIfStmt(ctx *IfStmtContext) interface{}

	// Visit a parse tree produced by TSwiftLanguage#IfChain.
	VisitIfChain(ctx *IfChainContext) interface{}

	// Visit a parse tree produced by TSwiftLanguage#ElseStmt.
	VisitElseStmt(ctx *ElseStmtContext) interface{}

	// Visit a parse tree produced by TSwiftLanguage#SwitchStmt.
	VisitSwitchStmt(ctx *SwitchStmtContext) interface{}

	// Visit a parse tree produced by TSwiftLanguage#SwitchCase.
	VisitSwitchCase(ctx *SwitchCaseContext) interface{}

	// Visit a parse tree produced by TSwiftLanguage#DefaultCase.
	VisitDefaultCase(ctx *DefaultCaseContext) interface{}

	// Visit a parse tree produced by TSwiftLanguage#WhileStmt.
	VisitWhileStmt(ctx *WhileStmtContext) interface{}

	// Visit a parse tree produced by TSwiftLanguage#ForStmt.
	VisitForStmt(ctx *ForStmtContext) interface{}

	// Visit a parse tree produced by TSwiftLanguage#NumericRange.
	VisitNumericRange(ctx *NumericRangeContext) interface{}

	// Visit a parse tree produced by TSwiftLanguage#GuardStmt.
	VisitGuardStmt(ctx *GuardStmtContext) interface{}

	// Visit a parse tree produced by TSwiftLanguage#ReturnStmt.
	VisitReturnStmt(ctx *ReturnStmtContext) interface{}

	// Visit a parse tree produced by TSwiftLanguage#BreakStmt.
	VisitBreakStmt(ctx *BreakStmtContext) interface{}

	// Visit a parse tree produced by TSwiftLanguage#ContinueStmt.
	VisitContinueStmt(ctx *ContinueStmtContext) interface{}

	// Visit a parse tree produced by TSwiftLanguage#FuncCall.
	VisitFuncCall(ctx *FuncCallContext) interface{}

	// Visit a parse tree produced by TSwiftLanguage#ArgList.
	VisitArgList(ctx *ArgListContext) interface{}

	// Visit a parse tree produced by TSwiftLanguage#FuncArg.
	VisitFuncArg(ctx *FuncArgContext) interface{}

	// Visit a parse tree produced by TSwiftLanguage#FuncDecl.
	VisitFuncDecl(ctx *FuncDeclContext) interface{}

	// Visit a parse tree produced by TSwiftLanguage#ParamList.
	VisitParamList(ctx *ParamListContext) interface{}

	// Visit a parse tree produced by TSwiftLanguage#FuncParam.
	VisitFuncParam(ctx *FuncParamContext) interface{}

	// Visit a parse tree produced by TSwiftLanguage#StructDecl.
	VisitStructDecl(ctx *StructDeclContext) interface{}

	// Visit a parse tree produced by TSwiftLanguage#StructAttr.
	VisitStructAttr(ctx *StructAttrContext) interface{}

	// Visit a parse tree produced by TSwiftLanguage#StructFunc.
	VisitStructFunc(ctx *StructFuncContext) interface{}

	// Visit a parse tree produced by TSwiftLanguage#StructVector.
	VisitStructVector(ctx *StructVectorContext) interface{}
}
