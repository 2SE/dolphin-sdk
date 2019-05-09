package mock

import (
	"errors"
	"github.com/2se/dolphin-sdk/scheduler"
)

var MkService = new(MockUserService)

type MockUserService struct {
}

func (s *MockUserService) GetUser_V2(request *GetUserRequest) (*User, error) {
	if request.UserId == 1 {
		scheduler.SendRequest(nil)
		return &User{UserId: 1, UserName: "Dolphin", Age: 25}, nil
	} else {
		return nil, errors.New("user not found")
	}
}
func (s *MockUserService) GetUser(request *GetUserRequest) (*User, error) {
	return nil, nil
}
func (s *MockUserService) getUser_V1_123(request *GetUserRequest) (*User, error) {
	return nil, nil
}
func (s *MockUserService) GetUser_V1_123(request *GetUserRequest) (*User, error) {
	return nil, nil
}
func (s *MockUserService) Get1(a int) (b int, err error) {
	return 0, nil
}
func (s *MockUserService) Get2(request *GetUserRequest) (b int, err error) {
	return 0, nil
}
func (s *MockUserService) Get3(*GetUserRequest) (*User, int) {
	return nil, 0
}

/*func (s *MockService) GetUser_V(request *GetUserRequest) (*User, error) {
	return nil, nil
}*/
