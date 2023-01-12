package main

import (
	"fmt"
	"strconv"
	"sync"
	"time"

	//"github.com/garyburd/redigo/redis"

	"context"
	clusterredis "github.com/go-redis/redis/v9"
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

	c.Close()
}

func test_go_redis_cluster(redisAddr string) {
	fmt.Println("\r\ntest cluster redis set get", redisAddr)
	var ctx = context.Background()
	rdb := clusterredis.NewClusterClient(&clusterredis.ClusterOptions{
		Addrs: []string{redisAddr},

		OnConnect: func(ctx context.Context, conn *clusterredis.Conn) error {
			fmt.Printf("~~~~~>>>>>>>>>>>>>>>>>new conn from pool, conn=%v\n", conn)
			return nil
		},

		MinIdleConns: 1,
		MaxIdleConns: 2,
		PoolSize:3,
	})

	var wg sync.WaitGroup

	for i := 0; i < 100; i ++ {
		wg.Add(1)
		idx := i
		go func() {
			info_cluster, _ := rdb.Info(ctx,"Cluster").Result()
			fmt.Print("info cluster return:", info_cluster)
			var set_val = "cluster value test~~~~" + strconv.Itoa(idx)
			var set_key = "go_redis_key" + strconv.Itoa(idx)
			for j := 0; j < 100; j ++ {
				fmt.Println("<<<<<<<<<<<<<<<<pool status:", rdb.PoolStats())

				err := rdb.Set(ctx, set_key, set_val, 0).Err()
				fmt.Println("set err", err)
				val, err := rdb.Get(ctx, set_key).Result()
				fmt.Println("get return", val, err)

				if set_val != val {
					panic("get set not match")
				}

				ret_intcmd := rdb.HSet(ctx,"go_redis_hset_key", "field", "field_val_cluster" + redisAddr)
				fmt.Println("hset return", ret_intcmd, " hset err:", ret_intcmd.Err())
				hgetval, err := rdb.HGet(ctx, "go_redis_hset_key", "field").Result()
				fmt.Println("hget return", hgetval, err)
			}
			wg.Done()
		}()
	}

	fmt.Println("<<<<<<<<<<<<<<<<pool status:", rdb.PoolStats())
	wg.Wait()

	for i := 0; i < 1000; i ++ {
		time.Sleep(time.Second * 3)
		fmt.Println("<<<<<<<<<<<<<<<<pool status:", rdb.PoolStats())
	}

	rdb.Close()
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

	rdb.Close()
}

func test_redis_scan(redisAddr string) {
	fmt.Println("\r\ntest cluster redis set get", redisAddr)
	var ctx = context.Background()
	rdb := clusterredis.NewClusterClient(&clusterredis.ClusterOptions{
		Addrs: []string{redisAddr},

		OnConnect: func(ctx context.Context, conn *clusterredis.Conn) error {
			fmt.Printf("~~~~~>>>>>>>>>>>>>>>>>new conn from pool, conn=%v\n", conn)
			return nil
		},

		MinIdleConns: 1,
		MaxIdleConns: 2,
		PoolSize:3,
	})

	keycout := 0
	fn := func(ctx context.Context, client *clusterredis.Client) error {
		var cursor uint64 = 0
		for	 {
			scan_cmd := client.Scan(ctx, cursor, "go_redis_key*", 10)

			keys, cursor_ret, err := scan_cmd.Result()

			fmt.Println(cursor_ret, err)
			fmt.Println(keys)

			keycout += len(keys)

			cursor = cursor_ret
			if cursor==0 {
				break
			}
		}
		return nil
	}
	rdb.ForEachMaster(ctx, fn)
	fmt.Println("redis cluster scan key count", keycout)

	rdb.Close()
}

func main() {
	ip := "192.168.15.146"

	for {
		test_redis_scan(ip + ":7001")

		test_redigo(ip + ":7001")
		test_redigo(ip + ":7002")
		test_redigo(ip + ":7003")
		test_redigo(ip + ":7004")
		test_redigo(ip + ":7005")
		test_redigo(ip + ":7006")

		test_go_redis_standalone("192.168.3.60:6379")
		test_go_redis_cluster(ip + ":7002")
	}

}
