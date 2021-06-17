package controller

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"gocode/chatroom/client/model"
	"gocode/chatroom/client/utils"
	"gocode/chatroom/common/message"
	"net"
)

type UserController struct {
}

func (this *UserController) Login(userId int, userPwd string) (err error) {
	// fmt.Printf("userId = %d userPwd = %s\n", userId, userPwd)

	// 1. connect to server
	conn, err := net.Dial("tcp", "localhost:8889")
	if err != nil {
		fmt.Println("net.Dial err=", err)
		return
	}

	// should close connection in logout
	// defer conn.Close()

	// 2. prepare message
	var mes message.Message
	mes.Type = message.LoginMesType

	// 3. create LoginMes struct
	var loginMes message.LoginMes
	loginMes.UserId = userId
	loginMes.UserPwd = userPwd

	// 4. serialize loginMes (convert to string)
	data, err := json.Marshal(loginMes)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return
	}

	// 5. assign data to mes.Data
	mes.Data = string(data)

	// 6. serialize mes
	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return
	}

	// 7. send message(data) to server
	// 7.1 send data content length (convert data content length -> byte[])
	var pkgLen uint32          // 0 - 4294967295
	pkgLen = uint32(len(data)) // cast int to uint32
	var buffer [4]byte         // 4 * 8 = 32 bytes
	binary.BigEndian.PutUint32(buffer[0:4], pkgLen)

	n, err := conn.Write(buffer[:4])
	if n != 4 || err != nil {
		fmt.Println("conn.Write(buffer) fail", err)
		return
	}

	// fmt.Printf("Client sent data content length=%d, mes=%s\n", len(data), string(data))

	// 7.2 send actual data(mes)
	_, err = conn.Write(data)
	if err != nil {
		fmt.Println("conn.Write(data) fail", err)
		return
	}

	// 8. manage response messsage from server
	tf := &utils.Transfer{
		Conn: conn,
	}
	mes, err = tf.ReadPkg()
	if err != nil {
		fmt.Println("readPkg err=", err)
		return
	}

	var loginResMes message.LoginResMes
	err = json.Unmarshal([]byte(mes.Data), &loginResMes)
	if loginResMes.Code == 200 {
		fmt.Println("Login successful")

		CurUser = &model.CurUser{
			Conn:       conn,
			UserId:     userId,
			UserStatus: message.UserOnline,
		}

		// show current online users
		fmt.Println("Current online users:")
		for _, v := range loginResMes.UserIds {

			// won't show your own id
			if v == userId {
				continue
			}

			fmt.Println("User id:\t", v)

			//init onlineUsers
			user := &message.User{
				UserId:     v,
				UserStatus: message.UserOnline,
			}
			onlineUsers[v] = user
		}
		fmt.Println()
	} else {
		fmt.Println(loginResMes.Error)
		conn.Close()
		return errors.New(loginResMes.Error)
	}

	return
}

func (this *UserController) Register(userId int, userPwd string, userName string) (err error) {
	// 1. connect to server
	conn, err := net.Dial("tcp", "localhost:8889")
	if err != nil {
		fmt.Println("net.Dial err=", err)
		return
	}

	defer conn.Close() // close connection after registration

	// 2. prepare message
	var mes message.Message
	mes.Type = message.RegisterMesType

	// 3. create RegisterMes struct
	var registerMes message.RegisterMes
	registerMes.User.UserId = userId
	registerMes.User.UserPwd = userPwd
	registerMes.User.UserName = userName

	// 4. serialize registerMes (convert to string)
	data, err := json.Marshal(registerMes)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return
	}

	// 5. assign data to mes.Data
	mes.Data = string(data)

	// 6. serialize mes
	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return
	}

	tf := &utils.Transfer{
		Conn: conn,
	}

	// 7. send data to server
	err = tf.WritePkg(data)
	if err != nil {
		fmt.Println("Register message fail to send err=", err)
	}

	mes, err = tf.ReadPkg() // mes is RegisterResMes
	if err != nil {
		fmt.Println("readPkg err=", err)
		return
	}

	var registerResMes message.RegisterResMes
	err = json.Unmarshal([]byte(mes.Data), &registerResMes)
	if registerResMes.Code == 200 {
		fmt.Println("Register successful, please login now")
	} else {
		fmt.Println(registerResMes.Error)
	}

	return

}

func (this *UserController) Logout() (err error) {
	// 1. prepare message
	var mes message.Message
	mes.Type = message.LogoutMesType

	// 2. create LogoutRes struct
	var logoutMes message.LogoutMes
	logoutMes.UserId = CurUser.UserId

	// 3. serialize logoutRes (convert to string)
	data, err := json.Marshal(logoutMes)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return
	}

	// 4. assign data to mes.Data
	mes.Data = string(data)

	// 5. serialize mes
	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return
	}

	tf := &utils.Transfer{
		Conn: CurUser.Conn,
	}

	// 7. send data to server
	err = tf.WritePkg(data)
	if err != nil {
		fmt.Println("Logout message fail to send err=", err)
	}

	return

}
