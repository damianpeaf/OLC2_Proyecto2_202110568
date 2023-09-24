package repl

import (
	"github.com/damianpeaf/OLC2_Proyecto2_202110568/value"
	"log"
	"strings"

	"github.com/antlr4-go/antlr/v4"
)

type BaseScope struct {
	name       string
	parent     *BaseScope
	children   []*BaseScope
	variables  map[string]*Variable
	functions  map[string]value.IVOR
	structs    map[string]*Struct
	isStruct   bool
	IsMutating bool
}

func (s *BaseScope) Name() string {
	return s.name
}

func (s *BaseScope) Parent() *BaseScope {
	return s.parent
}

func (s *BaseScope) Children() []*BaseScope {
	return s.children
}

func (s *BaseScope) ValidType(_type string) bool {

	_, isStructType := s.structs[_type]

	return value.IsPrimitiveType(_type) || isStructType
}

func (s *BaseScope) AddChild(child *BaseScope) {
	s.children = append(s.children, child)
	child.parent = s
}

func (s *BaseScope) variableExists(variable *Variable) bool {

	if _, ok := s.variables[variable.Name]; ok {
		return true
	}

	return false

}

func (s *BaseScope) AddVariable(name string, varType string, value value.IVOR, isConst bool, allowNil bool, token antlr.Token) (*Variable, string) {

	variable := &Variable{
		Name:     name,
		Value:    value,
		Type:     varType,
		IsConst:  isConst,
		AllowNil: allowNil,
		Token:    token,
	}

	if s.variableExists(variable) {
		return nil, "La variable " + name + " ya existe"
	}

	typesOk, msg := variable.TypeValidation()

	// even if the variable is not valid, we add it to the scope, (internally it will be nil)
	s.variables[name] = variable

	if !typesOk {
		// report error
		return nil, msg
	}

	return variable, ""
}

func (s *BaseScope) GetVariable(name string) *Variable {
	// verify if is refering to and object/struct function
	if strings.Contains(name, ".") {
		return s.searchObjectVariable(name, nil)
	}

	initialScope := s

	for {
		if variable, ok := initialScope.variables[name]; ok {

			// verify if is refering to a pointer
			if variable.Type == value.IVOR_POINTER {
				return variable.Value.(*PointerValue).AssocVariable // pointer of a pointer ?
			}

			return variable
		}

		if initialScope.parent == nil {
			break
		}

		initialScope = initialScope.parent
	}

	return nil
}

// obj1.obj2.prop1

func (s *BaseScope) searchObjectVariable(name string, lastObj value.IVOR) *Variable {

	// split name by dot
	parts := strings.Split(name, ".")

	if len(parts) == 0 {
		log.Fatal("idk what u did, cant split by dot")
		return nil
	}

	if len(parts) == 1 {
		obj, ok := lastObj.(*ObjectValue)

		if ok {
			return obj.InternalScope.GetVariable(name)
		}

		log.Fatal("idk what u did, cant convert to object")
		return nil
	}

	// then parts should be 2 or more

	if lastObj == nil {
		variable := s.GetVariable(parts[0])

		if variable == nil {
			return nil
		}

		obj := variable.Value

		// obj must be an object/struct or vector
		switch obj := obj.(type) {
		case *ObjectValue:
			lastObj = obj
		case *VectorValue:
			lastObj = obj.ObjectValue
		default:
			return nil
		}

		return s.searchObjectVariable(strings.Join(parts[1:], "."), lastObj)
	}

	obj, ok := lastObj.(*ObjectValue)

	if ok {
		lastObj = obj.InternalScope.GetVariable(parts[0]).Value

		return s.searchObjectVariable(strings.Join(parts[1:], "."), lastObj)
	} else {
		log.Fatal("idk what u did, cant convert to object")
		return nil
	}
}

func (s *BaseScope) AddFunction(name string, function value.IVOR) (bool, string) {
	// check if function already exists

	if _, ok := s.functions[name]; ok {
		return false, "La funcion " + name + " ya existe"
	}

	s.functions[name] = function

	return true, ""
}

func (s *BaseScope) GetFunction(name string) (value.IVOR, string) {

	// verify if is refering to and object/struct function
	if strings.Contains(name, ".") {
		return s.searchObjectFunction(name, nil)
	}

	initialScope := s

	for {
		if function, ok := initialScope.functions[name]; ok {
			return function, ""
		}

		if initialScope.parent == nil {
			break
		}

		initialScope = initialScope.parent
	}

	return nil, "La funcion " + name + " no existe"
}

// obj1.obj2.func1()

func (s *BaseScope) searchObjectFunction(name string, lastObj value.IVOR) (value.IVOR, string) {

	// split name by dot
	parts := strings.Split(name, ".")

	if len(parts) == 0 {
		log.Fatal("idk what u did, cant split by dot")
		return nil, ""
	}

	if len(parts) == 1 {
		obj, ok := lastObj.(*ObjectValue)

		if ok {
			return obj.InternalScope.GetFunction(name)
		}

		log.Fatal("idk what u did, cant convert to object")
		return nil, ""
	}

	// then parts should be 2 or more

	if lastObj == nil {
		variable := s.GetVariable(parts[0])

		if variable == nil {
			return nil, "No se puede acceder a la propiedad " + parts[0]
		}

		obj := variable.Value

		// obj must be an object/struct or vector

		switch obj := obj.(type) {
		case *ObjectValue:
			lastObj = obj
		case *VectorValue:
			lastObj = obj.ObjectValue
		default:
			return nil, "La propiedad '" + variable.Name + "' de tipo " + obj.Type() + " no tiene propiedades"
		}

		return s.searchObjectFunction(strings.Join(parts[1:], "."), lastObj)
	}

	obj, ok := lastObj.(*ObjectValue)

	if ok {
		lastObj = obj.InternalScope.GetVariable(parts[0]).Value

		return s.searchObjectFunction(strings.Join(parts[1:], "."), lastObj)
	} else {
		log.Fatal("idk what u did, cant convert to object")
		return nil, ""
	}
}

func (s *BaseScope) AddStruct(name string, structValue *Struct) (bool, string) {

	if _, ok := s.structs[name]; ok {
		return false, "La estructura " + name + " ya existe"
	}

	s.structs[name] = structValue
	return true, ""
}

func (s *BaseScope) GetStruct(name string) (*Struct, string) {

	initialScope := s

	for {
		if structValue, ok := initialScope.structs[name]; ok {
			return structValue, ""
		}

		if initialScope.parent == nil {
			break
		}

		initialScope = initialScope.parent
	}

	return nil, "La estructura " + name + " no existe"
}

func (s *BaseScope) Reset() {
	s.variables = make(map[string]*Variable)
	s.children = make([]*BaseScope, 0)
	s.functions = make(map[string]value.IVOR)
}

func (s *BaseScope) IsMutatingScope() bool {
	aux := s

	for {
		if aux.IsMutating {
			return true
		}

		if aux.parent == nil {
			break
		}

		aux = aux.parent
	}

	return false
}

func NewGlobalScope() *BaseScope {

	// register built-in functions

	funcs := make(map[string]value.IVOR)

	for k, v := range DefaultBuiltInFunctions {
		funcs[k] = v
	}

	return &BaseScope{
		name:      "global",
		variables: make(map[string]*Variable),
		children:  make([]*BaseScope, 0),
		structs:   make(map[string]*Struct),
		parent:    nil,
		functions: funcs,
	}
}

func NewLocalScope(name string) *BaseScope {
	return &BaseScope{
		name:      name,
		variables: make(map[string]*Variable),
		functions: make(map[string]value.IVOR),
		children:  make([]*BaseScope, 0),
		parent:    nil,
	}
}

type ScopeTrace struct {
	GlobalScope  *BaseScope
	CurrentScope *BaseScope
}

func (s *ScopeTrace) PushScope(name string) *BaseScope {

	newScope := NewLocalScope(name)
	s.CurrentScope.AddChild(newScope)
	s.CurrentScope = newScope

	return s.CurrentScope
}

func (s *ScopeTrace) PopScope() {
	s.CurrentScope = s.CurrentScope.Parent()
}

func (s *ScopeTrace) Reset() {
	s.CurrentScope = s.GlobalScope
}

func (s *ScopeTrace) AddVariable(name string, varType string, value value.IVOR, isConst bool, allowNil bool, token antlr.Token) (*Variable, string) {
	return s.CurrentScope.AddVariable(name, varType, value, isConst, allowNil, token)
}

func (s *ScopeTrace) GetVariable(name string) *Variable {
	return s.CurrentScope.GetVariable(name)
}

func (s *ScopeTrace) AddFunction(name string, function value.IVOR) (bool, string) {
	return s.CurrentScope.AddFunction(name, function)
}

func (s *ScopeTrace) GetFunction(name string) (value.IVOR, string) {
	return s.CurrentScope.GetFunction(name)
}

func (s *ScopeTrace) IsMutatingEnvironment() bool {
	return s.CurrentScope.IsMutatingScope()
}

func NewScopeTrace() *ScopeTrace {
	globalScope := NewGlobalScope()
	return &ScopeTrace{
		GlobalScope:  globalScope,
		CurrentScope: globalScope,
	}
}

func NewVectorScope() *BaseScope {
	var scope = &BaseScope{
		name:      "vector",
		variables: make(map[string]*Variable),
		children:  make([]*BaseScope, 0),
		functions: make(map[string]value.IVOR),
		parent:    nil,
	}

	// register object built-in functions

	return scope
}

func NewStructScope() *BaseScope {

	newGlobal := NewGlobalScope()

	return &BaseScope{
		name:      "struct",
		variables: make(map[string]*Variable),
		children:  make([]*BaseScope, 0),
		functions: make(map[string]value.IVOR),
		structs:   make(map[string]*Struct),
		parent:    newGlobal,
		isStruct:  true,
	}
}

// * Report

type ReportTable struct {
	GlobalScope ReportScope
}

type ReportScope struct {
	Name        string
	Vars        []ReportSymbol
	Funcs       []ReportSymbol
	Structs     []ReportSymbol
	ChildScopes []ReportScope
}

type ReportSymbol struct {
	Name   string
	Type   string
	Line   int
	Column int
}

func (s *ScopeTrace) Report() ReportTable {
	return ReportTable{
		GlobalScope: s.CurrentScope.Report(),
	}
}

func (s *BaseScope) Report() ReportScope {

	reportScope := ReportScope{
		Name:        s.name,
		Vars:        make([]ReportSymbol, 0),
		Funcs:       make([]ReportSymbol, 0),
		Structs:     make([]ReportSymbol, 0),
		ChildScopes: make([]ReportScope, 0),
	}

	for _, v := range s.variables {

		token := v.Token
		line := 0
		column := 0

		if token != nil {
			line = token.GetLine()
			column = token.GetColumn()
		}

		reportScope.Vars = append(reportScope.Vars, ReportSymbol{
			Name:   v.Name,
			Type:   v.Type,
			Line:   line,
			Column: column,
		})
	}

	for _, f := range s.functions {
		switch function := f.(type) {
		case *BuiltInFunction:
			reportScope.Funcs = append(reportScope.Funcs, ReportSymbol{
				Name:   function.Name,
				Type:   "Embebida: " + function.Name,
				Line:   0,
				Column: 0,
			})
		case *Function:

			line := 0
			column := 0

			if function.Token != nil {
				line = function.Token.GetLine()
				column = function.Token.GetColumn()
			}

			reportScope.Funcs = append(reportScope.Funcs, ReportSymbol{
				Name:   function.Name,
				Type:   function.ReturnType,
				Line:   line,
				Column: column,
			})
		case *ObjectBuiltInFunction:
			break
		default:
			log.Fatal("Function type not found")
		}
	}

	for _, v := range s.structs {
		reportScope.Structs = append(reportScope.Structs, ReportSymbol{
			Name:   v.Name,
			Type:   v.Name,
			Line:   v.Token.GetLine(),
			Column: v.Token.GetColumn(),
		})
	}

	for _, v := range s.children {
		reportScope.ChildScopes = append(reportScope.ChildScopes, v.Report())
	}

	return reportScope
}
