package tag

import (
	"testing"
)

func TestBase(t *testing.T) {
	b := Base{Autowired: true, Readonly: true}
	if token := b.Marshal(); token != "readonly,autowired" {
		t.Error("Marshal error 1")
	}

	b = Base{}
	if token := b.Marshal(); token != "" {
		t.Error("Marshal error 2")
	}

	b.Unmarshal("unique,postonly,belongTo=auth.User.Username")
	if token := b.Marshal(); token != "unique,postonly" {
		t.Error("Marshal error 3")
	}
}

func TestCanlendar(t *testing.T) {
	c := Calendar{ShowTime: true, ShowIcon: true}
	if token := c.Marshal(); token != "<calendar>showTime,showIcon" {
		t.Error("Marshal error 1")
	}

	c = Calendar{}
	c.Unmarshal("<calendar>showTime,readonly,unique,belongTo=auth.User.Username")
	if token := c.Marshal(); token != "<calendar>unique,readonly,showTime" {
		t.Error("Marshal error 2")
	}
}

func TestDropdown(t *testing.T) {
	d := Dropdown{}
	if d.Marshal() != "<dropdown>" {
		t.Error("Marshal error 1")
	}

	d = Dropdown{}
	d.Unmarshal("<dropdown>autowired,secret,hidden,belongTo=auth.User.Username")
	if token := d.Marshal(); token != "<dropdown>hidden,secret,autowired,belongTo=auth.User.Username" {
		t.Error("Marshal error 2")
	}
}

func TestNumber(t *testing.T) {
	min := 10
	max := 100
	n := Number{ShowButtons: true, Min: &min, Max: &max}
	if n.Marshal() != "<number>showButtons,min=10,max=100" {
		t.Error("Marshal error 1")
	}

	n = Number{}
	n.Unmarshal("<number>showTime,autowired,showButtons,secret,hidden,min=10,max=1000")
	if token := n.Marshal(); token != "<number>hidden,secret,autowired,showButtons,min=10,max=1000" {
		t.Error("Marshal error 2")
	}
}

func TestText(t *testing.T) {
	it := Text{}
	if it.Marshal() != "<text>" {
		t.Error("Marshal error 1")
	}

	it.Unmarshal("<text>autowired,secret,hidden")
	if token := it.Marshal(); token != "<text>hidden,secret,autowired" {
		t.Error("Marshal error 2")
	}
}

func TestPassword(t *testing.T) {
	u := Password{}
	if u.Marshal() != "<password>" {
		t.Error("Marshal error 1")
	}

	u.Unmarshal("<password>autowired,secret,hidden")
	if token := u.Marshal(); token != "<password>hidden,secret,autowired" {
		t.Error("Marshal error 2")
	}
}

func TestUuid(t *testing.T) {
	u := Uuid{}
	if u.Marshal() != "<uuid>" {
		t.Error("Marshal error 1")
	}

	u.Unmarshal("<uuid>autowired,secret,hidden")
	if token := u.Marshal(); token != "<uuid>hidden,secret,autowired" {
		t.Error("Marshal error 2")
	}
}

func TestSwitch(t *testing.T) {
	s := Switch{}
	if s.Marshal() != "<switch>" {
		t.Error("Marshal error 1")
	}

	s.Unmarshal("<switch>autowired,secret,hidden")
	if token := s.Marshal(); token != "<switch>hidden,secret,autowired" {
		t.Error("Marshal error 2")
	}
}
