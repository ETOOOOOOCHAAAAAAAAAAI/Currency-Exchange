package handler

import (
	"Currency-exchange/internal/models"
	"Currency-exchange/internal/repository"
	"encoding/json"
	"net/http"
)

type CurrencyHandler struct {
	repo *repository.CurrencyRepository
}

func NewCurrencyHandler(repo *repository.CurrencyRepository) *CurrencyHandler {
	return &CurrencyHandler{repo: repo}
}

func (h *CurrencyHandler) GetAllCurrencies(w http.ResponseWriter, r *http.Request) {
	currencies, err := h.repo.GetAll()
	if err != nil {
		http.Error(w, "Ошибка при получении валют из БД", http.StatusInternalServerError)
		return
	}
	jsonBytes, err := json.Marshal(currencies)
	if err != nil {
		http.Error(w, "Ошибка при формировании ответа", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

func (h *CurrencyHandler) GetCurrencyByCode(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("code")
	codeCurrency, err := h.repo.GetByCode(code)
	if err != nil {
		http.Error(w, "Валюта не найдена", http.StatusNotFound)
		return
	}
	jsonBytes, err := json.Marshal(codeCurrency)
	if err != nil {
		http.Error(w, "Ошибка при формировании ответа", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)

}

func (h *CurrencyHandler) CreateNewCurrency(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Ошибка при чтений данных", http.StatusInternalServerError)
		return
	}
	newCode := r.FormValue("code")
	newName := r.FormValue("name")
	newSign := r.FormValue("sign")
	if newCode == "" || newName == "" || newSign == "" {
		http.Error(w, "Неверный запрос", http.StatusBadRequest)
		return
	}
	c := models.Currency{Code: newCode, FullName: newName, Sign: newSign}
	result, err := h.repo.CreateCurrencies(c)
	if err != nil {
		http.Error(w, "Такая валюта уже существует", http.StatusConflict)
		return
	}
	jsonBytes, err := json.Marshal(result)
	if err != nil {
		http.Error(w, "Ошибка при формировании ответа", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonBytes)
}
