package iniconfig

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"
)

type kv map[string]string

type Config struct {
	sections map[string]kv
}

func NewConfig(rd io.Reader) Config {
	c := Config{map[string]kv{}}
	r := bufio.NewReader(rd)
	currentSection := ""
	c.sections[currentSection] = kv{}
	for {
		l, err := r.ReadString('\n')
		if err != nil {
			break
		}
		l = strings.TrimRight(l, "\r\n")
		l = strings.Split(l, "#")[0]
		l = strings.Split(l, ";")[0]
		l = strings.TrimSpace(l)
		if len(l) == 0 {
			continue
		}
		last := len(l) - 1
		if l[0] == '[' && l[last] == ']' {
			currentSection = l[1:last]
			c.sections[currentSection] = kv{}
		} else {
			keyval := strings.Split(l, "=")
			if len(keyval) < 2 {
				key := strings.TrimSpace(keyval[0])
				c.sections[currentSection][key] = ""
				continue
			}
			keys := keyval[:len(keyval)-1]
			val := strings.TrimSpace(keyval[len(keyval)-1])
			for _, k := range keys {
				k = strings.TrimSpace(k)
				c.sections[currentSection][k] = val
			}
		}
	}
	return c
}

func NewConfigString(s string) Config {
	return NewConfig(bytes.NewBufferString(s))
}

func (c Config) Get(section, key string) (string, bool) {
	keyvalues, exists := c.sections[section]
	if !exists {
		return "", exists
	}
	value, exists := keyvalues[key]
	if !exists {
		return value, exists
	}
	return value, true
}

func (c Config) Put(section, key, val string) {
	_, exists := c.sections[section]
	if !exists {
		c.sections[section] = make(kv)
	}
	c.sections[section][key] = val
}

func (c Config) Delete(section, key string) error {
	_, exists := c.sections[section]
	if !exists {
		return fmt.Errorf("section [%s] does not exist")
	}
	delete(c.sections[section], key)
	return nil
}

func (c Config) DeleteSection(section string) {
	delete(c.sections, section)
}

func (c Config) String() string {
	ret := ""
	keyvalues := c.sections[""]
	for key, val := range keyvalues {
		if val == "" {
			ret += key + "\n"
		} else {
			ret += key + " = " + val + "\n"
		}
	}
	ret += "\n"
	for section, keyvalues := range c.sections {
		if section == "" {
			continue
		}
		ret += "[" + section + "]\n"
		for key, val := range keyvalues {
			if val == "" {
				ret += key + "\n"
			} else {
				ret += key + " = " + val + "\n"
			}
		}
		ret += "\n"
	}
	return ret
}
