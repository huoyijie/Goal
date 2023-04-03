package tag

import "strings"

type Autowired bool

// Marshal implements Tag
func (a *Autowired) Marshal() (token string) {
	if *a {
		token = Key(a)
	}
	return
}

// Match implements Tag
func (a *Autowired) Match(token string) bool {
	return strings.Contains(token, Key(a))
}

// Unmarshal implements Tag
func (a *Autowired) Unmarshal(token string) {
	if a.Match(token) {
		*a = true
	}
}

type Secret bool

// Marshal implements Tag
func (s *Secret) Marshal() (token string) {
	if *s {
		token = Key(s)
	}
	return
}

// Match implements Tag
func (s *Secret) Match(token string) bool {
	return strings.Contains(token, Key(s))
}

// Unmarshal implements Tag
func (s *Secret) Unmarshal(token string) {
	if s.Match(token) {
		*s = true
	}
}

type Hidden bool

// Marshal implements Tag
func (h *Hidden) Marshal() (token string) {
	if *h {
		token = Key(h)
	}
	return
}

// Match implements Tag
func (h *Hidden) Match(token string) bool {
	return strings.Contains(token, Key(h))
}

// Unmarshal implements Tag
func (h *Hidden) Unmarshal(token string) {
	if h.Match(token) {
		*h = true
	}
}

type Postonly bool

// Marshal implements Tag
func (p *Postonly) Marshal() (token string) {
	if *p {
		token = Key(p)
	}
	return
}

// Match implements Tag
func (p *Postonly) Match(token string) bool {
	return strings.Contains(token, Key(p))
}

// Unmarshal implements Tag
func (p *Postonly) Unmarshal(token string) {
	if p.Match(token) {
		*p = true
	}
}

type Readonly bool

// Marshal implements Tag
func (r *Readonly) Marshal() (token string) {
	if *r {
		token = Key(r)
	}
	return
}

// Match implements Tag
func (r *Readonly) Match(token string) bool {
	return strings.Contains(token, Key(r))
}

// Unmarshal implements Tag
func (r *Readonly) Unmarshal(token string) {
	if r.Match(token) {
		*r = true
	}
}

type Primary bool

// Marshal implements Tag
func (p *Primary) Marshal() (token string) {
	if *p {
		token = Key(p)
	}
	return
}

// Match implements Tag
func (p *Primary) Match(token string) bool {
	return strings.Contains(token, Key(p))
}

// Unmarshal implements Tag
func (p *Primary) Unmarshal(token string) {
	if p.Match(token) {
		*p = true
	}
}

type Unique bool

// Marshal implements Tag
func (u *Unique) Marshal() (token string) {
	if *u {
		token = Key(u)
	}
	return
}

// Match implements Tag
func (u *Unique) Match(token string) bool {
	return strings.Contains(token, Key(u))
}

// Unmarshal implements Tag
func (u *Unique) Unmarshal(token string) {
	if u.Match(token) {
		*u = true
	}
}

type ShowTime bool

// Marshal implements Tag
func (s *ShowTime) Marshal() (token string) {
	if *s {
		token = Key(s)
	}
	return
}

// Match implements Tag
func (s *ShowTime) Match(token string) bool {
	return strings.Contains(token, Key(s))
}

// Unmarshal implements Tag
func (s *ShowTime) Unmarshal(token string) {
	if s.Match(token) {
		*s = true
	}
}

type ShowIcon bool

// Marshal implements Tag
func (s *ShowIcon) Marshal() (token string) {
	if *s {
		token = Key(s)
	}
	return
}

// Match implements Tag
func (s *ShowIcon) Match(token string) bool {
	return strings.Contains(token, Key(s))
}

// Unmarshal implements Tag
func (s *ShowIcon) Unmarshal(token string) {
	if s.Match(token) {
		*s = true
	}
}

type Filter bool

// Marshal implements Tag
func (f *Filter) Marshal() (token string) {
	if *f {
		token = Key(f)
	}
	return
}

// Match implements Tag
func (f *Filter) Match(token string) bool {
	return strings.Contains(token, Key(f))
}

// Unmarshal implements Tag
func (f *Filter) Unmarshal(token string) {
	if f.Match(token) {
		*f = true
	}
}

type ShowButtons bool

// Marshal implements Tag
func (s *ShowButtons) Marshal() (token string) {
	if *s {
		token = Key(s)
	}
	return
}

// Match implements Tag
func (s *ShowButtons) Match(token string) bool {
	return strings.Contains(token, Key(s))
}

// Unmarshal implements Tag
func (s *ShowButtons) Unmarshal(token string) {
	if s.Match(token) {
		*s = true
	}
}

type Float bool

// Marshal implements Tag
func (f *Float) Marshal() (token string) {
	if *f {
		token = Key(f)
	}
	return
}

// Match implements Tag
func (f *Float) Match(token string) bool {
	return strings.Contains(token, Key(f))
}

// Unmarshal implements Tag
func (f *Float) Unmarshal(token string) {
	if f.Match(token) {
		*f = true
	}
}

type Sortable bool

// Marshal implements Tag
func (s *Sortable) Marshal() (token string) {
	if *s {
		token = Key(s)
	}
	return
}

// Match implements Tag
func (s *Sortable) Match(token string) bool {
	return strings.Contains(token, Key(s))
}

// Unmarshal implements Tag
func (s *Sortable) Unmarshal(token string) {
	if s.Match(token) {
		*s = true
	}
}

type Asc bool

// Marshal implements Tag
func (a *Asc) Marshal() (token string) {
	if *a {
		token = Key(a)
	}
	return
}

// Match implements Tag
func (a *Asc) Match(token string) bool {
	return strings.Contains(token, Key(a))
}

// Unmarshal implements Tag
func (a *Asc) Unmarshal(token string) {
	if a.Match(token) {
		*a = true
	}
}

type Desc bool

// Marshal implements Tag
func (d *Desc) Marshal() (token string) {
	if *d {
		token = Key(d)
	}
	return
}

// Match implements Tag
func (d *Desc) Match(token string) bool {
	return strings.Contains(token, Key(d))
}

// Unmarshal implements Tag
func (d *Desc) Unmarshal(token string) {
	if d.Match(token) {
		*d = true
	}
}

var (
	_ Tag = (*Autowired)(nil)
	_ Tag = (*Secret)(nil)
	_ Tag = (*Hidden)(nil)
	_ Tag = (*Postonly)(nil)
	_ Tag = (*Readonly)(nil)
	_ Tag = (*Primary)(nil)
	_ Tag = (*Unique)(nil)
	_ Tag = (*ShowTime)(nil)
	_ Tag = (*ShowIcon)(nil)
	_ Tag = (*Filter)(nil)
	_ Tag = (*ShowButtons)(nil)
	_ Tag = (*Float)(nil)
	_ Tag = (*Sortable)(nil)
	_ Tag = (*Asc)(nil)
	_ Tag = (*Desc)(nil)
)
