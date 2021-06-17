package main

import (
	"bufio"
	"fmt"
	"gocode/chatroom/client/controller"
	"os"
	"strconv"
)

// to capture user input
var (
	userId   int
	userPwd  string
	userName string
)

func main() {

	// record user selection
	var key int
	// whether to show first page menu
	var loop = true
	// scan input from terminal
	scanner := bufio.NewScanner(os.Stdin)

	// first page
	for loop {

		fmt.Println()
		fmt.Println("---------------Welcome to chat system---------------")
		fmt.Println("\t\t1. Login to chat room")
		fmt.Println("\t\t2. Register")
		fmt.Println("\t\t3. Exit")
		fmt.Printf("Please select (1-3): ")
		fmt.Scanf("%d\n", &key)

		switch key {
		case 1:
			fmt.Println("Login......")
			fmt.Printf("Please input user id: ")
			scanner.Scan()
			userId, err := strconv.Atoi(scanner.Text()) // convert string to int
			if err != nil {
				fmt.Println("strconv.Atoi(scanner.Text() err=", err)
				continue
			}
			fmt.Printf("Please input password: ")
			scanner.Scan()
			userPwd = scanner.Text()

			uc := &controller.UserController{}
			err = uc.Login(userId, userPwd)
			if err == nil {
				showMenu()
			}

		case 2:
			fmt.Println("Registering......")
			fmt.Printf("Please input user id: ")
			scanner.Scan()
			userId, err := strconv.Atoi(scanner.Text())
			if err != nil {
				fmt.Println("strconv.Atoi err=", err)
				continue
			}
			fmt.Printf("Please input user password: ")
			scanner.Scan()
			userPwd = scanner.Text()
			fmt.Printf("Please input username(nickname): ")
			scanner.Scan()
			userName = scanner.Text()

			uc := &controller.UserController{}
			uc.Register(userId, userPwd, userName)

		case 3:
			fmt.Println("Exit")
			os.Exit(0)

		default:
			fmt.Println("Wrong input, please try again")
		}
	}
}
