package coverage

type Context struct {
	strategy Strategy
}

func NewContext(s Strategy) *Context {
	return &Context{strategy: s}
}

func (c *Context) SetStrategy(s Strategy) {
	c.strategy = s
}

func (c *Context) Check(req Request) (Decision, error) {
	return c.strategy.Check(req)
}
