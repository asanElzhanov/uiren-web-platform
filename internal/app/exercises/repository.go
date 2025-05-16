package exercises

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

func (r *exercisesRepository) getExercisesByCodes(ctx context.Context, codes []string) ([]Exercise, error) {
	var (
		collection = r.db.Collection(exercisesCollection)
		filter     = bson.M{
			"code": bson.M{
				"$in": codes,
			},
			"deleted_at": nil,
		}
		response []Exercise
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

func (r *exercisesRepository) getExercise(ctx context.Context, code string) (Exercise, error) {
	var (
		collection = r.db.Collection(exercisesCollection)
		filter     = bson.M{"code": code, "deleted_at": nil}
		response   Exercise
	)

	if err := collection.FindOne(ctx, filter).Decode(&response); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return Exercise{}, ErrNotFound
		}
	}

	return response, nil
}

func (r *exercisesRepository) createExercise(ctx context.Context, dto CreateExerciseDTO) (primitive.ObjectID, error) {
	var (
		collection = r.db.Collection(exercisesCollection)
	)

	res, err := collection.InsertOne(ctx, dto)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return primitive.NilObjectID, ErrCodeAlreadyExists
		}
		return primitive.NilObjectID, err
	}

	oid, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		logger.Error("exercisesRepository.createExercise type assertion InsertedID.(ObjectID)")
	}

	return oid, nil
}

func (r *exercisesRepository) deleteExercise(ctx context.Context, code string) error {
	var (
		collection = r.db.Collection(exercisesCollection)
		filter     = bson.M{"code": code, "deleted_at": nil}
		update     = bson.M{"$set": bson.M{"deleted_at": time.Now()}}
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

func (r *exercisesRepository) getExerciseType(ctx context.Context, code string) (string, error) {
	var (
		collection = r.db.Collection(exercisesCollection)
		filter     = bson.M{"code": code, "deleted_at": nil}
		projection = bson.M{"type": 1, "_id": 0}
	)

	raw, err := collection.FindOne(ctx, filter, options.FindOne().SetProjection(projection)).Raw()
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return "", ErrNotFound
		}
		return "", err
	}

	exerciseType, ok := raw.Lookup("type").StringValueOK()
	if !ok {
		return "", ErrNotFound
	}

	return exerciseType, nil
}

func (r *exercisesRepository) updateExercise(ctx context.Context, code string, dto UpdateExerciseDTO) error {
	var (
		collection = r.db.Collection(exercisesCollection)
		filter     = bson.M{"code": code, "deleted_at": nil}
		update     = bson.M{}
	)

	if dto.Question != nil {
		update["question"] = *dto.Question
	}
	if dto.Hints != nil {
		update["hints"] = dto.Hints
	}
	if dto.Explanation != nil {
		update["explanation"] = *dto.Explanation
	}
	if dto.Options != nil {
		update["options"] = dto.Options
	}
	if dto.CorrectAnswer != nil {
		update["correct_answer"] = *dto.CorrectAnswer
	}
	if dto.CorrectOrder != nil {
		update["correct_order"] = dto.CorrectOrder
	}
	if dto.Pairs != nil {
		update["pairs"] = dto.Pairs
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

func (r *exercisesRepository) getAllExercises(ctx context.Context) ([]Exercise, error) {
	var (
		collection = r.db.Collection(exercisesCollection)
		filter     = bson.M{"deleted_at": nil}
		result     []Exercise
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
