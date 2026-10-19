package flow

import "hospitaladmission/internal/admission"

type Context struct {
	strategy Strategy
}

func NewContext(s Strategy) *Context {
	return &Context{strategy: s}
}

func (c *Context) SetStrategy(s Strategy) {
	c.strategy = s
}

func (c *Context) Begin(a admission.Admission) (admission.Result, error) {
	return c.strategy.Begin(a)
}
