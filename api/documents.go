package api

import (
	"code.google.com/p/go-uuid/uuid"
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
	defaultIndex := c.MustGet("Index").(string)

	if c.BindJSON(&incoming) == nil {
		log.Printf("Request JSON: %+v", incoming)

		if incoming.Index == "" {
			incoming.Index = defaultIndex
		}

		resp, err = client.Index().Index(incoming.Index).
			Type(incoming.Type).Id(uuid.New()).BodyJson(incoming.Body).Do()

	} else if c.Bind(&incomingMessage) == nil {

		resp, err = client.Index().Index(defaultIndex).
			Type("message").Id(uuid.New()).BodyJson(incomingMessage).Do()
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
	defaultIndex := c.MustGet("Index").(string)

	requestedId := c.Query("docId")
	requestedType := c.Query("docType")
	requestedIndex := c.Query("index")

	if requestedIndex == "" {

		requestedIndex = defaultIndex
	}

	esResp, err := client.Get().Index(requestedIndex).Type(requestedType).Id(requestedId).Do()

	response := make(map[string]interface{})

	if err != nil {
		response["success"] = false
		response["error"] = err.Error()
		c.JSON(http.StatusBadRequest, response)
		return
	}

	if !esResp.Found {
		response["success"] = false
		response["error"] = "Document not found."
		c.JSON(http.StatusOK, response)
		return

	} else {
		response["success"] = true
		response["doc"] = esResp.Source
		c.JSON(http.StatusOK, response)
		return
	}
}

func fullTextSearch(c *gin.Context) {

	//term := c.Query("term")
	c.String(http.StatusNotImplemented, "Not Yet Implemented")
	return
}
