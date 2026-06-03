package service

import (
	"Currency-exchange/internal/models"
	"Currency-exchange/internal/repository"
)

type CurrencyService struct {
	repo *repository.CurrencyRepository
}

func NewCurrencyService(repo *repository.CurrencyRepository) *CurrencyService {
	return &CurrencyService{repo: repo}
}

func (s *CurrencyService) GetAllCurrencies() ([]models.Currency, error) {
	return s.repo.GetAll()
}

func (s *CurrencyService) GetCurrencyByCode(code string) (models.Currency, error) {
	return s.repo.GetByCode(code)
}

func (s *CurrencyService) CreateCurrency(c models.Currency) (models.Currency, error) {
	return s.repo.CreateCurrencies(c)
}
