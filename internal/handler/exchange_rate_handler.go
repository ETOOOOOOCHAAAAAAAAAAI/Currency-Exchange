package handler

import (
	"Currency-exchange/internal/models"
	"Currency-exchange/internal/repository"
	"encoding/json"
	"net/http"
	"strconv"
)

type ExchangeRateHandler struct {
	repo         *repository.ExchangeRateRepository
	currencyRepo *repository.CurrencyRepository
}

type ExchangeResponse struct {
	BaseCurrency    models.Currency `json:"baseCurrency"`
	TargetCurrency  models.Currency `json:"targetCurrency"`
	Rate            float64         `json:"rate"`
	Amount          float64         `json:"amount"`
	ConvertedAmount float64         `json:"convertedAmount"`
}

func NewExchangeRateHandler(repo *repository.ExchangeRateRepository, currencyRepo *repository.CurrencyRepository) *ExchangeRateHandler {
	return &ExchangeRateHandler{repo: repo,
		currencyRepo: currencyRepo,
	}

}

func (h *ExchangeRateHandler) GetAllExchangeRate(w http.ResponseWriter, r *http.Request) {
	rate, err := h.repo.GetAll()
	if err != nil {
		SendJSONError(w, "Ошибка при получений данных из обменника БД", http.StatusInternalServerError)
		return
	}
	jsonBytes, err := json.Marshal(rate)
	if err != nil {
		SendJSONError(w, "Ошибка при формираваний ответа", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)

}

func (h *ExchangeRateHandler) GetByCodeInExchangeRate(w http.ResponseWriter, r *http.Request) {
	codes := r.PathValue("codes")
	if len(codes) != 6 {
		SendJSONError(w, "Неверный формат кода валют", http.StatusBadRequest)
		return
	}
	baseCode := codes[:3]
	targetCode := codes[3:]
	if baseCode == targetCode {
		SendJSONError(w, "Одинаковый код валют", http.StatusBadRequest)
		return
	}
	rate, err := h.repo.GetRateByCode(baseCode, targetCode)
	if err != nil {
		SendJSONError(w, "Код валют не найден", http.StatusNotFound)
		return
	}
	jsonBytes, err := json.Marshal(rate)
	if err != nil {
		SendJSONError(w, "Ошибка при формировании ответа", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

func (h *ExchangeRateHandler) CreateNewExchangeRate(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		SendJSONError(w, "Ошибка при чтении данных", http.StatusInternalServerError)
		return
	}
	newBaseCurrencyCode := r.FormValue("baseCurrencyCode")
	newTargetCurrencyCode := r.FormValue("targetCurrencyCode")
	if newBaseCurrencyCode == newTargetCurrencyCode {
		SendJSONError(w, "Одинаковый код валют", http.StatusBadRequest)
		return
	}
	newRate := r.FormValue("rate")
	if newBaseCurrencyCode == "" || newTargetCurrencyCode == "" || newRate == "" {
		SendJSONError(w, "Неверный запрос", http.StatusBadRequest)
		return
	}
	baseСurrency, err := h.currencyRepo.GetByCode(newBaseCurrencyCode)
	if err != nil {
		SendJSONError(w, "Такой валюты не существует", http.StatusNotFound)
		return
	}
	targetCurrency, err := h.currencyRepo.GetByCode(newTargetCurrencyCode)
	if err != nil {
		SendJSONError(w, "Такой валюты не существует", http.StatusNotFound)
		return
	}
	rateFloat, err := strconv.ParseFloat(newRate, 64)
	if err != nil {
		SendJSONError(w, "Ошибка при конвертации соотношения", http.StatusBadRequest)
		return
	}
	e := models.ExcangeRate{BaseCurrency: baseСurrency, TargetCurrency: targetCurrency, Rate: rateFloat}
	result, err := h.repo.CreateNewExchangeRate(e)
	if err != nil {
		SendJSONError(w, "Такое соотношение уже существует", http.StatusConflict)
		return
	}
	jsonBytes, err := json.Marshal(result)
	if err != nil {
		SendJSONError(w, "Ошибка при формировании ответа", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonBytes)
}

func (h *ExchangeRateHandler) UpdateExchangeRate(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		SendJSONError(w, "Ошибка при чтении формы", http.StatusInternalServerError)
		return
	}
	codes := r.PathValue("codes")
	if len(codes) != 6 {
		SendJSONError(w, "Неверный формат кода валют", http.StatusBadRequest)
		return
	}
	baseCurrencyCode := codes[:3]
	targetCurrencyCode := codes[3:]
	if baseCurrencyCode == targetCurrencyCode {
		SendJSONError(w, "Одинаковый код валют", http.StatusBadRequest)
		return
	}
	baseCurrency, err := h.currencyRepo.GetByCode(baseCurrencyCode)
	if err != nil {
		SendJSONError(w, "Такой валюты не существует", http.StatusNotFound)
		return
	}
	targetCurrency, err := h.currencyRepo.GetByCode(targetCurrencyCode)
	if err != nil {
		SendJSONError(w, "Такой валюты не существует", http.StatusNotFound)
		return
	}
	newRate := r.FormValue("rate")
	rateFloat, err := strconv.ParseFloat(newRate, 64)
	if err != nil {
		SendJSONError(w, "Ошибка при конвертации соотнешия", http.StatusBadRequest)
		return
	}
	e := models.ExcangeRate{BaseCurrency: baseCurrency, TargetCurrency: targetCurrency, Rate: rateFloat}
	result, err := h.repo.UpdateExchangeRate(e)
	if err != nil {
		SendJSONError(w, "Соотшения не найден", http.StatusNotFound)
		return
	}
	jsonBytes, err := json.Marshal(result)
	if err != nil {
		SendJSONError(w, "Ошибка при формировании ответа", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

func (h *ExchangeRateHandler) ExchangeCalculator(w http.ResponseWriter, r *http.Request) {
	from := r.URL.Query().Get("from")
	to := r.URL.Query().Get("to")
	amountStr := r.URL.Query().Get("amount")
	if from == "" || to == "" || amountStr == "" {
		SendJSONError(w, "Отсутствуют необходимые параметры запроса", http.StatusBadRequest)
		return
	}
	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		SendJSONError(w, "Неверный формат суммы", http.StatusBadRequest)
		return
	}
	var finalRate float64
	rate, err := h.repo.GetRateByCode(from, to)
	if err == nil {
		finalRate = rate.Rate
	} else {
		rateReverse, err := h.repo.GetRateByCode(to, from)
		if err == nil {
			finalRate = 1.0 / rateReverse.Rate
		} else {
			rateUSDFrom, err1 := h.repo.GetRateByCode("USD", from)
			rateUSDTo, err2 := h.repo.GetRateByCode("USD", to)
			if err1 == nil && err2 == nil {
				finalRate = rateUSDTo.Rate / rateUSDFrom.Rate
			} else {
				SendJSONError(w, "Обменный курс для пары не найден", http.StatusNotFound)
				return
			}
		}
	}
	convertedAmount := amount * finalRate
	baseCurrency, err := h.currencyRepo.GetByCode(from)
	if err != nil {
		SendJSONError(w, "Базовая валюта не найдена", http.StatusNotFound)
		return
	}
	targetCurrency, err := h.currencyRepo.GetByCode(to)
	if err != nil {
		SendJSONError(w, "Целевая валюта не найдена", http.StatusNotFound)
		return
	}
	response := ExchangeResponse{
		BaseCurrency:    baseCurrency,
		TargetCurrency:  targetCurrency,
		Rate:            finalRate,
		Amount:          amount,
		ConvertedAmount: convertedAmount,
	}
	jsonBytes, err := json.Marshal(response)
	if err != nil {
		SendJSONError(w, "Ошибка при формировании ответа", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}
