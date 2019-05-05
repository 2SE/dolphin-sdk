package server

import (
	"errors"
	"fmt"
	"github.com/golang/protobuf/descriptor"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/stretchr/testify/assert"
	"gitlab.2se.com/hashhash/server-sdk/mock"
	"reflect"
	"testing"
)

var mockService = new(MockService)

type MockService struct {
}

func (s *MockService) GetUser_V2(request *mock.GetUserRequest) (*mock.User, error) {
	return nil, nil
}
func (s *MockService) GetUser(request *mock.GetUserRequest) (*mock.User, error) {
	return nil, nil
}
func (s *MockService) getUser_V1_123(request *mock.GetUserRequest) (*mock.User, error) {
	return nil, nil
}
func (s *MockService) GetUser_V1_123(request *mock.GetUserRequest) (*mock.User, error) {
	return nil, nil
}
func (s *MockService) Get1(a int) (b int, err error) {
	return 0, nil
}
func (s *MockService) Get2(request *mock.GetUserRequest) (b int, err error) {
	return 0, nil
}
func (s *MockService) Get3(*mock.GetUserRequest) (*mock.User, int) {
	return nil, 0
}

func TestReflectMethod(t *testing.T) {
	typ := reflect.TypeOf(mockService)
	//vtyp := reflect.ValueOf(mockService)
	fmt.Println(typ.Elem().Name())

	for i := 0; i < typ.NumMethod(); i++ {
		a := typ.Method(i)

		if a.Func.Type().NumIn() != 2 || a.Func.Type().NumOut() != 2 {
			continue
		}
		if a.Func.Type().In(1).Kind() != reflect.Ptr || a.Func.Type().Out(0).Kind() != reflect.Ptr || a.Func.Type().Out(1).Name() != "error" {
			continue
		}
		p1 := reflect.New(a.Func.Type().In(1).Elem()).Interface().(descriptor.Message)
		fd1, _ := descriptor.ForMessage(p1)
		if *fd1.Syntax != syntax {
			continue
		}
		val := getAny()
		tmp := reflect.New(a.Func.Type().In(1).Elem()).Interface().(descriptor.Message)

		fmt.Println(tmp)
		err := ptypes.UnmarshalAny(val, tmp)
		if err != nil {
			fmt.Println(err)
			return
		}
		v := parseVersion(a.Name)
		fmt.Println("version :", v)
		fmt.Println("-------------------")
	}
}

func Get(pb *mock.GetUserRequest) (*mock.User, error) {
	//return &mock.User{UserId: 1, UserName: "Jack", Age: 20}, nil
	return nil, errors.New("this is error")
}

func TestMethodCall(t *testing.T) {
	typ := reflect.ValueOf(Get)

	val := getAny()
	tmp := reflect.New(reflect.TypeOf(mock.GetUserRequest{})).Interface().(descriptor.Message)

	fmt.Println(tmp)
	err := ptypes.UnmarshalAny(val, tmp)
	if err != nil {
		fmt.Println(err)
		return
	}
	inputs := make([]reflect.Value, 1)
	inputs[0] = reflect.ValueOf(tmp)
	vals := typ.Call(inputs)

	if !vals[1].IsNil() {
		fmt.Println(500, vals[1].Interface().(error).Error())
	}
}

func getAny() *any.Any {

	pm := &mock.GetUserRequest{
		UserId: 10086,
	}
	object, _ := ptypes.MarshalAny(pm)
	return object
}
func TestPasreVersion(t *testing.T) {
	assert.Equal(t, parseVersion("getHash_V12334"), "V1")
	assert.Equal(t, parseVersion("getHash_V1_2334"), "V1_2334")
	assert.Equal(t, parseVersion("getHash_V_12334"), "V1")
	assert.Equal(t, parseVersion("getHash_V0_2334"), "V1")
	assert.Equal(t, parseVersion("getHash_V02334"), "V1")
	assert.Equal(t, parseVersion("getHash_V1233_4"), "V1")
	assert.Equal(t, parseVersion("getHash"), "V1")
	assert.Equal(t, parseVersion("getHashV1233_4"), "V1")
}
