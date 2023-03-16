package controller

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v4"
	"github.com/innocentkithinji/xmtest/entity"
	"github.com/innocentkithinji/xmtest/errors"
	"github.com/innocentkithinji/xmtest/service"
	"github.com/labstack/echo/v4"
)

type companyController struct {
	companyService service.Service
}

func (cc companyController) Create(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*entity.JWTClaims)

	var company = new(entity.Company)
	if err := c.Bind(company); err != nil {
		log.Printf("Failed to Bind company data: %s", err)
		return c.JSON(
			http.StatusBadRequest,
			errors.ServiceError{Message: "Error unpacking the json body"})
	}
	err := c.Validate(company)
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
	company.OwnerId = claims.UID
	newCompany, err := cc.companyService.Create(company)
	if err != nil {
		log.Printf("Failed to create company: %s", err)
		return c.JSON(
			http.StatusInternalServerError,
			errors.ServiceError{Message: err.Error()})
	}

	return c.JSON(http.StatusCreated, newCompany)
}

func (cc companyController) Retrieve(c echo.Context) error {
	uid := c.Param("uid")
	company, err := cc.companyService.Retrieve(uid)
	if err != nil {
		log.Printf("Failed to retrieve the given company: %s", err)
		return c.JSON(
			http.StatusNotFound,
			errors.ServiceError{Message: "Company was not found"})
	}

	return c.JSON(http.StatusOK, company)
}

func (cc companyController) Update(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*entity.JWTClaims)

	uid := c.Param("uid")
	var company = new(entity.Company)
	if err := c.Bind(company); err != nil {
		log.Printf("Failed to Bind Request data to user data: %s", err)
		return c.JSON(
			http.StatusBadRequest,
			errors.ServiceError{Message: "Error unpacking the json body"})
	}

	companyInfo, err := cc.companyService.Update(uid, company, claims.UID)
	if err != nil {
		log.Printf("Failed to Update Company Info: %v", err)
		return c.JSON(
			http.StatusInternalServerError,
			errors.ServiceError{Message: "Error Updating Company Info"})
	}

	return c.JSON(http.StatusAccepted, companyInfo)
}

func (cc companyController) Delete(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*entity.JWTClaims)

	uid := c.Param("uid")
	err := cc.companyService.Delete(uid, claims.UID)
	if err != nil {
		log.Printf("Failed to Delete company: %v", err)
		return c.JSON(
			http.StatusInternalServerError,
			errors.ServiceError{Message: "Error Deleting company info"})
	}

	return c.NoContent(http.StatusOK)
}

func NewCompanyController(s service.Service) Controller {
	return &companyController{companyService: s}
}
