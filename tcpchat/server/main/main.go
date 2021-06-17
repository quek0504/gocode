package main

import (
	"fmt"
	"gocode/chatroom/server/model"
	"net"
	"time"
)

func process(conn net.Conn) {
	defer conn.Close() // close connection when err occurs

	// keep connection alive with client, to read any message from client
	process := &Process{
		Conn: conn,
	}

	err := process.serve()
	if err != nil {
		// print error and close connection
		fmt.Println("Goroutine err=", err)
		return
	}
}

func init() {
	// initialize redis connection pool
	initPool("localhost:6379", 16, 0, 300*time.Second)
	initUserDao()
}

func initUserDao() {
	// note : this must run after initPool
	// MyUsreDao needs pool from redis.go
	model.MyUserDao = model.NewUserDao(pool)
}

func main() {
	fmt.Println("Server listening at port 8889...")
	listen, err := net.Listen("tcp", "0.0.0.0:8889")
	defer listen.Close()
	if err != nil {
		fmt.Println("net.Listen err=", err)
		return
	}

	for {
		fmt.Println("Waiting for client...")
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("listen.Accept err=", err)
		}

		// create goroutine to process each client request
		go process(conn)

	}

}
