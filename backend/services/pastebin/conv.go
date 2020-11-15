package pastebin

import (
	"encoding/json"

	"go.mongodb.org/mongo-driver/bson"
)

func AsSliceOfInterface(rawPasteBinItems *[]RawPasteBinItem) []interface{} {
	var result []interface{}
	for _, rb := range *rawPasteBinItems {
		result = append(result, rb)
	}
	return result
}

func AsJSON(item *bson.M) (PasteBinItem, error) {
	var result PasteBinItem
	b, err := json.Marshal(*item)
	if err != nil {
		return result, err
	}
	err = json.Unmarshal(b, &result)

	return result, err
}

func AsBSON(item *bson.M) (PasteBinItem, error) {
	var result PasteBinItem
	b, err := bson.Marshal(*item)
	if err != nil {
		return result, err
	}
	err = bson.Unmarshal(b, &result)

	return result, err
}

func AsSliceOfJSON(items *[]bson.M) ([]PasteBinItem, error) {
	var result []PasteBinItem
	for _, item := range *items {
		pb, err := AsJSON(&item)
		if err != nil {
			return result, err
		}
		result = append(result, pb)
	}

	return result, nil
}

func AsSliceOfBSON(items *[]bson.M) ([]PasteBinItem, error) {
	var result []PasteBinItem
	for _, item := range *items {
		pb, err := AsBSON(&item)
		if err != nil {
			return result, err
		}
		result = append(result, pb)
	}

	return result, nil
}
