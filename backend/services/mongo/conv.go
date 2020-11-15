package mongo

import (
	"go.mongodb.org/mongo-driver/bson"
)

func AsBsonRawValue(item *bson.M) (bson.RawValue, error) {
	var rawValue bson.RawValue
	b, _ := bson.Marshal(*item)
	err := bson.Unmarshal(b, &rawValue)

	return rawValue, err
}

func AsBsonM(item interface{}) (bson.M, error) {
	var bsonM bson.M
	b, _ := bson.Marshal(item)
	err := bson.Unmarshal(b, &bsonM)

	return bsonM, err
}

func AsSliceOfBsonM(items *[]interface{}) ([]bson.M, error) {
	var result []bson.M
	for _, item := range *items {
		val, err := AsBsonM(item)
		if err != nil {
			return result, err
		}
		result = append(result, val)
	}

	return result, nil
}
