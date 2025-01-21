package server

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type channel struct {
	Conn  map[string]*websocket.Conn // map[usesrname]conn
	Users map[string]bool            // map[username]bool
}

type Server struct {
	Channels map[int]channel //map[port]channel
	Upgrader websocket.Upgrader
}

func (s *Server) wsHanlder(w http.ResponseWriter, r *http.Request) {
	conn, err := s.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error while upgrading websocket: ", err.Error())
		return
	}

	defer conn.Close()

	//обработка сообщений
	for {
		messageType, message, err := conn.ReadMessage()

		//обработка отключения
		if err != nil {
			vars := mux.Vars(r)
			username := vars["username"]
			port, _ := strconv.Atoi(vars["port"])
			delete(s.Channels[port].Conn, username)
			delete(s.Channels[port].Users, username)

			for _, connection := range s.Channels[port].Conn {
				err = connection.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf(">---user %s disconnected---<", username)))
				if err != nil {
					fmt.Println("error:", err.Error())
				}
			}
			break
		}

		vars := mux.Vars(r)
		username := vars["username"]
		port, err := strconv.Atoi(vars["port"])
		if err != nil {
			fmt.Printf("invalid port: %s", vars["port"])
			break
		}

		if _, ok := s.Channels[port]; !ok {
			s.Channels[port] = channel{Users: map[string]bool{}, Conn: map[string]*websocket.Conn{}}
			log.Printf("new chat opened on port %d\n", port)
		}

		if _, ok := s.Channels[port].Users[username]; !ok {
			s.Channels[port].Users[username] = true
			s.Channels[port].Conn[username] = conn

			for _, connection := range s.Channels[port].Conn {
				err := connection.WriteMessage(messageType, []byte(fmt.Sprintf("<---user %s connected--->", username)))
				if err != nil {
					fmt.Println("error:", err.Error())
				}
			}
			continue
		}

		if len(message) > 1 {
			if string(message[:len(message)-1]) == "--users" {
				users := "<---connected users--->"
				for name := range s.Channels[port].Users {
					users += "\n\t" + name
				}
				users += "\n"
				message = []byte(users)
			} else {
				message = []byte(fmt.Sprintf("%d:%s -> %s", port, username, message))
			}
			for _, connection := range s.Channels[port].Conn {
				err = connection.WriteMessage(messageType, message)
				if err != nil {
					fmt.Println("error:", err.Error())
				}
			}
		}
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
	s.Channels = make(map[int]channel)
	s.Upgrader = upgrader

	r := mux.NewRouter()
	r.HandleFunc("/{port}/{username}", s.wsHanlder)

	fmt.Println("Server started on http://localhost:8080")
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		return fmt.Errorf("error while starting server: %s", err.Error())
	}
	return nil
}
