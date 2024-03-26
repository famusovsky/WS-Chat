package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
	"ws-server/internal/models"
)

type messagesDB struct {
	db *sql.DB
}

func (m *messagesDB) AddMessage(sender string, text string, creation time.Time) error {
	tx, err := m.db.Begin()
	if err != nil {
		return fmt.Errorf("error while starting transaction: %s", err.Error())
	}
	defer tx.Rollback()

	q := `INSERT INTO messages (text, sender, creation) VALUES ($1, $2, $3);`
	if err = tx.QueryRow(q, text, sender, creation).Err(); err != nil {
		return fmt.Errorf("error while adding message to the database: %s", err.Error())
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("error while committing transaction: %s", err.Error())
	}

	return nil
}

func (m *messagesDB) GetLastMessages(count int) ([]models.Message, error) {
	wrapErr := errors.New("error while getting messages from the database")

	q := `SELECT * FROM (
		SELECT * from messages ORDER BY creation DESC LIMIT $1
	) AS tmp ORDER BY creation ASC;`

	rows, err := m.db.Query(q, count)
	if err != nil {
		return nil, errors.Join(wrapErr, err)
	}

	res := []models.Message{}
	for rows.Next() {
		var message models.Message
		if err = rows.Scan(&message.Id, &message.Text, &message.Sender, &message.Creation); err != nil {
			return nil, errors.Join(wrapErr, err)
		}
		res = append(res, message)
	}

	return res, nil
}

func (m *messagesDB) GetLastMessagesByNickname(count int, nickname string) ([]models.Message, error) {
	wrapErr := errors.New("error while getting messages from the database")

	q := `SELECT * FROM (
		SELECT * from messages WHERE sender = $1 ORDER BY creation DESC LIMIT $2
	) AS tmp ORDER BY creation ASC;`

	rows, err := m.db.Query(q, nickname, count)
	if err != nil {
		return nil, errors.Join(wrapErr, err)
	}

	res := []models.Message{}
	for rows.Next() {
		var message models.Message
		if err = rows.Scan(&message); err != nil {
			return nil, errors.Join(wrapErr, err)
		}
		res = append(res, message)
	}

	return res, nil
}

func (m *messagesDB) GetMessage(id int) (models.Message, error) {
	wrapErr := errors.New("error while getting message from the database")

	q := `SELECT * from messages WHERE id = $1;`
	rows, err := m.db.Query(q, id)
	if err != nil {
		return models.Message{}, errors.Join(wrapErr, err)
	}

	var message models.Message
	if err = rows.Scan(&message); err != nil {
		return models.Message{}, errors.Join(wrapErr, err)
	}

	return message, nil
}

func overrideDB(db *sql.DB) error {
	q := `DROP TABLE IF EXISTS messages;`

	if _, err := db.Exec(q); err != nil {
		return fmt.Errorf("error while dropping tables: %s", err)
	}
	q = `CREATE TABLE messages (
		id SERIAL PRIMARY KEY,
		text TEXT,
		sender TEXT,
		creation TIMESTAMP
	);`

	if _, err := db.Exec(q); err != nil {
		return fmt.Errorf("error while creating tables: %s", err)
	}

	return nil
}

func checkDB(db *sql.DB) error {
	var (
		qMessages = `SELECT COUNT(*) = 4 AS properMessages
		FROM information_schema.columns
		WHERE table_schema = 'public'
		AND table_name = 'messages'
		AND (
			(column_name = 'id' AND data_type = 'integer')
			OR (column_name = 'text' AND data_type = 'text')
			OR (column_name = 'sender' AND data_type = 'text')
			OR (column_name = 'creation' AND data_type = 'timestamp without time zone')	
		);`
		properMessages bool
	)

	if err := db.QueryRow(qMessages).Scan(&properMessages); err != nil {
		return fmt.Errorf("error while checking 'messages' table: %s", err.Error())
	}

	if !properMessages {
		return errors.New("'messages' table is not ok: proper 'messages' table is { id INTEGER; text TEXT; sender TEXT; creation TIMESTAMP }")
	}

	return nil
}
