package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/viorel-d/binscrape/backend/common"
	"github.com/viorel-d/binscrape/backend/services/mongo"
	"github.com/viorel-d/binscrape/backend/services/pastebin"
)

type PasteBinItemDetailResponseBody struct {
	Data       pastebin.PasteBinItem `json: "data"`
	StatusCode int                   `json: "statusCode"`
}

func writeBinDetailResponse(
	w http.ResponseWriter,
	data *pastebin.PasteBinItem,
	statusCode int,
) {
	var pb pastebin.PasteBinItem
	if statusCode == common.HTTPStatusOk {
		pb = *data
	}

	resBody := PasteBinItemDetailResponseBody{pb, statusCode}
	body, err := json.Marshal(resBody)
	if err != nil {
		common.WriteHTTPResponse(w, common.HTTPStatusInternalError, nil, nil)
	}

	common.WriteHTTPResponse(w, statusCode, body, nil)
}

func GetPasteBinItemDetailsHandler(w http.ResponseWriter, r *http.Request) {
	pasteBinItemID := r.URL.Query().Get("id")
	if pasteBinItemID == "" {
		writeBinDetailResponse(w, nil, common.HTTPStatusBadRequest)
		return
	}
	col := mongo.GetCollection("bins")
	item, err := mongo.FindOne(col, pasteBinItemID, nil)
	if err != nil {
		writeBinDetailResponse(w, nil, common.HTTPStatusInternalError)
		return
	}

	bsonMItem, err := mongo.AsBsonM(item)
	if err != nil {
		writeBinDetailResponse(w, nil, common.HTTPStatusInternalError)
		return
	}

	jsonItem, err := pastebin.AsJSON(&bsonMItem)
	if err != nil {
		writeBinDetailResponse(w, nil, common.HTTPStatusInternalError)
		return
	}

	writeBinDetailResponse(w, &jsonItem, common.HTTPStatusOk)
}
