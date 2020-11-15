package services

import (
	"time"

	"github.com/viorel-d/binscrape/backend/config"
	"github.com/viorel-d/binscrape/backend/services/mongo"
	"github.com/viorel-d/binscrape/backend/services/pastebin"
	"go.mongodb.org/mongo-driver/bson"
)

func storePublicPasteBinItems(rawPasteBinItems *[]pastebin.RawPasteBinItem) {
	col := mongo.GetCollection(config.MongoItemsCollectionName)
	urlIndex := bson.M{"url": 1}
	canCreateURLIndex := mongo.EnsureIndex(col, "url")
	if canCreateURLIndex {
		mongo.CreateUniqueIndex(col, &urlIndex)
	}
	items := pastebin.AsSliceOfInterface(rawPasteBinItems)
	mongo.InsertMany(col, &items)
}

func ScrapePublicBins() {
	for {
		pbr := pastebin.NewPasteBinHTTPRequest()
		rawPasteBinItems := pbr.ParseResponse()
		rawPasteBinItemsContent := pastebin.GetRawPasteBinItemsContent(&rawPasteBinItems)
		storePublicPasteBinItems(&rawPasteBinItemsContent)
		time.Sleep(60 * time.Second)
	}
}
