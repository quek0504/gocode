package main

import (
	"fmt"

	"github.com/gomodule/redigo/redis"
)

func main() {
	conn, err := redis.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		fmt.Println("Redis connection error: ", err)
		return
	}

	defer conn.Close()

	_, err = conn.Do("Set", "key1", 998)
	if err != nil {
		fmt.Println(err)
		return
	}
	result1, err := redis.Int(conn.Do("Get", "key1"))
	if err != nil {
		fmt.Println("get key1 failed", err)
		return
	}
	fmt.Println("key1=", result1)

	_, err = conn.Do("HSet", "user01", "name", "tom")
	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = conn.Do("HSet", "user01", "age", "18")
	if err != nil {
		fmt.Println(err)
		return
	}

	name, err := redis.String(conn.Do("HGet", "user01", "name"))
	if err != nil {
		fmt.Println("hset err= ", err)
		return
	}
	fmt.Println("name=", name)

	age, err := redis.Int(conn.Do("HGet", "user01", "age"))
	if err != nil {
		fmt.Println("hset err= ", err)
		return
	}
	fmt.Println("age=", age)

	_, err = conn.Do("HMSet", "user02", "name", "john", "age", 19)
	if err != nil {
		fmt.Println(err)
		return
	}

	// HMGet example
	result2, err := redis.Strings(conn.Do("HMGet", "user02", "name", "age"))
	if err != nil {
		fmt.Println("hmset err= ", err)
		return
	}

	for i, v := range result2 {
		fmt.Printf("result2[%d]=%s\n", i, v)
	}

}
