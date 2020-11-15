package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/viorel-d/binscrape/backend/config"

	"github.com/viorel-d/binscrape/backend/common"
	"github.com/viorel-d/binscrape/backend/services/mongo"
	"github.com/viorel-d/binscrape/backend/services/pastebin"
)

type AllPasteBinItemsResponseBody struct {
	Data       []pastebin.PasteBinItem `json: "data"`
	StatusCode int                     `json: "statusCode"`
}

func writeAllPasteBinItemsResponse(
	w http.ResponseWriter,
	data *[]pastebin.PasteBinItem,
	statusCode int,
) {
	var pbItems []pastebin.PasteBinItem
	if statusCode == common.HTTPStatusOk {
		pbItems = *data
	}

	resBody := AllPasteBinItemsResponseBody{pbItems, statusCode}
	body, err := json.Marshal(resBody)
	if err != nil {
		common.WriteHTTPResponse(w, common.HTTPStatusInternalError, nil, nil)
		return
	}

	common.WriteHTTPResponse(w, statusCode, body, nil)
}

func GetAllPasteBinItemsHandler(w http.ResponseWriter, r *http.Request) {
	col := mongo.GetCollection(config.MongoItemsCollectionName)
	items := mongo.Find(col, nil)
	bsonItems, err := mongo.AsSliceOfBsonM(&items)
	if err != nil {
		writeAllPasteBinItemsResponse(w, nil, common.HTTPStatusInternalError)
		return
	}

	jsonPasteBinItems, err := pastebin.AsSliceOfJSON(&bsonItems)
	if err != nil {
		writeAllPasteBinItemsResponse(w, nil, common.HTTPStatusInternalError)
		return
	}

	writeAllPasteBinItemsResponse(w, &jsonPasteBinItems, common.HTTPStatusOk)
}
