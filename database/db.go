package database

import (
	"GPT_BOT/user"
	"database/sql"
	"errors"
	"fmt"
	"github.com/go-telegram/bot/models"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

type Storage struct {
	db *sql.DB
}

func InitDB(storagePath string) (*Storage, error) {
	var err error
	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, err
	}
	return &Storage{db: db}, err
}

func (s *Storage) Close() error {
	return s.db.Close()
}
func (s *Storage) AddUserDB(user *model.User) error {
	insertUserSQL := `
	INSERT OR IGNORE INTO users (id, name, apiKey, modelID) 
	VALUES (?, ?, ?, ?);`
	statement, err := s.db.Prepare(insertUserSQL)
	if err != nil {
		return err
	}
	defer statement.Close()

	_, err = statement.Exec(user.Id, user.Name, user.APIKey, user.ModelID)
	if err != nil {
		return err
	}

	log.Println("User added successfully!")
	return nil
}

func (s *Storage) GetUserDB(update *models.Update) (model.User, error) {
	getUserSQL := `
        SELECT *
        FROM users
        WHERE id = ?;`

	var user model.User
	err := s.db.QueryRow(getUserSQL, update.Message.From.ID).Scan(&user.Id, &user.Name, &user.Balance, &user.NumLama, &user.NumGPT, &user.APIKey, &user.ModelID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.User{}, fmt.Errorf("user not found for ID %d", update.Message.From.ID)
		}
		return model.User{}, err
	}

	return user, nil
}

func (s *Storage) ChangeInfo(id int64, fieldName, data string) error {
	changeModelSQL := fmt.Sprintf("UPDATE users SET %s = ? WHERE id = ?", fieldName)
	statement, err := s.db.Prepare(changeModelSQL)
	if err != nil {
		return err
	}
	defer statement.Close()

	_, err = statement.Exec(data, id)
	if err != nil {
		return err
	}

	log.Println(fmt.Sprintf("%s changed successfully!", fieldName))
	return nil
}
