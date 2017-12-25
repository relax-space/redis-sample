package main

import (
	"errors"
	"net/url"
	"time"

	"fmt"

	"github.com/garyburd/redigo/redis"
)

var (
	errorInvalidScheme = errors.New("invalid Redis database URI scheme")
	pool               *redis.Pool
	redisWaitingTime   = 1000 * time.Millisecond
	redisMaxIdle       = 5
)

type RedisContext struct {
	Pool *redis.Pool
}

func NewRedisContext(uri string) (*RedisContext, error) {
	c := &RedisContext{
		Pool: &redis.Pool{
			MaxIdle:     redisMaxIdle,
			IdleTimeout: 240 * time.Second,
			Dial: func() (redis.Conn, error) {
				return redisConnFromUri(uri)
			},
			TestOnBorrow: func(c redis.Conn, t time.Time) error {
				_, err := c.Do("PING")
				return err
			},
		},
	}

	redisConn := c.Conn()
	defer redisConn.Close()

	pong, err := redis.String(redisConn.Do("PING", "pong"))
	if err != nil {
		//log.Errorf("Fail to send ping message to redis.")
		return nil, err
	}
	//log.Debugln(pong)
	fmt.Printf(pong)
	return c, nil
}

func (c *RedisContext) Conn() redis.Conn {
	return c.Pool.Get()
}

func redisConnFromUri(uriString string) (redis.Conn, error) {
	uri, err := url.Parse(uriString)
	if err != nil {
		return nil, err
	}

	var network string
	var host string
	var password string
	var db string

	switch uri.Scheme {
	case "redis":
		network = "tcp"
		host = uri.Host
		if uri.User != nil {
			password, _ = uri.User.Password()
		}
		if len(uri.Path) > 1 {
			db = uri.Path[1:]
		}
	case "unix":
		network = "unix"
		host = uri.Path
	default:
		return nil, errorInvalidScheme
	}

	conn, err := redis.Dial(network, host)
	if err != nil {
		return nil, err
	}

	if password != "" {
		_, err := conn.Do("AUTH", password)
		if err != nil {
			conn.Close()
			return nil, err
		}
	}

	if db != "" {
		_, err := conn.Do("SELECT", db)
		if err != nil {
			conn.Close()
			return nil, err
		}
	}

	return conn, nil
}
