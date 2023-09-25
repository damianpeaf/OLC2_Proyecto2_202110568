package abstract

import (
	"github.com/damianpeaf/OLC2_Proyecto2_202110568/ir/tac"
)

// ***** BASE SCOPE *****
type BaseScope struct {
	Name     string
	Parent   *BaseScope
	Children []*BaseScope
	IVORS    map[string]*IVOR
	Factory  *tac.TACFactory
	// functions  map[string]value.IVOR
	// structs    map[string]*Struct
}

func (s *BaseScope) AddChild(child *BaseScope) {
	s.Children = append(s.Children, child)
	child.Parent = s
}

func (s *BaseScope) NewVariable(name string, val tac.SimpleValue, _type string) {
	if val == nil {
		return
	}
	s.Factory.AppendToBlock(s.Factory.NewComment().SetComment("Variable " + name + ": " + _type))
	stackAddress := s.Factory.Utility.SaveValOnStack(val)

	s.IVORS[name] = &IVOR{
		Name:    name,
		Type:    _type,
		Address: stackAddress,
	}
}

// TODO: GetVariable
// TODO: searchObjectVariable

// TODO: NewFunction
// TODO: GetFunction
// TODO: searchObjectFunction

// TODO: NewStruct
// TODO: GetStruct

func (s *BaseScope) Reset() {
	s.IVORS = make(map[string]*IVOR)
	// s.children = make([]*BaseScope, 0)
	// s.functions = make(map[string]value.IVOR)
}

// ***** SCOPE TRACE *****
type ScopeTrace struct {
	GlobalScope  *BaseScope
	CurrentScope *BaseScope
	Factory      *tac.TACFactory
}

func NewGlobalScope(factory *tac.TACFactory) *BaseScope {
	// TODO: register built-in functions
	return &BaseScope{
		Name:     "global",
		IVORS:    make(map[string]*IVOR),
		Children: make([]*BaseScope, 0),
		Parent:   nil,
	}
}

func NewLocalScope(name string, factory *tac.TACFactory) *BaseScope {
	return &BaseScope{
		Name:     name,
		IVORS:    make(map[string]*IVOR),
		Children: make([]*BaseScope, 0),
		Parent:   nil,
		Factory:  factory,
	}
}

func (s *ScopeTrace) PushScope(name string) *BaseScope {

	newScope := NewLocalScope(name, s.Factory)
	s.CurrentScope.AddChild(newScope)
	s.CurrentScope = newScope

	return s.CurrentScope
}

func (s *ScopeTrace) PopScope() {
	s.CurrentScope = s.CurrentScope.Parent
}

func (s *ScopeTrace) Reset() {
	s.CurrentScope = s.GlobalScope
}

func NewScopeTrace(factory *tac.TACFactory) *ScopeTrace {
	globalScope := NewGlobalScope(factory)
	globalScope.Factory = factory
	return &ScopeTrace{
		GlobalScope:  globalScope,
		CurrentScope: globalScope,
		Factory:      factory,
	}
}
