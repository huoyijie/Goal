package tag

// int | float
type Number struct {
	Base
	ShowButtons,
	Float,
	Uint bool `json:",omitempty"`
	Min,
	Max *int `json:",omitempty"`
}

// Head implements Component
func (n *Number) Head() string {
	return ComponentHead(n)
}

// Is implements Component
func (n *Number) Is(token string) bool {
	return IsComponent(n, token)
}

// Marshal implements Tag
func (n *Number) Marshal() string {
	return Marshal(n)
}

// Match implements Tag
func (n *Number) Match(token string) bool {
	return n.Is(token)
}

// Unmarshal implements Tag
func (n *Number) Unmarshal(token string) {
	Unmarshal(token, n)
}

var _ Tag = (*Number)(nil)
var _ Component = (*Number)(nil)
