package models

import (
	"time"
)

// Message - структура, представляющая сообщение.
type Message struct {
	Id       int       // Id - id сообщения.
	Text     string    // Text - текст сообщения.
	Sender   string    // Sender - имя отправителя сообщения.
	Creation time.Time // Creation - время создания сообщения.
}
