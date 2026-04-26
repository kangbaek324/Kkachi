package service

type Service interface{}

type currencyService struct{}

func New() Service {
	return &currencyService{}
}
