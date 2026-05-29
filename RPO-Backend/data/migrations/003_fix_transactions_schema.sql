-- +goose Up

PRAGMA foreign_keys = OFF;

CREATE TABLE transactions_new (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    amount REAL NOT NULL,
    card_id INTEGER NOT NULL,
    terminal_id INTEGER NOT NULL,
    operation TEXT NOT NULL DEFAULT 'withdraw',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (card_id) REFERENCES cards(id),
    FOREIGN KEY (terminal_id) REFERENCES terminals(id)
);

INSERT INTO transactions_new (id, amount, card_id, terminal_id, operation, created_at)
SELECT id, amount, card_id, terminal_id, 'withdraw', created_at
FROM transactions;

DROP TABLE transactions;
ALTER TABLE transactions_new RENAME TO transactions;

PRAGMA foreign_keys = ON;

-- +goose Down

PRAGMA foreign_keys = OFF;

CREATE TABLE transactions_old (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    amount REAL NOT NULL,
    card_id INTEGER NOT NULL,
    terminal_id INTEGER NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (card_id) REFERENCES transport_cards(id),
    FOREIGN KEY (terminal_id) REFERENCES terminals(id)
);

INSERT INTO transactions_old (id, amount, card_id, terminal_id, created_at)
SELECT id, amount, card_id, terminal_id, created_at
FROM transactions;

DROP TABLE transactions;
ALTER TABLE transactions_old RENAME TO transactions;

PRAGMA foreign_keys = ON;
