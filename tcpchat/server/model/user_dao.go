package model

import (
	"encoding/json"
	"fmt"
	"gocode/chatroom/common/message"

	"github.com/gomodule/redigo/redis"
)

// global var
var (
	MyUserDao *UserDao
)

type UserDao struct {
	pool *redis.Pool
}

// factory pattern to create UserDao
func NewUserDao(pool *redis.Pool) (userDao *UserDao) {
	userDao = &UserDao{
		pool: pool,
	}
	return
}

func (this *UserDao) getUserById(conn redis.Conn, id int) (user *User, err error) {
	res, err := redis.String(conn.Do("HGet", "users", id))
	if err != nil {
		if err == redis.ErrNil { // id not exist in users hash
			err = ERROR_USER_NOTEXISTS
		}
		return
	}

	// create user instance
	user = &User{}

	err = json.Unmarshal([]byte(res), user)
	if err != nil {
		fmt.Println("json.Unmarshal user err=", err)
		return
	}
	return
}

// login verification
func (this *UserDao) Login(userId int, userPwd string) (user *User, err error) {
	// get connection from redis pool
	conn := this.pool.Get()
	defer conn.Close()

	user, err = this.getUserById(conn, userId)
	if err != nil {
		return
	}

	// get user from redis, check its credential
	if user.UserPwd != userPwd {
		err = ERROR_USER_PWD
		return
	}
	return
}

func (this *UserDao) Register(user *message.User) (err error) {
	// get connection from redis pool
	conn := this.pool.Get()
	defer conn.Close()

	_, err = this.getUserById(conn, user.UserId)
	if err == nil {
		err = ERROR_USER_EXISTS
		return
	}

	// can't find given userId in redis, it is a new user, proceed to register
	data, err := json.Marshal(user)
	if err != nil {
		return
	}

	// save new user to redis
	_, err = conn.Do("HSet", "users", user.UserId, string(data))
	if err != nil {
		fmt.Println("Saving new user err=", err)
		return
	}
	return
}
