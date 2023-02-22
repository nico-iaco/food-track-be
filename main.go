package main

import (
	"context"
	"database/sql"
	firebase "firebase.google.com/go/v4"
	"food-track-be/controller"
	"food-track-be/repository"
	"food-track-be/service"
	"github.com/gin-contrib/cors"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	_ "food-track-be/docs"
)

//	@title			Food track be API
//	@version		1.0
//	@description	This is a sample server celler server.

//	@contact.name	Nicola Iacovelli
//	@contact.email	nicolaiacovelli98@gmail.com

//	@host		localhost:8080
//	@BasePath	/api/meal

//	@externalDocs.description	OpenAPI
//	@externalDocs.url			https://swagger.io/resources/open-api/
func main() {

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbTimeout, _ := strconv.ParseInt(os.Getenv("DB_TIMEOUT"), 10, 64)

	pgconn := pgdriver.NewConnector(
		pgdriver.WithAddr(dbHost+":"+dbPort),
		pgdriver.WithUser(dbUser),
		pgdriver.WithPassword(dbPassword),
		pgdriver.WithDatabase(dbName),
	)
	sqldb := sql.OpenDB(pgconn)

	db := bun.NewDB(sqldb, pgdialect.New())
	db.SetConnMaxIdleTime(time.Duration(dbTimeout) * time.Second)

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

	app, err := firebase.NewApp(context.Background(), nil)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}

	mr := repository.NewMealRepository(*db)
	fcr := repository.NewFoodConsumptionRepository(*db)
	gs := service.NewGroceryService()
	fcs := service.NewFoodConsumptionService(fcr, gs)
	ms := service.NewMealService(mr, fcs)
	mc := controller.NewMealController(ms, app)
	fcc := controller.NewFoodConsumptionController(fcs, app)

	r := gin.Default()
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowHeaders = append(corsConfig.AllowHeaders, "Authorization")
	//corsConfig.AllowHeaders = append(corsConfig.AllowHeaders, "iv-user")
	r.Use(cors.New(corsConfig))

	mealApi := r.Group("/api/meal")
	{
		mealApi.GET("/", mc.FindAllMeals)
		mealApi.GET(":mealId/", mc.FindMealById)
		mealApi.POST("/", mc.CreateMeal)
		mealApi.PATCH(":mealId/", mc.UpdateMeal)
		mealApi.DELETE(":mealId/", mc.DeleteMeal)
		mealApi.GET("/statistics/", mc.GetMealStatistics)

		mealApi.GET(":mealId/consumption/", fcc.FindAllConsumptionForMeal)
		mealApi.POST(":mealId/consumption/", fcc.AddFoodConsumption)
		mealApi.PATCH(":mealId/consumption/:consumptionId/", fcc.UpdateFoodConsumption)
		mealApi.DELETE(":mealId/consumption/:foodConsumptionId/", fcc.DeleteFoodConsumption)
	}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
