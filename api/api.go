package api

import (
	"code.google.com/p/go-uuid/uuid"
	"github.com/gin-gonic/gin"
	"gopkg.in/olivere/elastic.v2"
	"log"
	"net/http"
)

type IncomingDoc struct {
	Type string            `json:"type"`
	Body map[string]string `json:"body" binding:"required"`
}

func esInjector(ESClient *elastic.Client, index string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("ESClient", ESClient)
		c.Set("Index", index)
		c.Next()
	}
}

func indexer(c *gin.Context) {
	var incoming IncomingDoc
	client := c.MustGet("ESClient").(*elastic.Client)
	index := c.MustGet("Index").(string)

	if c.BindJSON(&incoming) == nil {
		log.Printf("Request JSON: %+v", incoming)
		resp, err := client.Index().Index(index).
			Type(incoming.Type).Id(uuid.New()).BodyJson(incoming.Body).Do()

		if err != nil {
			log.Println(err.Error())
			c.String(http.StatusBadRequest, err.Error())
			return
		}
		c.JSON(http.StatusOK, resp)
		return
	}
	c.String(http.StatusBadRequest, "Failed to bind JSON Request.")
}

func docGetter(c *gin.Context) {

	client := c.MustGet("ESClient").(*elastic.Client)
	index := c.MustGet("Index").(string)

	requestedId := c.Query("docId")
	requestedType := c.Query("docType")
	esResp, err := client.Get().Index(index).Type(requestedType).Id(requestedId).Do()

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

func statusResponder(c *gin.Context) {
	c.String(http.StatusOK, "All systems nominal :P")
}

func SetupAPI(ESClient *elastic.Client, index string) *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(esInjector(ESClient, index))

	//API v0 endpoints
	v0 := router.Group("/v0")
	{
		v0.GET("/status", statusResponder)
		v0.POST("/index_doc", indexer)
		v0.GET("/get_doc", docGetter)
	}
	gin.SetMode(gin.ReleaseMode)
	return router
}
