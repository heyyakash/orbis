package controllers

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorhill/cronexpr"
	"github.com/heyyakash/orbis/modals"
	"gorm.io/gorm"
)

type CronController struct {
	DB *gorm.DB
}

func GetCronController(db *gorm.DB) *CronController {
	return &CronController{
		DB: db,
	}
}

func (c *CronController) AddCronJob() gin.HandlerFunc {

	return func(ctx *gin.Context) {
		var job modals.JobRequest
		if err := ctx.BindJSON(&job); err != nil {
			log.Print(err)
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}
		expr, err := cronexpr.Parse(job.Schedule)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
		}
		nextTime := expr.Next(time.Now())
		cronjob := modals.CronJob{Command: job.Command, Schedule: job.Schedule, NextRun: nextTime}
		result := c.DB.Create(&cronjob)
		if result.Error != nil {
			log.Print(result.Error)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal server error",
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"message": "Added job with id : " + strconv.Itoa(int(cronjob.JobId)),
		})

	}
}

func (c *CronController) GetCronJob() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		idParam := ctx.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid id format",
			})
			return
		}
		job := modals.CronJob{JobId: uint(id)}
		c.DB.Where("job_id = ?", job.JobId).First(&job)
		ctx.JSON(http.StatusOK, job)

	}
}

func (c *CronController) GetAllCronJob() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		jobs := []modals.CronJob{}
		result := c.DB.Find(&jobs)
		if result.Error != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": result.Error.Error(),
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"data": jobs,
		})

	}
}

func (c *CronController) DeleteCronJob() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		idParam := ctx.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid id format",
			})
			return
		}
		job := modals.CronJob{JobId: uint(id)}
		c.DB.Where("job_id = ?", job.JobId).Delete(&job)
		ctx.JSON(http.StatusOK, gin.H{})

	}
}

func (c *CronController) DeleteAllCronJobs() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		c.DB.Unscoped().Where("1 = 1").Delete(&modals.CronJob{})
		ctx.JSON(http.StatusOK, gin.H{})
	}
}

func (c *CronController) UpdateJob() gin.HandlerFunc {
	return func(ctx *gin.Context) {

	}
}
