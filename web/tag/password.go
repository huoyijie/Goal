package tag

// string
type Password struct {
	Base
}

// Head implements Component
func (p *Password) Head() string {
	return ComponentHead(p)
}

// Is implements Component
func (p *Password) Is(token string) bool {
	return IsComponent(p, token)
}

// Marshal implements Tag
func (p *Password) Marshal() string {
	return Marshal(p)
}

// Match implements Tag
func (p *Password) Match(token string) bool {
	return p.Is(token)
}

// Unmarshal implements Tag
func (p *Password) Unmarshal(token string) {
	Unmarshal(token, p)
}

var _ Tag = (*Password)(nil)
var _ Component = (*Password)(nil)
