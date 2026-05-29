package repository

import (
	"database/sql"
	"fmt"
)

type TerminalRepository struct {
	DB *sql.DB
}

func NewTerminalRepository(db *sql.DB) *TerminalRepository {
	return &TerminalRepository{DB: db}
}

type Terminal struct {
	ID           int
	SerialNumber string
	Address      string
	Name         string
}

func (r *TerminalRepository) GetAllTerminals() ([]Terminal, error) {
	rows, err := r.DB.Query("select * from terminals")
	if err != nil {
		return nil, fmt.Errorf("error querying terminals: %w", err)
	}
	defer rows.Close()

	var terminals []Terminal
	for rows.Next() {
		var t Terminal
		err := rows.Scan(&t.ID, &t.SerialNumber, &t.Address, &t.Name)
		if err != nil {
			return nil, fmt.Errorf("error scanning terminal row: %w", err)
		}
		terminals = append(terminals, t)
	}

	return terminals, nil
}

func (r *TerminalRepository) GetTerminalByID(id int) (*Terminal, error) {
	row := r.DB.QueryRow("select * from terminals where id = ?", id)
	var t Terminal
	err := row.Scan(&t.ID, &t.SerialNumber, &t.Address, &t.Name)
	if err != nil {
		return nil, fmt.Errorf("error scanning terminal row: %w", err)
	}
	return &t, nil
}

func (r *TerminalRepository) CreateTerminal(serialNumber, address, name string) error {
	_, err := r.DB.Exec("insert into terminals (serial_number, address, name) values (?, ?, ?)", serialNumber, address, name)
	if err != nil {
		return fmt.Errorf("error creating terminal: %w", err)
	}
	return nil
}

func (r *TerminalRepository) UpdateTerminal(id int, serialNumber, address, name string) error {
	_, err := r.DB.Exec("update terminals set serial_number = ?, address = ?, name = ? where id = ?", serialNumber, address, name, id)
	if err != nil {
		return fmt.Errorf("error updating terminal: %w", err)
	}
	return nil
}

func (r *TerminalRepository) DeleteTerminalByID(id int) error {
	_, err := r.DB.Exec("delete from terminals where id = ?", id)
	if err != nil {
		return fmt.Errorf("error deleting terminal: %w", err)
	}
	return nil
}
