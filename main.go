package main

import (
	"database/sql"
	"food-track-be/controller"
	"food-track-be/repository"
	"food-track-be/service"
	"github.com/gin-contrib/cors"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {

	pgconn := pgdriver.NewConnector(
		pgdriver.WithAddr(os.Getenv("DB_HOST")+":"+os.Getenv("DB_PORT")),
		pgdriver.WithUser(os.Getenv("DB_USER")),
		pgdriver.WithPassword(os.Getenv("DB_PASSWORD")),
		pgdriver.WithDatabase(os.Getenv("DB_NAME")),
		pgdriver.WithTimeout(1*time.Second),
	)
	sqldb := sql.OpenDB(pgconn)

	db := bun.NewDB(sqldb, pgdialect.New())

	defer func(db *bun.DB) {
		err := db.Close()
		if err != nil {
			panic(err)
		}
	}(db)

	err := db.Ping()
	if err != nil {
		log.Println(err)
		panic(err)
	}

	mr := repository.NewMealRepository(*db)
	fcr := repository.NewFoodConsumptionRepository(*db)
	gs := service.NewGroceryService()
	fcs := service.NewFoodConsumptionService(fcr, gs)
	ms := service.NewMealService(mr, fcs)
	mc := controller.NewMealController(ms)
	fcc := controller.NewFoodConsumptionController(fcs)

	r := gin.Default()
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowHeaders = append(corsConfig.AllowHeaders, "iv-user")
	r.Use(cors.New(corsConfig))

	r.GET("/api/meal/", mc.FindAllMeals)
	r.GET("/api/meal/:mealId/", mc.FindMealById)
	r.POST("/api/meal/", mc.CreateMeal)
	r.PATCH("/api/meal/:mealId/", mc.UpdateMeal)
	r.DELETE("/api/meal/:mealId/", mc.DeleteMeal)
	r.GET("/api/meal/statistics/", mc.GetMealStatistics)

	r.GET("/api/meal/:mealId/consumption/", fcc.FindAllConsumptionForMeal)
	r.POST("/api/meal/:mealId/consumption/", fcc.AddFoodConsumption)
	r.PATCH("/api/meal/:mealId/consumption/:consumptionId/", fcc.UpdateFoodConsumption)
	r.DELETE("/api/meal/:mealId/consumption/:foodConsumptionId/", fcc.DeleteFoodConsumption)

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
