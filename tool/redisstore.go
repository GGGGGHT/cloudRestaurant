package tool

import (
	"fmt"
	"github.com/go-redis/redis"
	_ "github.com/go-redis/redis"
	"github.com/mojocn/base64Captcha"
	"time"
)

type RedisStore struct {
	client *redis.Client
	store  *base64Captcha.Store
}

var Rstore *RedisStore

func InitRedisStore() *RedisStore {
	config := GetConfig().RedisConfig
	client := redis.NewClient(&redis.Options{
		Addr:     config.Addr + ":" + config.Port,
		Password: config.Password,
		DB:       config.Db,
	})
	Rstore = &RedisStore{client: client}

	return Rstore
}

func (rs *RedisStore) Set(id string, value string) {
	if err := rs.client.Set(id, value, time.Minute*3).Err(); err != nil {
		fmt.Printf("%v\n", err)
	}
}

func (rs *RedisStore) Get(id string, clear bool) string {
	res, err := rs.client.Get(id).Result()
	if err != nil {
		fmt.Println(err)
		return ""
	}

	if clear {
		err := rs.client.Del(id).Err()
		if err != nil {
			fmt.Println(err)
			return ""
		}
	}

	return res
}

func (rs *RedisStore) Verify(id, answer string, clear bool) bool {
	v := rs.Get(id, clear)
	return v == answer
}
