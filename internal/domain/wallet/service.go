package wallet

type Service interface{}

type walletService struct{}

func NewService() Service {
	return &walletService{}
}
