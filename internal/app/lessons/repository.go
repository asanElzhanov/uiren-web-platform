package lessons

import (
	"context"
	"errors"
	"time"
	"uiren/pkg/logger"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func (r *lessonRepository) getLesson(ctx context.Context, code string) (lesson, error) {
	var (
		collection = r.db.Collection(lessonsCollection)
		filter     = bson.M{"code": code, "deleted_at": nil}
		response   lesson
	)

	res := collection.FindOne(ctx, filter)

	if err := res.Decode(&response); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return lesson{}, ErrNotFound
		}
		return lesson{}, err
	}

	return response, nil
}

func (r *lessonRepository) createLesson(ctx context.Context, dto CreateLessonDTO) (primitive.ObjectID, error) {
	var (
		collection = r.db.Collection(lessonsCollection)
	)

	result, err := collection.InsertOne(ctx, dto)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return primitive.NilObjectID, ErrCodeAlreadyExists
		}
		return primitive.NilObjectID, err
	}

	oid, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		logger.Error("lessonRepository.createLesson type assertion InsertedID.(ObjectID)")
	}

	return oid, nil
}

func (r *lessonRepository) updateLesson(ctx context.Context, code string, dto UpdateLessonDTO) error {
	var (
		collection = r.db.Collection(lessonsCollection)
		filter     = bson.M{
			"code":       code,
			"deleted_at": nil,
		}
		update = bson.M{}
	)

	if dto.Title != nil {
		update["title"] = *dto.Title
	}
	if dto.Description != nil {
		update["description"] = *dto.Description
	}
	if len(update) == 0 {
		return ErrNoFieldsToUpdate
	}
	update = bson.M{"$set": update}

	res, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return ErrNotFound
	}

	return nil
}

func (r *lessonRepository) deleteLesson(ctx context.Context, code string) error {
	var (
		collection = r.db.Collection(lessonsCollection)
		filter     = bson.M{
			"code":       code,
			"deleted_at": nil,
		}
		update = bson.M{"$set": bson.M{"deleted_at": time.Now()}}
	)

	res, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if res.MatchedCount == 0 {
		return ErrNotFound
	}

	return nil
}

func (r *lessonRepository) addExerciseToList(ctx context.Context, code, exerciseCode string) error {
	var (
		collection = r.db.Collection(lessonsCollection)
		filter     = bson.M{
			"code":       code,
			"deleted_at": nil,
		}
		update = bson.M{"$addToSet": bson.M{"exercises": exerciseCode}}
	)

	res, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if res.MatchedCount == 0 {
		return ErrNotFound
	}
	if res.ModifiedCount == 0 {
		return ErrExerciseAlreadyInSet
	}

	return nil
}

func (r *lessonRepository) deleteExerciseFromList(ctx context.Context, code, exerciseCode string) error {
	var (
		collection = r.db.Collection(lessonsCollection)
		filter     = bson.M{"code": code, "deleted_at": nil}
		update     = bson.M{"$pull": bson.M{"exercises": exerciseCode}}
	)

	res, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if res.MatchedCount == 0 {
		return ErrNotFound
	}
	if res.ModifiedCount == 0 {
		return ErrExerciseNotInList
	}

	return nil
}

func (r *lessonRepository) getAllLessons(ctx context.Context) ([]lesson, error) {
	var (
		collection = r.db.Collection(lessonsCollection)
		filter     = bson.M{"deleted_at": nil}
		result     []lesson
	)

	cur, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	if err := cur.All(ctx, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func (r *lessonRepository) lessonExists(ctx context.Context, code string) (bool, error) {
	var (
		collection = r.db.Collection(lessonsCollection)
		filter     = bson.M{"code": code, "deleted_at": nil}
	)

	count, err := collection.CountDocuments(ctx, filter, options.Count().SetLimit(1))
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
