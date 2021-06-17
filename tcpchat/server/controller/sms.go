package controller

import (
	"encoding/json"
	"fmt"
	"gocode/chatroom/common/message"
	"gocode/chatroom/server/utils"
	"net"
)

type SmsController struct {
}

func (this *SmsController) SendGroupMes(mes *message.Message) {
	// broadcast message to everyone who is online

	// deserialize mes to extract smsMes.UserId
	var smsMes message.SmsMes
	err := json.Unmarshal([]byte(mes.Data), &smsMes)
	if err != nil {
		fmt.Println("json.Unmarshal err=", err)
		return
	}

	// serialize again and forward to client
	data, err := json.Marshal(mes)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return
	}

	for id, uc := range userMgr.onlineUsers {
		// filter out sender/youself, dont send to youself
		if id == smsMes.UserId {
			continue
		}
		// start broadcasting
		this.SendMesToEachOnlineUser(data, uc.Conn)
	}
}

func (this *SmsController) SendMesToEachOnlineUser(data []byte, conn net.Conn) {
	tf := &utils.Transfer{
		Conn: conn,
	}

	err := tf.WritePkg(data)
	if err != nil {
		fmt.Println("Broadcast message err=", err)
	}

}
