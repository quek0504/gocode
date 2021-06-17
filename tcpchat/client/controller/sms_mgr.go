package controller

import (
	"encoding/json"
	"fmt"
	"gocode/chatroom/common/message"
)

func OutputGroupMes(mes *message.Message) { // SmsMes type
	// 1. deserialize mes
	var smsMes message.SmsMes
	err := json.Unmarshal([]byte(mes.Data), &smsMes)
	if err != nil {
		fmt.Println("json.Unmarshal err=", err.Error())
		return
	}

	// 2. show content, who saying what
	content := fmt.Sprintf("\nUSER %d saying: %s",
		smsMes.UserId, smsMes.Content)
	fmt.Println(content)
	fmt.Println()
}
