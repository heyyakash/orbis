package main

import (
	"log"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/heyyakash/orbis/db"
	"github.com/heyyakash/orbis/helpers"
	"github.com/heyyakash/orbis/modals"
	"github.com/heyyakash/orbis/routes"
)

func Init() {
	db.Init()
}

var JobsChannel = make(chan modals.CronJob)

func main() {
	Init()
	r := gin.Default()
	r.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "pong",
		})
	})

	pool, err := strconv.Atoi(helpers.GetString("POOL"))
	if err != nil {
		panic("POOL size missing")
	}

	for i := 1; i <= pool; i++ {
		go Worker(JobsChannel)
	}

	go Schedule()

	routes.CronRoutes(r)

	log.Print("Server Started on port 8080")
	r.Run(":8080")

}

func Schedule() {
	for {
		var jobs []modals.CronJob
		result := db.Store.DB.Where("next_run <= ?", time.Now()).Find(&jobs)
		if result.Error != nil {
			panic(result.Error)
		}
		for _, job := range jobs {
			JobsChannel <- job
		}
		time.Sleep(1 * time.Minute)
	}
}

func Worker(jobChannel <-chan modals.CronJob) {
	for job := range jobChannel {
		log.Print(job)
	}
}
