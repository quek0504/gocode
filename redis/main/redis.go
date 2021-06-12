package main

import (
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
)

func main() {
	conn, err := redis.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		fmt.Println("Redis connection error: ", err)
		return
	}

	// close connection
	defer conn.Close()

	// redis string example
	_, err = conn.Do("Set", "key1", 998)
	if err != nil {
		fmt.Println(err)
		return
	}

	// redis TTL example
	_, err = conn.Do("EXPIRE", "key1", 4)
	if err != nil {
		fmt.Println(err)
		return
	}
	ttl, err := redis.Int(conn.Do("TTL", "key1"))
	if err != nil {
		fmt.Println("get key1 ttl failed", err)
		return
	}
	fmt.Println("key1 ttl=", ttl)
	fmt.Println("sleeping for 4 sec...")
	time.Sleep(4 * time.Second) // key1 ttl 4 sec, try changing time to sleep
	result1, err := redis.Int(conn.Do("Get", "key1"))
	if err != nil {
		fmt.Println("get key1 failed", err)
		return
	}
	fmt.Println("key1=", result1)

	// redis hash HSet and HGet example
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

	// redis hash HMSet and HMGet example
	_, err = conn.Do("HMSet", "user02", "name", "john", "age", 19)
	if err != nil {
		fmt.Println(err)
		return
	}
	result2, err := redis.Strings(conn.Do("HMGet", "user02", "name", "age"))
	if err != nil {
		fmt.Println("hmset err= ", err)
		return
	}
	for i, v := range result2 {
		fmt.Printf("result2[%d]=%s\n", i, v)
	}

	// redis list example
	_, err = conn.Do("lpush", "animalList", "1:lion")
	if err != nil {
		fmt.Println(err)
		return
	}
	result3, err := redis.String(conn.Do("rpop", "animalList"))
	if err != nil {
		fmt.Println("rpop err=", err)
	}
	fmt.Println("result3=", result3)

}
