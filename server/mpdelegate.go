package server

import (
	"errors"
	"fmt"
	"github.com/2se/dolphin-sdk/pb"
	"github.com/golang/protobuf/descriptor"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/sirupsen/logrus"
	"reflect"
)

var (
	delegate = &mpdelegate{
		services:  make(map[string]reflect.Value),
		direction: make(map[string]map[string]map[string]*grpcMethod),
	}

	ErrOverloadNotSupported = errors.New("The registered service does not support overloading of version,resource,action")
	ErrParamNotSpecified    = errors.New("Parameter not specified")
)

const (
	panicCode = "1"
	panicStr  = "An unknowable error"
)

type grpcMethod struct {
	method reflect.Method
	numIn  int
	numOut int
	argin  reflect.Type
	argout reflect.Type
}

type mpdelegate struct {
	services  map[string]reflect.Value
	direction map[string]map[string]map[string]*grpcMethod
}

func (m *mpdelegate) registerService(resource string, value reflect.Value) {
	m.services[resource] = value
}
func (m *mpdelegate) registerMethod(version, resource, action string, mehtod reflect.Method, in, out reflect.Type, numIn, numOut int) error {
	if m.direction[resource] == nil {
		m.direction[resource] = make(map[string]map[string]*grpcMethod)
	}
	if m.direction[resource][version] == nil {
		m.direction[resource][version] = make(map[string]*grpcMethod)
	}
	if _, ok := m.direction[resource][version][action]; ok {
		return ErrOverloadNotSupported
	}

	m.direction[resource][version][action] = &grpcMethod{
		method: mehtod,
		argin:  in,
		argout: out,
		numIn:  numIn,
		numOut: numOut,
	}
	return nil
}
func catchPanic(req *pb.ClientComRequest) {
	if err := recover(); err != nil {
		logrus.WithFields(logrus.Fields{
			"resource": req.MethodPath.Resource,
			"version":  req.MethodPath.Revision,
			"action":   req.MethodPath.Action,
			"traceId":  req.TraceId,
			"panic":    panicCode,
		}).Error(err)

	}
}
func (m *mpdelegate) invoke(req *pb.ClientComRequest) *pb.ServerComResponse {
	defer catchPanic(req)
	response := &pb.ServerComResponse{
		Id:      req.Id,
		TraceId: req.TraceId,
		Code:    200,
	}
	fmt.Println(m)
	grpcM := m.direction[req.MethodPath.Resource][req.MethodPath.Revision][req.MethodPath.Action]
	inputs := make([]reflect.Value, grpcM.numIn)
	inputs[0] = m.services[req.MethodPath.Resource]
	if grpcM.numIn == 2 {
		tmp := reflect.New(grpcM.argin).Interface().(descriptor.Message)
		err := ptypes.UnmarshalAny(req.Params, tmp)
		if err != nil {
			response.Code = 400
			response.Text = fmt.Sprintf("%s,and the param type is %s", ErrParamNotSpecified.Error(), grpcM.argin.String())
			return response
		}
		inputs[1] = reflect.ValueOf(tmp)
	}

	vals := grpcM.method.Func.Call(inputs)

	errIndx := 0
	if len(vals) == 2 {
		errIndx = 1
		if vals[0].Elem().Type() == grpcM.argout && !vals[0].IsNil() {
			object, err := ptypes.MarshalAny(vals[0].Interface().(proto.Message))
			if err != nil {
				response.Code = 500
				response.Text = err.Error()
				return response
			}
			response.Body = object
		}
	}
	if !vals[errIndx].IsNil() {
		response.Code = 500
		response.Text = vals[errIndx].Interface().(error).Error()
	}
	return response
}
