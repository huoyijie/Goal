package tag

// time.Time
type File struct {
	Base
	// used by backend
	UploadTo *UploadTo `json:",omitempty"`
}

// Head implements Component
func (f *File) Head() string {
	return ComponentHead(f)
}

// Is implements Component
func (f *File) Is(token string) bool {
	return IsComponent(f, token)
}

// Marshal implements Tag
func (f *File) Marshal() string {
	return Marshal(f)
}

// Match implements Tag
func (f *File) Match(token string) bool {
	return f.Is(token)
}

// Unmarshal implements Tag
func (f *File) Unmarshal(token string) {
	Unmarshal(token, f)
}

var _ Tag = (*File)(nil)
var _ Component = (*File)(nil)
