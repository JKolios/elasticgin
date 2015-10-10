package api

import (
	"github.com/gin-gonic/gin"
	"gopkg.in/olivere/elastic.v2"
	"log"
	"net/http"
)

type IncomingDoc struct {
	Type string            `json:"type" binding:"required"`
	Id   string            `json:"id" binding:"required"`
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
			Type(incoming.Type).Id(incoming.Id).BodyJson(incoming.Body).Do()

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

func statusResponder(c *gin.Context) {
	c.String(http.StatusOK, "Hi")
}

func SetupAPI(ESClient *elastic.Client, index string) *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(esInjector(ESClient, index))
	router.POST("/index", indexer)
	router.GET("/status", statusResponder)
	return router
}
