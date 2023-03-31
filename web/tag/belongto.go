package tag

import (
	"fmt"
	"strings"
)

// relation
type BelongTo struct {
	Pkg   string `json:",omitempty"`
	Name  string `json:",omitempty"`
	Field string `json:",omitempty"`
}

// Zero Value
func (bt *BelongTo) Empty() bool {
	return bt.Pkg == "" || bt.Name == "" || bt.Field == ""
}

// Marshal implements Tag
func (bt *BelongTo) Marshal() (token string) {
	if bt.Empty() {
		return
	}
	return fmt.Sprintf("%s%s.%s.%s", Prefix(bt), bt.Pkg, bt.Name, bt.Field)
}

// Match implements Tag
func (bt *BelongTo) Match(token string) bool {
	return strings.Contains(token, Prefix(bt))
}

// Unmarshal implements Tag
func (bt *BelongTo) Unmarshal(token string) {
	if propVal := ParseString(bt, token); propVal != "" {
		parts := strings.Split(propVal, ".")
		if len(parts) == 3 {
			bt.Pkg = parts[0]
			bt.Name = parts[1]
			bt.Field = parts[2]
		}
	}
}

var _ Tag = (*BelongTo)(nil)
