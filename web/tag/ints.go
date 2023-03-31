package tag

import (
	"fmt"
	"strings"
)

type Min int

// Marshal implements Tag
func (m *Min) Marshal() string {
	return fmt.Sprintf("%s%d", Prefix(m), *m)
}

// Match implements Tag
func (m *Min) Match(token string) bool {
	return strings.Contains(token, Prefix(m))
}

// Unmarshal implements Tag
func (m *Min) Unmarshal(token string) {
	*m = Min(ParseInt(m, token))
}

type Max int

// Marshal implements Tag
func (m *Max) Marshal() string {
	return fmt.Sprintf("%s%d", Prefix(m), *m)
}

// Match implements Tag
func (m *Max) Match(token string) bool {
	return strings.Contains(token, Prefix(m))
}

// Unmarshal implements Tag
func (m *Max) Unmarshal(token string) {
	*m = Max(ParseInt(m, token))
}

var _ Tag = (*Min)(nil)
var _ Tag = (*Max)(nil)
