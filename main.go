package main

import (
	"log"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorhill/cronexpr"
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
		command := strings.Split(job.Command, " ")
		cmd := exec.Command(command[0], command[1:]...)
		output, err := cmd.CombinedOutput()
		if err != nil {
			log.Printf("Error executing command '%s': %s\n", job.Command, err)
		} else {
			log.Printf("Output of command '%s': %s\n", job.Command, output)
		}
		expr, err := cronexpr.Parse(job.Schedule)
		if err != nil {
			log.Print("Error parsing cron expression for job id = ", job.JobId)
		}
		nextTime := expr.Next(time.Now())
		db.Store.UpdateTimeById("job_id", &job, nextTime, &modals.CronJob{})
	}
}
