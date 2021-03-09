package main

import (
	"cloudRestaurant/controller"
	"cloudRestaurant/tool"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {

	app := gin.Default()
	registerRoute(app)
	config, err := tool.ParseConfig("./config/app.json")
	if err != nil {
		panic(err.Error())
	}

	if _, err := tool.OrmEngine(config); err != nil{
		panic(err.Error())
		return
	}

	if err := app.Run(config.AppHost + ":" + string(config.AppPort)); err != nil {
		log.Fatal(err.Error())
	}
}

func registerRoute(engine *gin.Engine) {
	new(controller.HelloController).Router(engine)
	new(controller.MemberController).Router(engine)
}
