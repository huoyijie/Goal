package tag

// string
type Uuid struct {
	Base
}

// Head implements Component
func (u *Uuid) Head() string {
	return ComponentHead(u)
}

// Is implements Component
func (u *Uuid) Is(token string) bool {
	return IsComponent(u, token)
}

// Marshal implements Tag
func (u *Uuid) Marshal() string {
	return Marshal(u)
}

// Match implements Tag
func (u *Uuid) Match(token string) bool {
	return u.Is(token)
}

// Unmarshal implements Tag
func (u *Uuid) Unmarshal(token string) {
	Unmarshal(token, u)
}

var _ Tag = (*Uuid)(nil)
var _ Component = (*Uuid)(nil)
