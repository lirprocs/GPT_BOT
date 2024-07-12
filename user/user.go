package model

import (
	"GPT_BOT/lama"
	"github.com/go-telegram/bot/models"
)

type User struct {
	Id      int64
	Name    string
	Balance int
	NumLama int
	NumGPT  int
	APIKey  string
	ModelID string
	storage Storage
}

type Storage interface {
	AddUserDB(user *User) error
	GetUserDB(update *models.Update) (User, error)
	ChangeInfo(id int64, fieldName, data string) error
}

func NewUser(storage Storage) *User {
	return &User{
		storage: storage,
	}
}

func (u *User) AddUser(update *models.Update) error {
	user := &User{
		Id: update.Message.From.ID, Name: update.Message.From.Username, APIKey: lama.GROQAPIKey, ModelID: lama.ModelID,
	}
	err := u.storage.AddUserDB(user)
	if err != nil {
		return err
	}
	return nil
}

func (u *User) GetUser(update *models.Update) (User, error) {
	data, err := u.storage.GetUserDB(update)
	if err != nil {
		return User{}, err
	}
	return data, nil
}

func (u *User) UpdateUser(update *models.Update, infoName, data string) error {
	id := update.Message.From.ID
	switch infoName {
	case "/api":
		err := u.storage.ChangeInfo(id, "apiKey", data)
		if err != nil {
			return err
		}
	case "/model":
		err := u.storage.ChangeInfo(id, "modelID", data)
		if err != nil {
			return err
		}
	}
	return nil
}
