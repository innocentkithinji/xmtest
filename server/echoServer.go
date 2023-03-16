package server

import (
	"fmt"
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/viper"
)

type EchoServer interface {
	Server
	AddGroup(path string) *echo.Group
}

type (
	s struct {
		echo *echo.Echo
	}

	echoValidator struct {
		validator *validator.Validate
	}
)

func (eV echoValidator) Validate(i interface{}) error {
	if err := eV.validator.Struct(i); err != nil {
		return err
	}
	return nil
}

func (s s) AddGroup(path string) *echo.Group {
	group := s.echo.Group(path)

	return group
}

func (s s) Serve(port string) {
	s.echo.Use(middleware.Recover())
	s.echo.Use(middleware.Logger())
	serviceName := viper.Get("service_name").(string)
	env := viper.Get("environment").(string)
	loggingPrefix := fmt.Sprintf("%s:%s", serviceName, env)
	s.echo.Logger.SetPrefix(loggingPrefix)

	startPort := fmt.Sprintf(":%s", port)
	log.Fatal(s.echo.Start(startPort))

}

func NewEchoServer() EchoServer {
	e := echo.New()
	e.Validator = &echoValidator{validator: validator.New()}
	return &s{echo: e}
}
