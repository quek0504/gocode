package main

import (
	"bufio"
	"fmt"
	"gocode/chatroom/client/controller"
	"gocode/chatroom/client/utils"
	"gocode/chatroom/common/message"
	"io"
	"net"
	"os"
)

// manage user session
var (
	session bool
)

// show menu after login
func showMenu() {
	sc := &controller.SmsController{}
	uc := &controller.UserController{}
	scanner := bufio.NewScanner(os.Stdin)
	var key int
	var content string
	session = true

	// keep connection alive with server, to read any message from server
	go processServerMes(controller.CurUser.Conn)

	for session {
		fmt.Printf("---------------Welcome back USER %d--------------\n", controller.CurUser.UserId)
		fmt.Println("\t\t1. Show online users")
		fmt.Println("\t\t2. Send message")
		fmt.Println("\t\t3. Conversation list")
		fmt.Println("\t\t4. Sign out")
		fmt.Printf("Please select (1-4): ")
		fmt.Scanf("%d\n", &key)
		switch key {
		case 1:
			controller.ShowOnlineUser()
			fmt.Println()
		case 2:
			fmt.Println("Please say something to all:)")
			scanner.Scan()
			content = scanner.Text()
			sc.SendGroupMes(content)
			fmt.Println()
		case 3:
			fmt.Println("Conversation list")
			fmt.Println()
		case 4:
			fmt.Println("Closing connection...")
			uc.Logout()
			session = false // control back to main
		default:
			fmt.Println("Wrong input, please try again")
			fmt.Println()
		}
	}
}

// manage connection with server
func processServerMes(conn net.Conn) {
	defer conn.Close()

	tf := &utils.Transfer{
		Conn: conn,
	}
	// read message from server
	for {
		mes, err := tf.ReadPkg()
		if err != nil {
			session = false
			if err == io.EOF {
				fmt.Println("Server close connection...")
				return
			} else {
				fmt.Println("readPkg err=", err)
				return
			}
		}

		// if got message, determine its mesType
		switch mes.Type {
		case message.NotifyUserStatusMesType:
			// server notify client someone has online
			// add the one who has just online to client managed list of online users, map[int]onlineUsers
			controller.UpdateUserStatus(&mes)
		case message.SmsMesType: // broadcast message
			controller.OutputGroupMes(&mes)
		case message.LogoutResMesType:
			controller.Logout()
			// not processing server message anymore, goroutine ends
			return
		default:
			fmt.Println("Unknown message")
		}
	}
}
