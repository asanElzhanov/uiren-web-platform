package modules

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	modulesCollection = "modules"
)

type modulesRepository struct {
	db *mongo.Database
}

func NewModulesRepository(db *mongo.Database) *modulesRepository {
	return &modulesRepository{
		db: db,
	}
}

func (r *modulesRepository) getModule(ctx context.Context, code string) (module, error) {
	var (
		collection = r.db.Collection(modulesCollection)
		filter     = bson.M{
			"code":       code,
			"deleted_at": nil,
		}
		response module
	)

	moduleEntity := collection.FindOne(ctx, filter)

	if err := moduleEntity.Decode(&response); err != nil {
		return module{}, err
	}

	return response, nil
}
