package tag

// bool
type Switch struct {
	Base
}

// Head implements Component
func (s *Switch) Head() string {
	return ComponentHead(s)
}

// Is implements Component
func (s *Switch) Is(token string) bool {
	return IsComponent(s, token)
}

// Marshal implements Tag
func (s *Switch) Marshal() string {
	return Marshal(s)
}

// Match implements Tag
func (s *Switch) Match(token string) bool {
	return s.Is(token)
}

// Unmarshal implements Tag
func (s *Switch) Unmarshal(token string) {
	Unmarshal(token, s)
}

var _ Tag = (*Switch)(nil)
var _ Component = (*Switch)(nil)
