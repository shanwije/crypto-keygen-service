package repositories

import (
	"context"
	"errors"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type KeyGenRepository interface {
	SaveKey(userID int, network string, address string, publicKey string, privateKey string) error
	GetKey(userID int, network string) (string, string, string, error)
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

func (r *MongoRepository) SaveKey(userID int, network string, address string, publicKey string, encryptedPrivateKey string) error {
	log.Printf("Saving encrypted private key: %s", encryptedPrivateKey)

	keyData := KeyData{
		UserID:              userID,
		Network:             network,
		Address:             address,
		PublicKey:           publicKey,
		EncryptedPrivateKey: encryptedPrivateKey,
	}
	filter := bson.M{"user_id": userID, "network": network}
	update := bson.M{"$set": keyData}
	_, err := r.collection.UpdateOne(context.Background(), filter, update, options.Update().SetUpsert(true))
	return err
}

func (r *MongoRepository) GetKey(userID int, network string) (string, string, string, error) {
	filter := bson.M{"user_id": userID, "network": network}
	var keyData KeyData
	err := r.collection.FindOne(context.Background(), filter).Decode(&keyData)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return "", "", "", errors.New("key not found")
		}
		return "", "", "", err
	}

	log.Printf("Retrieved encrypted private key: %s", keyData.EncryptedPrivateKey)

	return keyData.Address, keyData.PublicKey, keyData.EncryptedPrivateKey, nil
}

func (r *MongoRepository) KeyExists(userID int, network string) (bool, error) {
	filter := bson.M{"user_id": userID, "network": network}
	count, err := r.collection.CountDocuments(context.Background(), filter)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
