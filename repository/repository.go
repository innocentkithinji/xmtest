package repository

import (
	"github.com/google/martian/log"
	"github.com/innocentkithinji/xmtest/entity"
	"go.mongodb.org/mongo-driver/bson"
)

type Repository interface {
	Create(i *entity.Company) (*entity.Company, error)
	Get(id string) (*entity.Company, error)
	Update(company *entity.Company) (*entity.Company, error)
	Filter(filters map[string]interface{}) (*entity.Company, error)
	Delete(id string) error
}

func CreateUpdateBson(i interface{}) (bson.M, error) {
	updateBsonM, err := bson.Marshal(i)
	if err != nil {
		log.Errorf("Unable to Marshal item: %s", err)
		return nil, err
	}
	updateDoc := bson.M{}
	err = bson.Unmarshal(updateBsonM, &updateDoc)
	if err != nil {
		log.Errorf("Unable to unmarshall  bson: %s", err)
		return nil, err
	}

	// Create an update operation with the $set update modifier
	update := bson.M{
		"$set": updateDoc,
	}

	return update, nil
}
