package main

import (
	"errors"
	"fmt"
	"gocode/chatroom/common/message"
	"gocode/chatroom/server/controller"
	"gocode/chatroom/server/utils"
	"io"
	"net"
)

type Process struct {
	Conn net.Conn
}

// manage connection with client
func (this *Process) serve() (err error) {
	// read message from client
	for {
		tf := &utils.Transfer{
			Conn: this.Conn,
		}
		mes, err := tf.ReadPkg()
		if err != nil {
			if err == io.EOF {
				return errors.New("Client close connection...")
			} else {
				fmt.Println("readPkg err=", err)
				return err
			}
		}
		err = this.serverProcessMes(&mes) // match mes type and process mes accordingly
		if err != nil {
			return err
		}
	}
}

func (this *Process) serverProcessMes(mes *message.Message) (err error) {
	// for debugging, printing message from client
	fmt.Println("mes=", mes)

	switch mes.Type {
	case message.LoginMesType:
		uc := &controller.UserController{
			Conn: this.Conn,
		}
		err = uc.ServerProcessLogin(mes)
	case message.RegisterMesType:
		uc := &controller.UserController{
			Conn: this.Conn,
		}
		err = uc.ServerProcessRegister(mes)
	case message.LogoutMesType:
		uc := &controller.UserController{
			Conn: this.Conn,
		}
		uc.ServerProcessLogout(mes)
	case message.SmsMesType:
		sc := &controller.SmsController{}
		sc.SendGroupMes(mes)
	default:
		fmt.Println("Message type not matched..")
	}
	return
}
