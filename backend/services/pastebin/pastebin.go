package pastebin

import (
	"log"
	"strconv"
	"time"

	"github.com/viorel-d/binscrape/backend/common"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PasteBinItem struct {
	ID        string            `json: "id" bson:"_id, omitempty"`
	URL       string            `json: "url" bson:"url"`
	CreatedAt string            `json: "createdAt" bson:"created_at"`
	Content   string            `json: "content" bson:"content"`
	Meta      map[string]string `json: "meta" bson:"meta, omitempty"`
}

type PasteBinResponse struct {
	Body []byte
}

type RawPasteBinItem struct {
	ID        string
	URL       string
	Content   string
	CreatedAt string
	Meta      map[string]string
}

const PasteBinURL = "https://pastebin.com"

func NewPasteBinHTTPRequest() (pbr *PasteBinResponse) {
	body, err := common.MakeHTTPRequest(
		"GET",
		PasteBinURL,
		nil,
		nil,
	)
	if err != nil {
		log.Panic(err)
	}
	pbr = &PasteBinResponse{body}

	return
}

func (pbRes *PasteBinResponse) extractLinks() []string {
	linkRegexp := `href="(\/[a-zA-Z0-9]{8})"`
	strBody := string(pbRes.Body)
	links := common.RegexpFindAllString(linkRegexp, strBody)

	return links
}

func (pbRes *PasteBinResponse) extractPasteBinItemsTypes() []string {
	binTypeRegexp := `<span>([a-zA-Z0-9\+#\.]+)\s\|`
	strBody := string(pbRes.Body)
	pasteBinItemsTypes := common.RegexpFindAllString(binTypeRegexp, strBody)

	return pasteBinItemsTypes
}

func (pbRes *PasteBinResponse) ParseResponse() (rawPasteBinItems []RawPasteBinItem) {
	links := pbRes.extractLinks()
	pasteBinItemsTypes := pbRes.extractPasteBinItemsTypes()

	if len(links) != len(pasteBinItemsTypes) {
		log.Panicln("Response parse error: links and pasteBinItemsTypes differ in length")
	}

	for i, link := range links {
		pasteBinItemType := pasteBinItemsTypes[i]
		url := "/raw" + link
		pasteBinItemTypeMap := make(map[string]string)
		pasteBinItemTypeMap["pasteBinItemType"] = pasteBinItemType
		objectID := primitive.NewObjectID()
		rb := RawPasteBinItem{
			objectID.Hex(),
			url,
			"",
			"",
			pasteBinItemTypeMap,
		}
		rawPasteBinItems = append(rawPasteBinItems, rb)
	}

	return
}

func GetRawPasteBinItemsContent(
	rawPasteBinItems *[]RawPasteBinItem,
) (newRawPasteBinItemsContent []RawPasteBinItem) {
	for _, rpbi := range *rawPasteBinItems {
		finalURL := PasteBinURL + rpbi.URL
		resBody, err := common.MakeHTTPRequest(
			"GET",
			finalURL,
			nil,
			nil,
		)
		if err != nil {
			return
		}
		newRawPasteBinItem := RawPasteBinItem{
			rpbi.ID,
			rpbi.URL,
			string(resBody),
			strconv.FormatInt(time.Now().Unix(), 10),
			rpbi.Meta,
		}
		newRawPasteBinItemsContent = append(
			newRawPasteBinItemsContent,
			newRawPasteBinItem,
		)
		time.Sleep(1 * time.Second)
	}

	return
}
