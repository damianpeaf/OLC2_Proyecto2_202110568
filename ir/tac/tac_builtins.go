package tac

import (
	"strconv"
)

func (f *TACFactory) reserveParams(name string) []*Temp {

	switch name {
	case "__concat_str":
		return []*Temp{f.NewTemp(), f.NewTemp(), f.NewTemp()}
	case "__print_str":
		return []*Temp{f.NewTemp()}
	case "__zero_division":
		return []*Temp{f.NewTemp()}
	case "__compare_str":
		return []*Temp{f.NewTemp(), f.NewTemp(), f.NewTemp()}
	case "__print_bool":
		return []*Temp{f.NewTemp()}
	case "__and":
		return []*Temp{f.NewTemp(), f.NewTemp(), f.NewTemp()}
	case "__or":
		return []*Temp{f.NewTemp(), f.NewTemp(), f.NewTemp()}
	case "__not":
		return []*Temp{f.NewTemp(), f.NewTemp()}
	case "__alloc_frame":
		return []*Temp{f.NewTemp(), f.NewTemp()} // size, prevFrame
	case "__vector_append":
		return []*Temp{f.NewTemp(), f.NewTemp()} // vector, item
	case "__vector_remove_last":
		return []*Temp{f.NewTemp()} // vector
	case "__vector_remove":
		return []*Temp{f.NewTemp(), f.NewTemp()} // vector, index
	case "__string_to_int":
		return []*Temp{f.NewTemp()} // string and result
	case "__string_to_float":
		return []*Temp{f.NewTemp()} // string and result
	case "__int_to_string":
		return []*Temp{f.NewTemp()} // int and result
	case "__float_to_string":
		return []*Temp{f.NewTemp()} // float and result
	case "__bool_to_string":
		return []*Temp{f.NewTemp()} // bool and result
	}

	return nil
}

func (f *TACFactory) ConcatBuiltIn() *MethodDcl {
	params := f.GetBuiltinParams("__concat_str")

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
		Name: "__concat_str",
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

func (f *TACFactory) PrintBoolBuiltIn() *MethodDcl {
	params := f.GetBuiltinParams("__print_bool")
	block := make(TACBlock, 0)

	t1 := params[0]

	printFalseLabel := f.NewLabel()
	endPrintBoolLabel := f.NewLabel()

	condition := f.NewBoolExpression().SetLeft(t1).SetRight(f.NewLiteral().SetValue("0")).SetOp(EQ) // if(t1 == 0)
	conditional := f.NewConditionalJump().SetCondition(condition).SetTarget(printFalseLabel)        // goto print_false
	block = append(block, conditional)

	printTrue := f.Utility.PrintStringStream("true")
	block = append(block, printTrue...)

	jumpEnd := f.NewUnconditionalJump().SetTarget(endPrintBoolLabel) // goto end_print_bool
	block = append(block, jumpEnd)

	block = append(block, printFalseLabel)
	printFalse := f.Utility.PrintStringStream("false")
	block = append(block, printFalse...)

	block = append(block, endPrintBoolLabel)

	return &MethodDcl{
		Name:  "__print_bool",
		Block: block,
	}
}

func (f *TACFactory) PrintStrBuiltIn() *MethodDcl {

	params := f.GetBuiltinParams("__print_str")

	t1 := params[0]
	t2 := f.NewTemp()

	/*
		t1 contains the address of the string
		t2: aux temporal that will contain the char to print

		print_str:
			t2 = (int) heap[t1]
			if(t2 == 0) goto end_print_str
			print(t2)
			t1 = t1 + 1
			goto print_str
		end_print_str:
	*/

	printStr := f.NewLabel()
	endPrintStr := f.NewLabel()

	assign1 := f.NewSimpleAssignment().SetAssignee(t2).SetVal(f.NewHeapIndexed().SetIndex(t1)) // t2 = (int) heap[t1]

	condition := f.NewBoolExpression().SetLeft(t2).SetRight(f.NewLiteral().SetValue("0")).SetOp(EQ) // if(t2 == 0)
	conditional := f.NewConditionalJump().SetCondition(condition).SetTarget(endPrintStr)            // goto end_print_str

	print := f.NewPrint().SetVal(t2).SetMode(PRINT_CHAR).SetCast("int")

	assign2 := f.NewCompoundAssignment().SetAssignee(t1).SetLeft(t1).SetRight(f.NewLiteral().SetValue("1")).SetOperator(PLUS) // t1 = t1 + 1

	assign3 := f.NewUnconditionalJump().SetTarget(printStr) // goto print_str

	return &MethodDcl{
		Name: "__print_str",
		Block: []TACStmtI{
			printStr,
			assign1,
			conditional,
			print,
			assign2,
			assign3,
			endPrintStr,
		},
	}
}

func (f *TACFactory) ZeroDivisionBuiltIn() *MethodDcl {

	block := make(TACBlock, 0)
	params := f.GetBuiltinParams("__zero_division")

	t1 := params[0] // temporal that will contain the denominator

	/*
		if(t1 != 0) goto no_zero_division
		print(...) MathError
		t1 = 1 // error code
		goto end_zero_division
		no_zero_division:
		t1 = 0
		end_zero_division:
	*/

	noZeroDivision := f.NewLabel()
	endZeroDivision := f.NewLabel()

	condition := f.NewBoolExpression().SetLeft(t1).SetRight(f.NewLiteral().SetValue("0")).SetOp(NEQ).SetLeftCast("int") // if(t1 != 0)
	conditional := f.NewConditionalJump().SetCondition(condition).SetTarget(noZeroDivision)                             // goto no_zero_division
	block = append(block, conditional)

	prints := f.Utility.PrintStringStream("MathError")
	block = append(block, prints...)
	block = append(block, f.NewPrint().SetMode(PRINT_CHAR).SetVal(f.NewLiteral().SetValue(strconv.Itoa('\n'))))

	assign1 := f.NewSimpleAssignment().SetAssignee(t1).SetVal(f.NewLiteral().SetValue("1")) // t1 = 1
	block = append(block, assign1)

	assign2 := f.NewUnconditionalJump().SetTarget(endZeroDivision) // goto end_zero_division
	block = append(block, assign2)

	block = append(block, noZeroDivision)
	assign3 := f.NewSimpleAssignment().SetAssignee(t1).SetVal(f.NewLiteral().SetValue("0")) // t1 = 0
	block = append(block, assign3)

	block = append(block, endZeroDivision)

	return &MethodDcl{
		Name:  "__zero_division",
		Block: block,
	}
}

func (f *TACFactory) AndBuiltIn() *MethodDcl {

	params := f.GetBuiltinParams("__and")
	block := make(TACBlock, 0)
	left := params[0]
	right := params[1]
	result := params[2]

	/*
		if(left == 0) goto false
		if(right == 0) goto false
		temp = 1
		goto end
		false:
			temp = 0
	*/

	falseLabel := f.NewLabel()
	endLabel := f.NewLabel()

	condition1 := f.NewBoolExpression().SetLeft(left).SetRight(f.NewLiteral().SetValue("0")).SetOp(EQ)
	conditional1 := f.NewConditionalJump().SetCondition(condition1).SetTarget(falseLabel)

	condition2 := f.NewBoolExpression().SetLeft(right).SetRight(f.NewLiteral().SetValue("0")).SetOp(EQ)
	conditional2 := f.NewConditionalJump().SetCondition(condition2).SetTarget(falseLabel)

	assign1 := f.NewSimpleAssignment().SetAssignee(result).SetVal(f.NewLiteral().SetValue("1"))
	assign2 := f.NewUnconditionalJump().SetTarget(endLabel)

	block = append(block, conditional1)
	block = append(block, conditional2)
	block = append(block, assign1)
	block = append(block, assign2)
	block = append(block, falseLabel)
	block = append(block, f.NewSimpleAssignment().SetAssignee(result).SetVal(f.NewLiteral().SetValue("0")))
	block = append(block, endLabel)

	return &MethodDcl{
		Name:  "__and",
		Block: block,
	}
}

func (f *TACFactory) OrBuiltIn() *MethodDcl {

	params := f.GetBuiltinParams("__or")
	block := make(TACBlock, 0)

	left := params[0]
	right := params[1]
	result := params[2]

	/*
		if(left != 0) goto true
		if(right != 0) goto true
		temp = 0
		goto end
		true:
			temp = 1
	*/

	trueLabel := f.NewLabel()
	endLabel := f.NewLabel()

	condition1 := f.NewBoolExpression().SetLeft(left).SetRight(f.NewLiteral().SetValue("0")).SetOp(NEQ)
	conditional1 := f.NewConditionalJump().SetCondition(condition1).SetTarget(trueLabel)

	condition2 := f.NewBoolExpression().SetLeft(right).SetRight(f.NewLiteral().SetValue("0")).SetOp(NEQ)
	conditional2 := f.NewConditionalJump().SetCondition(condition2).SetTarget(trueLabel)

	assign1 := f.NewSimpleAssignment().SetAssignee(result).SetVal(f.NewLiteral().SetValue("0"))
	assign2 := f.NewUnconditionalJump().SetTarget(endLabel)

	block = append(block, conditional1)
	block = append(block, conditional2)
	block = append(block, assign1)
	block = append(block, assign2)
	block = append(block, trueLabel)
	block = append(block, f.NewSimpleAssignment().SetAssignee(result).SetVal(f.NewLiteral().SetValue("1")))
	block = append(block, endLabel)

	return &MethodDcl{
		Name:  "__or",
		Block: block,
	}
}

func (f *TACFactory) NotBuiltIn() *MethodDcl {
	params := f.GetBuiltinParams("__not")
	block := make(TACBlock, 0)

	left := params[0]
	result := params[1]

	/*
		if(left == 0) goto true
		temp = 0
		goto end
		true:
			temp = 1
	*/

	trueLabel := f.NewLabel()
	endLabel := f.NewLabel()

	condition := f.NewBoolExpression().SetLeft(left).SetRight(f.NewLiteral().SetValue("0")).SetOp(EQ)
	conditional := f.NewConditionalJump().SetCondition(condition).SetTarget(trueLabel)

	assign1 := f.NewSimpleAssignment().SetAssignee(result).SetVal(f.NewLiteral().SetValue("0"))
	assign2 := f.NewUnconditionalJump().SetTarget(endLabel)

	block = append(block, conditional)
	block = append(block, assign1)
	block = append(block, assign2)
	block = append(block, trueLabel)
	block = append(block, f.NewSimpleAssignment().SetAssignee(result).SetVal(f.NewLiteral().SetValue("1")))
	block = append(block, endLabel)

	return &MethodDcl{
		Name:  "__not",
		Block: block,
	}
}

func (f *TACFactory) CompareStrBuiltIn() *MethodDcl {

	block := make(TACBlock, 0)
	params := f.GetBuiltinParams("__compare_str")

	s1 := params[0] // temporal that will contain the first string
	s2 := params[1] // temporal that will contain the second string
	r := params[2]  // temporal that will contain the result

	// return values:
	// 0 -> s1 == s2
	// -1 -> s1 < s2
	// 1 -> s1 > s2

	/*
		t1: will contain the char of the first string
		t2: will contain the char of the second string

		cmp_str:
			t1 = (int) heap[s1]
			t2 = (int) heap[s2]

			if(t1 == 0) goto end_of_str_1
			if(t2 == 0) goto s1_greater_than_s2

			if(t1 < t2) goto s1_less_than_s2
			if(t1 > t2) goto s1_greater_than_s2

			next_char:
				s1 = s1 + 1
				s2 = s2 + 1
				goto cmp_str

		end_of_str_1:
			if(t2 == 0) goto equal_str
			goto s1_less_than_s2

		equal_str:
			r = 0
			goto end_cmp_str

		s1_less_than_s2:
			r = -1
			goto end_cmp_str

		s1_greater_than_s2:
			r = 1
			goto end_cmp_str

		end_cmp_str:
	*/

	// temps
	t1 := f.NewTemp()
	t2 := f.NewTemp()

	// labels:
	cmpStr := f.NewLabel()
	endOfStr1 := f.NewLabel()
	s1GreaterThanS2 := f.NewLabel()
	s1LessThanS2 := f.NewLabel()
	equalStr := f.NewLabel()
	endCmpStr := f.NewLabel()

	// cmp_str:
	block = append(block, cmpStr)

	// t1 = (int) heap[s1]
	assign1 := f.NewSimpleAssignment().SetAssignee(t1).SetVal(f.NewHeapIndexed().SetIndex(s1))
	block = append(block, assign1)

	// t2 = (int) heap[s2]
	assign2 := f.NewSimpleAssignment().SetAssignee(t2).SetVal(f.NewHeapIndexed().SetIndex(s2))
	block = append(block, assign2)

	// if(t1 == 0) goto end_of_str_1
	condition1 := f.NewBoolExpression().SetLeft(t1).SetLeftCast("int").SetRight(f.NewLiteral().SetValue("0")).SetOp(EQ)
	conditional1 := f.NewConditionalJump().SetCondition(condition1).SetTarget(endOfStr1)
	block = append(block, conditional1)

	// if(t2 == 0) goto s1_greater_than_s2
	condition2 := f.NewBoolExpression().SetLeft(t2).SetLeftCast("int").SetRight(f.NewLiteral().SetValue("0")).SetOp(EQ)
	conditional2 := f.NewConditionalJump().SetCondition(condition2).SetTarget(s1GreaterThanS2)
	block = append(block, conditional2)

	// if(t1 < t2) goto s1_less_than_s2
	condition3 := f.NewBoolExpression().SetLeft(t1).SetLeftCast("int").SetRight(t2).SetOp(LT).SetRightCast("int")
	conditional3 := f.NewConditionalJump().SetCondition(condition3).SetTarget(s1LessThanS2)
	block = append(block, conditional3)

	// if(t1 > t2) goto s1_greater_than_s2
	condition4 := f.NewBoolExpression().SetLeft(t1).SetLeftCast("int").SetRight(t2).SetOp(GT).SetRightCast("int")
	conditional4 := f.NewConditionalJump().SetCondition(condition4).SetTarget(s1GreaterThanS2)
	block = append(block, conditional4)

	// s1 = s1 + 1
	assign3 := f.NewCompoundAssignment().SetAssignee(s1).SetLeft(s1).SetRight(f.NewLiteral().SetValue("1")).SetOperator(PLUS)
	block = append(block, assign3)

	// s2 = s2 + 1
	assign4 := f.NewCompoundAssignment().SetAssignee(s2).SetLeft(s2).SetRight(f.NewLiteral().SetValue("1")).SetOperator(PLUS)
	block = append(block, assign4)

	// goto cmp_str
	block = append(block, f.NewUnconditionalJump().SetTarget(cmpStr))

	// end_of_str_1:
	block = append(block, endOfStr1)

	// if(t2 == 0) goto equal_str
	condition5 := f.NewBoolExpression().SetLeft(t2).SetRight(f.NewLiteral().SetValue("0")).SetOp(EQ)
	conditional5 := f.NewConditionalJump().SetCondition(condition5).SetTarget(equalStr)
	block = append(block, conditional5)

	// goto s1_less_than_s2
	block = append(block, f.NewUnconditionalJump().SetTarget(s1LessThanS2))

	// equal_str:
	block = append(block, equalStr)

	// r = 0
	assign5 := f.NewSimpleAssignment().SetAssignee(r).SetVal(f.NewLiteral().SetValue("0"))
	block = append(block, assign5)

	// goto end_cmp_str
	block = append(block, f.NewUnconditionalJump().SetTarget(endCmpStr))

	// s1_less_than_s2:
	block = append(block, s1LessThanS2)

	// r = -1
	assign6 := f.NewSimpleAssignment().SetAssignee(r).SetVal(f.NewLiteral().SetValue("-1"))
	block = append(block, assign6)

	// goto end_cmp_str
	block = append(block, f.NewUnconditionalJump().SetTarget(endCmpStr))

	// s1_greater_than_s2:
	block = append(block, s1GreaterThanS2)

	// r = 1
	assign7 := f.NewSimpleAssignment().SetAssignee(r).SetVal(f.NewLiteral().SetValue("1"))
	block = append(block, assign7)

	// end_cmp_str:
	block = append(block, endCmpStr)

	return &MethodDcl{
		Name:  "__compare_str",
		Block: block,
	}
}

func (f *TACFactory) AllocFrameBuiltIn() *MethodDcl {
	block := make(TACBlock, 0)
	params := f.GetBuiltinParams("__alloc_frame")

	size := params[0]
	prevFrame := params[1]
	auxCount := f.NewTemp()
	p := f.NewStackPtr()

	endLabel := f.NewLabel()
	allocLabel := f.NewLabel()

	/*
		save on stack
		1. prevFrame
		2. reserved 'size' bytes
	*/

	// 1. prevFrame
	assign1 := f.NewSimpleAssignment().SetAssignee(f.NewStackIndexed().SetIndex(p)).SetVal(prevFrame)                        //  stack[SP] = prevFrame
	increase := f.NewCompoundAssignment().SetAssignee(p).SetLeft(p).SetRight(f.NewLiteral().SetValue("1")).SetOperator(PLUS) // SP = SP + 1
	block = append(block, assign1)
	block = append(block, increase)

	// 2. reserved 'size' bytes
	// auxCount = 0
	assign2 := f.NewSimpleAssignment().SetAssignee(auxCount).SetVal(f.NewLiteral().SetValue("0"))
	block = append(block, assign2)
	block = append(block, allocLabel)
	// if auxCount == size goto end
	condition := f.NewBoolExpression().SetLeft(auxCount).SetRight(size).SetOp(EQ)
	conditional := f.NewConditionalJump().SetCondition(condition).SetTarget(endLabel)
	block = append(block, conditional)
	// stack[SP] = 0
	assign3 := f.NewSimpleAssignment().SetAssignee(f.NewStackIndexed().SetIndex(p)).SetVal(f.NewLiteral().SetValue("0"))
	block = append(block, assign3)
	// SP = SP + 1
	increase2 := f.NewCompoundAssignment().SetAssignee(p).SetLeft(p).SetRight(f.NewLiteral().SetValue("1")).SetOperator(PLUS)
	block = append(block, increase2)
	// auxCount = auxCount + 1
	assign4 := f.NewCompoundAssignment().SetAssignee(auxCount).SetLeft(auxCount).SetRight(f.NewLiteral().SetValue("1")).SetOperator(PLUS)
	block = append(block, assign4)
	// goto allocLabel
	block = append(block, f.NewUnconditionalJump().SetTarget(allocLabel))
	// endLabel
	block = append(block, endLabel)

	return &MethodDcl{
		Name:  "__alloc_frame",
		Block: block,
	}

}

func (f *TACFactory) VectorAppendBuiltIn() *MethodDcl {
	block := make(TACBlock, 0)
	params := f.GetBuiltinParams("__vector_append")

	vectorVarIndex := params[0]
	item := params[1]

	/*
		vectorHeapAddress = stack[vectorVarIndex] // prev address of the vector in the heap
		stack[vectorVarIndex] = heapPtr // new address of the vector

		// new size
		size = heap[vectorHeapAddress]
		size = size + 1

		heap[heapPtr] = size
		hp = hp + 1

		// clone the vector
		size = size - 1 // size of the previous vector

		current = 0
		vectorHeapAddress = vectorHeapAddress + 1 // skip the size

		clone:
		if(size == current) goto end_clone
		aux = heap[vectorHeapAddress]
		heap[hp] = aux

		hp = hp + 1
		vectorHeapAddress = vectorHeapAddress + 1
		current = current + 1
		goto clone
		end_clone:

		// add the new item
		heap[hp] = item
		hp = hp + 1


	*/

	vectorHeapAddress := f.NewTemp()
	size := f.NewTemp()
	heapPtr := f.NewHeapPtr()
	aux := f.NewTemp()
	current := f.NewTemp()

	cloneLabel := f.NewLabel()
	endCloneLabel := f.NewLabel()

	// vectorHeapAddress = stack[vectorVarIndex]
	assign1 := f.NewSimpleAssignment().SetAssignee(vectorHeapAddress).SetVal(f.NewStackIndexed().SetIndex(vectorVarIndex))
	block = append(block, assign1)

	// stack[vectorVarIndex] = heapPtr
	assign2 := f.NewSimpleAssignment().SetAssignee(f.NewStackIndexed().SetIndex(vectorVarIndex)).SetVal(heapPtr)
	block = append(block, assign2)

	// size = heap[vectorHeapAddress]
	assign3 := f.NewSimpleAssignment().SetAssignee(size).SetVal(f.NewHeapIndexed().SetIndex(vectorHeapAddress))
	block = append(block, assign3)

	// size = size + 1
	assign4 := f.NewCompoundAssignment().SetAssignee(size).SetLeft(size).SetRight(f.NewLiteral().SetValue("1")).SetOperator(PLUS)
	block = append(block, assign4)

	// heap[heapPtr] = size
	assign5 := f.NewSimpleAssignment().SetAssignee(f.NewHeapIndexed().SetIndex(heapPtr)).SetVal(size)
	block = append(block, assign5)

	// hp = hp + 1
	assign6 := f.NewCompoundAssignment().SetAssignee(f.NewHeapPtr()).SetLeft(f.NewHeapPtr()).SetRight(f.NewLiteral().SetValue("1")).SetOperator(PLUS)
	block = append(block, assign6)

	// size = size - 1
	assign7 := f.NewCompoundAssignment().SetAssignee(size).SetLeft(size).SetRight(f.NewLiteral().SetValue("1")).SetOperator(MINUS)
	block = append(block, assign7)

	// current = 0
	assign8 := f.NewSimpleAssignment().SetAssignee(current).SetVal(f.NewLiteral().SetValue("0"))
	block = append(block, assign8)

	// vectorHeapAddress = vectorHeapAddress + 1
	assign9 := f.NewCompoundAssignment().SetAssignee(vectorHeapAddress).SetLeft(vectorHeapAddress).SetRight(f.NewLiteral().SetValue("1")).SetOperator(PLUS)
	block = append(block, assign9)

	// clone:
	block = append(block, cloneLabel)

	// if(size == current) goto end_clone
	condition := f.NewBoolExpression().SetLeft(size).SetRight(current).SetOp(EQ).SetLeftCast("int").SetRightCast("int")
	conditional := f.NewConditionalJump().SetCondition(condition).SetTarget(endCloneLabel)
	block = append(block, conditional)

	// aux = heap[vectorHeapAddress]
	assign10 := f.NewSimpleAssignment().SetAssignee(aux).SetVal(f.NewHeapIndexed().SetIndex(vectorHeapAddress))
	block = append(block, assign10)

	// heap[hp] = aux
	assign11 := f.NewSimpleAssignment().SetAssignee(f.NewHeapIndexed().SetIndex(f.NewHeapPtr())).SetVal(aux)
	block = append(block, assign11)

	// hp = hp + 1
	assign12 := f.NewCompoundAssignment().SetAssignee(f.NewHeapPtr()).SetLeft(f.NewHeapPtr()).SetRight(f.NewLiteral().SetValue("1")).SetOperator(PLUS)
	block = append(block, assign12)

	// vectorHeapAddress = vectorHeapAddress + 1
	assign13 := f.NewCompoundAssignment().SetAssignee(vectorHeapAddress).SetLeft(vectorHeapAddress).SetRight(f.NewLiteral().SetValue("1")).SetOperator(PLUS)
	block = append(block, assign13)

	// current = current + 1
	assign14 := f.NewCompoundAssignment().SetAssignee(current).SetLeft(current).SetRight(f.NewLiteral().SetValue("1")).SetOperator(PLUS)
	block = append(block, assign14)

	// goto clone
	block = append(block, f.NewUnconditionalJump().SetTarget(cloneLabel))

	// end_clone:
	block = append(block, endCloneLabel)

	// heap[hp] = item
	assign15 := f.NewSimpleAssignment().SetAssignee(f.NewHeapIndexed().SetIndex(f.NewHeapPtr())).SetVal(item)
	block = append(block, assign15)

	// hp = hp + 1
	assign16 := f.NewCompoundAssignment().SetAssignee(f.NewHeapPtr()).SetLeft(f.NewHeapPtr()).SetRight(f.NewLiteral().SetValue("1")).SetOperator(PLUS)
	block = append(block, assign16)

	return &MethodDcl{
		Name:  "__vector_append",
		Block: block,
	}
}

func (f *TACFactory) VectorRemoveLastBuiltIn() *MethodDcl {
	block := make(TACBlock, 0)
	params := f.GetBuiltinParams("__vector_remove_last")

	vectorVarIndex := params[0]

	/*
		// just decrease the size of the vector
		vectorHeapAddress = stack[vectorVarIndex] // prev address of the vector in the heap
		size = heap[vectorHeapAddress]
		size = size - 1
		heap[vectorHeapAddress] = size
	*/

	vectorHeapAddress := f.NewTemp()
	size := f.NewTemp()

	// vectorHeapAddress = stack[vectorVarIndex]
	assign1 := f.NewSimpleAssignment().SetAssignee(vectorHeapAddress).SetVal(f.NewStackIndexed().SetIndex(vectorVarIndex))
	block = append(block, assign1)

	// size = heap[vectorHeapAddress]
	assign2 := f.NewSimpleAssignment().SetAssignee(size).SetVal(f.NewHeapIndexed().SetIndex(vectorHeapAddress))
	block = append(block, assign2)

	// size = size - 1
	assign3 := f.NewCompoundAssignment().SetAssignee(size).SetLeft(size).SetRight(f.NewLiteral().SetValue("1")).SetOperator(MINUS)
	block = append(block, assign3)

	// heap[vectorHeapAddress] = size
	assign4 := f.NewSimpleAssignment().SetAssignee(f.NewHeapIndexed().SetIndex(vectorHeapAddress)).SetVal(size)
	block = append(block, assign4)

	return &MethodDcl{
		Name:  "__vector_remove_last",
		Block: block,
	}
}

func (f *TACFactory) VectorRemoveBuiltIn() *MethodDcl {
	block := make(TACBlock, 0)
	params := f.GetBuiltinParams("__vector_remove")

	vectorVarIndex := params[0]
	removeIndex := params[1]

	/*
		// check if the index is valid

		vectorHeapAddress = stack[vectorVarIndex] // prev address of the vector in the heap
		size = heap[vectorHeapAddress]

		if(removeIndex < 0) goto error
		if(removeIndex >= size) goto error

		// remove the item
		stack[vectorVarIndex] = heapPtr // new address of the vector

		// new size
		size = size - 1

		heap[heapPtr] = size
		hp = hp + 1

		size = size + 1 // just to traverse all the items

		// clone the vector, except the item to remove
		current = 0
		vectorHeapAddress = vectorHeapAddress + 1 // skip the size

		clone:
		if(size == current) goto end_clone
		if(current == removeIndex) goto skip_item

		aux = heap[vectorHeapAddress]
		heap[hp] = aux

		hp = hp + 1

		skip_item:
		vectorHeapAddress = vectorHeapAddress + 1
		current = current + 1
		goto clone

		error:
		print("BoundsError")
		goto end_clone

		end_clone:
	*/

	vectorHeapAddress := f.NewTemp()
	size := f.NewTemp()
	aux := f.NewTemp()
	current := f.NewTemp()
	heapPtr := f.NewHeapPtr()

	cloneLabel := f.NewLabel()
	endCloneLabel := f.NewLabel()
	skipItemLabel := f.NewLabel()
	errorLabel := f.NewLabel()

	// vectorHeapAddress = stack[vectorVarIndex]
	assign1 := f.NewSimpleAssignment().SetAssignee(vectorHeapAddress).SetVal(f.NewStackIndexed().SetIndex(vectorVarIndex))
	block = append(block, assign1)

	// size = heap[vectorHeapAddress]
	assign2 := f.NewSimpleAssignment().SetAssignee(size).SetVal(f.NewHeapIndexed().SetIndex(vectorHeapAddress))
	block = append(block, assign2)

	// if(removeIndex < 0) goto error
	condition1 := f.NewBoolExpression().SetLeft(removeIndex).SetRight(f.NewLiteral().SetValue("0")).SetOp(LT).SetLeftCast("int").SetRightCast("int")
	conditional1 := f.NewConditionalJump().SetCondition(condition1).SetTarget(errorLabel)
	block = append(block, conditional1)

	// if(removeIndex >= size) goto error
	condition2 := f.NewBoolExpression().SetLeft(removeIndex).SetRight(size).SetOp(GTE).SetLeftCast("int").SetRightCast("int")
	conditional2 := f.NewConditionalJump().SetCondition(condition2).SetTarget(errorLabel)
	block = append(block, conditional2)

	// stack[vectorVarIndex] = heapPtr
	assign3 := f.NewSimpleAssignment().SetAssignee(f.NewStackIndexed().SetIndex(vectorVarIndex)).SetVal(heapPtr)
	block = append(block, assign3)

	// size = size - 1
	assign4 := f.NewCompoundAssignment().SetAssignee(size).SetLeft(size).SetRight(f.NewLiteral().SetValue("1")).SetOperator(MINUS)
	block = append(block, assign4)

	// heap[heapPtr] = size
	assign5 := f.NewSimpleAssignment().SetAssignee(f.NewHeapIndexed().SetIndex(heapPtr)).SetVal(size)
	block = append(block, assign5)

	// hp = hp + 1
	assign6 := f.NewCompoundAssignment().SetAssignee(f.NewHeapPtr()).SetLeft(f.NewHeapPtr()).SetRight(f.NewLiteral().SetValue("1")).SetOperator(PLUS)
	block = append(block, assign6)

	// size = size + 1
	assign7 := f.NewCompoundAssignment().SetAssignee(size).SetLeft(size).SetRight(f.NewLiteral().SetValue("1")).SetOperator(PLUS)
	block = append(block, assign7)

	// current = 0
	assign8 := f.NewSimpleAssignment().SetAssignee(current).SetVal(f.NewLiteral().SetValue("0"))
	block = append(block, assign8)

	// vectorHeapAddress = vectorHeapAddress + 1
	assign9 := f.NewCompoundAssignment().SetAssignee(vectorHeapAddress).SetLeft(vectorHeapAddress).SetRight(f.NewLiteral().SetValue("1")).SetOperator(PLUS)
	block = append(block, assign9)

	// clone:
	block = append(block, cloneLabel)

	// if(size == current) goto end_clone
	condition3 := f.NewBoolExpression().SetLeft(size).SetRight(current).SetOp(EQ).SetLeftCast("int").SetRightCast("int")
	conditional3 := f.NewConditionalJump().SetCondition(condition3).SetTarget(endCloneLabel)
	block = append(block, conditional3)

	// if(current == removeIndex) goto skip_item
	condition4 := f.NewBoolExpression().SetLeft(current).SetRight(removeIndex).SetOp(EQ).SetLeftCast("int").SetRightCast("int")
	conditional4 := f.NewConditionalJump().SetCondition(condition4).SetTarget(skipItemLabel)
	block = append(block, conditional4)

	// aux = heap[vectorHeapAddress]
	assign10 := f.NewSimpleAssignment().SetAssignee(aux).SetVal(f.NewHeapIndexed().SetIndex(vectorHeapAddress))
	block = append(block, assign10)

	// heap[hp] = aux
	assign11 := f.NewSimpleAssignment().SetAssignee(f.NewHeapIndexed().SetIndex(f.NewHeapPtr())).SetVal(aux)
	block = append(block, assign11)

	// hp = hp + 1
	assign12 := f.NewCompoundAssignment().SetAssignee(f.NewHeapPtr()).SetLeft(f.NewHeapPtr()).SetRight(f.NewLiteral().SetValue("1")).SetOperator(PLUS)
	block = append(block, assign12)

	// skip_item:
	block = append(block, skipItemLabel)

	// vectorHeapAddress = vectorHeapAddress + 1
	assign13 := f.NewCompoundAssignment().SetAssignee(vectorHeapAddress).SetLeft(vectorHeapAddress).SetRight(f.NewLiteral().SetValue("1")).SetOperator(PLUS)
	block = append(block, assign13)

	// current = current + 1
	assign14 := f.NewCompoundAssignment().SetAssignee(current).SetLeft(current).SetRight(f.NewLiteral().SetValue("1")).SetOperator(PLUS)
	block = append(block, assign14)

	// goto clone
	block = append(block, f.NewUnconditionalJump().SetTarget(cloneLabel))

	// error:
	block = append(block, errorLabel)

	// print("BoundsError")
	prints := f.Utility.PrintStringStream("BoundsError\n")
	block = append(block, prints...)

	// end_clone:
	block = append(block, endCloneLabel)

	return &MethodDcl{
		Name:  "__vector_remove",
		Block: block,
	}
}

func (f *TACFactory) StringToIntBuiltIn() *MethodDcl {
	params := f.GetBuiltinParams("__string_to_int")
	block := make(TACBlock, 0)
	stringHeapAddress := params[0]

	aux := f.NewTemp()

	/*
		 "0" is 48 in ascii, so we can use this to convert a string to an int
		 "." is 46 in ascii, so we can use this to check if the string is a float

		//  get the last char
		stringFinalAddress = stringHeapAddress
		aux = heap[stringFinalAddress]
		last:
		if aux == 0 goto lastEnd:
		if aux == 46 goto lastEnd:
		stringFinalAddress = stringFinalAddress + 1
		aux = heap[stringFinalAddress]
		goto last
		lastEnd:

		result = 0
		multiplier = 1

		// convert the string to int
		stringFinalAddress = stringFinalAddress - 1

		convert:
		if(stringFinalAddress < stringHeapAddress) goto endConvert
		aux = heap[stringFinalAddress]
		aux = aux - 48

		if aux < 0 goto error
		if aux > 9 goto error

		aux = aux * multiplier
		result = result + aux
		multiplier = multiplier * 10
		stringFinalAddress = stringFinalAddress - 1
		goto convert

		endConvert:
		stringHeapAddress = result
		goto end

		error:
		stringHeapAddress = nil

		end:
	*/

	stringFinalAddress := f.NewTemp()
	result := f.NewTemp()
	multiplier := f.NewTemp()

	lastLabel := f.NewLabel()
	lastEndLabel := f.NewLabel()
	errorLabel := f.NewLabel()
	endLabel := f.NewLabel()

	// stringFinalAddress = stringHeapAddress
	assign1 := f.NewSimpleAssignment().SetAssignee(stringFinalAddress).SetVal(stringHeapAddress)
	block = append(block, assign1)

	// aux = heap[stringFinalAddress]
	assign2 := f.NewSimpleAssignment().SetAssignee(aux).SetVal(f.NewHeapIndexed().SetIndex(stringFinalAddress))
	block = append(block, assign2)

	// last:
	block = append(block, lastLabel)

	// if aux == 0 goto lastEnd:
	condition := f.NewBoolExpression().SetLeft(aux).SetRight(f.NewLiteral().SetValue("0")).SetOp(EQ)
	conditional := f.NewConditionalJump().SetCondition(condition).SetTarget(lastEndLabel)
	block = append(block, conditional)

	// if aux == 46 goto lastEnd:
	condition2 := f.NewBoolExpression().SetLeft(aux).SetRight(f.NewLiteral().SetValue("46")).SetOp(EQ)
	conditional2 := f.NewConditionalJump().SetCondition(condition2).SetTarget(lastEndLabel)
	block = append(block, conditional2)

	// stringFinalAddress = stringFinalAddress + 1
	assign3 := f.NewCompoundAssignment().SetAssignee(stringFinalAddress).SetLeft(stringFinalAddress).SetRight(f.NewLiteral().SetValue("1")).SetOperator(PLUS)
	block = append(block, assign3)

	// aux = heap[stringFinalAddress]
	assign4 := f.NewSimpleAssignment().SetAssignee(aux).SetVal(f.NewHeapIndexed().SetIndex(stringFinalAddress))
	block = append(block, assign4)

	// goto last
	block = append(block, f.NewUnconditionalJump().SetTarget(lastLabel))

	// lastEnd:
	block = append(block, lastEndLabel)

	// result = 0
	assign5 := f.NewSimpleAssignment().SetAssignee(result).SetVal(f.NewLiteral().SetValue("0"))
	block = append(block, assign5)

	// multiplier = 1
	assign6 := f.NewSimpleAssignment().SetAssignee(multiplier).SetVal(f.NewLiteral().SetValue("1"))
	block = append(block, assign6)

	// stringFinalAddress = stringFinalAddress - 1
	assign7 := f.NewCompoundAssignment().SetAssignee(stringFinalAddress).SetLeft(stringFinalAddress).SetRight(f.NewLiteral().SetValue("1")).SetOperator(MINUS)
	block = append(block, assign7)

	// convert:
	convertLabel := f.NewLabel()
	endConvertLabel := f.NewLabel()
	block = append(block, convertLabel)

	// if(stringFinalAddress == stringHeapAddress) goto endConvert
	condition3 := f.NewBoolExpression().SetLeft(stringFinalAddress).SetRight(stringHeapAddress).SetOp(LT)
	conditional3 := f.NewConditionalJump().SetCondition(condition3).SetTarget(endConvertLabel)
	block = append(block, conditional3)

	// aux = heap[stringFinalAddress]
	assign8 := f.NewSimpleAssignment().SetAssignee(aux).SetVal(f.NewHeapIndexed().SetIndex(stringFinalAddress))
	block = append(block, assign8)

	// aux = aux - 48
	assign9 := f.NewCompoundAssignment().SetAssignee(aux).SetLeft(aux).SetRight(f.NewLiteral().SetValue("48")).SetOperator(MINUS)
	block = append(block, assign9)

	// if aux < 0 goto error
	condition4 := f.NewBoolExpression().SetLeft(aux).SetRight(f.NewLiteral().SetValue("0")).SetOp(LT)
	conditional4 := f.NewConditionalJump().SetCondition(condition4).SetTarget(errorLabel)
	block = append(block, conditional4)

	// if aux > 9 goto error
	condition5 := f.NewBoolExpression().SetLeft(aux).SetRight(f.NewLiteral().SetValue("9")).SetOp(GT)
	conditional5 := f.NewConditionalJump().SetCondition(condition5).SetTarget(errorLabel)
	block = append(block, conditional5)

	// aux = aux * multiplier
	assign10 := f.NewCompoundAssignment().SetAssignee(aux).SetLeft(aux).SetRight(multiplier).SetOperator(MULTIPLY)
	block = append(block, assign10)

	// result = result + aux
	assign11 := f.NewCompoundAssignment().SetAssignee(result).SetLeft(result).SetRight(aux).SetOperator(PLUS)
	block = append(block, assign11)

	// multiplier = multiplier * 10
	assign12 := f.NewCompoundAssignment().SetAssignee(multiplier).SetLeft(multiplier).SetRight(f.NewLiteral().SetValue("10")).SetOperator(MULTIPLY)
	block = append(block, assign12)

	// stringFinalAddress = stringFinalAddress - 1
	assign13 := f.NewCompoundAssignment().SetAssignee(stringFinalAddress).SetLeft(stringFinalAddress).SetRight(f.NewLiteral().SetValue("1")).SetOperator(MINUS)
	block = append(block, assign13)

	// goto convert
	block = append(block, f.NewUnconditionalJump().SetTarget(convertLabel))

	// endConvert:
	block = append(block, endConvertLabel)

	// return result
	assign14 := f.NewSimpleAssignment().SetAssignee(stringHeapAddress).SetVal(result)
	block = append(block, assign14)

	// goto end
	block = append(block, f.NewUnconditionalJump().SetTarget(endLabel))

	// error:
	block = append(block, errorLabel)

	nilValue := f.Utility.NilValue()
	assign15 := f.NewSimpleAssignment().SetAssignee(stringHeapAddress).SetVal(nilValue)
	block = append(block, assign15)

	// end:
	block = append(block, endLabel)

	return &MethodDcl{
		Name:  "__string_to_int",
		Block: block,
	}
}

func (f *TACFactory) StringToFloatBuiltIn() *MethodDcl {
	params := f.GetBuiltinParams("__string_to_float")
	block := make(TACBlock, 0)
	stringHeapAddress := params[0]

	intToStringParams := f.GetBuiltinParams("__string_to_int")
	stringParam := intToStringParams[0]

	// divide the string at ".", so we can convert the left and right parts to int, and then divide them

	/*
		    // call __string_to_int to convert the left part of the string to int
			stringParam = stringHeapAddress
			__string_to_int()
			left = stringParam

			// check if its nil
			if(left == nil) goto error

			// search for the "."
			decimalPartStart = stringHeapAddress
			aux = heap[decimalPartStart]
			decimalPartStartLabel:
			if(aux == 46) goto decimalPartEnd
			if(aux == 0) goto decimalPartEnd
			decimalPartStart = decimalPartStart + 1
			aux = heap[decimalPartStart]
			goto decimalPartStartLabel
			decimalPartEnd:

			// call __string_to_int to convert the right part of the string to int
			stringParam = decimalPartStart + 1
			__string_to_int()
			right = stringParam

			// check if its nil
			if(right == nil) goto error

			// Get the final address of the string
			stringFinalAddress = stringHeapAddress
			aux = heap[stringFinalAddress]
			last:
			if aux == 0 goto lastEnd:
			stringFinalAddress = stringFinalAddress + 1
			aux = heap[stringFinalAddress]
			goto last
			lastEnd:


			// convert the right part to float
			decimalSize = decimalPartStart - stringFinalAddress // .05 (end) -> 2
			divisor = 1
			current = 0
			convert:
			if(current == decimalSize) goto endConvert
			divisor = divisor * 10
			current = current + 1
			goto convert

			endConvert:
			right = right / divisor

			// convert the left part to float
			result = left + right

			// return result
			stringHeapAddress = result
			goto end

			error:
			stringHeapAddress = nil

			end:
	*/

	left := f.NewTemp()
	right := f.NewTemp()
	decimalPartStart := f.NewTemp()
	aux := f.NewTemp()
	stringFinalAddress := f.NewTemp()
	decimalSize := f.NewTemp()
	divisor := f.NewTemp()
	current := f.NewTemp()
	result := f.NewTemp()

	decimalPartEndLabel := f.NewLabel()
	decimalPartStartLabel := f.NewLabel()
	conversionLabel := f.NewLabel()
	endConvertLabel := f.NewLabel()
	lastLabel := f.NewLabel()
	lastEndLabel := f.NewLabel()
	errorLabel := f.NewLabel()
	endLabel := f.NewLabel()

	// stringParam = stringHeapAddress
	assign1 := f.NewSimpleAssignment().SetAssignee(stringParam).SetVal(stringHeapAddress)
	block = append(block, assign1)

	// __string_to_int()
	call1 := f.NewMethodCall("__string_to_int")
	block = append(block, call1)

	// left = stringParam
	assign2 := f.NewSimpleAssignment().SetAssignee(left).SetVal(stringParam)
	block = append(block, assign2)

	// if(left == nil) goto error
	condition := f.NewBoolExpression().SetLeft(left).SetRight(f.Utility.NilValue()).SetOp(EQ)
	conditional := f.NewConditionalJump().SetCondition(condition).SetTarget(errorLabel)
	block = append(block, conditional)

	// decimalPartStart = stringHeapAddress
	assign3 := f.NewSimpleAssignment().SetAssignee(decimalPartStart).SetVal(stringHeapAddress)
	block = append(block, assign3)

	// aux = heap[decimalPartStart]
	assign4 := f.NewSimpleAssignment().SetAssignee(aux).SetVal(f.NewHeapIndexed().SetIndex(decimalPartStart))
	block = append(block, assign4)

	// decimalPartStartLabel:
	block = append(block, decimalPartStartLabel)

	// if(aux == 46) goto decimalPartEnd
	condition2 := f.NewBoolExpression().SetLeft(aux).SetRight(f.NewLiteral().SetValue("46")).SetOp(EQ)
	conditional2 := f.NewConditionalJump().SetCondition(condition2).SetTarget(decimalPartEndLabel)
	block = append(block, conditional2)

	// if aux == 0 goto decimalPartEnd
	condition30 := f.NewBoolExpression().SetLeft(aux).SetRight(f.NewLiteral().SetValue("0")).SetOp(EQ)
	conditional30 := f.NewConditionalJump().SetCondition(condition30).SetTarget(decimalPartEndLabel)
	block = append(block, conditional30)

	// decimalPartStart = decimalPartStart + 1
	assign5 := f.NewCompoundAssignment().SetAssignee(decimalPartStart).SetLeft(decimalPartStart).SetRight(f.NewLiteral().SetValue("1")).SetOperator(PLUS)
	block = append(block, assign5)

	// aux = heap[decimalPartStart]
	assign6 := f.NewSimpleAssignment().SetAssignee(aux).SetVal(f.NewHeapIndexed().SetIndex(decimalPartStart))
	block = append(block, assign6)

	// goto decimalPartStartLabel
	block = append(block, f.NewUnconditionalJump().SetTarget(decimalPartStartLabel))

	// decimalPartEnd:
	block = append(block, decimalPartEndLabel)

	// decimalPartStart = decimalPartStart + 1
	assign7 := f.NewCompoundAssignment().SetAssignee(decimalPartStart).SetLeft(decimalPartStart).SetRight(f.NewLiteral().SetValue("1")).SetOperator(PLUS)
	block = append(block, assign7)

	// stringParam = decimalPartStart
	assign80 := f.NewSimpleAssignment().SetAssignee(stringParam).SetVal(decimalPartStart)
	block = append(block, assign80)

	// __string_to_int()
	call2 := f.NewMethodCall("__string_to_int")
	block = append(block, call2)

	// right = stringParam
	assign8 := f.NewSimpleAssignment().SetAssignee(right).SetVal(stringParam)
	block = append(block, assign8)

	// if(right == nil) goto error
	condition3 := f.NewBoolExpression().SetLeft(right).SetRight(f.Utility.NilValue()).SetOp(EQ)
	conditional3 := f.NewConditionalJump().SetCondition(condition3).SetTarget(errorLabel)
	block = append(block, conditional3)

	// stringFinalAddress = stringHeapAddress
	assign9 := f.NewSimpleAssignment().SetAssignee(stringFinalAddress).SetVal(stringHeapAddress)
	block = append(block, assign9)

	// aux = heap[stringFinalAddress]
	assign10 := f.NewSimpleAssignment().SetAssignee(aux).SetVal(f.NewHeapIndexed().SetIndex(stringFinalAddress))
	block = append(block, assign10)

	// last:
	block = append(block, lastLabel)

	// if aux == 0 goto lastEnd:
	condition4 := f.NewBoolExpression().SetLeft(aux).SetRight(f.NewLiteral().SetValue("0")).SetOp(EQ)
	conditional4 := f.NewConditionalJump().SetCondition(condition4).SetTarget(lastEndLabel)
	block = append(block, conditional4)

	// stringFinalAddress = stringFinalAddress + 1
	assign11 := f.NewCompoundAssignment().SetAssignee(stringFinalAddress).SetLeft(stringFinalAddress).SetRight(f.NewLiteral().SetValue("1")).SetOperator(PLUS)
	block = append(block, assign11)

	// aux = heap[stringFinalAddress]
	assign12 := f.NewSimpleAssignment().SetAssignee(aux).SetVal(f.NewHeapIndexed().SetIndex(stringFinalAddress))
	block = append(block, assign12)

	// goto last
	block = append(block, f.NewUnconditionalJump().SetTarget(lastLabel))

	// lastEnd:
	block = append(block, lastEndLabel)

	// decimalSize =  stringFinalAddress - decimalPartStart
	assign13 := f.NewCompoundAssignment().SetAssignee(decimalSize).SetLeft(stringFinalAddress).SetRight(decimalPartStart).SetOperator(MINUS)
	block = append(block, assign13)

	// divisor = 1
	assign14 := f.NewSimpleAssignment().SetAssignee(divisor).SetVal(f.NewLiteral().SetValue("1"))
	block = append(block, assign14)

	// current = 0
	assign15 := f.NewSimpleAssignment().SetAssignee(current).SetVal(f.NewLiteral().SetValue("0"))
	block = append(block, assign15)

	// convert:
	block = append(block, conversionLabel)

	// if(current >= decimalSize) goto endConvert
	condition5 := f.NewBoolExpression().SetLeft(current).SetRight(decimalSize).SetOp(GTE)
	conditional5 := f.NewConditionalJump().SetCondition(condition5).SetTarget(endConvertLabel)
	block = append(block, conditional5)

	// divisor = divisor * 10
	assign16 := f.NewCompoundAssignment().SetAssignee(divisor).SetLeft(divisor).SetRight(f.NewLiteral().SetValue("10")).SetOperator(MULTIPLY)
	block = append(block, assign16)

	// current = current + 1
	assign17 := f.NewCompoundAssignment().SetAssignee(current).SetLeft(current).SetRight(f.NewLiteral().SetValue("1")).SetOperator(PLUS)
	block = append(block, assign17)

	// goto convert
	block = append(block, f.NewUnconditionalJump().SetTarget(conversionLabel))

	// endConvert:
	block = append(block, endConvertLabel)

	// right = right / divisor
	assign18 := f.NewCompoundAssignment().SetAssignee(right).SetLeft(right).SetRight(divisor).SetOperator(DIVIDE)
	block = append(block, assign18)

	// result = left + right
	assign19 := f.NewCompoundAssignment().SetAssignee(result).SetLeft(left).SetRight(right).SetOperator(PLUS)
	block = append(block, assign19)

	// stringHeapAddress = result
	assign20 := f.NewSimpleAssignment().SetAssignee(stringHeapAddress).SetVal(result)
	block = append(block, assign20)

	// goto end
	block = append(block, f.NewUnconditionalJump().SetTarget(endLabel))

	// error:
	block = append(block, errorLabel)

	nilValue := f.Utility.NilValue()
	assign21 := f.NewSimpleAssignment().SetAssignee(stringHeapAddress).SetVal(nilValue)
	block = append(block, assign21)

	// end:
	block = append(block, endLabel)

	return &MethodDcl{
		Name:  "__string_to_float",
		Block: block,
	}
}

func (f *TACFactory) IntToStringBuiltIn() *MethodDcl {
	params := f.GetBuiltinParams("__int_to_string")
	block := make(TACBlock, 0)
	param := params[0]

	/*
		stringAddress = hp

		// check if the number is negative
		if(param >= 0) goto positive
		param = param * -1
		heap[hp] = 45 // "-"
		hp = hp + 1
		positive:

		if param != 0 goto next
		// if the number is 0, just print "0"
		heap[hp] = 48 // "0"
		hp = hp + 1
		heap[hp] = 0
		hp = hp + 1
		goto end

		next:
		// get the last address of the string
		stringFinalAddress = hp - 1
		aux = param
		last:
		if aux == 0 goto lastEnd:
		aux = aux / 10
		aux = (int) aux
		stringFinalAddress = stringFinalAddress + 1
		goto last
		lastEnd:

		hp = stringFinalAddress + 1

		// convert the number to string
		convert:
		if(param == 0) goto endConvert
		aux = param % 10
		aux = aux + 48
		heap[stringFinalAddress] = aux
		stringFinalAddress = stringFinalAddress - 1
		param = param / 10
		param = (int) param
		goto convert
		endConvert:

		heap[hp] = 0
		hp = hp + 1

		end:
		param = stringAddress


	*/

	stringAddress := f.NewTemp()
	stringFinalAddress := f.NewTemp()
	aux := f.NewTemp()

	positiveLabel := f.NewLabel()
	nextLabel := f.NewLabel()
	lastLabel := f.NewLabel()
	lastEndLabel := f.NewLabel()
	convertLabel := f.NewLabel()
	endConvertLabel := f.NewLabel()
	endLabel := f.NewLabel()

	// stringAddress = hp
	assign1 := f.NewSimpleAssignment().SetAssignee(stringAddress).SetVal(f.NewHeapPtr())
	block = append(block, assign1)

	// if(param >= 0) goto positive
	condition := f.NewBoolExpression().SetLeft(param).SetRight(f.NewLiteral().SetValue("0")).SetOp(GTE)
	conditional := f.NewConditionalJump().SetCondition(condition).SetTarget(positiveLabel)
	block = append(block, conditional)

	// param = param * -1
	assign2 := f.NewCompoundAssignment().SetAssignee(param).SetLeft(param).SetRight(f.NewLiteral().SetValue("-1")).SetOperator(MULTIPLY)
	block = append(block, assign2)

	// heap[hp] = 45 // "-"
	assign3 := f.NewSimpleAssignment().SetAssignee(f.NewHeapIndexed().SetIndex(f.NewHeapPtr())).SetVal(f.NewLiteral().SetValue("45"))
	block = append(block, assign3)

	// hp = hp + 1
	assign4 := f.NewCompoundAssignment().SetAssignee(f.NewHeapPtr()).SetLeft(f.NewHeapPtr()).SetRight(f.NewLiteral().SetValue("1")).SetOperator(PLUS)
	block = append(block, assign4)

	// positive:
	block = append(block, positiveLabel)

	// if param != 0 goto next
	condition2 := f.NewBoolExpression().SetLeft(param).SetRight(f.NewLiteral().SetValue("0")).SetOp(NEQ)
	conditional2 := f.NewConditionalJump().SetCondition(condition2).SetTarget(nextLabel)
	block = append(block, conditional2)

	// heap[hp] = 48 // "0"
	assign5 := f.NewSimpleAssignment().SetAssignee(f.NewHeapIndexed().SetIndex(f.NewHeapPtr())).SetVal(f.NewLiteral().SetValue("48"))
	block = append(block, assign5)

	// hp = hp + 1
	assign6 := f.NewCompoundAssignment().SetAssignee(f.NewHeapPtr()).SetLeft(f.NewHeapPtr()).SetRight(f.NewLiteral().SetValue("1")).SetOperator(PLUS)
	block = append(block, assign6)

	// heap[hp] = 0
	assign7 := f.NewSimpleAssignment().SetAssignee(f.NewHeapIndexed().SetIndex(f.NewHeapPtr())).SetVal(f.NewLiteral().SetValue("0"))
	block = append(block, assign7)

	// hp = hp + 1
	assign8 := f.NewCompoundAssignment().SetAssignee(f.NewHeapPtr()).SetLeft(f.NewHeapPtr()).SetRight(f.NewLiteral().SetValue("1")).SetOperator(PLUS)
	block = append(block, assign8)

	// goto end
	block = append(block, f.NewUnconditionalJump().SetTarget(endLabel))

	// next:
	block = append(block, nextLabel)

	// stringFinalAddress = hp - 1
	assign9 := f.NewCompoundAssignment().SetAssignee(stringFinalAddress).SetLeft(f.NewHeapPtr()).SetRight(f.NewLiteral().SetValue("1")).SetOperator(MINUS)
	block = append(block, assign9)

	// aux = param
	assign10 := f.NewSimpleAssignment().SetAssignee(aux).SetVal(param)
	block = append(block, assign10)

	// last:
	block = append(block, lastLabel)

	// if aux == 0 goto lastEnd:
	condition3 := f.NewBoolExpression().SetLeft(aux).SetRight(f.NewLiteral().SetValue("0")).SetOp(EQ)
	conditional3 := f.NewConditionalJump().SetCondition(condition3).SetTarget(lastEndLabel)
	block = append(block, conditional3)

	// aux = aux / 10
	assign11 := f.NewCompoundAssignment().SetAssignee(aux).SetLeft(aux).SetRight(f.NewLiteral().SetValue("10")).SetOperator(DIVIDE)
	block = append(block, assign11)

	// aux = (int) aux
	assign12 := f.NewSimpleAssignment().SetAssignee(aux).SetVal(aux).SetCast("int")
	block = append(block, assign12)

	// stringFinalAddress = stringFinalAddress + 1
	assign13 := f.NewCompoundAssignment().SetAssignee(stringFinalAddress).SetLeft(stringFinalAddress).SetRight(f.NewLiteral().SetValue("1")).SetOperator(PLUS)
	block = append(block, assign13)

	// goto last
	block = append(block, f.NewUnconditionalJump().SetTarget(lastLabel))

	// lastEnd:
	block = append(block, lastEndLabel)

	// hp = stringFinalAddress + 1
	assign14 := f.NewCompoundAssignment().SetAssignee(f.NewHeapPtr()).SetLeft(stringFinalAddress).SetRight(f.NewLiteral().SetValue("1")).SetOperator(PLUS)
	block = append(block, assign14)

	// convert:
	block = append(block, convertLabel)

	// if(param == 0) goto endConvert
	condition4 := f.NewBoolExpression().SetLeft(param).SetRight(f.NewLiteral().SetValue("0")).SetOp(EQ)
	conditional4 := f.NewConditionalJump().SetCondition(condition4).SetTarget(endConvertLabel)
	block = append(block, conditional4)

	// aux = param % 10
	assign15 := f.NewCompoundAssignment().SetAssignee(aux).SetLeft(param).SetRight(f.NewLiteral().SetValue("10")).SetOperator(MOD).SetLeftCast("int")
	block = append(block, assign15)

	// aux = aux + 48
	assign16 := f.NewCompoundAssignment().SetAssignee(aux).SetLeft(aux).SetRight(f.NewLiteral().SetValue("48")).SetOperator(PLUS)
	block = append(block, assign16)

	// heap[stringFinalAddress] = aux
	assign17 := f.NewSimpleAssignment().SetAssignee(f.NewHeapIndexed().SetIndex(stringFinalAddress)).SetVal(aux)
	block = append(block, assign17)

	// stringFinalAddress = stringFinalAddress - 1
	assign18 := f.NewCompoundAssignment().SetAssignee(stringFinalAddress).SetLeft(stringFinalAddress).SetRight(f.NewLiteral().SetValue("1")).SetOperator(MINUS)
	block = append(block, assign18)

	// param = param / 10
	assign19 := f.NewCompoundAssignment().SetAssignee(param).SetLeft(param).SetRight(f.NewLiteral().SetValue("10")).SetOperator(DIVIDE)
	block = append(block, assign19)

	// param = (int) param
	assign20 := f.NewSimpleAssignment().SetAssignee(param).SetVal(param).SetCast("int")
	block = append(block, assign20)

	// goto convert
	block = append(block, f.NewUnconditionalJump().SetTarget(convertLabel))

	// endConvert:
	block = append(block, endConvertLabel)

	// heap[hp] = 0
	assign21 := f.NewSimpleAssignment().SetAssignee(f.NewHeapIndexed().SetIndex(f.NewHeapPtr())).SetVal(f.NewLiteral().SetValue("0"))
	block = append(block, assign21)

	// hp = hp + 1
	assign22 := f.NewCompoundAssignment().SetAssignee(f.NewHeapPtr()).SetLeft(f.NewHeapPtr()).SetRight(f.NewLiteral().SetValue("1")).SetOperator(PLUS)
	block = append(block, assign22)

	// end:
	block = append(block, endLabel)

	// param = stringAddress
	assign23 := f.NewSimpleAssignment().SetAssignee(param).SetVal(stringAddress)
	block = append(block, assign23)

	return &MethodDcl{
		Name:  "__int_to_string",
		Block: block,
	}
}

func (f *TACFactory) FloatToStringBuiltIn() *MethodDcl {
	params := f.GetBuiltinParams("__float_to_string")
	intToStringParams := f.GetBuiltinParams("__int_to_string")
	stringParam := intToStringParams[0]
	block := make(TACBlock, 0)
	param := params[0]

	/*
		// get the integer part of the float
		integerPart = (int) param

		// get the decimal part of the float
		decimalPart = param - integerPart

		// convert the integer part to string
		stringParam = integerPart
		__int_to_string()
		integerPart = stringParam

		// convert the decimal part to string
		stringParam = decimalPart
		__int_to_string()
		decimalPart = stringParam

		// replace terminator char of the first string with "."
		param = decimalPart
		param = param - 1
		heap[param] = 46 // "."

		// return the final string
		param = integerPart
	*/

	integerPart := f.NewTemp()
	decimalPart := f.NewTemp()

	// integerPart = (int) param
	assign1 := f.NewSimpleAssignment().SetAssignee(integerPart).SetVal(param).SetCast("int")
	block = append(block, assign1)

	// decimalPart = param - integerPart
	assign2 := f.NewCompoundAssignment().SetAssignee(decimalPart).SetLeft(param).SetRight(integerPart).SetOperator(MINUS)
	block = append(block, assign2)

	// decimalPart = decimalPart * 1000
	assign21 := f.NewCompoundAssignment().SetAssignee(decimalPart).SetLeft(decimalPart).SetRight(f.NewLiteral().SetValue("10000")).SetOperator(MULTIPLY)
	block = append(block, assign21)

	// decimalPart = (int) decimalPart
	assign22 := f.NewSimpleAssignment().SetAssignee(decimalPart).SetVal(decimalPart).SetCast("int")
	block = append(block, assign22)

	// stringParam = integerPart
	assign3 := f.NewSimpleAssignment().SetAssignee(stringParam).SetVal(integerPart)
	block = append(block, assign3)

	// __int_to_string()
	call1 := f.NewMethodCall("__int_to_string")
	block = append(block, call1)

	// integerPart = stringParam
	assign4 := f.NewSimpleAssignment().SetAssignee(integerPart).SetVal(stringParam)
	block = append(block, assign4)

	// stringParam = decimalPart
	assign5 := f.NewSimpleAssignment().SetAssignee(stringParam).SetVal(decimalPart)
	block = append(block, assign5)

	// __int_to_string()
	call2 := f.NewMethodCall("__int_to_string")
	block = append(block, call2)

	// decimalPart = stringParam
	assign6 := f.NewSimpleAssignment().SetAssignee(decimalPart).SetVal(stringParam)
	block = append(block, assign6)

	// param = decimalPart
	assign7 := f.NewSimpleAssignment().SetAssignee(param).SetVal(decimalPart)
	block = append(block, assign7)

	// param = param - 1
	assign8 := f.NewCompoundAssignment().SetAssignee(param).SetLeft(param).SetRight(f.NewLiteral().SetValue("1")).SetOperator(MINUS)
	block = append(block, assign8)

	// heap[param] = 46 // "."
	assign9 := f.NewSimpleAssignment().SetAssignee(f.NewHeapIndexed().SetIndex(param)).SetVal(f.NewLiteral().SetValue("46"))
	block = append(block, assign9)

	// param = integerPart
	assign10 := f.NewSimpleAssignment().SetAssignee(param).SetVal(integerPart)
	block = append(block, assign10)

	return &MethodDcl{
		Name:  "__float_to_string",
		Block: block,
	}
}

func (f *TACFactory) BoolToStringBuiltIn() *MethodDcl {
	params := f.GetBuiltinParams("__bool_to_string")
	block := make(TACBlock, 0)
	param := params[0]

	address := f.NewTemp()
	falseLabel := f.NewLabel()
	endLabel := f.NewLabel()

	// address = hp
	assign1 := f.NewSimpleAssignment().SetAssignee(address).SetVal(f.NewHeapPtr())
	block = append(block, assign1)

	// if(param == 0) goto false
	condition := f.NewBoolExpression().SetLeft(param).SetRight(f.NewLiteral().SetValue("0")).SetOp(EQ)
	conditional := f.NewConditionalJump().SetCondition(condition).SetTarget(falseLabel)
	block = append(block, conditional)

	var trueStr string = "true"
	for i := 0; i < 4; i++ {
		code := int(trueStr[i])
		assign := f.NewSimpleAssignment().SetAssignee(f.NewHeapIndexed().SetIndex(f.NewHeapPtr())).SetVal(f.NewLiteral().SetValue(strconv.Itoa(code)))
		block = append(block, assign)
		assign2 := f.NewCompoundAssignment().SetAssignee(f.NewHeapPtr()).SetLeft(f.NewHeapPtr()).SetRight(f.NewLiteral().SetValue("1")).SetOperator(PLUS)
		block = append(block, assign2)
	}

	// goto end
	block = append(block, f.NewUnconditionalJump().SetTarget(endLabel))

	// false:
	block = append(block, falseLabel)

	var falseStr string = "false"
	for i := 0; i < 5; i++ {
		code := int(falseStr[i])
		assign := f.NewSimpleAssignment().SetAssignee(f.NewHeapIndexed().SetIndex(f.NewHeapPtr())).SetVal(f.NewLiteral().SetValue(strconv.Itoa(code)))
		block = append(block, assign)
		assign2 := f.NewCompoundAssignment().SetAssignee(f.NewHeapPtr()).SetLeft(f.NewHeapPtr()).SetRight(f.NewLiteral().SetValue("1")).SetOperator(PLUS)
		block = append(block, assign2)
	}

	// end:
	block = append(block, endLabel)

	// heap[hp] = 0
	assign3 := f.NewSimpleAssignment().SetAssignee(f.NewHeapIndexed().SetIndex(f.NewHeapPtr())).SetVal(f.NewLiteral().SetValue("0"))
	block = append(block, assign3)

	// hp = hp + 1
	assign4 := f.NewCompoundAssignment().SetAssignee(f.NewHeapPtr()).SetLeft(f.NewHeapPtr()).SetRight(f.NewLiteral().SetValue("1")).SetOperator(PLUS)
	block = append(block, assign4)

	// param = address
	assign5 := f.NewSimpleAssignment().SetAssignee(param).SetVal(address)
	block = append(block, assign5)

	return &MethodDcl{
		Name:  "__bool_to_string",
		Block: block,
	}
}
