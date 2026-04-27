package user

type Service interface{}

type userService struct{}

func NewService() Service {
	return &userService{}
}
