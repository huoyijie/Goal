package tag

// define a iface that return itself
type IBase interface {
	Get() *Base
}

// base
type Base struct {
	Primary,
	Unique,
	Postonly,
	Readonly,
	Sortable,
	Asc,
	Desc,
	GlobalSearch,
	Filter,
	Hidden,
	Secret,
	Autowired bool `json:",omitempty"`
}

// Get implements IBase
func (b *Base) Get() *Base {
	return b
}

// Marshal implements Tag
func (b *Base) Marshal() string {
	return Marshal(b)
}

// Match implements Tag
func (b *Base) Match(token string) bool {
	return true
}

// Unmarshal implements Tag
func (b *Base) Unmarshal(token string) {
	Unmarshal(token, b)
}

var _ Tag = (*Base)(nil)
var _ IBase = (*Base)(nil)
