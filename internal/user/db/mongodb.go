package db

import (
	"REST_API/internal/apperror"
	"REST_API/internal/user"
	"REST_API/pkg/logging"
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type db struct {
	collection *mongo.Collection
	logger     *logging.Logger
}

func (d *db) Create(ctx context.Context, user user.User) (string, error) {
	d.logger.Debug("Creating user")
	result, err := d.collection.InsertOne(ctx, user)
	if err != nil {
		return "", fmt.Errorf("failed to create user due to error: %v", err)
	}
	d.logger.Debug("Convert InsertedID to ObjectID")
	oid, ok := result.InsertedID.(bson.ObjectID)
	if ok {
		return oid.Hex(), nil
	}
	d.logger.Trace(user)
	return "", fmt.Errorf("failed to convert objectid to hex. probably oid: %s", oid)
}

func (d *db) FindAll(ctx context.Context) (u []user.User, err error) {

	cursor, err := d.collection.Find(ctx, bson.M{})
	if cursor.Err() != nil {
		return u, fmt.Errorf("failed to find one user by id. id: due to error %v", err)
	}

	if err = cursor.All(ctx, &u); err != nil {
		return u, fmt.Errorf("failed to read all documents from cursor. error: %v", err)
	}

	return u, nil
}

func (d *db) FindOne(ctx context.Context, id string) (u user.User, err error) {
	oid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return u, fmt.Errorf("failed to convert hex to objectid. hex: %s", id)
	}
	filter := bson.M{"_id": oid}

	result := d.collection.FindOne(ctx, filter)
	if result.Err() != nil {
		if errors.Is(result.Err(), mongo.ErrNoDocuments) {
			return u, apperror.ErrNotFound
		}
		return u, fmt.Errorf("failed to find one user by id. id: %s due to error %v", id, err)
	}
	if err = result.Decode(&u); err != nil {
		return u, fmt.Errorf("failed to decode user(id:%s) from DB. due to error %v", id, err)
	}
	return u, nil
}

func (d *db) Update(ctx context.Context, user user.User) error {
	objectId := user.ID
	//if err != nil {
	//	return fmt.Errorf("failed to convert user ID to objectid. ID: %s", user.ID)
	//}
	filter := bson.M{"_id": objectId}

	userBytes, err := bson.Marshal(user)
	if err != nil {
		return fmt.Errorf("failed to marshal user, error: %s", err.Error())
	}

	var updateUserObj bson.M
	err = bson.Unmarshal(userBytes, &updateUserObj)
	if err != nil {
		return fmt.Errorf("failed to unmarshal user bytes, error: %s", err.Error())
	}

	delete(updateUserObj, "_id")

	update := bson.M{
		"$set": updateUserObj,
	}

	result, err := d.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to ececute update user query. error: %v", err)
	}
	if result.MatchedCount == 0 {
		return apperror.ErrNotFound
	}

	d.logger.Tracef("Matched %d documentd and Modified %d documents", result.MatchedCount, result.ModifiedCount)

	return nil
}

func (d *db) Delete(ctx context.Context, id string) error {
	objectId, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("failed to convert user ID to objectid. ID: %s", id)
	}
	filter := bson.M{"_id": objectId}

	result, err := d.collection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to execute  query. error %v", err)
	}
	if result.DeletedCount == 0 {
		return apperror.ErrNotFound
	}

	d.logger.Tracef("Deleted %d documents ", result.DeletedCount)

	return nil
}

func NewStorage(database *mongo.Database, collection string, logger *logging.Logger) user.Storage {
	return &db{
		collection: database.Collection(collection),
		logger:     logger,
	}
}
