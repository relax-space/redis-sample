package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/garyburd/redigo/redis"
)

var (
	redisContext *RedisContext
	redisUrl     = flag.String("REDIS_CONN", os.Getenv("REDIS_CONN"), "REDIS_CONN")
)

func main() {

	redisContext, err := NewRedisContext(*redisUrl)
	if err != nil {
		log.Panic(err)
	}
	redisConn := redisContext.Conn()
	defer redisConn.Close()

	RedisSet(redisConn)

	openId, err := RedisGet(redisConn)
	fmt.Println(openId, err)
}

func RedisGet(conn redis.Conn) (openId string, err error) {
	openIdObj, err := redis.String(conn.Do("GET", "WXMall_Sale_Service:Test:xiao"))
	fmt.Printf("%T,%v", openIdObj, openIdObj)
	if err != nil {
		return
	}
	return
}

func RedisSet(conn redis.Conn) (openId string, err error) {
	openIdObj, err := redis.String(conn.Do("SET", "WXMall_Sale_Service:Test:xiao", "xinmiao"))
	fmt.Printf("%T,%v", openIdObj, openIdObj)
	if err != nil {
		return
	}
	return
}
