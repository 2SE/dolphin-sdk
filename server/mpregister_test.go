package server

import (
	"errors"
	"fmt"
	"github.com/2se/dolphin-sdk/mock"
	"github.com/2se/dolphin-sdk/mock/pb"
	"github.com/golang/protobuf/descriptor"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/stretchr/testify/assert"
	"os"
	"path"
	"reflect"
	"strings"
	"testing"
)

var mockService = new(MockService)

type MockService struct {
}

func (s *MockService) GetUser_V2(request *pb.GetUserRequest) (*pb.User, error) {
	return nil, nil
}
func (s *MockService) GetUser(request *pb.GetUserRequest) (*pb.User, error) {
	return nil, nil
}
func (s *MockService) getUser_V1_123(request *pb.GetUserRequest) (*pb.User, error) {
	return nil, nil
}
func (s *MockService) GetUser_V1_123(request *pb.GetUserRequest) (*pb.User, error) {
	return nil, nil
}
func (s *MockService) Get1(a int) (b int, err error) {
	return 0, nil
}
func (s *MockService) Get2(request *pb.GetUserRequest) (b int, err error) {
	return 0, nil
}
func (s *MockService) Get3(*pb.GetUserRequest) (*pb.User, int) {
	return nil, 0
}

func TestReflectMethod(t *testing.T) {

	typ := reflect.TypeOf(mockService)
	//vtyp := reflect.ValueOf(mockService)
	fmt.Println(typ.Elem().Name())
	fmt.Println(typ.Elem().Kind() == reflect.Struct)

	fmt.Println("filepath", typ.Elem().PkgPath())

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
		v, n := parseVersion(a.Name)
		fmt.Printf("action: %s version :%s\n", n, v)
		fmt.Println("-------------------")
	}
}

func Get(pb *pb.GetUserRequest) (*pb.User, error) {
	//return &mock.User{UserId: 1, UserName: "Jack", Age: 20}, nil
	return nil, errors.New("this is error")
}

func TestMethodCall(t *testing.T) {
	typ := reflect.ValueOf(Get)

	val := getAny()
	tmp := reflect.New(reflect.TypeOf(pb.GetUserRequest{})).Interface().(descriptor.Message)

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

	pm := &pb.GetUserRequest{
		UserId: 10086,
	}
	object, _ := ptypes.MarshalAny(pm)
	return object
}
func TestPasreResource(t *testing.T) {
	s := "UserService"
	b, n := parseResouce(s)
	assert.Equal(t, b, true)
	assert.Equal(t, n, "User")

	s = "Userservice"
	b, n = parseResouce(s)
	assert.Equal(t, b, false)
	assert.Equal(t, n, "")
}
func TestPasreVersion(t *testing.T) {
	v, m := parseVersion("getHash_V12334")
	assert.Equal(t, v, "v12334")
	assert.Equal(t, m, "getHash")
	v, m = parseVersion("getHash_V1_2334")
	assert.Equal(t, v, "v1.2334")
	assert.Equal(t, m, "getHash")
	v, m = parseVersion("getHash_V_12334")
	assert.Equal(t, v, "v1")
	assert.Equal(t, m, "getHash")
	v, m = parseVersion("getHash_V0_2334")
	assert.Equal(t, v, "v1")
	assert.Equal(t, m, "getHash")
	v, m = parseVersion("getHash_V02334")
	assert.Equal(t, v, "v1")
	assert.Equal(t, m, "getHash")
	v, m = parseVersion("getHash_V1233_4")
	assert.Equal(t, v, "v1233.4")
	assert.Equal(t, m, "getHash")
	v, m = parseVersion("getHash")
	assert.Equal(t, v, "v1")
	assert.Equal(t, m, "getHash")
	v, m = parseVersion("getHashV1233_4")
	assert.Equal(t, v, "v1")
	assert.Equal(t, m, "getHashV1233")
}
func TestVersionReplace(t *testing.T) {
	version := "V1_123"
	newv := strings.ToLower(strings.ReplaceAll(version, "_", "."))
	fmt.Println(newv)

	fmt.Println(reflect.TypeOf(pb.GetUserRequest{}).String())
}

func TestGenDoc(t *testing.T) {
	dir, _ := os.Getwd()
	dirTmp := path.Join(getParentDirectory(dir), "mock")
	//fmt.Println(os.Getwd())
	GenDoc("appName", []string{dirTmp}, mock.MkService)
}

func substr(s string, pos, length int) string {
	runes := []rune(s)
	l := pos + length
	if l > len(runes) {
		l = len(runes)
	}
	return string(runes[pos:l])
}

func getParentDirectory(dirctory string) string {
	return substr(dirctory, 0, strings.LastIndex(dirctory, "/"))
}
