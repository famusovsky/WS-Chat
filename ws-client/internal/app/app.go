package app

import (
	"bufio"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"

	"github.com/gorilla/websocket"
	"golang.org/x/sync/errgroup"
)

const exitComm = ":exit"

// Run - запуск приложения клиента.
func Run(host string, nickname string, scanner *bufio.Scanner) error {
	u := url.URL{Scheme: "ws", Host: host, Path: "/"}

	h := http.Header{}
	h.Add("Nickname", nickname)
	c, _, err := websocket.DefaultDialer.Dial(u.String(), h)
	if err != nil {
		return err
	}
	defer c.Close()

	connOk := atomic.Bool{}
	connOk.Store(true)

	wg := new(sync.WaitGroup)
	wg.Add(1)
	go grace(c, wg)

	go func() {
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				connOk.Store(false)
				return
			}
			fmt.Printf("%s\n", message)
		}
	}()

	go func() {
		fmt.Printf("To exit the chat input '%s'\n", exitComm)

		for {
			if !connOk.Load() {
				err = errors.New("connection is bad")
				break
			}
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
				break
			}
		}

		wg.Done()
	}()

	wg.Wait()
	return err
}

// grace - отслеживание и обработка принудительного выхода.
func grace(c *websocket.Conn, wg *sync.WaitGroup) {
	sigQuit := make(chan os.Signal, 2)
	signal.Notify(sigQuit, syscall.SIGINT, syscall.SIGTERM)

	eg := new(errgroup.Group)

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
	wg.Done()
}
