package db

import "context"

type KeyData struct {
	UserID              int    `bson:"user_id" json:"user_id"`
	Network             string `bson:"network" json:"network"`
	Address             string `bson:"address" json:"address"`
	PublicKey           string `bson:"public_key" json:"public_key"`
	EncryptedPrivateKey string `bson:"private_key" json:"private_key"`
}

type Database interface {
	SaveKey(ctx context.Context, keyData KeyData) error
	GetKey(ctx context.Context, userID int, network string) (KeyData, error)
	KeyExists(ctx context.Context, userID int, network string) (bool, error)
	CreateIndexes(ctx context.Context) error
}
