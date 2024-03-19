package app

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/gorilla/websocket"
	"golang.org/x/sync/errgroup"
)

const exitComm = ":exit"

// Run - запуск приложения клиента.
//
// Принимает: адрес сервера, имя пользователя, сканнер ввода.
//
// Возвращает: ошибку.
func Run(host string, nickname string, scanner *bufio.Scanner) error {
	u := url.URL{Scheme: "ws", Host: host, Path: "/"}

	h := http.Header{}
	h.Add("Nickname", nickname)
	c, _, err := websocket.DefaultDialer.Dial(u.String(), h)
	if err != nil {
		return err
	}
	defer c.Close()

	connOk := true
	mu := sync.Mutex{}

	go grace(c)

	go func(ok *bool, mu *sync.Mutex) {
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				mu.Lock()
				connOk = false
				mu.Unlock()
				return
			}
			fmt.Println(string(message))
		}
	}(&connOk, &mu)

	fmt.Printf("To exit the chat input '%s'\n", exitComm)

	for {
		mu.Lock()
		if !connOk {
			return errors.New("connection is bad")
		}
		mu.Unlock()
		ok := scanner.Scan()
		if !ok {
			continue
		}
		message := scanner.Text()
		if message == exitComm {
			break
		}
		err = c.WriteMessage(websocket.TextMessage, []byte(message))
		if err != nil {
			return err
		}
	}

	return nil
}

// grace - отслеживание и обработка принудительного выхода.
//
// Принимает: соединение.
func grace(c *websocket.Conn) {
	sigQuit := make(chan os.Signal, 2)
	signal.Notify(sigQuit, syscall.SIGINT, syscall.SIGTERM)

	eg, _ := errgroup.WithContext(context.Background())

	eg.Go(func() error {
		select {
		case s := <-sigQuit:
			return fmt.Errorf("captured signal: %v", s)
		}
	})

	if err := eg.Wait(); err != nil {
		fmt.Printf("gracefully shutting down the client: %v", err)
	}

	err := c.Close()
	if err != nil {
		fmt.Printf("error whil closing the connection: %v", err)
	}
	os.Exit(1) // XXX Я не придумал, как нормально закрыть цикл считывания сообщений, но вроде и так всё ок.
}
