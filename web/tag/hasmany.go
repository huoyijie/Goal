package tag

import (
	"fmt"
	"strings"
)

// relation
type HasMany struct {
	Pkg,
	Name string `json:",omitempty"`
}

// Zero Value
func (hm *HasMany) Empty() bool {
	return hm.Pkg == "" || hm.Name == ""
}

// Marshal implements Tag
func (hm *HasMany) Marshal() (token string) {
	if hm.Empty() {
		return
	}
	return fmt.Sprintf("%s%s.%s", Prefix(hm), hm.Pkg, hm.Name)
}

// Match implements Tag
func (hm *HasMany) Match(token string) bool {
	return strings.Contains(token, Prefix(hm))
}

// Unmarshal implements Tag
func (hm *HasMany) Unmarshal(token string) {
	if propVal := ParseString(hm, token); propVal != "" {
		parts := strings.Split(propVal, ".")
		if len(parts) == 2 {
			hm.Pkg = parts[0]
			hm.Name = parts[1]
		}
	}
}

var _ Tag = (*HasMany)(nil)
