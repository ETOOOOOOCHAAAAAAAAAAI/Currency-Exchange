package service

import (
	"Currency-exchange/internal/models"
	"Currency-exchange/internal/repository"
	"errors"
	"math"
)

type ExchangeService struct {
	repo         *repository.ExchangeRateRepository
	currencyRepo *repository.CurrencyRepository
}

func NewExchangeService(repo *repository.ExchangeRateRepository, currencyRepo *repository.CurrencyRepository) *ExchangeService {
	return &ExchangeService{repo: repo, currencyRepo: currencyRepo}
}

func (s *ExchangeService) GetAllExchangeRates() ([]models.ExcangeRate, error) {
	return s.repo.GetAll()
}

func (s *ExchangeService) GetRateByCode(base, target string) (models.ExcangeRate, error) {
	return s.repo.GetRateByCode(base, target)
}

func (s *ExchangeService) CreateExchangeRate(e models.ExcangeRate) (models.ExcangeRate, error) {
	return s.repo.CreateNewExchangeRate(e)
}

func (s *ExchangeService) UpdateExchangeRate(e models.ExcangeRate) (models.ExcangeRate, error) {
	return s.repo.UpdateExchangeRate(e)
}

func (s *ExchangeService) CalculateExchange(from, to string, amount float64) (models.ExchangeResponse, error) {
	var finalRate float64
	var response models.ExchangeResponse
	rate, err := s.repo.GetRateByCode(from, to)
	if err == nil {
		finalRate = rate.Rate
	} else {
		rateReverse, err := s.repo.GetRateByCode(to, from)
		if err == nil {
			finalRate = 1.0 / rateReverse.Rate
		} else {
			rateUSDFrom, err1 := s.repo.GetRateByCode("USD", from)
			rateUSDTo, err2 := s.repo.GetRateByCode("USD", to)
			if err1 == nil && err2 == nil {
				finalRate = rateUSDTo.Rate / rateUSDFrom.Rate
			} else {
				return response, errors.New("обменный курс для пары не найден")
			}
		}
	}
	baseCurrency, err := s.currencyRepo.GetByCode(from)
	if err != nil {
		return response, errors.New("базовая валюта не найдена")
	}
	targetCurrency, err := s.currencyRepo.GetByCode(to)
	if err != nil {
		return response, errors.New("целевая валюта не найдена")
	}

	convertedAmount := amount * finalRate
	finalRate = math.Round(finalRate*100) / 100
	convertedAmount = math.Round(convertedAmount*100) / 100
	response = models.ExchangeResponse{
		BaseCurrency:    baseCurrency,
		TargetCurrency:  targetCurrency,
		Rate:            finalRate,
		Amount:          amount,
		ConvertedAmount: convertedAmount,
	}
	return response, nil
}
