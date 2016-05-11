package api

import (
	"log"
	"net/http"

	"github.com/JKolios/elasticgin/config"
	"github.com/JKolios/elasticgin/es_requests"
	"github.com/gin-gonic/gin"
	"gopkg.in/olivere/elastic.v3"
)

type genericDoc struct {
	Type  string                 `json:"type" binding:"required"`
	Index string                 `json:"index"`
	Body  map[string]interface{} `json:"body" binding:"required"`
}

type termSearchPayload struct {
	Key  string		`json:"key" binding:"required"`
	Value string	`json:"value" binding:"required"`
	Index string 	`json:"index"`
	From  int		`json:"from"`
	MaxHits int		`json:"maxHits"` 
}

func indexDoc(c *gin.Context) {
	var incoming genericDoc
	var resp *elastic.IndexResponse
	var err error

	client := c.MustGet("ESClient").(*elastic.Client)
	defaultIndex := c.MustGet("Config").(*config.Config).DefaultIndex

	err = c.Bind(&incoming)
	if err != nil {
		log.Println(err.Error())
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	log.Printf("Incoming document struct: %+v", incoming)

	if incoming.Index == "" {
		incoming.Index = defaultIndex
	}
	resp, err = es_requests.IndexDocMapping(client, incoming.Index, incoming.Type, incoming.Body)

	if err != nil {
		log.Println(err.Error())
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, resp)
	return

}

func getDoc(c *gin.Context) {

	client := c.MustGet("ESClient").(*elastic.Client)
	defaultIndex := c.MustGet("Config").(*config.Config).DefaultIndex

	requestedID := c.Query("docID")
	requestedType := c.Query("docType")
	requestedIndex := c.Query("index")

	if requestedIndex == "" {

		requestedIndex = defaultIndex
	}

	esResp, err := es_requests.GetDoc(client, requestedIndex, requestedType, requestedID)

	responseBody := make(map[string]interface{})

	if err != nil {
		responseBody["success"] = false
		responseBody["error"] = err.Error()
		c.JSON(http.StatusBadRequest, responseBody)
		return
	}

	if !esResp.Found {
		responseBody["success"] = false
		responseBody["error"] = "Document not found."
		c.JSON(http.StatusOK, responseBody)
		return
	}

	responseBody["success"] = true
	responseBody["doc"] = esResp.Source
	c.JSON(http.StatusOK, responseBody)
	return

}

func termQuery(c *gin.Context) {

	client := c.MustGet("ESClient").(*elastic.Client)
	requestPayload := &termSearchPayload{From: 0, MaxHits:10, Index:c.MustGet("Config").(*config.Config).DefaultIndex}
	
	c.Bind(&requestPayload)

	log.Printf("Term Search initiated for index:%s key: %s value: %s maxHits: %d from: %d \n",
	 requestPayload.Index, requestPayload.Key, requestPayload.Value, requestPayload.MaxHits, requestPayload.From)

	Query := elastic.NewTermQuery(requestPayload.Key, requestPayload.Value)

	searchResult, err := client.Search().
		Index(requestPayload.Index).
		Query(Query).
		From(requestPayload.From).Size(requestPayload.MaxHits).
		Pretty(true).
		Do()

	if err != nil {
		log.Println(err.Error())
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	log.Printf("Term Search returned %d documents\n", searchResult.TotalHits())
	c.JSON(http.StatusOK, searchResult.Hits)

	return
}
