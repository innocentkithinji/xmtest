package repository

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/Shopify/sarama"
	"github.com/innocentkithinji/xmtest/entity"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type companyRepo struct {
	client     *mongo.Client
	db         *mongo.Database
	collection *mongo.Collection
	kafkaUri   string
}

type CompanyRepository interface {
	Create(i *entity.Company) (*entity.Company, error)
	Get(id string) (*entity.Company, error)
	Update(company *entity.Company) (*entity.Company, error)
	Filter(filters map[string]interface{}) (*entity.Company, error)
	Delete(id string) error
}

func (c companyRepo) Filter(filters map[string]interface{}) (*entity.Company, error) {
	filter := bson.M{}
	for key, value := range filters {
		filter[key] = value
	}

	var company entity.Company
	err := c.collection.FindOne(context.Background(), filter).Decode(&company)
	if err != nil {
		log.Println("Couldn't find the given company")
		return nil, err
	}

	return &company, nil
}

func (c companyRepo) Emit(kafkaURI string, event string, value string) error {
	if kafkaURI == "" {
		kafkaURI = viper.Get("KAFKA_URI").(string)
	}
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer([]string{c.kafkaUri}, config)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := producer.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	message := &sarama.ProducerMessage{
		Topic: viper.Get("COMPANY_EVENT_TOPIC").(string),
		Key:   sarama.StringEncoder(event),
		Value: sarama.StringEncoder(value),
	}

	_, _, err = producer.SendMessage(message)
	if err != nil {
		log.Printf("Unable to push message")
		return err
	}

	return nil
}

func (c companyRepo) Create(company *entity.Company) (*entity.Company, error) {
	insertResult, err := c.collection.InsertOne(context.Background(), company)
	if err != nil {
		return nil, err
	}
	log.Println("Created Company (REPO)")

	_ = c.Emit(c.kafkaUri, "creation", insertResult.InsertedID.(string))
	log.Printf("Inserted New Company In DB. ID: %v", insertResult.InsertedID)
	return company, nil
}

func (c companyRepo) Get(id string) (*entity.Company, error) {
	var company entity.Company
	if err := c.collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&company); err != nil {
		log.Printf("Company with ID(%s) was not found", id)
		return nil, errors.New(fmt.Sprintf("Company with ID(%s) was not found", id))
	}

	return &company, nil
}

func (c companyRepo) Update(company *entity.Company) (*entity.Company, error) {
	filter := bson.M{"_id": company.ID}

	update, err := CreateUpdateBson(company)
	if err != nil {
		log.Printf("Error creating update Item: %s", err)
		return nil, err
	}

	_, err = c.collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Printf("Unable to update company info: %s", err)
		return nil, err
	}

	_ = c.Emit(c.kafkaUri, "update", company.ID)

	return company, nil
}

func (c companyRepo) Delete(id string) error {
	filter := bson.M{"_id": id}
	_, err := c.collection.DeleteOne(context.Background(), filter)
	if err != nil {
		log.Printf("Error deleting company info: %s", err)
		return err
	}

	_ = c.Emit(c.kafkaUri, "deletion", id)

	return nil
}

func NewCompanyRepo(dbURI string, kafkaURI string) CompanyRepository {
	log.Println("Initializing Company Repo")
	clientOptions := options.Client().ApplyURI(dbURI)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Printf("Failed to connect to database: %s", err)
	}

	db := client.Database(viper.Get("DB_NAME").(string))
	collection := db.Collection(viper.Get("COMPANY_COLLECTION").(string))
	log.Printf("Connected to DB")

	return &companyRepo{
		client:     client,
		db:         db,
		collection: collection,
		kafkaUri:   kafkaURI,
	}
}
