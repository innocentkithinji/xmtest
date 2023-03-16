package service

import "github.com/innocentkithinji/xmtest/entity"

type Service interface {
	Create(company *entity.Company) (*entity.Company, error)
	Retrieve(uid string) (*entity.Company, error)
	Update(uid string, company *entity.Company, userUID string) (*entity.Company, error)
	Delete(uid string, userUID string) error
}
