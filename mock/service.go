package mock

import (
	"errors"
	pb2 "github.com/2se/dolphin-sdk/mock/pb"
	"github.com/2se/dolphin-sdk/pb"
)

var MkService = new(MockUserService)

type MockUserService struct {
}

//GetUser v2
func (s *MockUserService) GetUser_V2(request *pb2.GetUserRequest) (*pb2.User, error) {
	//SendRequest()
	//s.sendRequest()
	if request.UserId == 1 {
		return &pb2.User{UserId: 1, UserName: "Dolphin", Age: 25}, nil
	} else {
		return nil, errors.New("user not found")
	}
}

//GetUser v1
func (s *MockUserService) GetUser(request *pb2.GetUserRequest) (*pb2.User, error) {
	return &pb2.User{UserId: 2, UserName: "Jack", Age: 30}, nil
}

//GetUser v1.123
func (s *MockUserService) GetUser_V1_123(request *pb2.GetUserRequest) (*pb2.User, error) {
	return nil, nil
}

func (s *MockUserService) getUser_V1_123(request *pb2.GetUserRequest) (*pb2.User, error) {
	return nil, nil
}
func (s *MockUserService) Get1(a int) (b int, err error) {
	return 0, nil
}
func (s *MockUserService) GetUser_V3(info *pb.CurrentInfo, request *pb2.GetUserRequest) (err error) {
	return nil
}
func (s *MockUserService) Get3() (*pb2.User, int) {
	return nil, 0
}
func (s *MockUserService) NotParam() error {
	return nil
}
func (s *MockUserService) sendRequest() {
	/*	rep, err := server.SendGrpcRequest(
			&pb.MethodPath{
				Resource: "MockUser",
				Revision: "v1",
				Action:   "GetUser",
			},
		, //&pb2.GetUserRequest{},

		)
		if err != nil {
			fmt.Println("GetUser_V2 err:", err)
			return
		}
		if rep.Code == 200 {
			pmu := &pb2.User{}
			ptypes.UnmarshalAny(rep.Body, pmu)
			fmt.Println("response body :", pmu.String())
		}*/

}
