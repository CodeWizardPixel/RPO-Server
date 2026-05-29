package repository

import (
	"database/sql"
	"fmt"
)

type UserRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{DB: db}
}

type User struct {
	ID           int
	Login        string
	Name         string
	PasswordHash string
	IsAdmin      int
}

func (r *UserRepository) GetAllUsers() ([]User, error) {
	rows, err := r.DB.Query("select id, login, name, password_hash, is_admin from users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		err := rows.Scan(&u.ID, &u.Login, &u.Name, &u.PasswordHash, &u.IsAdmin)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, nil
}

func (r *UserRepository) GetUserByID(id int) (*User, error) {
	row := r.DB.QueryRow("select id, login, name, password_hash, is_admin from users where id = ?", id)
	var u User
	err := row.Scan(&u.ID, &u.Login, &u.Name, &u.PasswordHash, &u.IsAdmin)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) GetUserByLogin(login string) (*User, error) {
	row := r.DB.QueryRow("select id, login, name, password_hash, is_admin from users where login = ?", login)
	var u User
	err := row.Scan(&u.ID, &u.Login, &u.Name, &u.PasswordHash, &u.IsAdmin)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}
	return &u, nil
}

func (r *UserRepository) CreateUser(login, name, passwordHash string, isAdmin int) error {
	_, err := r.DB.Exec("insert into users (login, name, password_hash, is_admin) values (?, ?, ?, ?)", login, name, passwordHash, isAdmin)
	return err
}

func (r *UserRepository) UpdateUser(id int, name, passwordHash string, isAdmin int) error {
	_, err := r.DB.Exec("update users set name = ?, password_hash = ?, is_admin = ? where id = ?", name, passwordHash, isAdmin, id)
	return err
}

func (r *UserRepository) DeleteUserByID(id int) error {
	_, err := r.DB.Exec("delete from users where id = ?", id)
	return err
}