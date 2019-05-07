package server

import (
	"errors"
	"fmt"
	"github.com/golang/protobuf/descriptor"
	"github.com/sirupsen/logrus"
	"reflect"
	"regexp"
	"strings"
)

const (
	syntax         = "proto3"
	Pattern        = "^[A-Z1-9]{5}$"
	actNamePattern = "(_V[1-9]+)+((_[0-9]*)|())"
	defaultVersion = "V1"
	strEmpty       = ""
	line           = "_"
	dot            = "."
	service        = "Service"
	slen           = 7
)

var (
	reg, _        = regexp.Compile(actNamePattern)
	serviceReg, _ = regexp.Compile(service)

	ErrServiceType       = errors.New("service type must be struct")
	ErrServiceUseless    = errors.New("There is no public method in the service")
	ErrServiceNameSuffix = errors.New("The service name must end with the Service suffix")
	ErrServicesEmpty     = errors.New("The services cannot be empty.")
)

func parseResouce(serviceName string) (bool, string) {
	b := strings.Contains(serviceName, service)
	if !b {
		return b, strEmpty
	}
	return b, serviceName[0 : len(serviceName)-7]
}
func parseVersion(methodName string) (version, action string) {
	v := reg.FindString(methodName)
	if v == strEmpty {
		v = defaultVersion
	} else {
		v = v[1:]
	}
	arr := strings.Split(methodName, line)
	return strings.ToLower(strings.ReplaceAll(v, line, dot)), arr[0]
}
func parseServices(services ...interface{}) error {
	if len(services) == 0 {
		return ErrServicesEmpty
	}
	for k, s := range services {
		idx := k + 1
		typ := reflect.TypeOf(s)
		if typ.Elem().Kind() != reflect.Struct {
			return fmt.Errorf("the service index of %d is err:%v ", idx, ErrServiceType)
		}
		resource := typ.Elem().Name()
		b, r := parseResouce(resource)
		if !b {
			return fmt.Errorf("the service index of %d is err:%v ", idx, ErrServiceNameSuffix)
		}
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
			version, action := parseVersion(fnm)
			appInfo.registerMethod(version, r, action)
			err := delegate.registerMethod(version, r, action, a, in, out)
			if err != nil {
				return fmt.Errorf("the service index of %d and the method %s is err:%v the ", idx, fnm, err)
			}
			md.appendMethod(version, r, action, in, out)

			f = true
		}
		if !f {
			return fmt.Errorf("the service index of %d is err:%v ", idx, ErrServiceUseless)
		}
	}
	base.readyGo()
	md.genDoc()
	logrus.Info("The service group registered successfully.")
	return nil
}
