package visitor

import (
	"strconv"

	"github.com/damianpeaf/OLC2_Proyecto2_202110568/compiler"
	"github.com/damianpeaf/OLC2_Proyecto2_202110568/ir/abstract"
	"github.com/damianpeaf/OLC2_Proyecto2_202110568/ir/value"
)

func (v *IrVisitor) VisitStructDecl(ctx *compiler.StructDeclContext) interface{} {
	if v.ScopeTrace.CurrentScope != v.ScopeTrace.GlobalScope {
		return nil
	}

	name := ctx.ID().GetText()

	mainTrace := v.ScopeTrace
	structTrace := abstract.NewScopeTrace(false, 0, v.Factory)

	v.ScopeTrace = structTrace

	for _, prop := range ctx.AllStruct_prop() {
		v.Visit(prop)
	}
	v.ScopeTrace = mainTrace

	v.ScopeTrace.CurrentScope.NewStruct(name, structTrace)
	return nil
}

func (v *IrVisitor) VisitStructAttr(ctx *compiler.StructAttrContext) interface{} {

	if ctx.Type_() != nil && ctx.Expr() != nil {
		return nil
	}

	varName := ctx.ID().GetText()

	var varValue *value.ValueWrapper = v.GetNilVW()
	explicitType := ""
	implicitType := ""
	finalType := ""

	// value is defined
	if ctx.Expr() != nil {
		varValue = v.Visit(ctx.Expr()).(*value.ValueWrapper)
		implicitType = varValue.Metadata
	}

	if ctx.Type_() != nil {
		explicitType = ctx.Type_().GetText()
	}

	// explicit type and implicit type are defined
	if explicitType != "" && implicitType != "" {
		if explicitType != implicitType {
			return nil
		}
	}

	// only explicit type is defined
	if explicitType != "" && implicitType == "" {
		finalType = explicitType
	} else {
		// only implicit type is defined
		finalType = implicitType
	}

	v.ScopeTrace.NewProp(varName, finalType, varValue)

	return nil
}

func (v *IrVisitor) VisitStructFunc(ctx *compiler.StructFuncContext) interface{} {
	// TODO: implement
	return nil
}

func (v *IrVisitor) VisitStructVector(ctx *compiler.StructVectorContext) interface{} {
	// TODO: implement
	return nil
}

func (v *IrVisitor) VisitStructVectorExp(ctx *compiler.StructVectorExpContext) interface{} {
	return v.Visit(ctx.Struct_vector())
}

func (v *IrVisitor) VisitVectorFuncExp(ctx *compiler.VectorFuncExpContext) interface{} {
	return v.Visit(ctx.Vector_func())
}

func (v *IrVisitor) VisitVectorPropExp(ctx *compiler.VectorPropExpContext) interface{} {
	return v.Visit(ctx.Vector_prop())
}

func (v *IrVisitor) VisitVectorProp(ctx *compiler.VectorPropContext) interface{} {
	// TODO: implement
	return nil
}

func (v *IrVisitor) VisitVectorFunc(ctx *compiler.VectorFuncContext) interface{} {
	// TODO: implement
	return nil
}

// v.BuildStruct(canditateName, funcObj, args)
func (v *IrVisitor) BuildStruct(name string, scope *abstract.BaseScope, args []*abstract.Argument) *value.ValueWrapper {

	// create args map
	argMap := make(map[string]*abstract.Argument)
	for _, arg := range args {
		// repeat arg
		if _, ok := argMap[arg.Name]; ok {
			return v.GetNilVW()
		}
		argMap[arg.Name] = arg
	}

	// * create struct

	v.Utility.Comment("create struct")

	// 1. save the size for the struct on the heap
	structAddress := v.Factory.NewTemp()
	structSize := v.Factory.NewLiteral().SetValue(strconv.Itoa(scope.ScopeTrace.Correlative))

	v.Factory.AppendToBlock(v.Factory.NewSimpleAssignment().SetAssignee(v.Factory.NewHeapIndexed().SetIndex(v.Factory.NewHeapPtr())).SetVal(structSize))
	v.Factory.AppendToBlock(v.Factory.NewSimpleAssignment().SetAssignee(structAddress).SetVal(v.Factory.NewHeapPtr()))
	v.Utility.IncreaseHeapPtr()

	// save each field on the heap
	for i := 0; i < scope.ScopeTrace.Correlative; i++ {
		prop := scope.GetByAddress(i)

		vw := v.GetNilVW()
		arg := argMap[prop.Name]

		if arg != nil {
			vw = arg.Wrapper
		} else {
			if prop.DefaultValue != nil {
				vw = prop.DefaultValue
			}
		}

		v.Factory.AppendToBlock(v.Factory.NewSimpleAssignment().SetAssignee(v.Factory.NewHeapIndexed().SetIndex(v.Factory.NewHeapPtr())).SetVal(vw.Val))
		v.Utility.IncreaseHeapPtr()
	}

	return &value.ValueWrapper{
		Val:      structAddress,
		Metadata: name,
	}
}
