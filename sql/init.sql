CREATE TABLE IF NOT EXISTS Currencies (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    code VARCHAR UNIQUE,
    full_name VARCHAR,
    sign VARCHAR
);

CREATE TABLE IF NOT EXISTS Excange_Rates (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    base_currency_id INTEGER,
    target_currency_id INTEGER,
    rate DECIMAL(6),
    UNIQUE (base_currency_id, target_currency_id),
    FOREIGN KEY (base_currency_id) REFERENCES Currencies(id),
    FOREIGN KEY (target_currency_id) REFERENCES Currencies(id)
);