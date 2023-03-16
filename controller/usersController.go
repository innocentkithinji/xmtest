package controller

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/innocentkithinji/xmtest/entity"
	"github.com/innocentkithinji/xmtest/errors"
	"github.com/innocentkithinji/xmtest/service"
	"github.com/labstack/echo/v4"
)

type (
	UserController interface {
		Register(c echo.Context) error
		Login(c echo.Context) error
		Update(c echo.Context) error
		Retrieve(c echo.Context) error
	}

	userController struct {
		userService service.UserSvc
	}
)

func (u userController) Register(c echo.Context) error {
	var user = new(entity.User)
	if err := c.Bind(user); err != nil {
		log.Printf("Failed to Bind user data: %s", err)
		return c.JSON(
			http.StatusBadRequest,
			errors.ServiceError{Message: "Error unpacking the json body"})
	}
	err := c.Validate(user)
	if err != nil {
		log.Printf("Invalid data passed: %v", err.(validator.ValidationErrors))
		errs := map[string]string{}
		for _, e := range err.(validator.ValidationErrors) {
			key := e.Namespace()
			er := fmt.Sprintf("Error:Field validation for '%s' failed on the '%s' tag", e.Field(), e.Tag())

			errs[key] = er
		}
		return c.JSON(
			http.StatusBadRequest,
			errs)
	}

	userCreds, err := u.userService.Register(user)
	if err != nil {
		log.Printf("Failed to register user: %s", err)
		return c.JSON(
			http.StatusInternalServerError,
			errors.ServiceError{Message: err.Error()})
	}

	return c.JSON(http.StatusCreated, userCreds)

}

func (u userController) Login(c echo.Context) error {
	var user = new(entity.User)
	if err := c.Bind(user); err != nil {
		log.Printf("Failed to Bind user data: %s", err)
		return c.JSON(
			http.StatusBadRequest,
			errors.ServiceError{Message: "Error unpacking the json body"})
	}
	err := c.Validate(user)
	if err != nil {
		log.Printf("Invalid data passed: %v", err.(validator.ValidationErrors))
		errs := map[string]string{}
		for _, e := range err.(validator.ValidationErrors) {
			key := e.Namespace()
			er := fmt.Sprintf("Error:Field validation for '%s' failed on the '%s' tag", e.Field(), e.Tag())

			errs[key] = er
		}
		return c.JSON(
			http.StatusBadRequest,
			errs)
	}

	userCreds, err := u.userService.Login(user)
	if err != nil {
		log.Printf("Failed to register user: %s", err)
		return c.JSON(
			http.StatusInternalServerError,
			errors.ServiceError{Message: err.Error()})
	}

	return c.JSON(http.StatusCreated, userCreds)
}

func (u userController) Update(c echo.Context) error {
	uid := c.Param("uid")
	var user = new(entity.User)
	if err := c.Bind(user); err != nil {
		log.Printf("Failed to Bind Request data to user data: %s", err)
		return c.JSON(
			http.StatusBadRequest,
			errors.ServiceError{Message: "Error unpacking the json body"})
	}

	userInfo, err := u.userService.Update(uid, user)
	if err != nil {
		log.Printf("Failed to Update User Info: %v", err)
		return c.JSON(
			http.StatusInternalServerError,
			errors.ServiceError{Message: "Error Updating User Info"})
	}

	return c.JSON(http.StatusAccepted, userInfo)
}

func (u userController) Retrieve(c echo.Context) error {
	uid := c.Param("uid")
	user, err := u.userService.Retrieve(uid)
	if err != nil {
		log.Printf("Failed to retrieve the given user: %s", err)
		return c.JSON(
			http.StatusNotFound,
			errors.ServiceError{Message: "User was not found"})
	}

	return c.JSON(http.StatusOK, user)
}

func NewUserController(s service.UserSvc) UserController {
	return &userController{userService: s}
}
