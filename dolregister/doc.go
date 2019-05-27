package dolregister

import (
	"bytes"
	"fmt"
	"github.com/2se/dolphin-sdk/pb"
	"io/ioutil"
	"reflect"
)

var curType = reflect.TypeOf(pb.CurrentInfo{})

type Doc interface {
	SetTitle(appname string)
	AppendMethod(version, resource, action, doc string, ins []reflect.Type, out reflect.Type, numIn, numOut int)
	//生成document
	GenDoc()
}
type markdown struct {
	content  *bytes.Buffer
	resource map[string]struct{}
	rindex   int
	mindex   int
}

func (m *markdown) SetTitle(appname string) {
	m.content.WriteString(fmt.Sprintf("## %s\n", appname))
}
func (m *markdown) AppendMethod(version, resource, action, doc string, ins []reflect.Type, out reflect.Type, numIn, numOut int) {
	_, ok := m.resource[resource]
	if !ok {
		m.resource[resource] = struct{}{}
		m.rindex++
		m.mindex = 0
		m.content.WriteString(fmt.Sprintf("### %d. resource: %s\n", m.rindex, resource))
		ok = true
	}
	if ok {
		m.mindex++
		m.content.WriteString(fmt.Sprintf("#### %d. action: %s \t version:%s\n", m.mindex, action, version))
		m.content.WriteString("```\n")
		if doc != "" {
			m.content.WriteString("describe：\n")
			m.content.WriteString(doc)
			m.content.WriteString("\n")
		}
		if numIn > 1 {
			for _, v := range ins {
				if v != curType {
					m.content.WriteString(fmt.Sprintf("input param:%s\n", v.Name()))
				}
			}
		}
		if numOut == 2 {
			m.content.WriteString(fmt.Sprintf("output param:%s\n", out.Name()))
		}
		m.content.WriteString("```\n")
	}
}

func (m *markdown) GenDoc() {
	ioutil.WriteFile("./document.md", m.content.Bytes(), 0644)
}
