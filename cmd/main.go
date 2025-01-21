package main

import (
	"broadcast-server/internal/server"
	"broadcast-server/internal/user"
	"flag"
	"fmt"
	"os"
)

func start() {
	var s server.Server
	err := s.StartServer()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

func connect(port int, username string) {
	if port < 1000 {
		fmt.Println("port must be from 1000 to 9999")
		return
	} else if port == 8080 {
		fmt.Println("port 8080 can't be used as it reserved by server")
		return
	}

	err := user.Connect(port, username)
	if err != nil {
		fmt.Printf("Error while connection to the server: %w", err.Error())
		return
	}
}

func main() {
	if len(os.Args) == 1 {
		fmt.Println("use `help` command")
	} else if os.Args[1] == "start" {
		start()
	} else if os.Args[1] == "connect" {
		// according to the source files, flag.Parse() parses os.Args, so i'm cutting os.Args to avoid parsing command
		os.Args = os.Args[1:]
		PORT := flag.Int("port", 0, "chat port")
		USERNAME := flag.String("username", "", "username that will be displayed")
		flag.Parse()

		connect(*PORT, *USERNAME)
	} else if os.Args[1] == "help" {
		fmt.Println("to start server use ->\tbroadcast-server start")
		fmt.Println("to connect to chat with use ->\tbroadcast-server connect --port <port> --username <username>")
	} else {
		fmt.Println("unknown command\nuse `help` command")
	}
}
