package iniconfig

import (
	"testing"
)

var (
	test_ini = `
hoge # comment
foo = bar
[foo] ; comment with ;
bar = baz # comment on key value
foobar
# commnet
`
)

func TestBasic(t *testing.T) {
	c := NewConfigString(test_ini)

	baz, exist := c.Get("foo", "bar")
	if !exist {
		t.Errorf("[foo] bar not exists")
	}

	if baz != "baz" {
		t.Errorf("[foo] bar must be baz")
	}
	baz, exist = c.Get("fooo", "bar")
	if exist {
		t.Errorf("[fooo] bar exists")
	}
	baz, exist = c.Get("foo", "barr")
	if exist {
		t.Errorf("[foo] barr exists")
	}

	c.Put("foo", "barr", "buzz")
	buzz, exist := c.Get("foo", "barr")
	if !exist {
		t.Errorf("[foo] barr not exists")
	}
	if buzz != "buzz" {
		t.Errorf("[foo] barr must be buzz")
	}

	c.Put("fooo", "barr", "buzz")
	buzz, exist = c.Get("fooo", "barr")
	if !exist {
		t.Errorf("[fooo] barr not exists")
	}
	if buzz != "buzz" {
		t.Errorf("[fooo] barr must be buzz")
	}

	err := c.Delete("foooo", "barr")
	if err == nil {
		t.Errorf("Delete to unknown section must be error")
	}

	err = c.Delete("fooo", "barr")
	if err != nil {
		t.Errorf("Delete To known section must not be error %v", err)
	}

	_, exist = c.Get("fooo", "barr")
	if exist {
		t.Errorf("Deleted key still exists")
	}

	c.DeleteSection("fooo")
	_, exist = c.Get("fooo", "barr")
	if exist {
		t.Errorf("Deleted key still exists")
	}

	print(c.String())

}
