package tac

// ** MethodDcl
type MethodDcl struct {
	Name  string
	Block TACBlock
}

func (m *MethodDcl) String() string {
	blockStr := ""
	for _, stmt := range m.Block {
		blockStr += "\t" + stmt.String() + "\n"
	}
	blockStr += "return;\n"
	return "\nvoid " + m.Name + "() {\n" + blockStr + "}\n"
}

// builder utils
func (m *MethodDcl) AppendStmt(stmt TACStmtI) *MethodDcl {
	m.Block = append(m.Block, stmt)
	return m
}

func (m *MethodDcl) AppendStmts(stmts TACBlock) *MethodDcl {
	for _, stmt := range stmts {
		m.AppendStmt(stmt)
	}
	return m
}

func (m *MethodDcl) SetName(name string) *MethodDcl {
	m.Name = name
	return m
}

// ** MethodCall
type MethodCall struct {
	Name string // ? use a pointer to MethodDcl?
}

func (m *MethodCall) String() string {
	return m.Name + "();"
}

// builder utils
func (m *MethodCall) SetName(name string) *MethodCall {
	m.Name = name
	return m
}
