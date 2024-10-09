package service

import (
	"database/sql"
	"fmt"
	"test-golang/internal/helpers"
	"test-golang/internal/models"
)

type UserService struct {
	db *sql.DB
}

func NewUserService(db *sql.DB) *UserService {
	return &UserService{
		db: db,
	}
}

func (s *UserService) GetByID(id int) (models.User, error) {
	if s.db == nil {
		return models.User{}, fmt.Errorf("nil pointer: UserService.db")
	}

	var user models.User
	var createdAt, updatedAt []uint8

	err := s.db.QueryRow(`SELECT id, name, age, address, created_at, updated_at FROM users WHERE id = ?`, id).Scan(&user.ID, &user.Name, &user.Age, &user.Address, &createdAt, &updatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return models.User{}, fmt.Errorf("user with ID %d not found", id)
		}
		return models.User{}, fmt.Errorf("database error: %w", err)
	}

	user.CreatedAt, err = helpers.ParseTimestamp(createdAt)
	if err != nil {
		return models.User{}, fmt.Errorf("unable to parse timestamp: %w", err)
	}

	user.UpdatedAt, err = helpers.ParseTimestamp(updatedAt)
	if err != nil {
		return models.User{}, fmt.Errorf("unable to parse timestamp: %w", err)
	}

	return user, nil
}

func (s *UserService) GetAll() ([]models.User, error) {
	rows, err := s.db.Query(`SELECT id, name, age, address, created_at, updated_at FROM users ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		var createdAt, updatedAt []uint8

		err := rows.Scan(&user.ID, &user.Name, &user.Age, &user.Address, &createdAt, &updatedAt)
		if err != nil {
			return nil, err
		}

		user.CreatedAt, err = helpers.ParseTimestamp(createdAt)
		if err != nil {
			return nil, fmt.Errorf("unable to parse timestamp: %w", err)
		}

		user.UpdatedAt, err = helpers.ParseTimestamp(updatedAt)
		if err != nil {
			return nil, fmt.Errorf("unable to parse timestamp: %w", err)
		}

		users = append(users, user)

	}

	return users, nil
}

func (s *UserService) Create(user models.User) (models.User, error) {
	if s.db == nil {
		return models.User{}, fmt.Errorf("nil pointer: UserService.db")
	}

	result, err := s.db.Exec(`INSERT INTO users (name, age, address) VALUES (?, ?, ?)`, user.Name, user.Age, user.Address)
	if err != nil {
		return models.User{}, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return models.User{}, err
	}

	return s.GetByID(int(id))
}

func (s *UserService) Update(user models.User) (models.User, error) {
	if s.db == nil {
		return models.User{}, fmt.Errorf("nil pointer: UserService.db")
	}

	result, err := s.db.Exec(`UPDATE users SET name = ?, age = ?, address = ? WHERE id = ?`, user.Name, user.Age, user.Address, user.ID)
	if err != nil {
		return models.User{}, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return models.User{}, err
	}
	if rowsAffected == 0 {
		return models.User{}, fmt.Errorf("user with ID %d not found", user.ID)
	}

	return s.GetByID(user.ID)
}

func (s *UserService) Delete(id int) error {
	if s.db == nil {
		return fmt.Errorf("nil pointer: UserService.db")
	}

	result, err := s.db.Exec(`DELETE FROM users WHERE id = ?`, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user with ID %d not found", id)
	}

	return nil
}
