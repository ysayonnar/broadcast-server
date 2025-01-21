package user

import (
	"bufio"
	"fmt"
	"os"

	"github.com/gorilla/websocket"
)

func Connect(port int, username string) error {
	serverAddr := fmt.Sprintf("ws://localhost:8080/%d/%s", port, username)

	conn, _, err := websocket.DefaultDialer.Dial(serverAddr, nil)
	if err != nil {
		return err
	}
	defer conn.Close()

	err = conn.WriteMessage(websocket.TextMessage, []byte(""))
	if err != nil {
		return err
	}

	fmt.Println("Welcome to chat on port", port)
	_, message, err := conn.ReadMessage()
	if err != nil {
		fmt.Println("error: ", err.Error())
	}
	fmt.Println(string(message))

	go func() {
		for {
			reader := bufio.NewReader(os.Stdin)
			msg, _ := reader.ReadString('\n')
			fmt.Print("\033[F")
			fmt.Print("\033[K")
			msg = msg[:len(msg)-1] // обрезаем переход на новую строку
			if msg == "" {
				continue
			} else {
				conn.WriteMessage(websocket.TextMessage, []byte(msg))
			}
		}
	}()

	for {
		_, message, err := conn.ReadMessage()
		//обработка отключения
		if err != nil {
			conn.Close()
			fmt.Println("Server stopped.")
			os.Exit(1)
			break
		}
		if err != nil {
			fmt.Println("error: ", err.Error())
		}
		fmt.Println(string(message))
	}

	return nil
}
