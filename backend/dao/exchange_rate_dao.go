package dao

import (
	"exchangeapp/global"
	"exchangeapp/models"

	"gorm.io/gorm"
)

type ExchangeRateDAO interface {
	CreateExchangeRate(exchangeRate *models.ExchangeRate) error
	GetExchangeRates() ([]models.ExchangeRate, error)
	GetExchangeRateByID(id uint) (*models.ExchangeRate, error)
	GetExchangeRateByCurrency(fromCurrency, toCurrency string) (*models.ExchangeRate, error)
	UpdateExchangeRate(exchangeRate *models.ExchangeRate) error
	DeleteExchangeRate(id uint) error
}

type exchangeRateDAO struct {
	db *gorm.DB
}

func NewExchangeRateDAO() ExchangeRateDAO {
	return &exchangeRateDAO{
		db: global.Db,
	}
}

func (erd *exchangeRateDAO) CreateExchangeRate(exchangeRate *models.ExchangeRate) error {
	return erd.db.Create(exchangeRate).Error
}

func (erd *exchangeRateDAO) GetExchangeRates() ([]models.ExchangeRate, error) {
	var exchangeRates []models.ExchangeRate
	if err := erd.db.Find(&exchangeRates).Error; err != nil {
		return nil, err
	}
	return exchangeRates, nil
}

func (erd *exchangeRateDAO) GetExchangeRateByID(id uint) (*models.ExchangeRate, error) {
	var exchangeRate models.ExchangeRate
	if err := erd.db.Where("id = ?", id).First(&exchangeRate).Error; err != nil {
		return nil, err
	}
	return &exchangeRate, nil
}

func (erd *exchangeRateDAO) GetExchangeRateByCurrency(fromCurrency, toCurrency string) (*models.ExchangeRate, error) {
	var exchangeRate models.ExchangeRate
	if err := erd.db.Where("from_currency = ? AND to_currency = ?", fromCurrency, toCurrency).
		Order("created_at DESC").
		First(&exchangeRate).Error; err != nil {
		return nil, err
	}
	return &exchangeRate, nil
}

func (erd *exchangeRateDAO) UpdateExchangeRate(exchangeRate *models.ExchangeRate) error {
	return erd.db.Save(exchangeRate).Error
}

func (erd *exchangeRateDAO) DeleteExchangeRate(id uint) error {
	return erd.db.Delete(&models.ExchangeRate{}, id).Error
}
