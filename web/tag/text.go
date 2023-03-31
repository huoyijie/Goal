package tag

// string
type Text struct {
	Base
}

// Head implements Component
func (t *Text) Head() string {
	return ComponentHead(t)
}

// Is implements Component
func (t *Text) Is(token string) bool {
	return IsComponent(t, token)
}

// Marshal implements Tag
func (t *Text) Marshal() string {
	return Marshal(t)
}

// Match implements Tag
func (t *Text) Match(token string) bool {
	return t.Is(token)
}

// Unmarshal implements Tag
func (t *Text) Unmarshal(token string) {
	Unmarshal(token, t)
}

var _ Tag = (*Text)(nil)
var _ Component = (*Text)(nil)
