package tag

// time.Time
type Calendar struct {
	Base
	ShowTime,
	ShowIcon bool `json:",omitempty"`
}

// Head implements Component
func (c *Calendar) Head() string {
	return ComponentHead(c)
}

// Is implements Component
func (c *Calendar) Is(token string) bool {
	return IsComponent(c, token)
}

// Marshal implements Tag
func (c *Calendar) Marshal() string {
	return Marshal(c)
}

// Match implements Tag
func (c *Calendar) Match(token string) bool {
	return c.Is(token)
}

// Unmarshal implements Tag
func (c *Calendar) Unmarshal(token string) {
	Unmarshal(token, c)
}

var _ Tag = (*Calendar)(nil)
var _ Component = (*Calendar)(nil)
