package modules

import (
	"context"
	"errors"
	"time"
	"uiren/pkg/logger"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
		if errors.Is(err, mongo.ErrNoDocuments) {
			return module{}, ErrNotFound
		}
		return module{}, err
	}

	return response, nil
}

func (r *modulesRepository) createModule(ctx context.Context, dto CreateModuleDTO) (primitive.ObjectID, error) {
	var (
		collection = r.db.Collection(modulesCollection)
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
		logger.Error("modulesRepository.createModule type assertion InsertedID.(ObjectID)")
	}

	return oid, nil
}

func (r *modulesRepository) deleteModule(ctx context.Context, code string) error {
	var (
		collection = r.db.Collection(modulesCollection)
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

func (r *modulesRepository) updateModule(ctx context.Context, code string, dto UpdateModuleDTO) error {
	var (
		collection = r.db.Collection(modulesCollection)
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
	if dto.Goal != nil {
		update["goal"] = *dto.Goal
	}
	if dto.Difficulty != nil {
		update["goal"] = *dto.Difficulty
	}
	if dto.UnlockReq != nil {
		update["unlock_requirements"] = *dto.UnlockReq
	}
	if dto.Reward != nil {
		update["reward"] = *dto.Reward
	}

	if len(update) == 0 {
		return ErrNoFieldsToUpdate
	}

	update = bson.M{
		"$set": update,
	}
	res, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if res.MatchedCount == 0 {
		return ErrNotFound
	}

	return nil
}

func (r *modulesRepository) addLessonToList(ctx context.Context, code, lessonCode string) error {
	var (
		collection = r.db.Collection(modulesCollection)
		filter     = bson.M{"code": code, "deleted_at": nil}
		update     = bson.M{"$addToSet": bson.M{
			"lessons": lessonCode,
		}}
	)

	res, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if res.MatchedCount == 0 {
		return ErrNotFound
	}
	if res.ModifiedCount == 0 {
		return ErrLessonAlreadyInSet
	}

	return nil
}

func (r *modulesRepository) deleteLessonFromList(ctx context.Context, code, lessonCode string) error {
	var (
		collection = r.db.Collection(modulesCollection)
		filter     = bson.M{"code": code, "deleted_at": nil}
		update     = bson.M{"$pull": bson.M{
			"lessons": lessonCode,
		}}
	)

	res, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if res.MatchedCount == 0 {
		return ErrNotFound
	}
	if res.ModifiedCount == 0 {
		return ErrLessonNotInList
	}

	return nil
}
