package application

import (
	"sync"

	"github.com/gorilla/websocket"
)

// User - структура, представляющая пользователя.
type User struct {
	nickname   string
	connection *websocket.Conn
	query      [][]byte
	mu         sync.Mutex
	started    bool
}

// GetUser - создание нового пользователя.
//
// Принимает: имя пользователя, соединение пользователя.
//
// Возвращает: пользователя.
func GetUser(nick string, conn *websocket.Conn) *User {
	return &User{
		nickname:   nick,
		connection: conn,
		query:      [][]byte{},
		mu:         sync.Mutex{},
		started:    false,
	}
}

// GetNickname - получение имени пользователя.
//
// Возвращает: имя.
func (u *User) GetNickname() string {
	return u.nickname
}

// AddMessage - добавление сообщения пользователю.
//
// Принимает: сообщение.
func (u *User) AddMessage(msg []byte) {
	u.mu.Lock()
	u.query = append(u.query, msg)
	u.mu.Unlock()
}

// Start - начало обработки пользователя.
func (u *User) Start() {
	u.mu.Lock()
	u.started = true
	u.mu.Unlock()

	go func(u *User) {
		for {
			u.mu.Lock()
			if !u.started {
				break
			}
			if len(u.query) > 0 {
				u.connection.WriteMessage(websocket.TextMessage, u.query[0])

				u.query = u.query[1:]
			}

			u.mu.Unlock()
		}
	}(u)
}

// Start - конец обработки пользователя.
func (u *User) Stop() {
	u.mu.Lock()
	if u.started {
		u.started = false
	}
	u.mu.Unlock()
}
