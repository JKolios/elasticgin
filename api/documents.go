package api

import (
	"github.com/Jkolios/elasticgin/config"
	"github.com/Jkolios/elasticgin/es_requests"
	"github.com/gin-gonic/gin"
	"gopkg.in/olivere/elastic.v2"
	"log"
	"net/http"
)

type IncomingDoc struct {
	Type  string            `json:"type" binding:"required"`
	Index string            `json:"index"`
	Body  map[string]string `json:"body" binding:"required"`
}

func indexDoc(c *gin.Context) {
	var incoming IncomingDoc
	var incomingMessage string
	var resp *elastic.IndexResult
	var err error

	client := c.MustGet("ESClient").(*elastic.Client)
	defaultIndex := c.MustGet("Config").(config.Config).DefaultIndex

	if c.BindJSON(&incoming) == nil {
		log.Printf("Request JSON: %+v", incoming)

		if incoming.Index == "" {
			incoming.Index = defaultIndex
		}

		resp, err = es_requests.IndexDocMapping(client, incoming.Index, incoming.Type, incoming.Body)
	} else if c.Bind(&incomingMessage) == nil {

		resp, err = es_requests.IndexDocMessage(client, incoming.Index, "message", incomingMessage)
	} else {
		c.String(http.StatusBadRequest, "Cannot process this document")
		return
	}

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
	defaultIndex := c.MustGet("Config").(config.Config).DefaultIndex

	requestedId := c.Query("docId")
	requestedType := c.Query("docType")
	requestedIndex := c.Query("index")

	if requestedIndex == "" {

		requestedIndex = defaultIndex
	}

	esResp, err := es_requests.GetDoc(client, requestedIndex, requestedType, requestedId)

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

	} else {
		responseBody["success"] = true
		responseBody["doc"] = esResp.Source
		c.JSON(http.StatusOK, responseBody)
		return
	}
}

func fullTextSearch(c *gin.Context) {

	//term := c.Query("term")
	c.String(http.StatusNotImplemented, "Not Yet Implemented")
	return
}
