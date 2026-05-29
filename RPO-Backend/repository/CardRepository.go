package repository

import (
	"database/sql"
	"fmt"
)

type CardRepository struct {
	DB *sql.DB
}

func NewCardRepository(db *sql.DB) *CardRepository {
	return &CardRepository{DB: db}
}

type Card struct {
	ID         int
	CardNumber string
	Balance    float64
	IsBlocked  int
	OwnerName  string
	KeyID      *int
}

func (r *CardRepository) GetAllCards() ([]Card, error) {
	rows, err := r.DB.Query("select * from cards")
	if err != nil {
		return nil, fmt.Errorf("error querying  cards: %w", err)
	}
	defer rows.Close()

	var cards []Card
	for rows.Next() {
		var tc Card
		err := rows.Scan(&tc.ID, &tc.CardNumber, &tc.Balance, &tc.IsBlocked, &tc.OwnerName, &tc.KeyID)
		if err != nil {
			return nil, fmt.Errorf("error scanning  card row: %w", err)
		}
		cards = append(cards, tc)
	}

	return cards, nil
}

func (r *CardRepository) GetCardByID(id int) (*Card, error) {
	row := r.DB.QueryRow("select * from cards where id = ?", id)
	var tc Card
	err := row.Scan(&tc.ID, &tc.CardNumber, &tc.Balance, &tc.IsBlocked, &tc.OwnerName, &tc.KeyID)
	if err != nil {
		return nil, fmt.Errorf("error scanning  card row: %w", err)
	}
	return &tc, nil
}

func (r *CardRepository) CreateCard(cardNumber string, balance float64, isBlocked int, ownerName string, keyID *int) error {
	_, err := r.DB.Exec("insert into cards (card_number, balance, is_blocked, owner_name, key_id) values (?, ?, ?, ?, ?)",
		cardNumber, balance, isBlocked, ownerName, keyID)
	if err != nil {
		return fmt.Errorf("error creating  card: %w", err)
	}
	return nil
}

func (r *CardRepository) UpdateCard(id int, balance float64, isBlocked int, ownerName string, keyID *int) error {
	_, err := r.DB.Exec("update cards set balance = ?, is_blocked = ?, owner_name = ?, key_id = ? where id = ?",
		balance, isBlocked, ownerName, keyID, id)
	if err != nil {
		return fmt.Errorf("error updating  card: %w", err)
	}
	return nil
}

func (r *CardRepository) UpdateCardBalance(id int, balance float64) error {
	_, err := r.DB.Exec("update cards set balance = ? where id = ?", balance, id)
	if err != nil {
		return fmt.Errorf("error updating  card balance: %w", err)
	}
	return nil
}

func (r *CardRepository) DeleteCardByID(id int) error {
	_, err := r.DB.Exec("delete from cards where id = ?", id)
	if err != nil {
		return fmt.Errorf("error deleting  card: %w", err)
	}
	return nil
}
