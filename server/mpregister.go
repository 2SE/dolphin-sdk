package server

import (
	"errors"
	"fmt"
	"github.com/2se/dolphin-sdk/dolregister"
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
	reg, _          = regexp.Compile(actNamePattern)
	serviceReg, _   = regexp.Compile(service)
	registerManager = dolregister.NewRegisterManager()

	ErrServiceType       = errors.New("service type must be struct")
	ErrServiceUseless    = errors.New("There is no public method in the service")
	ErrServiceNameSuffix = errors.New("The service name must end with the Service suffix")
	ErrServicesEmpty     = errors.New("The services cannot be empty.")
)

func parseResouce(serviceName string) (bool, string) {
	l := len(serviceName)
	if l <= slen {
		return false, strEmpty
	}
	suf := serviceName[l-slen:]
	if suf != service {
		return false, strEmpty
	}
	return true, serviceName[:len(serviceName)-7]
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
			numIn := a.Func.Type().NumIn()
			numOut := a.Func.Type().NumOut()
			if numIn > 3 || numOut > 2 {
				continue
			}
			var (
				out reflect.Type
				ins = make([]reflect.Type, numIn-1)
			)
			if numIn > 1 {
				fl := false
				for i := 1; i < numIn; i++ {
					if a.Func.Type().In(i).Kind() != reflect.Ptr {
						fl = true
						break
					}
					ins[i-1] = a.Func.Type().In(i).Elem()
					p1 := reflect.New(ins[i-1]).Interface().(descriptor.Message)
					fd1, _ := descriptor.ForMessage(p1)
					if *fd1.Syntax != syntax {
						fl = true
						break
					}
				}
				if fl {
					continue
				}
			}
			if numOut == 0 {

			} else if numOut == 1 {
				if a.Func.Type().Out(0).Name() != "error" {
					continue
				}
			} else {
				if a.Func.Type().Out(0).Kind() != reflect.Ptr || a.Func.Type().Out(1).Name() != "error" {
					continue
				}
				out = a.Func.Type().Out(0).Elem()
				p2 := reflect.New(out).Interface().(descriptor.Message)
				fd2, _ := descriptor.ForMessage(p2)
				if *fd2.Syntax != syntax {
					continue
				}
			}
			fnm := a.Name
			version, action := parseVersion(fnm)
			registerManager.RegisterMethod(version, r, action)
			err := delegate.registerMethod(version, r, action, a, ins, out, numIn, numOut)
			if err != nil {
				return fmt.Errorf("the service index of %d and the method %s is err:%v the ", idx, fnm, err)
			}
			registerManager.AppendMethod(version, r, action, ins, out, numIn, numOut)
			f = true
		}
		if !f {
			return fmt.Errorf("the service index of %d is err:%v ", idx, ErrServiceUseless)
		}
		delegate.registerService(r, reflect.ValueOf(s))
	}
	base.readyGo()
	registerManager.GenDoc()
	logrus.Info("The service group registered successfully.")
	return nil
}
