package tag

import (
	"fmt"
	"strings"
)

// relation
type HasOne struct {
	Pkg,
	Name,
	Field string `json:",omitempty"`
}

// Zero Value
func (ho *HasOne) Empty() bool {
	return ho.Pkg == "" || ho.Name == "" || ho.Field == ""
}

// Marshal implements Tag
func (ho *HasOne) Marshal() (token string) {
	if ho.Empty() {
		return
	}
	return fmt.Sprintf("%s%s.%s.%s", Prefix(ho), ho.Pkg, ho.Name, ho.Field)
}

// Match implements Tag
func (ho *HasOne) Match(token string) bool {
	return strings.Contains(token, Prefix(ho))
}

// Unmarshal implements Tag
func (ho *HasOne) Unmarshal(token string) {
	if propVal := ParseString(ho, token); propVal != "" {
		parts := strings.Split(propVal, ".")
		if len(parts) == 3 {
			ho.Pkg = parts[0]
			ho.Name = parts[1]
			ho.Field = parts[2]
		}
	}
}

var _ Tag = (*HasOne)(nil)
