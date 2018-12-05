package redis

import (
    "github.com/go-redis/redis"
    "PentagoServer/dao"
)

type redisDAO struct {
    redis *redis.Client
}

func Connect() (dao.CacheDAO, error){
    client := redis.NewClient(&redis.Options{
        Addr:     "localhost:6379",
        Password: "",
        DB:       0,
    })

    _, err := client.Ping().Result()
    if err != nil {
        panic(err)
    }
    return &redisDAO{redis: client}, nil
}
