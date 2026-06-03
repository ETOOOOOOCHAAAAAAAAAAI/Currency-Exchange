package main

import (
	"Currency-exchange/internal/handler"
	"Currency-exchange/internal/middleware"
	"Currency-exchange/internal/repository"
	"Currency-exchange/internal/service"
	"log"
	"net/http"
)

func main() {
	db, err := repository.NewSQLiteDB("./sqllite.db")
	if err != nil {
		log.Fatalf("Ошибка подключений к БД: %v", err)
	}
	defer db.Close()
	currencyRepo := repository.NewCurrencyRepository(db)
	exchangeRateRepo := repository.NewExchangeRateRepository(db)
	currencyService := service.NewCurrencyService(currencyRepo)
	exchangeService := service.NewExchangeService(exchangeRateRepo, currencyRepo)
	currencyHandler := handler.NewCurrencyHandler(currencyService)
	exchangeRateHandler := handler.NewExchangeRateHandler(exchangeService, currencyService)
	mux := http.NewServeMux()
	mux.HandleFunc("GET /currencies", currencyHandler.GetAllCurrencies)
	mux.HandleFunc("GET /currency/{code}", currencyHandler.GetCurrencyByCode)
	mux.HandleFunc("POST /currencies", currencyHandler.CreateNewCurrency)
	mux.HandleFunc("GET /exchangeRates", exchangeRateHandler.GetAllExchangeRate)
	mux.HandleFunc("GET /exchangeRate/{codes}", exchangeRateHandler.GetByCodeInExchangeRate)
	mux.HandleFunc("POST /exchangeRates", exchangeRateHandler.CreateNewExchangeRate)
	mux.HandleFunc("PATCH /exchangeRate/{codes}", exchangeRateHandler.UpdateExchangeRate)
	mux.HandleFunc("GET /exchange", exchangeRateHandler.ExchangeCalculator)
	protectedMux := middleware.CORS(mux)
	log.Println("Сервер запущен на порту 8080...")
	err = http.ListenAndServe(":8080", protectedMux)
	if err != nil {
		log.Fatalf("Ошибка сервера: %v", err)
	}
}
