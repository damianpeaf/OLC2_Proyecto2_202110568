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

	// labels:
	cmpStr := f.NewLabel()
	endOfStr1 := f.NewLabel()
	s1GreaterThanS2 := f.NewLabel()
	s1LessThanS2 := f.NewLabel()
	nextChar := f.NewLabel()
	equalStr := f.NewLabel()
	endCmpStr := f.NewLabel()

	// cmp_str:
	block = append(block, cmpStr)

	// t1 = (int) heap[s1]
	assign1 := f.NewSimpleAssignment().SetAssignee(f.NewTemp()).SetVal(f.NewHeapIndexed().SetIndex(s1))
	block = append(block, assign1)

	// t2 = (int) heap[s2]
	assign2 := f.NewSimpleAssignment().SetAssignee(f.NewTemp()).SetVal(f.NewHeapIndexed().SetIndex(s2))
	block = append(block, assign2)

	// if(t1 == 0) goto end_of_str_1
	condition1 := f.NewBoolExpression().SetLeft(f.NewTemp()).SetRight(f.NewLiteral().SetValue("0")).SetOp(EQ)
	conditional1 := f.NewConditionalJump().SetCondition(condition1).SetTarget(endOfStr1)
	block = append(block, conditional1)

	// if(t2 == 0) goto s1_greater_than_s2
	condition2 := f.NewBoolExpression().SetLeft(f.NewTemp()).SetRight(f.NewLiteral().SetValue("0")).SetOp(EQ)
	conditional2 := f.NewConditionalJump().SetCondition(condition2).SetTarget(s1GreaterThanS2)
	block = append(block, conditional2)

	// if(t1 < t2) goto s1_less_than_s2
	condition3 := f.NewBoolExpression().SetLeft(f.NewTemp()).SetRight(f.NewTemp()).SetOp(LT)
	conditional3 := f.NewConditionalJump().SetCondition(condition3).SetTarget(s1LessThanS2)
	block = append(block, conditional3)

	// if(t1 > t2) goto s1_greater_than_s2
	condition4 := f.NewBoolExpression().SetLeft(f.NewTemp()).SetRight(f.NewTemp()).SetOp(GT)
	conditional4 := f.NewConditionalJump().SetCondition(condition4).SetTarget(s1GreaterThanS2)
	block = append(block, conditional4)

	// next_char:
	block = append(block, nextChar)

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
	condition5 := f.NewBoolExpression().SetLeft(f.NewTemp()).SetRight(f.NewLiteral().SetValue("0")).SetOp(EQ)
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

	// goto end_cmp_str
	block = append(block, assign6)

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
