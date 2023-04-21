package tag

// Object
type Dropdown struct {
	Base
	Strings,
	Ints,
	Uints,
	Floats,
	DynamicStrings,
	DynamicInts,
	DynamicUints,
	DynamicFloats bool `json:",omitempty"`
	BelongTo *BelongTo `json:",omitempty"`
	HasOne   *HasOne   `json:",omitempty"`
}

// Head implements Component
func (d *Dropdown) Head() string {
	return ComponentHead(d)
}

// Is implements Component
func (d *Dropdown) Is(token string) bool {
	return IsComponent(d, token)
}

// Marshal implements Tag
func (d *Dropdown) Marshal() string {
	return Marshal(d)
}

// Match implements Tag
func (d *Dropdown) Match(token string) bool {
	return d.Is(token)
}

// Unmarshal implements Tag
func (d *Dropdown) Unmarshal(token string) {
	Unmarshal(token, d)
}

var _ Tag = (*Dropdown)(nil)
var _ Component = (*Dropdown)(nil)
