package service

type Service interface{}

type walletService struct{}

func New() Service {
	return &walletService{}
}
