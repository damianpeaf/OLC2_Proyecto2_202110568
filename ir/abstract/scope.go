package abstract

import "fmt"

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
	s.Variables[name] = &IVOR{
		Name:    name,
		Type:    _type,
		Address: s.ScopeTrace.Correlative,
	}
	s.ScopeTrace.Correlative++
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
	Correlative  int
}

func NewGlobalScope(trace *ScopeTrace) *BaseScope {
	// TODO: register built-in functions

	initialFuncs := make(map[string]*Function)
	initialFuncs["print"] = &Function{
		Name: "print",
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

func (s *ScopeTrace) NewFunction(name string, params []*Param) *Function {
	return s.CurrentScope.NewFunction(name, params)
}

func (s *ScopeTrace) GetFunction(name string) *Function {
	return s.CurrentScope.GetFunction(name)
}

func NewScopeTrace() *ScopeTrace {
	trace := &ScopeTrace{
		GlobalScope:  nil,
		CurrentScope: nil,
		Correlative:  0,
	}
	globalScope := NewGlobalScope(trace)
	trace.GlobalScope = globalScope
	trace.CurrentScope = globalScope
	return trace
}
