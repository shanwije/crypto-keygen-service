package repositories

import (
	"context"
	"errors"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type KeyGenRepository interface {
	SaveKey(userID int, network string, keyData KeyData) error
	GetKey(userID int, network string) (KeyData, error)
	KeyExists(userID int, network string) (bool, error)
	CreateIndexes() error
}

type MongoRepository struct {
	collection *mongo.Collection
}

type KeyData struct {
	UserID              int    `bson:"user_id"`
	Network             string `bson:"network"`
	Address             string `bson:"address"`
	PublicKey           string `bson:"public_key"`
	EncryptedPrivateKey string `bson:"private_key"`
}

func NewMongoRepository(client *mongo.Client, dbName, collectionName string) *MongoRepository {
	collection := client.Database(dbName).Collection(collectionName)
	repo := &MongoRepository{collection: collection}
	err := repo.CreateIndexes()
	if err != nil {
		log.Fatalf("Failed to create indexes: %v", err)
	}
	return repo
}

func (r *MongoRepository) CreateIndexes() error {
	indexModel := mongo.IndexModel{
		Keys: bson.D{
			{Key: "user_id", Value: 1},
			{Key: "network", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	}
	_, err := r.collection.Indexes().CreateOne(context.Background(), indexModel)
	return err
}

func (r *MongoRepository) SaveKey(userID int, network string, keyData KeyData) error {
	log.WithFields(log.Fields{
		"user_id": userID,
		"network": network,
	}).Info("Saving keys to repository")

	filter := bson.M{"user_id": userID, "network": network}
	update := bson.M{"$set": keyData}
	_, err := r.collection.UpdateOne(context.Background(), filter, update, options.Update().SetUpsert(true))
	if err != nil {
		log.WithFields(log.Fields{
			"user_id": userID,
			"network": network,
		}).WithError(err).Error("Failed to save keys to repository")
	}
	return err
}

func (r *MongoRepository) GetKey(userID int, network string) (KeyData, error) {
	filter := bson.M{"user_id": userID, "network": network}
	var keyData KeyData
	err := r.collection.FindOne(context.Background(), filter).Decode(&keyData)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return KeyData{}, errors.New("key not found")
		}
		log.WithError(err).Error("Failed to retrieve key from repository")
		return KeyData{}, err
	}

	log.WithFields(log.Fields{
		"user_id": userID,
		"network": network,
	}).Info("Retrieved keys from repository")

	return keyData, nil
}

func (r *MongoRepository) KeyExists(userID int, network string) (bool, error) {
	filter := bson.M{"user_id": userID, "network": network}
	count, err := r.collection.CountDocuments(context.Background(), filter)
	if err != nil {
		log.WithError(err).Error("Failed to check if key exists in repository")
		return false, err
	}
	return count > 0, nil
}
