package currency

type Service interface{}

type currencyService struct{}

func NewService() Service {
	return &currencyService{}
}
