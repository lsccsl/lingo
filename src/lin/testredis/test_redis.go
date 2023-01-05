package main

import (
	"fmt"

	"github.com/garyburd/redigo/redis"
)

func test_redigo() {
	c, err := redis.Dial("tcp", "127.0.0.1:6379")
	fmt.Println("redis dial:", c, " err", err)
	//r, err := c.Do("AUTH", "")
	//fmt.Println("redis auth:", r, " err", err)
	if err != nil {
		fmt.Println("Connect to redis error", err)
		return
	} else {
		fmt.Println("connect to redis suc")
	}

	defer c.Close()

	c.Do("SET", "go_redis_key", "go_redis_val")
	val, err := redis.String(c.Do("GET", "go_redis_key"))
	fmt.Println("redis key:", val)

	c.Do("HSET", "go_redis_hset_key", "field", "field_val")
	val, err = redis.String(c.Do("HGET", "go_redis_hset_key", "field"))
	fmt.Println("redis mkey:", val)
}

func test_go_redis() {

}

func main() {
	test_redigo()
	test_go_redis()
}
