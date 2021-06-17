package controller

import (
	"encoding/json"
	"fmt"
	"gocode/chatroom/client/utils"
	"gocode/chatroom/common/message"
)

type SmsController struct {
}

func (this *SmsController) SendGroupMes(content string) (err error) {
	// 1. create a mes
	var mes message.Message
	mes.Type = message.SmsMesType

	// 2. create SmsMes instant
	var smsMes message.SmsMes
	smsMes.Content = content
	smsMes.UserId = CurUser.UserId
	smsMes.UserStatus = CurUser.UserStatus

	// 3. serialize smsMes
	data, err := json.Marshal(smsMes)
	if err != nil {
		fmt.Println("SendGroupMes json.Marshal fail=", err.Error())
		return
	}

	mes.Data = string(data)

	// 4. serialize mes
	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("SendGroupMes json.Marshal fail=", err.Error())
		return
	}

	// 5. start sending to server
	tf := &utils.Transfer{
		Conn: CurUser.Conn,
	}

	err = tf.WritePkg(data)
	if err != nil {
		fmt.Println("SendGroupMes err=", err.Error())
		return
	}

	return

}
