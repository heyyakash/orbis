package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/heyyakash/orbis/controllers"
	"github.com/heyyakash/orbis/db"
)

func CronRoutes(c *gin.Engine) {
	cronController := controllers.GetCronController(db.Store.DB)
	// Routes
	c.POST("set", cronController.AddCronJob())
	c.GET(":id", cronController.GetCronJob())
	c.GET("all", cronController.GetAllCronJob())
	c.DELETE(":id", cronController.DeleteCronJob())
	c.DELETE("all", cronController.DeleteAllCronJobs())
}
