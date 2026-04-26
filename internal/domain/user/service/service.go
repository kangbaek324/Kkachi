package service

type Service interface{}

type userService struct{}

func New() Service {
	return &userService{}
}
