package server

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type user struct {
	Username string
	Port     int
}

type channel struct {
	Port  int
	Users []user
}

type Server struct {
	Channels []channel
	Upgrader websocket.Upgrader
}

func (s *Server) wsHanlder(w http.ResponseWriter, r *http.Request) {
	conn, err := s.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error while upgrading websocket")
		return
	}

	defer conn.Close()

	for {
		// здесь обрабатывать сообщения
	}
}

func (s *Server) StartServer() error {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	s.Upgrader = upgrader

	r := mux.NewRouter()
	r.HandleFunc("/{id}", s.wsHanlder)

	fmt.Println("Server started on http://localhost:8080")
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		return fmt.Errorf("Error while starting server: %s", err.Error())
	}
	return nil
}
