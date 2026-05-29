package repository

import (
	"database/sql"
	"fmt"
	"time"
)

type TransactionRepository struct {
	DB *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{DB: db}
}

type Transaction struct {
	ID         int
	Amount     float64
	CardID     int
	TerminalID int
	Operation  string
	CreatedAt  time.Time
}

func (r *TransactionRepository) GetAllTransactions() ([]Transaction, error) {
	rows, err := r.DB.Query("select id, amount, card_id, terminal_id, operation, created_at from transactions")
	if err != nil {
		return nil, fmt.Errorf("error querying transactions: %w", err)
	}
	defer rows.Close()

	var transactions []Transaction
	for rows.Next() {
		var t Transaction
		err := rows.Scan(&t.ID, &t.Amount, &t.CardID, &t.TerminalID, &t.Operation, &t.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("error scanning transaction row: %w", err)
		}
		transactions = append(transactions, t)
	}

	return transactions, nil
}

func (r *TransactionRepository) GetTransactionByID(id int) (*Transaction, error) {
	row := r.DB.QueryRow("select id, amount, card_id, terminal_id, operation, created_at from transactions where id = ?", id)
	var t Transaction
	err := row.Scan(&t.ID, &t.Amount, &t.CardID, &t.TerminalID, &t.Operation, &t.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("error scanning transaction row: %w", err)
	}
	return &t, nil
}

func (r *TransactionRepository) CreateTransaction(amount float64, cardID, terminalID int, operation string) error {
	if operation == "" {
		operation = "withdraw"
	}

	_, err := r.DB.Exec("insert into transactions (amount, card_id, terminal_id, operation) values (?, ?, ?, ?)",
		amount, cardID, terminalID, operation)
	if err != nil {
		return fmt.Errorf("error creating transaction: %w", err)
	}
	return nil
}

func (r *TransactionRepository) ProcessCardOperation(cardNumber, terminalSerialNumber string, amount float64, operation string) (*Card, error) {
	tx, err := r.DB.Begin()
	if err != nil {
		return nil, fmt.Errorf("error starting transaction: %w", err)
	}
	defer tx.Rollback()

	var card Card
	err = tx.QueryRow(
		"select id, card_number, balance, is_blocked, owner_name, key_id from cards where card_number = ?",
		cardNumber,
	).Scan(&card.ID, &card.CardNumber, &card.Balance, &card.IsBlocked, &card.OwnerName, &card.KeyID)
	if err != nil {
		return nil, fmt.Errorf("error finding card by number: %w", err)
	}

	var terminalID int
	err = tx.QueryRow("select id from terminals where serial_number = ?", terminalSerialNumber).Scan(&terminalID)
	if err != nil {
		return nil, fmt.Errorf("error finding terminal by serial number: %w", err)
	}

	if card.IsBlocked == 1 {
		return nil, fmt.Errorf("card is blocked")
	}

	switch operation {
	case "withdraw":
		if card.Balance < amount {
			return nil, fmt.Errorf("insufficient funds. required: %.2f, available: %.2f", amount, card.Balance)
		}
		card.Balance -= amount
	case "deposit":
		card.Balance += amount
	default:
		return nil, fmt.Errorf("invalid operation: must be withdraw or deposit")
	}

	_, err = tx.Exec("update cards set balance = ? where id = ?", card.Balance, card.ID)
	if err != nil {
		return nil, fmt.Errorf("error updating card balance: %w", err)
	}

	_, err = tx.Exec(
		"insert into transactions (amount, card_id, terminal_id, operation) values (?, ?, ?, ?)",
		amount,
		card.ID,
		terminalID,
		operation,
	)
	if err != nil {
		return nil, fmt.Errorf("error creating transaction: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("error committing transaction: %w", err)
	}

	return &card, nil
}

func (r *TransactionRepository) DeleteTransactionByID(id int) error {
	_, err := r.DB.Exec("delete from transactions where id = ?", id)
	if err != nil {
		return fmt.Errorf("error deleting transaction: %w", err)
	}
	return nil
}
