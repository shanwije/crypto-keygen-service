package repositories

import (
	"context"
	"crypto-keygen-service/internal/db"
)

type KeyGenRepository struct {
	database db.Database
}

func NewKeyGenRepository(database db.Database) *KeyGenRepository {
	return &KeyGenRepository{database: database}
}

func (r *KeyGenRepository) SaveKey(ctx context.Context, keyData db.KeyData) error {
	return r.database.SaveKey(ctx, keyData)
}

func (r *KeyGenRepository) GetKey(ctx context.Context, userID int, network string) (db.KeyData, error) {
	return r.database.GetKey(ctx, userID, network)
}

func (r *KeyGenRepository) KeyExists(ctx context.Context, userID int, network string) (bool, error) {
	return r.database.KeyExists(ctx, userID, network)
}

func (r *KeyGenRepository) CreateIndexes(ctx context.Context) error {
	return r.database.CreateIndexes(ctx)
}
