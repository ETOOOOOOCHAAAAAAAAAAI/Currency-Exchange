package handler

import (
	"Currency-exchange/internal/models"
	"Currency-exchange/internal/service"
	"encoding/json"
	"net/http"
	"strconv"
)

type ExchangeRateHandler struct {
	exchangeService *service.ExchangeService
	currencyService *service.CurrencyService
}

func NewExchangeRateHandler(exchangeService *service.ExchangeService, currencyService *service.CurrencyService) *ExchangeRateHandler {
	return &ExchangeRateHandler{exchangeService: exchangeService, currencyService: currencyService}
}

func (h *ExchangeRateHandler) GetAllExchangeRate(w http.ResponseWriter, r *http.Request) {
	rates, err := h.exchangeService.GetAllExchangeRates()
	if err != nil {
		SendJSONError(w, "Ошибка при получений данных", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(rates)
}

func (h *ExchangeRateHandler) GetByCodeInExchangeRate(w http.ResponseWriter, r *http.Request) {
	codes := r.PathValue("codes")
	if len(codes) != 6 {
		SendJSONError(w, "Неверный формат кода валют", http.StatusBadRequest)
		return
	}
	baseCode, targetCode := codes[:3], codes[3:]
	if baseCode == targetCode {
		SendJSONError(w, "Одинаковый код валют", http.StatusBadRequest)
		return
	}
	rate, err := h.exchangeService.GetRateByCode(baseCode, targetCode)
	if err != nil {
		SendJSONError(w, "Код валют не найден", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(rate)
}

func (h *ExchangeRateHandler) CreateNewExchangeRate(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		SendJSONError(w, "Ошибка при чтении данных", http.StatusInternalServerError)
		return
	}
	newBaseCode := r.FormValue("baseCurrencyCode")
	newTargetCode := r.FormValue("targetCurrencyCode")
	newRate := r.FormValue("rate")
	if newBaseCode == "" || newTargetCode == "" || newRate == "" {
		SendJSONError(w, "Неверный запрос", http.StatusBadRequest)
		return
	}
	if newBaseCode == newTargetCode {
		SendJSONError(w, "Одинаковый код валют", http.StatusBadRequest)
		return
	}
	baseCurrency, err := h.currencyService.GetCurrencyByCode(newBaseCode)
	if err != nil {
		SendJSONError(w, "Базовая валюта не существует", http.StatusNotFound)
		return
	}
	targetCurrency, err := h.currencyService.GetCurrencyByCode(newTargetCode)
	if err != nil {
		SendJSONError(w, "Целевая валюта не существует", http.StatusNotFound)
		return
	}
	rateFloat, err := strconv.ParseFloat(newRate, 64)
	if err != nil {
		SendJSONError(w, "Ошибка при конвертации соотношения", http.StatusBadRequest)
		return
	}
	e := models.ExcangeRate{BaseCurrency: baseCurrency, TargetCurrency: targetCurrency, Rate: rateFloat}
	result, err := h.exchangeService.CreateExchangeRate(e)
	if err != nil {
		SendJSONError(w, "Такое соотношение уже существует", http.StatusConflict)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(result)
}

func (h *ExchangeRateHandler) UpdateExchangeRate(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		SendJSONError(w, "Ошибка при чтении формы", http.StatusInternalServerError)
		return
	}
	codes := r.PathValue("codes")
	if len(codes) != 6 {
		SendJSONError(w, "Неверный формат кода валют", http.StatusBadRequest)
		return
	}
	baseCode, targetCode := codes[:3], codes[3:]
	if baseCode == targetCode {
		SendJSONError(w, "Одинаковый код валют", http.StatusBadRequest)
		return
	}
	baseCurrency, err := h.currencyService.GetCurrencyByCode(baseCode)
	if err != nil {
		SendJSONError(w, "Базовая валюта не существует", http.StatusNotFound)
		return
	}
	targetCurrency, err := h.currencyService.GetCurrencyByCode(targetCode)
	if err != nil {
		SendJSONError(w, "Целевая валюта не существует", http.StatusNotFound)
		return
	}
	newRate := r.FormValue("rate")
	rateFloat, err := strconv.ParseFloat(newRate, 64)
	if err != nil {
		SendJSONError(w, "Ошибка при конвертации", http.StatusBadRequest)
		return
	}
	e := models.ExcangeRate{BaseCurrency: baseCurrency, TargetCurrency: targetCurrency, Rate: rateFloat}
	result, err := h.exchangeService.UpdateExchangeRate(e)
	if err != nil {
		SendJSONError(w, "Соотношение не найдено", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

func (h *ExchangeRateHandler) ExchangeCalculator(w http.ResponseWriter, r *http.Request) {
	from, to, amountStr := r.URL.Query().Get("from"), r.URL.Query().Get("to"), r.URL.Query().Get("amount")
	if from == "" || to == "" || amountStr == "" {
		SendJSONError(w, "Отсутствуют необходимые параметры", http.StatusBadRequest)
		return
	}
	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		SendJSONError(w, "Неверный формат суммы", http.StatusBadRequest)
		return
	}
	response, err := h.exchangeService.CalculateExchange(from, to, amount)
	if err != nil {
		SendJSONError(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
