package controller

import (
	"encoding/json"
	"fmt"
	"gocode/chatroom/client/model"
	"gocode/chatroom/common/message"
)

var (
	onlineUsers map[int]*message.User = make(map[int]*message.User, 10) // who is currently online
	CurUser     *model.CurUser                                          // after user has login, init CurUser
)

// update who is currently online
func UpdateUserStatus(mes *message.Message) { // NotifyUserStatus type
	// deserialize mes
	var notifyUserStatusMes message.NotifyUserStatusMes
	err := json.Unmarshal([]byte(mes.Data), &notifyUserStatusMes)
	if err != nil {
		fmt.Println("json.Unmarshal err=", err.Error())
		return
	}

	// check onlineUsers
	user, ok := onlineUsers[notifyUserStatusMes.UserId]
	if !ok { // initially not in onlineUsers
		user = &message.User{
			UserId: notifyUserStatusMes.UserId,
		}
		fmt.Printf("\nUSER %d has just online\n", notifyUserStatusMes.UserId)
	}

	user.UserStatus = notifyUserStatusMes.Status
	onlineUsers[notifyUserStatusMes.UserId] = user

	ShowOnlineUser()
}

// show current online user in terminal
func ShowOnlineUser() {
	fmt.Println("Current online users:")
	for id, user := range onlineUsers {
		if user.UserStatus == 0 {
			fmt.Println("User id:\t", id)
		}
	}
}

func Logout() {
	// close connection and clean up CurUser
	CurUser.Conn.Close()
	CurUser = nil
}
