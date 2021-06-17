package controller

import (
	"encoding/json"
	"fmt"
	"gocode/chatroom/common/message"
	"gocode/chatroom/server/model"
	"gocode/chatroom/server/utils"
	"net"
)

type UserController struct {
	Conn   net.Conn
	UserId int // user connection
}

// broadcast to every currently online users that I am online now (after login)
func (this *UserController) notifyOthersOnlineUser(userId int) {
	for id, uc := range userMgr.onlineUsers {
		// yourself
		if id == userId {
			continue
		}
		// start notifying others
		uc.notifyIAmOnline(userId)
	}
}

func (this *UserController) notifyIAmOnline(userId int) {
	// prepare mes
	var mes message.Message
	mes.Type = message.NotifyUserStatusMesType

	// prepare notifyUserStatusMes
	var notifyUserStatusMes message.NotifyUserStatusMes
	notifyUserStatusMes.UserId = userId
	notifyUserStatusMes.Status = message.UserOnline

	// serialize notifyUserStatusMes
	data, err := json.Marshal(notifyUserStatusMes)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return
	}

	mes.Data = string(data)

	// serialize mes
	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return
	}

	// send message
	tf := &utils.Transfer{
		Conn: this.Conn,
	}

	err = tf.WritePkg(data)
	if err != nil {
		fmt.Println("NotifyIAmOnline err=", err)
		return
	}

}

func (this *UserController) ServerProcessLogin(mes *message.Message) (err error) {
	// 1. get mes.Data
	var loginMes message.LoginMes
	err = json.Unmarshal([]byte(mes.Data), &loginMes)
	if err != nil {
		fmt.Println("json.Unmarshal fail err=", err)
		return
	}

	// 2. prepare response message
	var resMes message.Message
	resMes.Type = message.LoginResMesType

	var loginResMes message.LoginResMes

	// 3. verify credentials
	user, err := model.MyUserDao.Login(loginMes.UserId, loginMes.UserPwd)

	if err != nil {
		if err == model.ERROR_USER_NOTEXISTS {
			loginResMes.Code = 401 // user not exist
			loginResMes.Error = err.Error()
		} else if err == model.ERROR_USER_PWD {
			loginResMes.Code = 300
			loginResMes.Error = err.Error()
		} else {
			loginResMes.Code = 500 // unknown error
			loginResMes.Error = "Internal server error"
		}
	} else {
		loginResMes.Code = 200

		this.UserId = loginMes.UserId
		// update onlineUsers
		userMgr.AddOnlineUser(this)
		this.notifyOthersOnlineUser(loginMes.UserId)
		// put currently online users to loginResMes.UsersId
		for id, _ := range userMgr.onlineUsers {
			loginResMes.UserIds = append(loginResMes.UserIds, id)
		}
		fmt.Println("USER", user.UserId, "login successful")
	}

	// 4. serialize loginResMes
	data, err := json.Marshal(loginResMes)
	if err != nil {
		fmt.Println("json.Marshal(loginResMes) fail", err)
		return
	}

	// 5. serialize response message
	resMes.Data = string(data)

	data, err = json.Marshal(resMes)
	if err != nil {
		fmt.Println("json.Marshal(resMes) fail", err)
		return
	}

	tf := &utils.Transfer{
		Conn: this.Conn,
	}

	// 6. send response message
	err = tf.WritePkg(data)
	return

}

func (this *UserController) ServerProcessRegister(mes *message.Message) (err error) {
	// 1. get mes.Data
	var registerMes message.RegisterMes
	err = json.Unmarshal([]byte(mes.Data), &registerMes)
	if err != nil {
		fmt.Println("json.Unmarshal fail err=", err)
		return
	}

	// 2. prepare response message
	var resMes message.Message
	resMes.Type = message.RegisterResMesType

	var registerResMes message.RegisterResMes

	// 3. register logic
	err = model.MyUserDao.Register(&registerMes.User)
	if err != nil {
		if err == model.ERROR_USER_EXISTS {
			registerResMes.Code = 403
			registerResMes.Error = err.Error()
		} else {
			registerResMes.Code = 500
			registerResMes.Error = "Internal Server Error"
		}
	} else {
		registerResMes.Code = 200
	}

	// 4. serialize registerResMes
	data, err := json.Marshal(registerResMes)
	if err != nil {
		fmt.Println("json.Marshal(loginResMes) fail", err)
		return
	}

	// 5. serialize response message
	resMes.Data = string(data)

	data, err = json.Marshal(resMes)
	if err != nil {
		fmt.Println("json.Marshal(resMes) fail", err)
		return
	}

	tf := &utils.Transfer{
		Conn: this.Conn,
	}

	// 5. send response message
	err = tf.WritePkg(data)
	return

}

func (this *UserController) ServerProcessLogout(mes *message.Message) (err error) {
	// 1. get mes.Data
	var logoutMes message.LogoutMes
	err = json.Unmarshal([]byte(mes.Data), &logoutMes)
	if err != nil {
		fmt.Println("json.Unmarshal fail err=", err)
		return
	}

	// 2. prepare response message
	var resMes message.Message
	resMes.Type = message.LogoutResMesType
	var logoutResMes message.LogoutResMes

	// 3. logout logic, update onlineUsers
	err = userMgr.DeleteOnlineUser(logoutMes.UserId)
	if err != nil {
		logoutResMes.Code = 500
		logoutResMes.Error = err.Error()
	} else {
		logoutResMes.Code = 200
	}

	// 4. serialize logoutResMes
	data, err := json.Marshal(logoutResMes)
	if err != nil {
		fmt.Println("json.Marshal(logoutResMes) fail", err)
		return
	}

	// 5. serialize response message
	resMes.Data = string(data)

	data, err = json.Marshal(resMes)
	if err != nil {
		fmt.Println("json.Marshal(resMes) fail", err)
		return
	}

	tf := &utils.Transfer{
		Conn: this.Conn,
	}

	// 6. send response message
	err = tf.WritePkg(data)

	return

}
