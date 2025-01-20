package main

import (
	"flag"
	"fmt"
	"os"
)

func start() {
	fmt.Println("Server started")
}

func connect(port int, username string) {
	fmt.Printf("PORT: %d, username: %s", port, username)
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
