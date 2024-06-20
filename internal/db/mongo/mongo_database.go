package mongo

import (
	"context"
	dbi "crypto-keygen-service/internal/db"
	"errors"
	"time"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDatabase struct {
	Collection *mongo.Collection
	Client     *mongo.Client
}

func NewMongoDatabase(mongoURI, dbName, collectionName string) (*MongoDatabase, error) {
	clientOptions := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	collection := client.Database(dbName).Collection(collectionName)
	db := &MongoDatabase{Collection: collection, Client: client}
	err = db.CreateIndexes(context.Background())
	if err != nil {
		return nil, err
	}
	return db, nil
}

func (db *MongoDatabase) CreateIndexes(ctx context.Context) error {
	indexModel := mongo.IndexModel{
		Keys: bson.D{
			{Key: "user_id", Value: 1},
			{Key: "network", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	}
	_, err := db.Collection.Indexes().CreateOne(ctx, indexModel)
	return err
}

func (db *MongoDatabase) SaveKey(ctx context.Context, keyData dbi.KeyData) error {
	log.WithFields(log.Fields{
		"user_id": keyData.UserID,
		"network": keyData.Network,
	}).Info("Saving keys to repository")

	filter := bson.M{"user_id": keyData.UserID, "network": keyData.Network}
	update := bson.M{"$set": keyData}
	_, err := db.Collection.UpdateOne(ctx, filter, update, options.Update().SetUpsert(true))
	if err != nil {
		log.WithFields(log.Fields{
			"user_id": keyData.UserID,
			"network": keyData.Network,
		}).WithError(err).Error("Failed to save keys to repository")
	}
	return err
}

func (db *MongoDatabase) GetKey(ctx context.Context, userID int, network string) (dbi.KeyData, error) {
	filter := bson.M{"user_id": userID, "network": network}
	var keyData dbi.KeyData
	err := db.Collection.FindOne(ctx, filter).Decode(&keyData)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return dbi.KeyData{}, errors.New("key not found")
		}
		log.WithError(err).Error("Failed to retrieve key from repository")
		return dbi.KeyData{}, err
	}

	log.WithFields(log.Fields{
		"user_id": userID,
		"network": network,
	}).Info("Retrieved keys from repository")

	return keyData, nil
}

func (db *MongoDatabase) KeyExists(ctx context.Context, userID int, network string) (bool, error) {
	filter := bson.M{"user_id": userID, "network": network}
	count, err := db.Collection.CountDocuments(ctx, filter)
	if err != nil {
		log.WithError(err).Error("Failed to check if key exists in repository")
		return false, err
	}
	return count > 0, nil
}
