package message

// common 'protocol' used by client and server to communicate/data transfering

// const
const (
	LoginMesType            = "LoginMes"
	LoginResMesType         = "LoginResMes"
	RegisterMesType         = "RegisterMes"
	RegisterResMesType      = "RegisterResMes"
	LogoutMesType           = "LogoutMes"
	LogoutResMesType        = "LogoutResMes"
	NotifyUserStatusMesType = "NotifyUserStatusMes"
	SmsMesType              = "SmsMes"
)

const (
	UserOnline = iota
	UserOffline
	UserBusy
	UserSnooze
)

type Message struct {
	Type string `json:"type"` // message type
	Data string `json:"data"` // message content
}

// login message
type LoginMes struct {
	UserId   int    `json:"userId"`
	UserPwd  string `json:"userPwd"`
	UserName string `json:"userName"`
}

// login response
type LoginResMes struct {
	Code    int    `json:"code"`    // status code
	Error   string `json:"error"`   // error message
	UserIds []int  `json:"userIds"` // userId slice
}

// register message
type RegisterMes struct {
	User User `json:"user"`
}

// register response
type RegisterResMes struct {
	Code  int    `json:"code"`  // status code
	Error string `json:"error"` // error message
}

// logout message
type LogoutMes struct {
	UserId int `json:"userId"`
}

// logout response mesage
type LogoutResMes struct {
	Code  int    `json:"code"`  // status code
	Error string `json:"error"` // error message
}

// used by server to forward or push message
type NotifyUserStatusMes struct {
	UserId int `json:"userId"`
	Status int `json:"status"` // online/offline
}

// used by client to send message
type SmsMes struct {
	Content string `json:"content"`
	User           // User
}
