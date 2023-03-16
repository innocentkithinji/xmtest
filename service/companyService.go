package service

import (
	"errors"
	"log"
	"strings"

	"github.com/google/uuid"
	"github.com/innocentkithinji/xmtest/entity"
	"github.com/innocentkithinji/xmtest/repository"
)

type companyService struct {
	companyRepo repository.Repository
}

func (c companyService) ValidateName(companyName string) error {
	if companyName != "" {
		nameIdentifier := getNameIdentifier(companyName)
		filter := map[string]interface{}{
			"nameidentifier": nameIdentifier,
		}

		company, _ := c.companyRepo.Filter(filter)
		if company != nil {
			log.Println("Company with Similar Name Found")
			return errors.New("company with similar name was found")
		}
	}
	return nil
}

func (c companyService) Create(company *entity.Company) (*entity.Company, error) {
	if err := c.ValidateName(company.Name); err != nil {
		return nil, err
	}

	company.NameIdentifier = getNameIdentifier(company.Name)

	company.ID = uuid.New().String()
	return c.companyRepo.Create(company)
}

func (c companyService) Retrieve(uid string) (*entity.Company, error) {
	return c.companyRepo.Get(uid)
}

func (c companyService) Update(uid string, update *entity.Company, userUID string) (*entity.Company, error) {
	company, err := c.companyRepo.Get(uid)
	if err != nil {
		log.Printf("Could not find company with ID: %s", uid)
		return nil, err
	}

	if company.OwnerId != userUID {
		log.Printf("Unauthorised updated for the user")
		return nil, errors.New("Unauthorised users")
	}
	if update.Name != "" {
		if err := c.ValidateName(update.Name); err != nil {
			return nil, err
		}
		company.Name = update.Name
	}
	if update.Description != "" {
		company.Description = update.Description
	}
	if update.Employees != 0 {
		company.Employees = update.Employees
	}
	if update.Registered != company.Registered {
		company.Registered = update.Registered
	}
	if update.Type != company.Type {
		company.Type = update.Type
	}

	return c.companyRepo.Update(company)
}

func (c companyService) Delete(uid string, userUID string) error {
	company, err := c.companyRepo.Get(uid)
	if err != nil {
		log.Printf("Could not find company with ID: %s", uid)
		return err
	}
	if company.OwnerId != userUID {
		log.Printf("Unauthorised updated for the user")
		return errors.New("Unauthorised users")
	}
	return c.companyRepo.Delete(uid)
}

func NewCompanyService(companyRepo repository.Repository) Service {
	return companyService{companyRepo: companyRepo}
}

func getNameIdentifier(name string) string {
	return strings.ToLower(strings.ReplaceAll(name, " ", ""))
}
