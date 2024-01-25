package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"redis-server/handler"
	"redis-server/parser"
	"strings"
)

func main() {
	ln, err := net.Listen("tcp", ":6379")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	defer ln.Close()
	fmt.Println("Server is up on port 6379.")
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v", err)
		}
		go handle(conn)
	}
}

func handle(conn net.Conn) {
	defer conn.Close()
	for {
		buf := make([]byte, 1024)
		_, err := conn.Read(buf)
		if err != nil {
			if err != io.EOF {
				log.Printf("Failed to read: %v", err)
			}
			break
		}
		data, _ := parser.Parse(buf)
		if handler.StringCommands[strings.ToLower(data.Array[0].Str)] {
			conn.Write(handler.HandleString(data))
		} else if handler.ListCommands[strings.ToLower(data.Array[0].Str)] {
			conn.Write(handler.HandleList(data))
		} else {
			conn.Write(handler.ErrResp(errors.New("command not supported")))
		}
	}
}
