package tool

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
)

func InitSession(engine *gin.Engine) {
	config := GetConfig().RedisConfig
	store, err := redis.NewStore(10, "tcp", config.Addr+":"+config.Port, "", []byte("secret"))
	if err != nil {
		fmt.Println(err.Error())
	}

	engine.Use(sessions.Sessions("mysession", store))
}

func SetSession(context *gin.Context, key, value interface{}) error {
	session := sessions.Default(context)
	if nil == session {
		return nil
	}
	session.Set(key, value)

	return session.Save()
}

func GetSession(context *gin.Context, key interface{}) interface{} {
	session := sessions.Default(context)
	if nil == session {
		return nil
	}

	return session.Get(key)
}
