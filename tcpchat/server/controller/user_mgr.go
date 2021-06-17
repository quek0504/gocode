package controller

import "fmt"

// global var
var (
	userMgr *UserMgr
)

// init UserMgr
func init() {
	userMgr = &UserMgr{
		onlineUsers: make(map[int]*UserController, 1024),
	}
}

type UserMgr struct {
	onlineUsers map[int]*UserController // used by server to record user status
}

func (this *UserMgr) AddOnlineUser(uc *UserController) {
	val, ok := this.onlineUsers[uc.UserId]
	if ok {
		this.DeleteOnlineUser(uc.UserId)
		// user login again, kick out previous user session
		val.Conn.Close()
	}
	this.onlineUsers[uc.UserId] = uc
}

func (this *UserMgr) DeleteOnlineUser(userId int) (err error) {
	_, ok := this.onlineUsers[userId]
	if !ok {
		err = fmt.Errorf("User %d is not connected", userId)
		return
	}
	delete(this.onlineUsers, userId)
	return
}

func (this *UserMgr) GetAllOnlineUser() map[int]*UserController {
	return this.onlineUsers
}

func (this *UserMgr) GetOnlineUserById(userId int) (up *UserController, err error) {
	_, ok := this.onlineUsers[userId]
	if !ok {
		// user is offline
		err = fmt.Errorf("User %d is offline", userId)
		return
	}
	return
}
