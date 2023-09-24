package repl

import "github.com/damianpeaf/OLC2_Proyecto2_202110568/value"

const (
	BreakItem    = "break"
	ContinueItem = "continue"
	ReturnItem   = "return"
)

type CallStackItem struct {
	ReturnValue value.IVOR
	Type        []string
	Action      string
}

func (csi *CallStackItem) IsType(t string) bool {

	for _, i := range csi.Type {
		if i == t {
			return true
		}
	}

	return false
}

func (csi *CallStackItem) IsAction(a string) bool {
	return csi.Action == a
}

func (csi *CallStackItem) ResetAction() {
	csi.Action = ""
}

type CallStack struct {
	Items []*CallStackItem
}

func (cs *CallStack) Push(item *CallStackItem) {
	cs.Items = append(cs.Items, item)
}

func (cs *CallStack) Pop() *CallStackItem {
	item := cs.Items[len(cs.Items)-1]
	cs.Items = cs.Items[:len(cs.Items)-1]
	return item
}

func (cs *CallStack) Peek() *CallStackItem {
	return cs.Items[len(cs.Items)-1]
}

func (cs *CallStack) In(item *CallStackItem) bool {
	for _, i := range cs.Items {
		if i == item {
			return true
		}
	}
	return false
}

// Remove items from the stack until the item is found
func (cs *CallStack) Clean(item *CallStackItem) {

	if !cs.In(item) {
		return
	}

	for {
		peek := cs.Pop()

		if peek == item {
			break
		}
	}

}

func (cs *CallStack) IsContinueEnv() (bool, *CallStackItem) {

	// continue can be only in a loop
	// so it can be in break env like switch
	// but cannot interfer a function call that is a return env
	start := len(cs.Items) - 1

	for i := start; i >= 0; i-- {
		if cs.Items[i].IsType(ContinueItem) {
			return true, cs.Items[i]
		}

		if cs.Items[i].IsType(ReturnItem) {
			return false, nil
		}
	}

	return false, nil

}

func (cs *CallStack) IsBreakEnv() (bool, *CallStackItem) {
	// break item must be the peek of the stack, cannot interrupt a function call or a loop

	if len(cs.Items) == 0 {
		return false, nil
	}

	if cs.Items[len(cs.Items)-1].IsType(BreakItem) {
		return true, cs.Items[len(cs.Items)-1]
	}

	return false, nil
}

func (cs *CallStack) IsReturnEnv() (bool, *CallStackItem) {

	// return item can interfer with any other item

	for i := len(cs.Items) - 1; i >= 0; i-- {
		if cs.Items[i].IsType(ReturnItem) {
			return true, cs.Items[i]
		}
	}

	return false, nil
}

func (cs *CallStack) Len() int {
	return len(cs.Items)
}

func NewCallStack() *CallStack {
	return &CallStack{Items: []*CallStackItem{}}
}
