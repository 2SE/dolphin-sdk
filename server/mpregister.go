package server

import (
	"errors"
	"github.com/golang/protobuf/descriptor"
	"reflect"
	"regexp"
)

const (
	syntax         = "proto3"
	Pattern        = "^[A-Z1-9]{5}$"
	actNamePattern = "(_V[1-9]{1})+((_[0-9]*)|())"
	defaultVersion = "V1"
	strEmpty       = ""
)

var (
	reg, _               = regexp.Compile(actNamePattern)
	ErrServiceType       = errors.New("service type must be struct")
	ErrServiceNotThrough = errors.New("The service does not comply with the specification")
)

func parseVersion(methodName string) string {
	v := reg.FindString(methodName)
	if v == strEmpty {
		v = defaultVersion
	} else {
		v = v[1:]
	}
	return v
}
func parseService(s interface{}) error {
	typ := reflect.TypeOf(s)
	if typ.Kind() != reflect.Struct {
		return ErrServiceType
	}
	resource := typ.Elem().Name()
	f := false
	for i := 0; i < typ.NumMethod(); i++ {
		a := typ.Method(i)
		if a.Func.Type().NumIn() != 2 || a.Func.Type().NumOut() != 2 {
			continue
		}
		if a.Func.Type().In(1).Kind() != reflect.Ptr || a.Func.Type().Out(0).Kind() != reflect.Ptr || a.Func.Type().Out(1).Name() != "error" {
			continue
		}
		in := a.Func.Type().In(1).Elem()
		p1 := reflect.New(in).Interface().(descriptor.Message)
		fd1, _ := descriptor.ForMessage(p1)
		if *fd1.Syntax != syntax {
			continue
		}
		out := a.Func.Type().Out(0).Elem()
		p2 := reflect.New(out).Interface().(descriptor.Message)
		fd2, _ := descriptor.ForMessage(p2)
		if *fd2.Syntax != syntax {
			continue
		}
		fnm := a.Name
		version := parseVersion(fnm)
		appInfo.registerMethod(version, resource, fnm)
		delegate.registerMethod(version, resource, fnm, a, in, out)
		f = true
	}
	if f {
		return nil
	}
	return ErrServiceNotThrough
}
