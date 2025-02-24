package lessons

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	lessonsCollection = "lessons"
)

type lessonRepository struct {
	db *mongo.Database
}

func NewLessonRepository(db *mongo.Database) *lessonRepository {
	return &lessonRepository{
		db: db,
	}
}

func (r *lessonRepository) getLessonsByCodes(ctx context.Context, codes []string) ([]lesson, error) {
	var (
		collection = r.db.Collection(lessonsCollection)
		filter     = bson.M{
			"code": bson.M{
				"$in": codes,
			},
			"deleted_at": nil,
		}
		response []lesson
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
