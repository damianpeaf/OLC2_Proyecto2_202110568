package abstract

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/damianpeaf/OLC2_Proyecto2_202110568/ir/tac"
	"github.com/damianpeaf/OLC2_Proyecto2_202110568/ir/utils"
)

// ***** BASE SCOPE *****
type BaseScope struct {
	Name       string
	Parent     *BaseScope
	Children   []*BaseScope
	Variables  map[string]*IVOR
	Functions  map[string]*Function
	ScopeTrace *ScopeTrace
	innerOrder int // TODO: need to find a way to reset this for recursive functions
}

func (s *BaseScope) AddChild(child *BaseScope) {
	s.Children = append(s.Children, child)
	child.Parent = s

	fmt.Printf("%s -> %s\n", s.Name, child.Name)
}

func (s *BaseScope) NewVariable(name string, _type string) {

	offset := 0
	if s.ScopeTrace.FrameRelative {
		offset = s.ScopeTrace.ParamOffset + 1 // headers + params
	}

	s.Variables[name] = &IVOR{
		Name:          name,
		Type:          _type,
		Address:       s.ScopeTrace.Correlative,
		FrameRelative: s.ScopeTrace.FrameRelative,
		Offset:        offset,
	}
	s.ScopeTrace.Correlative++
}

func (s *BaseScope) DirectVariable(variable *IVOR) {
	s.Variables[variable.Name] = variable
}

func (s *BaseScope) GetVariable(pattern string) *IVOR {

	if strings.Contains(pattern, ".") {
		paternParts := strings.Split(pattern, ".")
		mainObjName := paternParts[0]
		mainObj := s.GetVariable(mainObjName)

		if mainObj == nil {
			return nil
		}

		var scope *BaseScope

		if IsPrimitiveType(mainObj.Type) {
			return nil
		}

		// search for its scope
		if utils.IsVectorType(mainObj.Type) {
			scope = DefaultVectorScope
		}
		// TODO: get struct scope
		factory := s.ScopeTrace.Factory
		structPointer := factory.NewTemp()
		factory.AppendToBlock(factory.NewSimpleAssignment().SetAssignee(structPointer).SetVal(mainObj.GetStackStmt(factory)))
		// SP = stack[x] --> structPointer = struct_address

		return s.searchObjectProp(structPointer, scope, strings.Join(paternParts[1:], "."))
	}

	// TODO: check if pattern refers to a pointer
	aux := s
	for aux != nil {
		if aux.Variables[pattern] != nil {
			return aux.Variables[pattern]
		}
		aux = aux.Parent
	}
	return nil
}

func (s *BaseScope) searchObjectProp(structPointer *tac.Temp, scope *BaseScope, pattern string) *IVOR {
	factory := s.ScopeTrace.Factory
	patternParts := strings.Split(pattern, ".")
	propName := patternParts[0]

	if scope == DefaultVectorScope {
		switch propName {

		case "count":
			factory.AppendToBlock(factory.NewSimpleAssignment().SetAssignee(structPointer).SetVal(factory.NewHeapIndexed().SetIndex(structPointer)))
			return &IVOR{
				Name:          pattern,
				Type:          IVOR_INT,
				Address:       -1,
				FrameRelative: false,
				Offset:        -1,
				Temp:          structPointer,
			}

		case "isEmpty":
			computation := factory.NewTemp()
			factory.AppendToBlock(factory.NewSimpleAssignment().SetAssignee(structPointer).SetVal(factory.NewHeapIndexed().SetIndex(structPointer)))
			endLabel := factory.NewLabel()
			isEmptyLabel := factory.NewLabel()

			// if computation == 0 goto isEmpty
			// structPointer == 0
			// goto end
			// isEmpty:
			// computation == 1
			// end:

			condition := factory.NewBoolExpression().SetLeft(structPointer).SetLeftCast("int").SetRight(factory.NewLiteral().SetValue("0")).SetOp(tac.EQ)
			factory.AppendToBlock(factory.NewConditionalJump().SetCondition(condition).SetTarget(isEmptyLabel))
			factory.AppendToBlock(factory.NewSimpleAssignment().SetAssignee(computation).SetVal(factory.NewLiteral().SetValue("0")))
			factory.AppendToBlock(factory.NewUnconditionalJump().SetTarget(endLabel))
			factory.AppendToBlock(isEmptyLabel)
			factory.AppendToBlock(factory.NewSimpleAssignment().SetAssignee(computation).SetVal(factory.NewLiteral().SetValue("1")))
			factory.AppendToBlock(endLabel)

			return &IVOR{
				Name:          pattern,
				Type:          IVOR_BOOL,
				Address:       -1,
				FrameRelative: false,
				Offset:        -1,
				Temp:          computation,
			}

		default:
			return nil
		}
	}

	prop := scope.GetVariable(propName)

	if prop == nil {
		return nil
	}

	heapIndex := prop.Address
	heapIndex++
	// SP = SP + prop.offset
	factory.AppendToBlock(factory.NewCompoundAssignment().SetAssignee(structPointer).SetLeft(structPointer).SetRight(factory.NewLiteral().SetValue(strconv.Itoa(heapIndex))).SetOperator("+"))

	// SP = heap[SP]
	factory.AppendToBlock(factory.NewSimpleAssignment().SetAssignee(structPointer).SetVal(factory.NewHeapIndexed().SetIndex(structPointer)))

	// keep searching
	if !strings.Contains(pattern, ".") {
		return &IVOR{
			Name:          pattern,
			Type:          prop.Type,
			Address:       -1,
			FrameRelative: false,
			Offset:        -1,
			Temp:          structPointer,
		}
	}

	var newScope *BaseScope

	if IsPrimitiveType(prop.Type) {
		return nil
	}

	// search for its scope
	if utils.IsVectorType(prop.Type) {
		newScope = DefaultVectorScope
	}

	// TODO: get struct scope

	return s.searchObjectProp(structPointer, newScope, strings.Join(patternParts[1:], "."))
}

func (s *BaseScope) NewFunction(name string, f *Function) *Function {
	s.Functions[name] = f
	return f
}

func (s *BaseScope) GetFunction(name string) *Function {
	if strings.Contains(name, ".") {
		patternParts := strings.Split(name, ".")
		mainObjName := patternParts[0]
		mainObj := s.GetVariable(mainObjName)

		if mainObj == nil {
			return nil
		}

		var scope *BaseScope

		if IsPrimitiveType(mainObj.Type) {
			return nil
		}

		// search for its scope
		if utils.IsVectorType(mainObj.Type) {
			scope = DefaultVectorScope
		}

		// TODO: get struct scope
		factory := s.ScopeTrace.Factory
		structPointer := factory.NewTemp()
		factory.AppendToBlock(factory.NewSimpleAssignment().SetAssignee(structPointer).SetVal(mainObj.GetStackStmt(factory)))
		// SP = stack[x] --> structPointer = struct_address

		objFunc := s.searchObjectFunction(structPointer, scope, strings.Join(patternParts[1:], "."))

		if objFunc != nil {
			objFunc.StructPointer = structPointer
			objFunc.StructRef = mainObj
		}
		return objFunc
	}
	aux := s
	for aux != nil {
		fmt.Println("funcs: ", aux.Functions)
		if aux.Functions[name] != nil {
			return aux.Functions[name]
		}
		aux = aux.Parent
	}
	return nil
}

func (s *BaseScope) searchObjectFunction(structPointer *tac.Temp, scope *BaseScope, pattern string) *Function {
	factory := s.ScopeTrace.Factory

	patternParts := strings.Split(pattern, ".")
	funcName := patternParts[0]

	if scope == DefaultVectorScope {
		switch funcName {
		case "append":
			return VectorAppendFunc
		case "removeLast":
			return VectorRemoveLastFunc
		case "remove":
			return VectorRemoveFunc
		default:
			return nil
		}
	}

	// if its the last part of the pattern, return the function
	if !strings.Contains(pattern, ".") {
		return scope.GetFunction(funcName)
	}

	// traverse the object
	prop := scope.GetVariable(funcName)

	if prop == nil {
		return nil
	}

	heapIndex := prop.Address
	heapIndex++

	// SP = SP + prop.offset
	factory.AppendToBlock(factory.NewCompoundAssignment().SetAssignee(structPointer).SetLeft(structPointer).SetRight(factory.NewLiteral().SetValue(strconv.Itoa(heapIndex))).SetOperator("+"))

	// SP = heap[SP]
	factory.AppendToBlock(factory.NewSimpleAssignment().SetAssignee(structPointer).SetVal(factory.NewHeapIndexed().SetIndex(structPointer)))

	var newScope *BaseScope

	if IsPrimitiveType(prop.Type) {
		return nil
	}

	// search for its scope
	if utils.IsVectorType(prop.Type) {
		newScope = DefaultVectorScope
	}

	// TODO: get struct scope

	return s.searchObjectFunction(structPointer, newScope, strings.Join(patternParts[1:], "."))
}

// TODO: NewStruct
// TODO: GetStruct

func (s *BaseScope) Reset() {
	s.Variables = make(map[string]*IVOR)
	// s.children = make([]*BaseScope, 0)
	s.Functions = make(map[string]*Function)
}

// ***** SCOPE TRACE *****
type ScopeTrace struct {
	GlobalScope   *BaseScope
	CurrentScope  *BaseScope
	Correlative   int
	FrameRelative bool
	ParamOffset   int
	Factory       *tac.TACFactory
}

func NewGlobalScope(trace *ScopeTrace) *BaseScope {
	initialFuncs := make(map[string]*Function)
	initialFuncs["print"] = &Function{
		Name: "print",
		Type: BUILTIN_FUNCTION,
	}
	initialFuncs["Int"] = &Function{
		Name: "Int",
		Type: BUILTIN_FUNCTION,
	}
	initialFuncs["Float"] = &Function{
		Name: "Float",
		Type: BUILTIN_FUNCTION,
	}
	initialFuncs["String"] = &Function{
		Name: "String",
		Type: BUILTIN_FUNCTION,
	}

	return &BaseScope{
		Name:       "global",
		Variables:  make(map[string]*IVOR),
		Functions:  initialFuncs,
		Children:   make([]*BaseScope, 0),
		Parent:     nil,
		ScopeTrace: trace,
	}
}

func NewLocalScope(name string, trace *ScopeTrace) *BaseScope {
	return &BaseScope{
		Name:       name,
		Variables:  make(map[string]*IVOR),
		Functions:  make(map[string]*Function),
		Children:   make([]*BaseScope, 0),
		Parent:     nil,
		ScopeTrace: trace,
	}
}

func (s *ScopeTrace) PushScope(name string) *BaseScope {

	newScope := NewLocalScope(name, s)
	s.CurrentScope.AddChild(newScope)
	s.CurrentScope = newScope

	return s.CurrentScope
}

func (s *ScopeTrace) NextScope() *BaseScope {

	fmt.Println("NEXT SCOPE _1", s.CurrentScope.Name)

	if len(s.CurrentScope.Children) > 0 {

		if s.CurrentScope.innerOrder >= len(s.CurrentScope.Children) {
			fmt.Println(s.CurrentScope.Name, "has no more children")
		}

		fmt.Println("NEXT SCOPE _2", s.CurrentScope.Children[s.CurrentScope.innerOrder].Name)
		fmt.Println("")

		prevScope := s.CurrentScope
		s.CurrentScope = s.CurrentScope.Children[s.CurrentScope.innerOrder]
		prevScope.innerOrder++
		return s.CurrentScope
	}
	return nil
}

func (s *ScopeTrace) PopScope() {
	s.CurrentScope = s.CurrentScope.Parent
}

func (s *ScopeTrace) PrevScope() *BaseScope {

	fmt.Println("PREV SCOPE _1", s.CurrentScope.Name)

	if s.CurrentScope.Parent != nil {

		fmt.Println("PREV SCOPE _2", s.CurrentScope.Parent.Name)
		fmt.Println("")

		s.CurrentScope = s.CurrentScope.Parent
		return s.CurrentScope
	}
	return nil
}

func (s *ScopeTrace) Reset() {
	s.CurrentScope = s.GlobalScope
}

func (s *ScopeTrace) NewVariable(name string, _type string) {
	s.CurrentScope.NewVariable(name, _type)
}

func (s *ScopeTrace) GetVariable(pattern string) *IVOR {
	return s.CurrentScope.GetVariable(pattern)
}

func (s *ScopeTrace) NewFunction(name string, f *Function) *Function {
	return s.CurrentScope.NewFunction(name, f)
}

func (s *ScopeTrace) GetFunction(name string) *Function {
	return s.CurrentScope.GetFunction(name)
}

func NewScopeTrace(frameRelative bool, paramOffset int, factory *tac.TACFactory) *ScopeTrace {
	trace := &ScopeTrace{
		GlobalScope:   nil,
		CurrentScope:  nil,
		Correlative:   0,
		FrameRelative: frameRelative,
		ParamOffset:   paramOffset,
		Factory:       factory,
	}
	globalScope := NewGlobalScope(trace)
	trace.GlobalScope = globalScope
	trace.CurrentScope = globalScope
	return trace
}

func IsPrimitiveType(t string) bool {
	switch t {
	case IVOR_BOOL, IVOR_INT, IVOR_FLOAT, IVOR_STRING, IVOR_NIL, IVOR_CHARACTER:
		return true
	default:
		return false
	}
}
