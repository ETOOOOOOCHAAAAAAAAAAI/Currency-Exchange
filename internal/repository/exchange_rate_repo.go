package repository

import (
	"Currency-exchange/internal/models"
	"database/sql"
	"errors"
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

// создать POST для exchangeRate потом идти как обычно по плану repo -> handler ну и т.д
func (r *ExchangeRateRepository) CreateNewExchangeRate(e models.ExcangeRate) (models.ExcangeRate, error) {
	result, err := r.db.Exec(`INSERT INTO ExchangeRates (base_currency_id, target_currency_id, rate) 
VALUES (?,?,?)`, e.BaseCurrency.ID, e.TargetCurrency.ID, e.Rate)
	if err != nil {
		return e, fmt.Errorf("Такая соотношения уже существует: %w", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return e, fmt.Errorf("Ошибка при получений id: %w", err)
	}
	e.ID = id
	return e, nil
}

// сделать обновление существующего курса PATCH
func (r *ExchangeRateRepository) UpdateExchangeRate(e models.ExcangeRate) (models.ExcangeRate, error) {
	result, err := r.db.Exec(`UPDATE ExchangeRates 
		SET rate = ? 
		WHERE base_currency_id = ? AND target_currency_id = ?`, e.Rate, e.BaseCurrency, e.TargetCurrency)

	if err != nil {
		return e, fmt.Errorf("Таких валют не существует: %w", err)
	}
	count, err := result.RowsAffected()
	if err != nil {
		return e, fmt.Errorf("Ошибка при получении данных: %w", err)
	}
	if count == 0 {
		return e, errors.New("курс для этой пары не найден")
	}
	return e, nil
}

// написать самое важное - калькулятор
