package handler

import (
	"Currency-exchange/internal/repository"
	"encoding/json"
	"net/http"
)

type ExchangeRateHandler struct {
	repo *repository.ExchangeRateRepository
}

func NewExchangeRateHandler(repo *repository.ExchangeRateRepository) *ExchangeRateHandler {
	return &ExchangeRateHandler{repo: repo}
}

func (h *ExchangeRateHandler) GetAllExchangeRate(w http.ResponseWriter, r *http.Request) {
	rate, err := h.repo.GetAll()
	if err != nil {
		http.Error(w, "Ошибка при получений данных из обменника БД", http.StatusInternalServerError)
		return
	}
	jsonBytes, err := json.Marshal(rate)
	if err != nil {
		http.Error(w, "Ошибка при формираваний ответа", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)

}

func (h *ExchangeRateHandler) GetByCodeInExchangeRate(w http.ResponseWriter, r *http.Request) {
	codes := r.PathValue("codes")
	if len(codes) != 6 {
		http.Error(w, "Неверный формат кода валют", http.StatusBadRequest)
		return
	}
	baseCode := codes[:3]
	targetCode := codes[3:]
	if baseCode == targetCode {
		http.Error(w, "Одинаковый код валют", http.StatusBadRequest)
		return
	}
	rate, err := h.repo.GetRateByCode(baseCode, targetCode)
	if err != nil {
		http.Error(w, "Код валют не найден", http.StatusNotFound)
		return
	}
	jsonBytes, err := json.Marshal(rate)
	if err != nil {
		http.Error(w, "Ошибка при формировании ответа", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}
