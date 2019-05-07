package server

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"reflect"
)

type markdown struct {
	content  *bytes.Buffer
	resource map[string]struct{}
	rindex   int
	mindex   int
}

var md = &markdown{
	content:  new(bytes.Buffer),
	resource: make(map[string]struct{}),
}

func (m *markdown) setTitle(appname string) {
	m.content.WriteString(fmt.Sprintf("## %s\n", appname))
}
func (m *markdown) appendMethod(version, resource, action string, in, out reflect.Type) {
	_, ok := m.resource[resource]
	if !ok {
		m.resource[resource] = struct{}{}
		m.rindex++
		m.mindex = 0
		m.content.WriteString(fmt.Sprintf("### %d. resource: %s\n", m.rindex, resource))
	}
	if ok {
		m.mindex++
		m.content.WriteString(fmt.Sprintf("#### %d. action: %s \t version:%s\n", m.mindex, action, version))
		m.content.WriteString("```\n")
		m.content.WriteString(fmt.Sprintf("input param:%s\n", in.Name()))
		m.content.WriteString(fmt.Sprintf("output param:%s\n", out.Name()))
		m.content.WriteString("```\n")
	}
}

func (m *markdown) genDoc() {
	ioutil.WriteFile("./document.md", m.content.Bytes(), 0644)
}
