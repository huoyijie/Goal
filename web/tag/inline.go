package tag

type Inline struct {
	Base
	HasOne  *HasOne  `json:",omitempty"`
	HasMany *HasMany `json:",omitempty"`
}

// Head implements Component
func (i *Inline) Head() string {
	return ComponentHead(i)
}

// Is implements Component
func (i *Inline) Is(token string) bool {
	return IsComponent(i, token)
}

// Marshal implements Tag
func (i *Inline) Marshal() string {
	return Marshal(i)
}

// Match implements Tag
func (i *Inline) Match(token string) bool {
	return i.Is(token)
}

// Unmarshal implements Tag
func (i *Inline) Unmarshal(token string) {
	Unmarshal(token, i)
}

var _ Tag = (*Inline)(nil)
var _ Component = (*Inline)(nil)
