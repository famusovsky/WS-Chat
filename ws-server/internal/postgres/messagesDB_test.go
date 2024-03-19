// package postgres

// import (
// 	"testing"
// 	"time"

// 	"github.com/DATA-DOG/go-sqlmock"
// 	"github.com/stretchr/testify/assert"
// )

// func Test_messagesDB_AddMessage(t *testing.T) {
// 	db, mock, err := sqlmock.New()
// 	if err != nil {
// 		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
// 	}
// 	defer db.Close()

// 	messages := messagesDB{db: db}
// 	sender := "test_sender"
// 	text := "test_message"
// 	creation := time.Now()

// 	mock.ExpectBegin()
// 	mock.ExpectExec("INSERT INTO messages").
// 		WithArgs(text, sender, creation).
// 		WillReturnResult(sqlmock.NewResult(1, 1))
// 	mock.ExpectCommit()

// 	err = messages.AddMessage(sender, text, creation)
// 	assert.NoError(t, err)
// 	assert.NoError(t, mock.ExpectationsWereMet())
// }

// func Test_messagesDB_GetLastMessages(t *testing.T) {
// 	db, mock, err := sqlmock.New()
// 	if err != nil {
// 		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
// 	}
// 	defer db.Close()

// 	messages := messagesDB{db: db}
// 	count := 5

// 	rows := sqlmock.NewRows([]string{"id", "text", "sender", "creation"}).
// 		AddRow(1, "text1", "sender1", time.Now()).
// 		AddRow(2, "text2", "sender2", time.Now()).
// 		AddRow(3, "text3", "sender3", time.Now())

// 	mock.ExpectQuery("SELECT .* FROM messages.*").
// 		WithArgs(count).
// 		WillReturnRows(rows)

// 	result, err := messages.GetLastMessages(count)
// 	assert.NoError(t, err)
// 	assert.Len(t, result, 3)
// 	assert.NoError(t, mock.ExpectationsWereMet())
// }

// func Test_messagesDB_GetLastMessagesByNickname(t *testing.T) {
// 	db, mock, err := sqlmock.New()
// 	if err != nil {
// 		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
// 	}
// 	defer db.Close()

// 	messages := messagesDB{db: db}
// 	count := 5
// 	nickname := "test_sender"

// 	rows := sqlmock.NewRows([]string{"id", "text", "sender", "creation"}).
// 		AddRow(1, "text1", "sender1", time.Now()).
// 		AddRow(2, "text2", "sender1", time.Now()).
// 		AddRow(3, "text3", "sender1", time.Now())

// 	mock.ExpectQuery("SELECT .* FROM messages.*").
// 		WithArgs(nickname, count).
// 		WillReturnRows(rows)

// 	result, err := messages.GetLastMessagesByNickname(count, nickname)
// 	assert.NoError(t, err)
// 	assert.Len(t, result, 3)
// 	assert.NoError(t, mock.ExpectationsWereMet())
// }

// func Test_messagesDB_GetMessage(t *testing.T) {
// 	db, mock, err := sqlmock.New()
// 	if err != nil {
// 		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
// 	}
// 	defer db.Close()

// 	messages := messagesDB{db: db}
// 	messageID := 1

// 	rows := sqlmock.NewRows([]string{"id", "text", "sender", "creation"}).
// 		AddRow(1, "text1", "sender1", time.Now())

// 	mock.ExpectQuery("SELECT .* FROM messages.*").
// 		WithArgs(messageID).
// 		WillReturnRows(rows)

// 	result, err := messages.GetMessage(messageID)
// 	assert.NoError(t, err)
// 	assert.NotNil(t, result)
// 	assert.NoError(t, mock.ExpectationsWereMet())
// }

// func Test_overrideDB(t *testing.T) {
// 	db, mock, err := sqlmock.New()
// 	if err != nil {
// 		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
// 	}
// 	defer db.Close()

// 	mock.ExpectExec("DROP TABLE IF EXISTS messages").WillReturnResult(sqlmock.NewResult(0, 0))
// 	mock.ExpectExec("CREATE TABLE messages").WillReturnResult(sqlmock.NewResult(0, 0))

// 	err = overrideDB(db)
// 	assert.NoError(t, err)
// 	assert.NoError(t, mock.ExpectationsWereMet())
// }

// func Test_checkDB(t *testing.T) {
// 	db, mock, err := sqlmock.New()
// 	if err != nil {
// 		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
// 	}
// 	defer db.Close()

// 	rows := sqlmock.NewRows([]string{"properMessages"}).AddRow(true)
// 	mock.ExpectQuery("SELECT .* FROM information_schema.columns.*").WillReturnRows(rows)

// 	err = checkDB(db)
// 	assert.NoError(t, err)
// 	assert.NoError(t, mock.ExpectationsWereMet())
// }
