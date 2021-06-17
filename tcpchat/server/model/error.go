package model

import (
	"errors"
)

var (
	ERROR_USER_NOTEXISTS = errors.New("User not exists")
	ERROR_USER_EXISTS    = errors.New("User already registered, please try another user id")
	ERROR_USER_PWD       = errors.New("Password incorrect")
)
