package database

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func UpdateOrAddItem(collection string, item map[string]interface{}, filter map[string]interface{}) error {
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

func GetMappedStringItems(collection string, key string) ([]string, error) {
	db := GetDatabase()
	collected, err := db.Collection(collection).Find(context.TODO(), bson.D{})
	if err != nil {
		return nil, err
	}

	var items []bson.M
	err = collected.All(context.TODO(), &items)
	if err != nil {
		return nil, err
	}

	values := []string{}
	for _, item := range items {
		values = append(values, item[key].(string))
	}

	return values, nil
}

func GetItem[T any](collection string, filter map[string]interface{}) (*T, error) {
	db := GetDatabase()

	var item *T
	err := db.Collection(collection).FindOne(context.TODO(), filter).Decode(&item)
	if err != nil {
		return item, err
	}

	return item, nil
}
