package database

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetAndUpdateOrAddItem(collection string, item map[string]interface{}, filter map[string]interface{}) error {
	db := GetDatabase()
	dbFilter := bson.D{}

	for k, v := range filter {
		dbFilter = append(dbFilter, primitive.E{
			Key:   k,
			Value: v,
		})
	}

	var foundItem bson.M
	err := db.Collection(collection).FindOne(context.TODO(), dbFilter).Decode(&foundItem)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			_, err = db.Collection("fuel_prices").InsertOne(context.TODO(), item)
			if err != nil {
				return err
			}
			return nil
		} else {
			return err
		}
	}

	updateFields := bson.D{}

	for k, v := range item {
		updateFields = append(updateFields, primitive.E{
			Key:   k,
			Value: v,
		})
	}

	updateItem := bson.D{
		{Key: "$set", Value: updateFields},
	}

	_, err = db.Collection(collection).UpdateOne(context.TODO(), dbFilter, updateItem)
	if err != nil {
		return err
	}

	return nil
}
