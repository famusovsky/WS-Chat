package application

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

// connect - обработка подключения к новому пользователю.
func (app *App) connect(w http.ResponseWriter, r *http.Request) {
	var nickname string
	if nicks, ok := r.Header["Nickname"]; ok && len(nicks) > 0 {
		nickname = nicks[0]
	} else {
		nickname = "Unknown User"
	}

	connection, _ := app.upgrader.Upgrade(w, r, nil)
	defer connection.Close()

	connMsg := fmt.Sprintf("User connected with nickname: %s", nickname)
	exitMsg := fmt.Sprintf("user %s left the chat", nickname)

	app.infoLog.Print(connMsg)
	user := GetUser(nickname, connection)
	app.mu.Lock()
	app.users[user] = struct{}{}
	app.mu.Unlock()

	defer app.mu.Unlock()
	defer delete(app.users, user)
	defer app.mu.Lock()
	defer user.Stop()
	user.Start()

	history, err := app.db.GetLastMessages(10)
	if err != nil {
		user.AddMessage([]byte("Error while getting chat history has occured"))
		app.errorLog.Printf("Obtaining history error: %v", err)
	} else {
		for i := 0; i < len(history); i++ {
			msg := history[i]
			user.AddMessage(formatMessage(msg.Creation, []byte(msg.Text), []byte(msg.Sender)))
		}
	}

	app.WriteMessage(time.Now(), []byte(connMsg), "")

	for {
		mt, msg, err := connection.ReadMessage()

		if err != nil || mt == websocket.CloseMessage {
			if mt != websocket.CloseMessage {
				app.errorLog.Printf("connection error: %v\n", err)
			}

			app.infoLog.Print(exitMsg)
			app.WriteMessage(time.Now(), []byte(exitMsg), "")
			break
		}

		creation := time.Now()

		go app.WriteMessage(creation, msg, nickname)

		go app.handleMessage(creation, msg, nickname)
	}
}

// WriteMessage - отправление сообщения всем подключенным пользователям.
//
// Принимает: время создания сообщения, текст сообщения, имя отправителя.
func (app *App) WriteMessage(creation time.Time, message []byte, nickname string) {
	for u := range app.users {
		// FIXME не отправлять сообщение отправителю -- проблема в том, что у пользователей могут быть одинаковые ники.
		// if nickname == u.GetNickname() {
		// 	continue
		// }
		go u.AddMessage(formatMessage(creation, message, []byte(nickname)))
	}
}

// handleMessage - обработка сообщения.
//
// Принимает: время создания сообщения, текст сообщения, имя отправителя.
func (app *App) handleMessage(creation time.Time, message []byte, nickname string) {
	app.infoLog.Printf("Message:\n%s\nsent by %s\n", string(message), nickname)

	err := app.db.AddMessage(nickname, string(message), creation)
	if err != nil {
		app.errorLog.Printf("error while sending message by %s: %v\n", nickname, err)
	}
}

// formatMessage - форматирование данных о сообщении в текстовый формат.
//
// Принимает: время создания сообщения, текст сообщения, имя отправителя.
func formatMessage(creation time.Time, message, nickname []byte) []byte {
	return []byte(fmt.Sprintf("%s\t%s:\t%s", creation.Format("Jan 2 15:04"), nickname, message))
}
