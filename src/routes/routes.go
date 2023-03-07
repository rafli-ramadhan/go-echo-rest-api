package routes

import (
	//"fmt"
	"os"

	echo "github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/joho/godotenv"
	"go-rest-api/docs"
	"go-rest-api/src/connection"
	"gorm.io/gorm"

	authController "go-rest-api/src/controller/v1/auth"
	accountController "go-rest-api/src/controller/v1/account"
	attendanceController "go-rest-api/src/controller/v1/attendance"
	locationController "go-rest-api/src/controller/v1/location"

	accountRepository "go-rest-api/src/repository/v1/account"
	attendanceRepository "go-rest-api/src/repository/v1/attendance"
	locationRepository "go-rest-api/src/repository/v1/location"

	accountService "go-rest-api/src/service/v1/account"
	attendanceService "go-rest-api/src/service/v1/attendance"
	locationService "go-rest-api/src/service/v1/location"

	echoSwagger "github.com/swaggo/echo-swagger"
)

var master *gorm.DB
var router = echo.New()

type DB struct {
	Master *gorm.DB
}

func Run() {	
	godotenv.Load()
	RouterSetup()
	router.Logger.Fatal(router.Start(":5000"))
	// router.Start(fmt.Sprintf(":%s", os.Getenv("SERVER_PORT")))
}

func RouterSetup() {
	// set up
	router.Use(middleware.Logger())
	router.Use(middleware.CORS())

	// swagger
	docs.SwaggerInfo.Title = "Phincon Attendance App Rest API"
	docs.SwaggerInfo.Description = "Phincon Attendance App Rest API"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = os.Getenv("SWAGGER_HOST")
	docs.SwaggerInfo.Schemes = []string{"http", "https"}
	router.GET("/swagger/*any", func(c echo.Context) error {
		echoSwagger.WrapHandler(c)
		return nil
	})

	// database connection (type *gorm.DB)
	master = connection.DBMaster()

	// repository
	accountRepo := accountRepository.NewRepository(connection.DB{
		Master: master,
	})
	attendanceRepo := attendanceRepository.NewRepository(connection.DB{
		Master: master,
	})
	locationRepo := locationRepository.NewRepository(connection.DB{
		Master: master,
	})

	// service
	accountSvc := accountService.NewService(accountRepo)
	locationSvc := locationService.NewService(locationRepo)
	attendanceSvc := attendanceService.NewService(attendanceRepo, accountSvc, locationSvc)
	
	// controller
	authController := authController.NewController(accountSvc)
	accountController := accountController.NewController(accountSvc)
	attendanceController := attendanceController.NewController(attendanceSvc)
	locationController := locationController.NewController(locationSvc)

	// endpoint v1
	v1 := router.Group("/v1")

	auth := v1.Group("/auth")
	auth.POST("", func(c echo.Context) error {
		authController.Login(c)
		return nil
	})
	auth.PATCH("/forgot", func(c echo.Context) error {
		authController.ForgotPassword(c)
		return nil
	})

	accounts := v1.Group("/accounts")
	accounts.GET("", func(c echo.Context) error {
		accountController.Get(c)
		return nil
	})
	accounts.POST("/register", func(c echo.Context) error {
		accountController.Register(c)
		return nil
	})
	accounts.PATCH("", func(c echo.Context) error {
		accountController.Update(c)
		return nil
	})
	accounts.DELETE("", func(c echo.Context) error {
		accountController.Delete(c)
		return nil
	})

	attendance := v1.Group("/attendance")
	attendance.GET("/history", func(c echo.Context) error {
		attendanceController.Get(c)
		return nil
	})
	attendance.GET("/locations",  func(c echo.Context) error {
		attendanceController.GetByLocation(c)
		return nil
	})
	attendance.POST("", func(c echo.Context) error {
		attendanceController.Add(c)
		return nil
	})

	location := v1.Group("/locations")
	location.GET("", func(c echo.Context) error {
		locationController.Get(c)
		return nil
	})
	location.POST("", func(c echo.Context) error {
		locationController.Create(c)
		return nil
	})
	location.PATCH("", func(c echo.Context) error {
		locationController.Update(c)
		return nil
	})
	location.DELETE("", func(c echo.Context) error {
		locationController.Delete(c)
		return nil
	})

	// endpoint v2

	// endpoint v3

}
