package main

import (
	"cloudRestaurant/controller"
	"cloudRestaurant/tool"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"strings"
)

func main() {
	config, err := tool.ParseConfig("./config/app.json")
	if err != nil {
		panic(err.Error())
	}
	if _, err := tool.OrmEngine(config); err != nil {
		panic(err.Error())
		return
	}

	tool.InitRedisStore()

	app := gin.Default()
	tool.InitSession(app)
	registerRoute(app)
	// 设置全局跨域访问
	app.Use(Cors())
	if err := app.Run(config.AppHost + ":" + string(config.AppPort)); err != nil {
		log.Fatal(err.Error())
	}
}

func registerRoute(engine *gin.Engine) {
	new(controller.HelloController).Router(engine)
	new(controller.MemberController).Router(engine)
}

func Cors() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		method := ctx.Request.Method
		origin := ctx.Request.Header.Get("Origin")
		var headerKeys []string
		for k, _ := range ctx.Request.Header {
			headerKeys = append(headerKeys, k)
		}

		headerStr := strings.Join(headerKeys, ",")
		if headerStr != "" {
			fmt.Println(headerStr)
		}

		if origin != "" {
			ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		}

		if method == "OPTIONS" {
			ctx.JSON(200, "Options Request!")
		}

		ctx.Next()
	}
}
