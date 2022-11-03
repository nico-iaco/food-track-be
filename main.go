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
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {

	//dsn := os.ExpandEnv("postgres://$DB_USER:$DB_PASSWORD@$DB_HOST:$DB_PORT/$DB_NAME")

	//log.Println("connecting to database ", dsn)
	pgconn := pgdriver.NewConnector(
		pgdriver.WithAddr(os.Getenv("DB_HOST")+":"+os.Getenv("DB_PORT")),
		pgdriver.WithUser(os.Getenv("DB_USER")),
		pgdriver.WithPassword(os.Getenv("DB_PASSWORD")),
		pgdriver.WithDatabase(os.Getenv("DB_NAME")),
		pgdriver.WithTimeout(5*time.Second),
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

	r.GET("/meal", mc.FindAllMeals)
	r.GET("/meal/:mealId", mc.FindMealById)
	r.POST("/meal", mc.CreateMeal)
	r.PATCH("/meal/:mealId", mc.UpdateMeal)
	r.DELETE("/meal/:mealId", mc.DeleteMeal)
	r.GET("/meal/statistics", mc.GetMealStatistics)

	r.GET("/meal/:mealId/consumption", fcc.FindAllConsumptionForMeal)
	r.POST("/meal/:mealId/consumption", fcc.AddFoodConsumption)
	r.PATCH("/meal/:mealId/consumption/:consumptionId", fcc.UpdateFoodConsumption)
	r.DELETE("/meal/:mealId/consumption/:foodConsumptionId", fcc.DeleteFoodConsumption)

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
