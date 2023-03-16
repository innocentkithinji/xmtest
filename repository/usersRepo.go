package repository

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/innocentkithinji/xmtest/entity"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type usersRepo struct {
	client     *mongo.Client
	db         *mongo.Database
	collection *mongo.Collection
}

func (u usersRepo) Create(user *entity.User) (*entity.User, error) {
	insertResult, err := u.collection.InsertOne(context.Background(), user)
	if err != nil {
		return nil, err
	}
	log.Printf("Inserted New Company In DB. ID: %v", insertResult.InsertedID)
	return user, nil
}

func (u usersRepo) Get(id string) (*entity.User, error) {
	var user entity.User
	if err := u.collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&user); err != nil {
		log.Printf("User with ID(%s) was not found", id)
		return nil, errors.New(fmt.Sprintf("User with ID(%s) was not found", id))
	}

	return &user, nil
}

func (u usersRepo) Update(user *entity.User) (*entity.User, error) {
	filter := bson.M{"email": user}

	update, err := CreateUpdateBson(user)
	if err != nil {
		log.Printf("Error creating update Item: %s", err)
		return nil, err
	}

	_, err = u.collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Printf("Unable to update user info: %s", err)
		return nil, err
	}

	return user, nil
}

func (u usersRepo) Filter(filters map[string]interface{}) (*entity.User, error) {
	filter := bson.M{}
	for key, value := range filters {
		filter[key] = value
	}

	var user entity.User
	err := u.collection.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		log.Println("Couldn't find the given user")
		return nil, err
	}

	return &user, nil
}

func (u usersRepo) Delete(id string) error {
	filter := bson.M{"_id": id}
	_, err := u.collection.DeleteOne(context.Background(), filter)
	if err != nil {
		log.Printf("Error deleting user info: %s", err)
		return err
	}

	return nil
}

type UsersRepository interface {
	Create(user *entity.User) (*entity.User, error)
	Get(id string) (*entity.User, error)
	Update(user *entity.User) (*entity.User, error)
	Filter(filters map[string]interface{}) (*entity.User, error)
	Delete(id string) error
}

func NewUsersRepo(dbURI string) UsersRepository {
	log.Println("Initializing Company Repo")
	clientOptions := options.Client().ApplyURI(dbURI)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Printf("Failed to connect to database: %s", err)
	}

	db := client.Database(viper.Get("DB_NAME").(string))
	collection := db.Collection(viper.Get("USERS_COLLECTION").(string))
	log.Printf("Connected to DB")
	return &usersRepo{
		client:     client,
		db:         db,
		collection: collection,
	}
}
