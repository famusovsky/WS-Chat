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
		return errors.New("error while starting transaction: " + err.Error())
	}
	defer tx.Rollback()

	q := `INSERT INTO messages (text, sender, creation) VALUES ($1, $2, $3);`
	err = tx.QueryRow(q, text, sender, creation).Err()
	if err != nil {
		return errors.New("error while adding message to the database: " + err.Error())
	}

	err = tx.Commit()
	if err != nil {
		return errors.New("error while committing transaction: " + err.Error())
	}

	return nil
}

func (m *messagesDB) GetLastMessages(count int) ([]models.Message, error) {
	q := `SELECT * FROM (
		SELECT * from messages ORDER BY creation DESC LIMIT $1
	  ) AS tmp ORDER BY creation ASC;`
	rows, err := m.db.Query(q, count)
	if err != nil {
		return nil, errors.New("error while getting messages from the database: " + err.Error())
	}

	res := []models.Message{}
	for rows.Next() {
		var message models.Message
		err = rows.Scan(&message.Id, &message.Text, &message.Sender, &message.Creation)
		if err != nil {
			return nil, fmt.Errorf("error while getting messages from the database: %s", err.Error())
		}
		res = append(res, message)
	}

	return res, nil
}

func (m *messagesDB) GetLastMessagesByNickname(count int, nickname string) ([]models.Message, error) {
	q := `SELECT * FROM (
		SELECT * from messages WHERE sender = $1 ORDER BY creation DESC LIMIT $2
	  ) AS tmp ORDER BY creation ASC;`
	rows, err := m.db.Query(q, nickname, count)
	if err != nil {
		return nil, errors.New("error while getting messages from the database: " + err.Error())
	}

	res := []models.Message{}
	for rows.Next() {
		var message models.Message
		err = rows.Scan(&message)
		if err != nil {
			return nil, fmt.Errorf("error while getting messages from the database: %s", err.Error())
		}
		res = append(res, message)
	}

	return res, nil
}

func (m *messagesDB) GetMessage(id int) (models.Message, error) {
	q := `SELECT * from messages WHERE id = $1;`
	rows, err := m.db.Query(q, id)
	if err != nil {
		return models.Message{}, errors.New("error while getting message from the database: " + err.Error())
	}

	var message models.Message
	err = rows.Scan(&message)
	if err != nil {
		return models.Message{}, fmt.Errorf("error while getting message from the database: %s", err.Error())
	}

	return message, nil
}

func overrideDB(db *sql.DB) error {
	q := `DROP TABLE IF EXISTS messages;`

	_, err := db.Exec(q)
	if err != nil {
		return fmt.Errorf("error while dropping tables: %s", err)
	}
	q = `CREATE TABLE messages (
		id SERIAL PRIMARY KEY,
		text TEXT,
		sender TEXT,
		creation TIMESTAMP
	);`

	_, err = db.Exec(q)
	if err != nil {
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

	var err error = nil

	err = db.QueryRow(qMessages).Scan(&properMessages)
	if err != nil {
		return errors.Join(errors.New("error while checking 'messages' table"), err)
	}

	if !properMessages {
		err = errors.Join(err, errors.New(
			"'messages' table is not ok: proper 'messages' table is { id INTEGER; text TEXT; sender TEXT; creation TIMESTAMP }"))
	}

	return err
}
