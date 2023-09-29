package tac

func (f *TACFactory) reserveParams(name string) []*Temp {

	switch name {
	case "__concat":
		return []*Temp{f.NewTemp(), f.NewTemp(), f.NewTemp()}

	}

	return nil
}

func (f *TACFactory) ConcatBuiltIn() *MethodDcl {
	params := f.GetBuiltinParams("__concat")

	t1 := params[0]
	t2 := params[1]
	t3 := params[2]
	t4 := f.NewTemp()

	/*
		t1 contains the address of the first string
		t2 contains the address of the second string
		t3 contains the address of the new string

		t4: aux temporal

		1. save the address of the heap pointer on t3
		t3 = P

		2. save all the chars of the first string on heap

		save_s1:
			t4 = (int) heap[t1]
			if(t4 == 0) goto end_save_s1
			heap[P] = t4
			P = P + 1
			t1 = t1 + 1
			goto save_s1
		end_save_s1:

		3. save all the chars of the second string on heap

		save_s2:
			t4 = (int) heap[t2]
			if(t4 == 0) goto end_save_s2
			heap[P] = t4
			P = P + 1
			t2 = t2 + 1
			goto save_s2
		end_save_s2:

		4. save the null char on heap
		heap[P] = 0
		P = P + 1
	*/

	// 1. save the address of the heap pointer on t3
	assign1 := f.NewSimpleAssignment().SetAssignee(t3).SetVal(f.NewHeapPtr()) // t3 = P

	// 2. save all the chars of the first string on heap
	saveS1 := f.NewLabel()
	endSaveS1 := f.NewLabel()

	assign2 := f.NewSimpleAssignment().SetAssignee(t4).SetVal(f.NewHeapIndexed().SetIndex(t1)) // t4 = (int) heap[t1]

	condition := f.NewBoolExpression().SetLeft(t4).SetRight(f.NewLiteral().SetValue("0")).SetOp(EQ) // if(t4 == 0)
	conditional1 := f.NewConditionalJump().SetCondition(condition).SetTarget(endSaveS1)             // goto end_save_s1

	assign3 := f.NewSimpleAssignment().SetAssignee(f.NewHeapIndexed().SetIndex(f.NewHeapPtr())).SetVal(t4)                                            // heap[P] = t4
	assign4 := f.NewCompoundAssignment().SetAssignee(f.NewHeapPtr()).SetLeft(f.NewHeapPtr()).SetRight(f.NewLiteral().SetValue("1")).SetOperator(PLUS) // P = P + 1
	assign5 := f.NewCompoundAssignment().SetAssignee(t1).SetLeft(t1).SetRight(f.NewLiteral().SetValue("1")).SetOperator(PLUS)                         // t1 = t1 + 1
	assign6 := f.NewUnconditionalJump().SetTarget(saveS1)                                                                                             // goto save_s1

	// 3. save all the chars of the second string on heap

	saveS2 := f.NewLabel()
	endSaveS2 := f.NewLabel()

	assign7 := f.NewSimpleAssignment().SetAssignee(t4).SetVal(f.NewHeapIndexed().SetIndex(t2)) // t4 = (int) heap[t2]

	condition2 := f.NewBoolExpression().SetLeft(t4).SetRight(f.NewLiteral().SetValue("0")).SetOp(EQ) // if(t4 == 0)
	conditional2 := f.NewConditionalJump().SetCondition(condition2).SetTarget(endSaveS2)             // goto end_save_s2

	assign8 := f.NewSimpleAssignment().SetAssignee(f.NewHeapIndexed().SetIndex(f.NewHeapPtr())).SetVal(t4)                                            // heap[P] = t4
	assign9 := f.NewCompoundAssignment().SetAssignee(f.NewHeapPtr()).SetLeft(f.NewHeapPtr()).SetRight(f.NewLiteral().SetValue("1")).SetOperator(PLUS) // P = P + 1
	assign10 := f.NewCompoundAssignment().SetAssignee(t2).SetLeft(t2).SetRight(f.NewLiteral().SetValue("1")).SetOperator(PLUS)                        // t2 = t2 + 1
	assign11 := f.NewUnconditionalJump().SetTarget(saveS2)                                                                                            // goto save_s2

	// 4. save the null char on heap
	assign12 := f.NewSimpleAssignment().SetAssignee(f.NewHeapIndexed().SetIndex(f.NewHeapPtr())).SetVal(f.NewLiteral().SetValue("0"))                  // heap[P] = 0
	assign13 := f.NewCompoundAssignment().SetAssignee(f.NewHeapPtr()).SetLeft(f.NewHeapPtr()).SetRight(f.NewLiteral().SetValue("1")).SetOperator(PLUS) // P = P + 1

	return &MethodDcl{
		Name: "__concat",
		Block: []TACStmtI{
			assign1,
			saveS1,
			assign2,
			conditional1,
			assign3,
			assign4,
			assign5,
			assign6,
			endSaveS1,
			saveS2,
			assign7,
			conditional2,
			assign8,
			assign9,
			assign10,
			assign11,
			endSaveS2,
			assign12,
			assign13,
		},
	}
}
