package handler

import (
	"Currency-exchange/internal/models"
	"Currency-exchange/internal/service"
	"encoding/json"
	"net/http"
)

type CurrencyHandler struct {
	service *service.CurrencyService
}

func NewCurrencyHandler(service *service.CurrencyService) *CurrencyHandler {
	return &CurrencyHandler{service: service}
}

func (h *CurrencyHandler) GetAllCurrencies(w http.ResponseWriter, r *http.Request) {
	currencies, err := h.service.GetAllCurrencies()
	if err != nil {
		SendJSONError(w, "Ошибка при получении валют из БД", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(currencies)
}

func (h *CurrencyHandler) GetCurrencyByCode(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("code")
	codeCurrency, err := h.service.GetCurrencyByCode(code)
	if err != nil {
		SendJSONError(w, "Валюта не найдена", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(codeCurrency)
}

func (h *CurrencyHandler) CreateNewCurrency(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		SendJSONError(w, "Ошибка при чтений данных", http.StatusInternalServerError)
		return
	}
	newCode, newName, newSign := r.FormValue("code"), r.FormValue("name"), r.FormValue("sign")
	if newCode == "" || newName == "" || newSign == "" {
		SendJSONError(w, "Неверный запрос", http.StatusBadRequest)
		return
	}
	c := models.Currency{Code: newCode, FullName: newName, Sign: newSign}
	result, err := h.service.CreateCurrency(c)
	if err != nil {
		SendJSONError(w, "Такая валюта уже существует", http.StatusConflict)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(result)
}
