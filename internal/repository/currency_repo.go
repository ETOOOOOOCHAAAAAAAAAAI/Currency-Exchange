package repository

import (
	"Currency-exchange/internal/models"
	"database/sql"
	"fmt"
)

type CurrencyRepository struct {
	db *sql.DB
}

func NewCurrencyRepository(db *sql.DB) *CurrencyRepository {
	return &CurrencyRepository{db: db}
}

func (r *CurrencyRepository) GetAll() ([]models.Currency, error) {
	rows, err := r.db.Query("SELECT id, code, full_name, sign FROM Currencies")
	if err != nil {
		return nil, fmt.Errorf("ошибка при запросе валют: %w", err)
	}
	defer rows.Close()
	var currencies []models.Currency
	for rows.Next() {
		var c models.Currency
		err := rows.Scan(&c.ID, &c.Code, &c.FullName, &c.Sign)
		if err != nil {
			return nil, fmt.Errorf("ошибка при запросе валют: %w", err)
		}
		currencies = append(currencies, c)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при чтении строк: %w", err)
	}
	return currencies, nil
}

func (r *CurrencyRepository) GetByCode(code string) (models.Currency, error) {
	var c models.Currency
	err := r.db.QueryRow("SELECT id, code, full_name, sign FROM Currencies WHERE code = ?", code).Scan(&c.ID, &c.Code, &c.FullName, &c.Sign)
	if err != nil {
		return c, fmt.Errorf("Ошибка при запросе кода валют: %w", err)
	}
	return c, nil
}

func (r *CurrencyRepository) CreateCurrencies(c models.Currency) (models.Currency, error) {
	result, err := r.db.Exec("INSERT INTO Currencies (code, full_name, sign) VALUES (?,?,?)", c.Code, c.FullName, c.Sign)
	if err != nil {
		return c, fmt.Errorf("Такая валюта уже существует: %w", err)

	}
	id, err := result.LastInsertId()
	if err != nil {
		return c, fmt.Errorf("Ошибка при получений id: %w", err)

	}
	c.ID = id
	return c, nil
}
