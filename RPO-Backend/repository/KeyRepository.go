package repository

import (
	"database/sql"
	"fmt"
)

type KeyRepository struct {
	DB *sql.DB
}

func NewKeyRepository(db *sql.DB) *KeyRepository {
	return &KeyRepository{DB: db}
}

type Key struct {
	ID    int
	Value string
}

func (r *KeyRepository) GetAllKeys() ([]Key, error) {
	rows, err := r.DB.Query("select * from keys")
	if err != nil {
		return nil, fmt.Errorf("error querying keys: %w", err)
	}
	defer rows.Close()

	var keys []Key
	for rows.Next() {
		var k Key
		err := rows.Scan(&k.ID, &k.Value)
		if err != nil {
			return nil, fmt.Errorf("error scanning key row: %w", err)
		}
		keys = append(keys, k)
	}

	return keys, nil
}

func (r *KeyRepository) GetKeyByID(id int) (*Key, error) {
	row := r.DB.QueryRow("select * from keys where id = ?", id)
	var k Key
	err := row.Scan(&k.ID, &k.Value)
	if err != nil {
		return nil, fmt.Errorf("error scanning key row: %w", err)
	}
	return &k, nil
}

func (r *KeyRepository) CreateKey(value string) error {
	_, err := r.DB.Exec("insert into keys (value) values (?)", value)
	if err != nil {
		return fmt.Errorf("error creating key: %w", err)
	}
	return nil
}

func (r *KeyRepository) UpdateKey(id int, value string) error {
	_, err := r.DB.Exec("update keys set value = ? where id = ?", value, id)
	if err != nil {
		return fmt.Errorf("error updating key: %w", err)
	}
	return nil
}

func (r *KeyRepository) DeleteKeyByID(id int) error {
	_, err := r.DB.Exec("delete from keys where id = ?", id)
	if err != nil {
		return fmt.Errorf("error deleting key: %w", err)
	}
	return nil
}
