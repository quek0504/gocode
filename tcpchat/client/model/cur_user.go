package model

import (
	"net"
)

// current logged in user/yourself
type CurUser struct {
	Conn       net.Conn
	UserId     int
	UserStatus int
}
