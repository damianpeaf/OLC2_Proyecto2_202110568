package visitor

import (
	"fmt"
	"strconv"

	"github.com/damianpeaf/OLC2_Proyecto2_202110568/compiler"
	"github.com/damianpeaf/OLC2_Proyecto2_202110568/ir/abstract"
	"github.com/damianpeaf/OLC2_Proyecto2_202110568/ir/tac"
	"github.com/damianpeaf/OLC2_Proyecto2_202110568/ir/utils"
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

func (v *IrVisitor) saveVectorSize(temp *tac.Temp, size tac.SimpleValue) {
	// - this can be optimized
	v.Factory.AppendToBlock(v.Factory.NewSimpleAssignment().SetAssignee(temp).SetVal(v.Factory.NewHeapPtr()))
	v.Utility.IncreaseHeapPtr()
	heapAddres := v.Factory.NewHeapIndexed().SetIndex(temp)
	v.Factory.AppendToBlock(v.Factory.NewSimpleAssignment().SetAssignee(heapAddres).SetVal(size))
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

	var vectorItems []*value.ValueWrapper
	temp := v.Factory.NewTemp()

	if len(ctx.AllExpr()) == 0 {
		size := v.Factory.NewLiteral().SetValue("0")
		v.saveVectorSize(temp, size)
		return &value.ValueWrapper{
			Val:      temp,
			Metadata: "[]",
		}
	}

	for _, item := range ctx.AllExpr() {
		vectorItems = append(vectorItems, v.Visit(item).(*value.ValueWrapper))
	}

	// save start direction on heap
	size := v.Factory.NewLiteral().SetValue(strconv.Itoa(len(vectorItems)))
	v.saveVectorSize(temp, size)
	// save on heap
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

	indexes := []*value.ValueWrapper{}

	for _, expr := range ctx.AllExpr() {

		val := v.Visit(expr).(*value.ValueWrapper)

		if val.Metadata != abstract.IVOR_INT {
			return v.GetNilVW()
		}

		indexes = append(indexes, val)
	}

	if structType == abstract.IVOR_VECTOR && utils.IsVectorType(variable.Type) {
		// TODO  post dynamic check
		return &VectorItemReference{
			Vector: variable,
			Index:  index,
		}
	} else if structType == abstract.IVOR_MATRIX && utils.IsMatrixType(variable.Type) {
		return &MatrixItemReference{
			Matrix:  variable,
			Indexes: indexes,
		}

	}
	return v.GetNilVW()
}

func (v *IrVisitor) VisitVector_type(ctx *compiler.Vector_typeContext) interface{} {
	return ctx.GetText()
}

func (v *IrVisitor) VisitVectorAssign(ctx *compiler.VectorAssignContext) interface{} {
	// TODO: implement
	panic("implement me :(")
	// return v.GetNilVW()
}

func (v *IrVisitor) VisitVectorItemExp(ctx *compiler.VectorItemExpContext) interface{} {
	switch itemRef := v.Visit(ctx.Vector_item()).(type) {
	case *VectorItemReference:

		itemTemp := v.Factory.NewTemp()
		offset := itemRef.Vector.GetStackStmt(v.Factory) // stack[variable]
		offsetTemp := v.Factory.NewTemp()
		v.Factory.AppendToBlock(v.Factory.NewSimpleAssignment().SetAssignee(offsetTemp).SetVal(offset)) // t = stack[variable]
		return v.GetVectorvalue(itemTemp, offsetTemp, itemRef.Index, itemRef.Vector.Type)
	case *MatrixItemReference:

		itemTemp := v.Factory.NewTemp()
		offset := itemRef.Matrix.GetStackStmt(v.Factory) // stack[variable]
		offsetTemp := v.Factory.NewTemp()
		v.Factory.AppendToBlock(v.Factory.NewSimpleAssignment().SetAssignee(offsetTemp).SetVal(offset)) // t = stack[variable]
		_type := itemRef.Matrix.Type

		fmt.Println("MATRIX ITEM REFERENCE")
		fmt.Println(itemRef.Indexes)
		fmt.Println("MATRIX ITEM REFERENCE")

		for _, index := range itemRef.Indexes {
			vw := v.GetVectorvalue(itemTemp, offsetTemp, index, _type)
			// assign itemTemp to offsetTemp
			v.Factory.AppendToBlock(v.Factory.NewSimpleAssignment().SetAssignee(offsetTemp).SetVal(vw.Val))
			_type = vw.Metadata
		}

		return &value.ValueWrapper{
			Val:      itemTemp,
			Metadata: _type,
		}
	}

	return v.GetNilVW()
}

func (v *IrVisitor) GetVectorvalue(itemTemp, offsetTemp *tac.Temp, index *value.ValueWrapper, _type string) *value.ValueWrapper {

	endLabel := v.Factory.NewLabel()
	errorLabel := v.Factory.NewLabel()
	relativeAddress := index.Val

	// check size
	size := v.Factory.NewTemp()
	// size = heap[t]
	v.Factory.AppendToBlock(v.Factory.NewSimpleAssignment().SetAssignee(size).SetVal(v.Factory.NewHeapIndexed().SetIndex(offsetTemp)))

	//  relativeAddress >= size
	condition := v.Factory.NewBoolExpression().SetLeft(relativeAddress).SetRight(size).SetOp(tac.GTE).SetLeftCast("int")
	v.Factory.AppendToBlock(v.Factory.NewConditionalJump().SetCondition(condition).SetTarget(errorLabel))

	// if relativeAddress < 0
	condition = v.Factory.NewBoolExpression().SetLeft(relativeAddress).SetRight(v.Factory.NewLiteral().SetValue("0")).SetOp(tac.LT).SetLeftCast("int")
	v.Factory.AppendToBlock(v.Factory.NewConditionalJump().SetCondition(condition).SetTarget(errorLabel))

	// increase offset by 1, to skip size
	v.Factory.AppendToBlock(v.Factory.NewCompoundAssignment().SetAssignee(offsetTemp).SetLeft(offsetTemp).SetRight(v.Factory.NewLiteral().SetValue("1")).SetOperator("+"))

	absoluteAddress := v.Factory.NewTemp()
	v.Factory.AppendToBlock(v.Factory.NewCompoundAssignment().SetAssignee(absoluteAddress).SetLeft(relativeAddress).SetRight(offsetTemp).SetOperator("+"))
	// absoluteAddress = relativeAddress + offsetTemp

	heapIndexed := v.Factory.NewHeapIndexed().SetIndex(absoluteAddress)
	v.Factory.AppendToBlock(v.Factory.NewSimpleAssignment().SetAssignee(itemTemp).SetVal(heapIndexed))

	v.Factory.AppendToBlock(v.Factory.NewUnconditionalJump().SetTarget(endLabel))

	v.Factory.AppendToBlock(errorLabel)

	v.Factory.AppendToBlock(v.Factory.NewSimpleAssignment().SetAssignee(itemTemp).SetVal(v.GetNilVW().Val))
	v.Factory.AppendBlock(v.Utility.PrintStringStream("BoundsError\n"))

	v.Factory.AppendToBlock(endLabel)

	return &value.ValueWrapper{
		Val:      itemTemp,
		Metadata: utils.RemoveOneLevelOfBrackets(_type),
	}
}

func (v *IrVisitor) VisitRepeating(ctx *compiler.RepeatingContext) interface{} {

	if ctx.ID(0).GetText() != "repeating" {
		return v.GetNilVW()
	}

	if ctx.ID(1).GetText() != "count" {
		return v.GetNilVW()
	}

	reapeating_val := v.Visit(ctx.Expr(0)).(*value.ValueWrapper)
	count_val := v.Visit(ctx.Expr(1)).(*value.ValueWrapper)

	if count_val.Metadata != abstract.IVOR_INT {
		return v.GetNilVW()
	}

	_type := ""

	if ctx.Vector_type() != nil {
		_type = ctx.Vector_type().GetText()
		primitive_type := utils.RemoveBrackets(_type)

		if primitive_type != reapeating_val.Metadata {
			return v.GetNilVW()
		}
	} else if ctx.Matrix_type() != nil {

		_type = ctx.Matrix_type().GetText()

		if !(utils.IsMatrixType(reapeating_val.Metadata) || utils.IsVectorType(reapeating_val.Metadata)) {
			return v.GetNilVW()
		}
	}

	// save the size (count)
	vectAddress := v.Factory.NewTemp()
	v.saveVectorSize(vectAddress, count_val.Val)

	// save the items
	count := v.Factory.NewTemp()
	v.Factory.AppendToBlock(v.Factory.NewSimpleAssignment().SetAssignee(count).SetVal(v.Factory.NewLiteral().SetValue("0")))
	loopLabel := v.Factory.NewLabel()
	endLabel := v.Factory.NewLabel()

	v.Factory.AppendToBlock(loopLabel)

	condition := v.Factory.NewBoolExpression().SetLeft(count).SetRight(count_val.Val).SetOp(tac.GTE).SetLeftCast("int")
	v.Factory.AppendToBlock(v.Factory.NewConditionalJump().SetCondition(condition).SetTarget(endLabel))

	heapAddres := v.Factory.NewHeapIndexed().SetIndex(v.Factory.NewHeapPtr())
	// ? if the val is a vector or matrix, we need to copy it. But i dont give a fuck
	v.Factory.AppendToBlock(v.Factory.NewSimpleAssignment().SetAssignee(heapAddres).SetVal(reapeating_val.Val)) // ! <--
	v.Utility.IncreaseHeapPtr()

	v.Factory.AppendToBlock(v.Factory.NewCompoundAssignment().SetAssignee(count).SetLeft(count).SetRight(v.Factory.NewLiteral().SetValue("1")).SetOperator("+"))

	v.Factory.AppendToBlock(v.Factory.NewUnconditionalJump().SetTarget(loopLabel))

	v.Factory.AppendToBlock(endLabel)

	return &value.ValueWrapper{
		Val:      vectAddress,
		Metadata: _type,
	}
}

func (v *IrVisitor) VisitRepeatingExp(ctx *compiler.RepeatingExpContext) interface{} {
	return v.Visit(ctx.Repeating())
}
