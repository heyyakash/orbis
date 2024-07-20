package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/heyyakash/orbis/db"
)

func Init() {
	db.Init()
}

func main() {
	Init()
	r := gin.Default()
	r.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "pong",
		})
	})

	log.Print("Server Started on port 8080")
	r.Run(":8080")

}
