package exercises

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	exercisesCollection = "exercises"
)

type exercisesRepository struct {
	db *mongo.Database
}

func NewExercisesRepository(db *mongo.Database) *exercisesRepository {
	return &exercisesRepository{
		db: db,
	}
}

func (r *exercisesRepository) getExercisesByCodes(ctx context.Context, codes []string) ([]exercise, error) {
	var (
		collection = r.db.Collection(exercisesCollection)
		filter     = bson.M{
			"code": bson.M{
				"$in": codes,
			},
			"deleted_at": nil,
		}
		response []exercise
	)

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &response); err != nil {
		return nil, err
	}

	return response, nil
}
