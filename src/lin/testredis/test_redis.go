package main

import (
	"fmt"

	//"github.com/garyburd/redigo/redis"

	"context"
	clusterredis "github.com/go-redis/redis/v8"
	"github.com/gomodule/redigo/redis"
)

func test_redigo(redisAddr string) {
	c, err := redis.Dial("tcp", redisAddr)
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

	c.Do("SET", "go_redis_key", "go_redis_val_" + redisAddr)
	val, err := redis.String(c.Do("GET", "go_redis_key"))
	fmt.Println("redis key:", val, err)

	c.Do("HSET", "go_redis_hset_key", "field", "field_val_" + redisAddr)
	val, err = redis.String(c.Do("HGET", "go_redis_hset_key", "field"))
	fmt.Println("redis mkey:", val, err)
}

func test_go_redis_cluster(redisAddr string) {
	fmt.Println("\r\ntest cluster redis set get", redisAddr)
	var ctx = context.Background()
	rdb := clusterredis.NewClusterClient(&clusterredis.ClusterOptions{
		Addrs: []string{redisAddr},
	})

	info_cluster, _ := rdb.Info(ctx,"Cluster").Result()
	fmt.Print("info cluster return:", info_cluster)

	err := rdb.Set(ctx, "go_redis_key", "cluster value test~~~~", 0).Err()
	fmt.Println("set err", err)
	val, err := rdb.Get(ctx, "go_redis_key").Result()
	fmt.Println("get return", val, err)

	ret_intcmd := rdb.HSet(ctx,"go_redis_hset_key", "field", "field_val_cluster" + redisAddr)
	fmt.Println("hset return", ret_intcmd, " hset err:", ret_intcmd.Err())
	hgetval, err := rdb.HGet(ctx, "go_redis_hset_key", "field").Result()
	fmt.Println("hget return", hgetval, err)
}

func test_go_redis_standalone(redisAddr string) {
	fmt.Println("\r\ntest standalone redis set get", redisAddr)
	var ctx = context.Background()
	rdb := clusterredis.NewClient(&clusterredis.Options{
		Addr: redisAddr,
	})

	info_stringcmd := rdb.Info(ctx, "Cluster")
	fmt.Println("info:", info_stringcmd,
		"fullname:",info_stringcmd.FullName())
	fmt.Println(info_stringcmd.Result())

	iscluster_mode, cmd_err := info_stringcmd.Bool()
	fmt.Print(iscluster_mode, cmd_err)

	err := rdb.Set(ctx, "go_redis_key", "standalone value test~~~~", 0).Err()
	fmt.Println("standalone set err", err)

	val, err := rdb.Get(ctx, "go_redis_key").Result()
	fmt.Println("standalone get return", val, err)

	ret_intcmd := rdb.HSet(ctx,"go_redis_hset_key", "field", "field_val_standalone" + redisAddr)
	fmt.Println("standalone hset return", ret_intcmd, " hset err:", ret_intcmd.Err())
	hgetval, err := rdb.HGet(ctx, "go_redis_hset_key", "field").Result()
	fmt.Println("standalone hget return", hgetval, err)
}

func main() {
	ip := "192.168.15.146"
	test_redigo(ip + ":7001")
	test_redigo(ip + ":7002")
	test_redigo(ip + ":7003")
	test_redigo(ip + ":7004")
	test_redigo(ip + ":7005")
	test_redigo(ip + ":7006")

	test_go_redis_cluster(ip + ":7002")

	test_go_redis_standalone("192.168.3.60:6379")
}
