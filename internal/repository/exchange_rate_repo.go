package repository

import (
	"Currency-exchange/internal/models"
	"database/sql"
	"fmt"
)

type ExchangeRateRepository struct {
	db *sql.DB
}

func NewExchangeRateRepository(db *sql.DB) *ExchangeRateRepository {
	return &ExchangeRateRepository{db: db}
}

func (r *ExchangeRateRepository) GetAll() ([]models.ExcangeRate, error) {
	rows, err := r.db.Query(`SELECT
		    er.id, er.rate,
		    bc.id, bc.code, bc.full_name, bc.sign,
		    tc.id, tc.code, tc.full_name, tc.sign
	    FROM ExchangeRates er
		JOIN Currencies bc ON er.base_currency_id = bc.id
		JOIN Currencies tc ON er.target_currency_id = tc.id`)
	if err != nil {
		return nil, fmt.Errorf("Ошибка при запросе обменника: %w", err)
	}
	defer rows.Close()
	var exchangeRate []models.ExcangeRate
	for rows.Next() {
		var rate models.ExcangeRate
		err := rows.Scan(
			&rate.ID,
			&rate.Rate,
			&rate.BaseCurrency.ID,
			&rate.BaseCurrency.Code,
			&rate.BaseCurrency.FullName,
			&rate.BaseCurrency.Sign,
			&rate.TargetCurrency.ID,
			&rate.TargetCurrency.Code,
			&rate.TargetCurrency.FullName,
			&rate.TargetCurrency.Sign,
		)
		if err != nil {
			return nil, fmt.Errorf("Ошибка при запросе обменника: %w", err)
		}
		exchangeRate = append(exchangeRate, rate)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("Ошибка при чтений строк: %w", err)
	}
	return exchangeRate, nil
}

func (r *ExchangeRateRepository) GetRateByCode(baseCode, targetCode string) (models.ExcangeRate, error) {
	var rate models.ExcangeRate
	err := r.db.QueryRow(`SELECT
		    er.id, er.rate,
		    bc.id, bc.code, bc.full_name, bc.sign,
		    tc.id, tc.code, tc.full_name, tc.sign
	    FROM ExchangeRates er
		JOIN Currencies bc ON er.base_currency_id = bc.id
		JOIN Currencies tc ON er.target_currency_id = tc.id
		WHERE bc.code = ? and tc.code = ?`, baseCode, targetCode).Scan(&rate.ID,
		&rate.Rate,
		&rate.BaseCurrency.ID,
		&rate.BaseCurrency.Code,
		&rate.BaseCurrency.FullName,
		&rate.BaseCurrency.Sign,
		&rate.TargetCurrency.ID,
		&rate.TargetCurrency.Code,
		&rate.TargetCurrency.FullName,
		&rate.TargetCurrency.Sign,
	)
	if err != nil {
		return rate, fmt.Errorf("Ошибка при получении данных из БД: %w", err)
	}
	return rate, nil
}
