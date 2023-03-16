package main

import (
	"log"

	"github.com/golang-jwt/jwt/v4"
	"github.com/innocentkithinji/xmtest/config"
	"github.com/innocentkithinji/xmtest/controller"
	"github.com/innocentkithinji/xmtest/entity"
	"github.com/innocentkithinji/xmtest/repository"
	"github.com/innocentkithinji/xmtest/server"
	"github.com/innocentkithinji/xmtest/service"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
)

func init() {
	log.Println("Initializing Service")
	config.InitializeConfig()
}

func main() {

	// Configure middleware with the custom claims type
	authConfig := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(entity.JWTClaims)
		},
		SigningKey: []byte("secret"),
	}

	DBURI := viper.Get("MONGODB_URI").(string)
	KAFKAURI := viper.Get("KAFKA_URI").(string)

	CompanyRepo := repository.NewCompanyRepo(DBURI, KAFKAURI)
	companyService := service.NewCompanyService(CompanyRepo)
	companController := controller.NewCompanyController(companyService)
	echoServer := server.NewEchoServer()

	companyRoute := echoServer.AddGroup("/companies")

	companyRoute.POST("/", companController.Create, echojwt.WithConfig(authConfig))
	companyRoute.GET("/:uid", companController.Retrieve)
	companyRoute.PATCH("/:uid", companController.Update, echojwt.WithConfig(authConfig))
	companyRoute.DELETE("/:uid", companController.Delete, echojwt.WithConfig(authConfig))

	userRepo := repository.NewUsersRepo(DBURI)
	userService := service.NewUsersService(userRepo)
	userController := controller.NewUserController(userService)

	authRoute := echoServer.AddGroup("/auth")

	authRoute.POST("/register", userController.Register)
	authRoute.POST("/login", userController.Login)
	authRoute.PATCH("/:uid", userController.Update, echojwt.WithConfig(authConfig))
	authRoute.GET("/:uid", userController.Retrieve)

	port := viper.Get("port").(string)
	echoServer.Serve(port)
}
