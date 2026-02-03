package service

import (
	"exchangeapp/dao"
	"exchangeapp/models"
)

type ExchangeRateService interface {
	CreateExchangeRate(exchangeRate *models.ExchangeRate) error
	GetExchangeRates() ([]models.ExchangeRate, error)
	GetExchangeRateByID(id uint) (*models.ExchangeRate, error)
	GetExchangeRateByCurrency(fromCurrency, toCurrency string) (*models.ExchangeRate, error)
	UpdateExchangeRate(exchangeRate *models.ExchangeRate) error
	DeleteExchangeRate(id uint) error
}

type exchangeRateService struct {
	exchangeRateDAO dao.ExchangeRateDAO
}

func NewExchangeRateService(exchangeRateDAO dao.ExchangeRateDAO) ExchangeRateService {
	return &exchangeRateService{
		exchangeRateDAO: exchangeRateDAO,
	}
}

func (ers *exchangeRateService) CreateExchangeRate(exchangeRate *models.ExchangeRate) error {
	// 业务逻辑：验证汇率数据
	if exchangeRate.Rate <= 0 {
		return &AppError{Code: "INVALID_RATE", Message: "汇率必须大于0"}
	}

	if exchangeRate.FromCurrency == exchangeRate.ToCurrency {
		return &AppError{Code: "SAME_CURRENCY", Message: "源币种和目标币种不能相同"}
	}

	return ers.exchangeRateDAO.CreateExchangeRate(exchangeRate)
}

func (ers *exchangeRateService) GetExchangeRates() ([]models.ExchangeRate, error) {
	return ers.exchangeRateDAO.GetExchangeRates()
}

func (ers *exchangeRateService) GetExchangeRateByID(id uint) (*models.ExchangeRate, error) {
	return ers.exchangeRateDAO.GetExchangeRateByID(id)
}

func (ers *exchangeRateService) GetExchangeRateByCurrency(fromCurrency, toCurrency string) (*models.ExchangeRate, error) {
	return ers.exchangeRateDAO.GetExchangeRateByCurrency(fromCurrency, toCurrency)
}

func (ers *exchangeRateService) UpdateExchangeRate(exchangeRate *models.ExchangeRate) error {
	// 业务逻辑：验证汇率数据
	if exchangeRate.Rate <= 0 {
		return &AppError{Code: "INVALID_RATE", Message: "汇率必须大于0"}
	}

	return ers.exchangeRateDAO.UpdateExchangeRate(exchangeRate)
}

func (ers *exchangeRateService) DeleteExchangeRate(id uint) error {
	return ers.exchangeRateDAO.DeleteExchangeRate(id)
}
