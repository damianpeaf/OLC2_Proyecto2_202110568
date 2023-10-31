package tac

import (
	"fmt"
	"strconv"
)

type Lifesaver struct {
	block         *[]TACStmtI
	funcName      string
	table         *LifeSaverTable
	factory       *TACFactory
	initialOffest int
	Reserve       int
}

type LifeSaverTable struct {
	table      []*LifeSaverTableEntry
	factory    *TACFactory
	cases      []*LifeSaverCase
	returnTemp *Temp
}

func NewLifeSaverTable(factory *TACFactory, returnTemp *Temp) *LifeSaverTable {
	return &LifeSaverTable{
		table:      make([]*LifeSaverTableEntry, 0),
		factory:    factory,
		returnTemp: returnTemp,
	}
}

func (lst *LifeSaverTable) AddEntry(entry *LifeSaverTableEntry) {

	if entry.temp != nil && lst.factory._framePointer != nil && lst.factory._framePointer == entry.temp {
		return
	}

	if entry.temp != nil && lst.returnTemp != nil && lst.returnTemp == entry.temp {
		return
	}

	lst.table = append(lst.table, entry)
}

func (lst *LifeSaverTable) ReduceToCases() {

	for i, entry := range lst.table {
		problematic := false
		problematicEntries := make([]*LifeSaverTableEntry, 0)

		if entry.appearanceType == "decl" {

			var centinel TACStmtI
			// search for its usage
			for j := i + 1; j < len(lst.table); j++ {
				if lst.table[j].appearanceType == "decl" && lst.table[j].temp == entry.temp {
					break
				}

				if lst.table[j].appearanceType == "return" {
					break
				}

				if lst.table[j].appearanceType == "recursive" {
					problematic = true
					centinel = lst.table[j].stmt
				}
				if problematic && lst.table[j].appearanceType == "used" && lst.table[j].temp == entry.temp {
					lst.table[j].centinel = centinel
					problematicEntries = append(problematicEntries, lst.table[j])
				}
			}
		}

		// create a case for this temp
		if problematic && len(problematicEntries) > 0 {
			newCase := &LifeSaverCase{
				temp: entry.temp,
				decl: entry,
				used: problematicEntries,
			}
			lst.cases = append(lst.cases, newCase)
		}
	}
}

func (lst *LifeSaverTable) Print() {
	fmt.Println("")
	fmt.Println("Lifesaver Table:")
	for _, entry := range lst.table {
		if entry.appearanceType == "decl" || entry.appearanceType == "used" {
			fmt.Println("type: ", entry.appearanceType, " temp: ", entry.temp.String(), " stmt: ", entry.stmt.String())
		} else {
			fmt.Println("type: ", entry.appearanceType)
		}
	}
	fmt.Println("")
}

func insertStmtsAfter(block *[]TACStmtI, stmts []TACStmtI, after TACStmtI) {

	for i, stmt := range *block {
		if stmt == after {
			*block = append((*block)[:i+1], append(stmts, (*block)[i+1:]...)...)
			break
		}
	}
}

// func insertStmtsBefore(block *[]TACStmtI, stmts []TACStmtI, before TACStmtI) {

// 	for i, stmt := range *block {
// 		if stmt == before {
// 			*block = append((*block)[:i], append(stmts, (*block)[i:]...)...)
// 			break
// 		}
// 	}
// }

func (lst *Lifesaver) FixCases() {
	for _, c := range lst.table.cases {

		// save val on stack
		addressTemp := lst.factory.NewTemp()
		offset := lst.factory.NewLiteral().SetValue(strconv.Itoa(lst.initialOffest))
		getAddress := lst.factory.NewCompoundAssignment().SetAssignee(addressTemp).SetLeft(lst.factory._framePointer).SetRight(offset).SetOperator("+")
		saveVal := lst.factory.NewSimpleAssignment().SetAssignee(lst.factory.NewStackIndexed().SetIndex(addressTemp)).SetVal(c.temp)
		// increaseStack := lst.factory.NewCompoundAssignment().SetAssignee(lst.factory.NewStackPtr()).SetLeft(lst.factory.NewStackPtr()).SetRight(lst.factory.NewLiteral().SetValue("1")).SetOperator("+")

		/*
			t1 = F + offset
			stack[t1] = t2
		*/

		newDclStmt := []TACStmtI{
			getAddress,
			saveVal,
		}

		insertStmtsAfter(lst.block, newDclStmt, c.decl.stmt)

		// load val from stack
		for _, u := range c.used {
			// access stmts
			addressTemp := lst.factory.NewTemp()
			valTemp := lst.factory.NewTemp()
			offset := lst.factory.NewLiteral().SetValue(strconv.Itoa(lst.initialOffest))
			getAddress := lst.factory.NewCompoundAssignment().SetAssignee(addressTemp).SetLeft(lst.factory._framePointer).SetRight(offset).SetOperator("+")
			saveVal := lst.factory.NewSimpleAssignment().SetAssignee(valTemp).SetVal(lst.factory.NewStackIndexed().SetIndex(addressTemp))

			newUsedStmt := []TACStmtI{
				getAddress,
				saveVal,
			}

			insertStmtsAfter(lst.block, newUsedStmt, u.centinel)

			fmt.Println("")
			fmt.Println("inserting after: ", u.centinel.String())

			// and replace the temp
			switch s := u.stmt.(type) {
			case *SimpleAssignment:
				s.Val = valTemp
			case *CompoundAssignment:
				if s.Left == u.temp {
					s.Left = valTemp
				}
				if s.Right == u.temp {
					s.Right = valTemp
				}
			}
		}
		lst.initialOffest++
		lst.Reserve++
	}

	// then set the new size of the frame
	allocParams := lst.factory.GetBuiltinParams("__alloc_frame")
	frameSizeTemp := allocParams[0]

	for _, stmt := range *lst.block {
		switch s := stmt.(type) {
		case *SimpleAssignment:
			if s.Assignee == frameSizeTemp {
				prevVal, ok := strconv.Atoi(s.Val.String())
				if ok != nil {
					panic("Error parsing frame size")
				}
				s.Val = lst.factory.NewLiteral().SetValue(strconv.Itoa(prevVal + lst.Reserve))
			}
		}
	}
}

func (lst *LifeSaverTable) PrintCases() {

	fmt.Println("")
	fmt.Println("PROBLEMATIC CASES: ... ")

	for _, c := range lst.cases {
		fmt.Println("")
		fmt.Println("Case for temp: ", c.temp.String())
		fmt.Println("Decl: ", c.decl.stmt.String())
		fmt.Println("")
		for _, u := range c.used {
			fmt.Println("Used: ", u.stmt.String())
		}
		fmt.Println("")
	}

}

type LifeSaverCase struct {
	temp *Temp
	decl *LifeSaverTableEntry
	used []*LifeSaverTableEntry
}

type LifeSaverTableEntry struct {
	temp           *Temp
	appearanceType string
	stmt           TACStmtI
	centinel       TACStmtI
}

func NewDeclEntry(temp *Temp, stmt TACStmtI) *LifeSaverTableEntry {
	return &LifeSaverTableEntry{
		temp:           temp,
		appearanceType: "decl",
		stmt:           stmt,
	}
}

func NewUsedEntry(temp *Temp, stmt TACStmtI) *LifeSaverTableEntry {
	return &LifeSaverTableEntry{
		temp:           temp,
		appearanceType: "used",
		stmt:           stmt,
	}
}

func NewReturnJumpEntry() *LifeSaverTableEntry {
	return &LifeSaverTableEntry{
		temp:           nil,
		appearanceType: "return",
		stmt:           nil,
	}
}

func NewBlockLifesaver(block *[]TACStmtI, funcName string, factory *TACFactory, returnTemp *Temp, initialOffest int) *Lifesaver {
	return &Lifesaver{
		block:         block,
		funcName:      funcName,
		table:         NewLifeSaverTable(factory, returnTemp),
		factory:       factory,
		initialOffest: initialOffest,
	}
}

func NewRecuriveEntry(stmt TACStmtI) *LifeSaverTableEntry {
	return &LifeSaverTableEntry{
		temp:           nil,
		appearanceType: "recursive",
		stmt:           stmt,
	}
}

func (ls *Lifesaver) EvalBlock() {
	// Get the temps that area affected by a recursive call

	for _, stmt := range *ls.block {
		switch s := stmt.(type) {
		case *SimpleAssignment:

			// t1 = t2
			if temp, ok := s.Val.(*Temp); ok {
				used := NewUsedEntry(temp, s)
				ls.table.AddEntry(used)
			}

			if temp, ok := s.Assignee.(*Temp); ok {
				decl := NewDeclEntry(temp, s)
				ls.table.AddEntry(decl)
			}
			// there is also the case with StackIndexed and HeapIndexed, but we are going to pretend that they are not there :p
		case *CompoundAssignment:
			// t1 = t2 + t3

			if temp, ok := s.Left.(*Temp); ok {
				used := NewUsedEntry(temp, s)
				ls.table.AddEntry(used)
			}

			if temp, ok := s.Right.(*Temp); ok {
				used := NewUsedEntry(temp, s)
				ls.table.AddEntry(used)
			}

			if temp, ok := s.Assignee.(*Temp); ok {
				decl := NewDeclEntry(temp, s)
				ls.table.AddEntry(decl)
			}
		case *MethodCall:
			// function();

			if s.Name == ls.funcName {
				// recursive call
				recursive := NewRecuriveEntry(s)
				ls.table.AddEntry(recursive)
			}
		}
	}

	ls.table.ReduceToCases()
	ls.table.PrintCases()
	ls.FixCases()
}
