package tag

import (
	"fmt"
	"strings"
)

// relation
type Many2Many struct {
	Pkg,
	Name,
	Field string `json:",omitempty"`
}

// Zero Value
func (m2m *Many2Many) Empty() bool {
	return m2m.Pkg == "" || m2m.Name == "" || m2m.Field == ""
}

// Marshal implements Tag
func (m2m *Many2Many) Marshal() (token string) {
	if m2m.Empty() {
		return
	}
	return fmt.Sprintf("%s%s.%s.%s", Prefix(m2m), m2m.Pkg, m2m.Name, m2m.Field)
}

// Match implements Tag
func (m2m *Many2Many) Match(token string) bool {
	return strings.Contains(token, Prefix(m2m))
}

// Unmarshal implements Tag
func (m2m *Many2Many) Unmarshal(token string) {
	if propVal := ParseString(m2m, token); propVal != "" {
		parts := strings.Split(propVal, ".")
		if len(parts) == 3 {
			m2m.Pkg = parts[0]
			m2m.Name = parts[1]
			m2m.Field = parts[2]
		}
	}
}

var _ Tag = (*Many2Many)(nil)
