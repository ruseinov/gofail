package main

import (
	"github.com/garyburd/redigo/internal"
	"github.com/garyburd/redigo/redis"

	"fmt"
	"time"
)

const (
	idleTimeout = 240
)

type Redis struct {
	*redis.Pool
}

func InitRedis(server string) *Redis {
	return &Redis{initRedisPool(server)}
}

func initRedisPool(server string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: idleTimeout * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}

func (r *Redis) SetValue(key string, value string, expiration int64) error {
	conn := r.Get()
	defer conn.Close()

	conn.Send("SET", key, value)
	conn.Send("EXPIRE", key, expiration)
	err := conn.Flush()

	return err
}

func (r *Redis) GetValue(key string) (string, error) {
	conn := r.Get()
	defer conn.Close()

	return redis.String(conn.Do("GET", key))
}

func (r *Redis) Del(key string) error {
	conn := r.Get()
	defer conn.Close()

	_, err := conn.Do("DEL", key)

	return err
}

func main() {
	redis := InitRedis("blah")
	redis.Del("das")
	fmt.Println(internal.WatchState)
}
