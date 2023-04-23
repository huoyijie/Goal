package tag

type MultiSelect struct {
	Base
	Many2Many *Many2Many `json:",omitempty"`
}

// Head implements Component
func (ms *MultiSelect) Head() string {
	return ComponentHead(ms)
}

// Is implements Component
func (ms *MultiSelect) Is(token string) bool {
	return IsComponent(ms, token)
}

// Marshal implements Tag
func (ms *MultiSelect) Marshal() string {
	return Marshal(ms)
}

// Match implements Tag
func (ms *MultiSelect) Match(token string) bool {
	return ms.Is(token)
}

// Unmarshal implements Tag
func (ms *MultiSelect) Unmarshal(token string) {
	Unmarshal(token, ms)
}

var _ Tag = (*MultiSelect)(nil)
var _ Component = (*MultiSelect)(nil)
