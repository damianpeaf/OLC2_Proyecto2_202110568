package visitor

import (
	"strconv"

	"regexp"
	"strings"

	"github.com/damianpeaf/OLC2_Proyecto2_202110568/compiler"
	"github.com/damianpeaf/OLC2_Proyecto2_202110568/ir/abstract"
	"github.com/damianpeaf/OLC2_Proyecto2_202110568/ir/tac"
	"github.com/damianpeaf/OLC2_Proyecto2_202110568/ir/value"
)

type VectorItemReference struct {
	Vector *abstract.IVOR
	Index  *value.ValueWrapper
}

type MatrixItemReference struct {
	Matrix  *abstract.IVOR
	Indexes []*value.ValueWrapper
}

func IsVectorType(_type string) bool {

	// Vector starts with only one [ and ends with only one ]
	vectorPattern := "^\\[.*\\]$"

	// Matrix starts with AT LEAST two [[ and ends with at least two ]]
	matrixPattern := "^\\[\\[.*\\]\\](\\[.*\\]\\])*$"

	// match vector pattern but not matrix pattern

	match, _ := regexp.MatchString(vectorPattern, _type)
	match2, _ := regexp.MatchString(matrixPattern, _type)

	return match && !match2
}

func RemoveBrackets(s string) string {
	return strings.Trim(s, "[]")
}

func IsMatrixType(_type string) bool {

	// Matrix starts with AT LEAST two [[ and ends with at least two ]]
	matrixPattern := "^\\[\\[.*\\]\\](\\[.*\\]\\])*$"

	match, _ := regexp.MatchString(matrixPattern, _type)

	return match
}

func (v *IrVisitor) saveVectorSize(temp *tac.Temp, size int) {
	// - this can be optimized
	heapAddres := v.Factory.NewHeapIndexed().SetIndex(temp)
	v.Factory.AppendToBlock(v.Factory.NewSimpleAssignment().SetAssignee(heapAddres).SetVal(v.Factory.NewLiteral().SetValue(strconv.Itoa(size))))
}

func (v *IrVisitor) saveVectorItems(items []*value.ValueWrapper) {
	for _, item := range items {
		h := v.Factory.NewHeapPtr()
		heapAddres := v.Factory.NewHeapIndexed().SetIndex(h)
		v.Factory.AppendToBlock(v.Factory.NewSimpleAssignment().SetAssignee(heapAddres).SetVal(item.Val))
		v.Utility.IncreaseHeapPtr()
	}
}

func (v *IrVisitor) VisitVectorItemList(ctx *compiler.VectorItemListContext) interface{} {
	// save start direction on heap
	temp := v.Factory.NewTemp()
	v.Factory.AppendToBlock(v.Factory.NewSimpleAssignment().SetAssignee(temp).SetVal(v.Factory.NewHeapPtr()))
	v.Utility.IncreaseHeapPtr()

	var vectorItems []*value.ValueWrapper

	if len(ctx.AllExpr()) == 0 {
		v.saveVectorSize(temp, 0)
		return &value.ValueWrapper{
			Val: temp,
		}
	}

	for _, item := range ctx.AllExpr() {
		vectorItems = append(vectorItems, v.Visit(item).(*value.ValueWrapper))
	}

	// save on heap
	v.saveVectorSize(temp, len(vectorItems))
	v.saveVectorItems(vectorItems)

	var itemType = abstract.IVOR_NIL

	if ctx.Expr(0) != nil {
		itemType = vectorItems[0].Metadata
	}

	_type := "[" + itemType + "]"

	return &value.ValueWrapper{
		Val:      temp,
		Metadata: _type,
	}
}

func (v *IrVisitor) VisitVectorItem(ctx *compiler.VectorItemContext) interface{} {

	varName := ctx.Id_pattern().GetText()
	variable := v.ScopeTrace.GetVariable(varName)

	if variable == nil {
		return v.GetNilVW()
	}

	structType := abstract.IVOR_VECTOR

	index := v.Visit(ctx.Expr(0)).(*value.ValueWrapper)

	if len(ctx.AllExpr()) != 1 {
		structType = abstract.IVOR_MATRIX
	}

	indexes := []*value.ValueWrapper{index}

	for _, expr := range ctx.AllExpr() {

		val := v.Visit(expr).(*value.ValueWrapper)

		if val.Metadata != abstract.IVOR_INT {
			return v.GetNilVW()
		}

		indexes = append(indexes, val)
	}

	if structType == abstract.IVOR_VECTOR && IsVectorType(variable.Type) {
		// TODO  post dynamic check
		return &VectorItemReference{
			Vector: variable,
			Index:  index,
		}
	} else if structType == abstract.IVOR_MATRIX && IsMatrixType(variable.Type) {
		return &MatrixItemReference{
			Matrix:  variable,
			Indexes: indexes,
		}

	}
	return v.GetNilVW()
}

func (v *IrVisitor) VisitVectorProp(ctx *compiler.VectorPropContext) interface{} {
	return v.GetNilVW()
}

func (v *IrVisitor) VisitVectorFunc(ctx *compiler.VectorFuncContext) interface{} {
	return v.GetNilVW()
}

func (v *IrVisitor) VisitVector_type(ctx *compiler.Vector_typeContext) interface{} {
	return ctx.GetText()
}

func (v *IrVisitor) VisitVectorAssign(ctx *compiler.VectorAssignContext) interface{} {
	panic("implement me :(")
	// return v.GetNilVW()
}

func (v *IrVisitor) VisitStructVectorExp(ctx *compiler.StructVectorExpContext) interface{} {
	return v.Visit(ctx.Struct_vector())
}

func (v *IrVisitor) VisitVectorPropExp(ctx *compiler.VectorPropExpContext) interface{} {
	return v.Visit(ctx.Vector_prop())
}

func (v *IrVisitor) VisitVectorItemExp(ctx *compiler.VectorItemExpContext) interface{} {
	switch itemRef := v.Visit(ctx.Vector_item()).(type) {
	case *VectorItemReference:
		// TODO: index dynamic check

		itemTemp := v.Factory.NewTemp()
		relativeAddress := itemRef.Index.Val

		offset := itemRef.Vector.GetStackStmt(v.Factory) // stack[variable]
		offsetTemp := v.Factory.NewTemp()
		v.Factory.AppendToBlock(v.Factory.NewSimpleAssignment().SetAssignee(offsetTemp).SetVal(offset)) // t = stack[variable]
		// increase offset by 1, to skip size
		v.Factory.AppendToBlock(v.Factory.NewCompoundAssignment().SetAssignee(offsetTemp).SetLeft(offsetTemp).SetRight(v.Factory.NewLiteral().SetValue("1")).SetOperator("+"))

		absoluteAddress := v.Factory.NewTemp()
		v.Factory.AppendToBlock(v.Factory.NewCompoundAssignment().SetAssignee(absoluteAddress).SetLeft(relativeAddress).SetRight(offsetTemp).SetOperator("+"))
		// absoluteAddress = relativeAddress + offsetTemp

		heapIndexed := v.Factory.NewHeapIndexed().SetIndex(absoluteAddress)
		v.Factory.AppendToBlock(v.Factory.NewSimpleAssignment().SetAssignee(itemTemp).SetVal(heapIndexed))

		return &value.ValueWrapper{
			Val:      itemTemp,
			Metadata: RemoveBrackets(itemRef.Vector.Type),
		}
	case *MatrixItemReference:
		// TODO: implement
	}

	return v.GetNilVW()
}

func (v *IrVisitor) VisitVectorFuncExp(ctx *compiler.VectorFuncExpContext) interface{} {
	return v.Visit(ctx.Vector_func())
}

func (v *IrVisitor) VisitVectorExp(ctx *compiler.VectorExpContext) interface{} {
	return v.Visit(ctx.Vector_expr())
}

func (v *IrVisitor) VisitStructVector(ctx *compiler.StructVectorContext) interface{} {
	panic("implement structs, and then implement me :(")
	// _type := ctx.ID().GetText()

	// stc, msg := v.ScopeTrace.GlobalScope.GetStruct(_type)

	// if stc == nil {
	// 	v.ErrorTable.NewSemanticError(ctx.GetStart(), msg)
	// 	return value.DefaultNilValue
	// }

	// return NewVectorValue(make([]value.IVOR, 0), "["+_type+"]", _type)
}
