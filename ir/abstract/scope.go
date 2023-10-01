package abstract

import (
	"github.com/damianpeaf/OLC2_Proyecto2_202110568/ir/tac"
)

// ***** BASE SCOPE *****
type BaseScope struct {
	Name      string
	Parent    *BaseScope
	Children  []*BaseScope
	Variables map[string]*IVOR
	Factory   *tac.TACFactory
	Functions map[string]*Function
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

	s.Variables[name] = &IVOR{
		Name:    name,
		Type:    _type,
		Address: stackAddress,
	}
}

func (s *BaseScope) GetVariable(pattern string) *IVOR {
	// TODO: check if pattern refers to a struct field
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

// TODO: searchObjectVariable

func (s *BaseScope) NewFunction(name string, params []*Param) *Function {

	createdFunc := &Function{
		Name:  name,
		Param: params,
	}

	s.Functions[name] = createdFunc

	return createdFunc
}

func (s *BaseScope) GetFunction(name string) *Function {
	// TODO: check if function is a method
	aux := s
	for aux != nil {
		if aux.Functions[name] != nil {
			return aux.Functions[name]
		}
		aux = aux.Parent
	}
	return nil
}

// TODO: searchObjectFunction

// TODO: NewStruct
// TODO: GetStruct

func (s *BaseScope) Reset() {
	s.Variables = make(map[string]*IVOR)
	// s.children = make([]*BaseScope, 0)
	s.Functions = make(map[string]*Function)
}

// ***** SCOPE TRACE *****
type ScopeTrace struct {
	GlobalScope  *BaseScope
	CurrentScope *BaseScope
	Factory      *tac.TACFactory
}

func NewGlobalScope(factory *tac.TACFactory) *BaseScope {
	// TODO: register built-in functions

	initialFuncs := make(map[string]*Function)
	initialFuncs["print"] = &Function{
		Name: "print",
		Type: BUILTIN_FUNCTION,
	}

	return &BaseScope{
		Name:      "global",
		Variables: make(map[string]*IVOR),
		Functions: initialFuncs,
		Children:  make([]*BaseScope, 0),
		Parent:    nil,
	}
}

func NewLocalScope(name string, factory *tac.TACFactory) *BaseScope {
	return &BaseScope{
		Name:      name,
		Variables: make(map[string]*IVOR),
		Functions: make(map[string]*Function),
		Children:  make([]*BaseScope, 0),
		Parent:    nil,
		Factory:   factory,
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

func (s *ScopeTrace) NewVariable(name string, val tac.SimpleValue, _type string) {
	s.CurrentScope.NewVariable(name, val, _type)
}

func (s *ScopeTrace) GetVariable(pattern string) *IVOR {
	return s.CurrentScope.GetVariable(pattern)
}

func (s *ScopeTrace) NewFunction(name string, params []*Param) *Function {
	return s.CurrentScope.NewFunction(name, params)
}

func (s *ScopeTrace) GetFunction(name string) *Function {
	return s.CurrentScope.GetFunction(name)
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
