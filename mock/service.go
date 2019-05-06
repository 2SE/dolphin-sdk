package mock

var MockService = new(mockService)

type mockService struct {
}

func (s *mockService) GetUser_V2(request *GetUserRequest) (*User, error) {
	return nil, nil
}
func (s *mockService) GetUser(request *GetUserRequest) (*User, error) {
	return nil, nil
}
func (s *mockService) getUser_V1_123(request *GetUserRequest) (*User, error) {
	return nil, nil
}
func (s *mockService) GetUser_V1_123(request *GetUserRequest) (*User, error) {
	return nil, nil
}
func (s *mockService) Get1(a int) (b int, err error) {
	return 0, nil
}
func (s *mockService) Get2(request *GetUserRequest) (b int, err error) {
	return 0, nil
}
func (s *mockService) Get3(*GetUserRequest) (*User, int) {
	return nil, 0
}
