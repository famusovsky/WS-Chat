package postgres

import (
	"database/sql"
	"time"
	"ws-server/internal/models"
)

type MessagesDBProcessor interface {
	AddMessage(sender, text string, creation time.Time) error
	GetMessage(id int) (models.Message, error)
	GetLastMessages(count int) ([]models.Message, error)
	GetLastMessagesByNickname(count int, nickname string) ([]models.Message, error)
}

func Get(db *sql.DB, overrideTables bool) (MessagesDBProcessor, error) {
	if overrideTables {
		err := overrideDB(db)
		if err != nil {
			return nil, err
		}
	}

	err := checkDB(db)
	if err != nil {
		return nil, err
	}

	return &messagesDB{db}, nil
}
