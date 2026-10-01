package config

import (
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

type Redis struct {
	HOST string
	PORT string
	USER string
	PWD  string
}

func NewRedis(host, port, user, password string) *Redis {
	return &Redis{
		HOST: host,
		PORT: port,
		USER: user,
		PWD:  password,
	}
}

func (r *Redis) Connect() *redis.Client {
	redis := redis.NewClient(&redis.Options{
		Username: r.USER,
		Password: r.PWD,
		Addr:     fmt.Sprintf("%s:%s", r.HOST, r.PORT),
	})

	if res, err := redis.Ping(context.Background()).Result(); err != nil {
		log.Printf("Close Redis : %s\n\n", err.Error())
	} else {
		log.Printf("%s ! Conected redis successfully\n\n", res)

	}
	return redis
}
