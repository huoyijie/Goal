package tag

import (
	"fmt"
	"strings"
)

// string
type UploadTo struct {
	Path string `json:",omitempty"`
}

// Marshal implements Tag
func (ut *UploadTo) Marshal() (token string) {
	if ut.Path == "" {
		return
	}
	return fmt.Sprintf("%s%s", Prefix(ut), ut.Path)
}

// Match implements Tag
func (ut *UploadTo) Match(token string) bool {
	return strings.Contains(token, Prefix(ut))
}

// Unmarshal implements Tag
func (ut *UploadTo) Unmarshal(token string) {
	ut.Path = ParseString(ut, token)
}

var _ Tag = (*UploadTo)(nil)
