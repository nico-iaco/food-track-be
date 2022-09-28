package main

import (
	"database/sql"
	"food-track-be/controller"
	"food-track-be/repository"
	"food-track-be/service"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {

	dsn := os.ExpandEnv("postgres://$DB_USER:$DB_PASSWORD@$DB_HOST:$DB_PORT/$DB_NAME")
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	db := bun.NewDB(sqldb, pgdialect.New())

	mr := repository.NewMealRepository(*db)
	fcr := repository.NewFoodConsumptionRepository(*db)
	ms := service.NewMealService(mr)
	gs := service.NewGroceryService()
	fcs := service.NewFoodConsumptionService(fcr, gs)
	mc := controller.NewMealController(ms)
	fcc := controller.NewFoodConsumptionController(fcs)
	r := gin.Default()

	r.GET("/meals", mc.FindAllMeals)
	r.GET("/meals/:id", mc.FindMealById)
	r.POST("/meals", mc.CreateMeal)
	r.PATCH("/meals", mc.UpdateMeal)
	r.DELETE("/meals/:id", mc.DeleteMeal)

	r.GET("/meals/:mealsId/consumption", fcc.FindAllConsumptionForMeal)
	r.POST("/meals/:mealsId/consumption", fcc.AddFoodConsumption)
	r.PATCH("/meals/:mealsId/consumption/:consumptionId", fcc.UpdateFoodConsumption)
	r.DELETE("/meals/:mealsId/consumption/:foodConsumptionId", fcc.DeleteFoodConsumption)

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
