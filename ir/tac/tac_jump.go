package tac

// ** ConditionalJump
type ConditionalJump struct {
	Condition BoolExpression
	Target    *Label
}

func (c *ConditionalJump) String() string {
	return "if (" + c.Condition.String() + ") goto " + c.Target.LabelName() + ";"
}

// builder utils
func (c *ConditionalJump) SetCondition(condition BoolExpression) *ConditionalJump {
	c.Condition = condition
	return c
}

func (c *ConditionalJump) SetTarget(target *Label) *ConditionalJump {
	c.Target = target
	return c
}

// ** UnconditionalJump
type UnconditionalJump struct {
	Target *Label
}

func (u *UnconditionalJump) String() string {
	return "goto " + u.Target.LabelName() + ";"
}

// builder utils
func (u *UnconditionalJump) SetTarget(target *Label) *UnconditionalJump {
	u.Target = target
	return u
}
