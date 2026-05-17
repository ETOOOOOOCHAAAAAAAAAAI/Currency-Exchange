INSERT INTO Currencies (code, full_name, sign) VALUES ('USD', 'United States dollar', '$');
INSERT INTO Currencies (code, full_name, sign) VALUES ('EUR', 'Euro', '€');
INSERT INTO Currencies (code, full_name, sign) VALUES ('RUB', 'Russian Ruble', '₽');

INSERT INTO ExchangeRates (base_currency_id, target_currency_id, rate) VALUES (1, 2, 0.95);
INSERT INTO ExchangeRates (base_currency_id, target_currency_id, rate) VALUES (1, 3, 90.50);
INSERT INTO ExchangeRates (base_currency_id, target_currency_id, rate) VALUES (2, 3, 95.00);