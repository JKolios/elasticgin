package api

import (
	"github.com/gin-gonic/gin"
	"gopkg.in/olivere/elastic.v2"
)

func esInjector(ESClient *elastic.Client, index string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("ESClient", ESClient)
		c.Set("Index", index)
		c.Next()
	}
}

func SetupAPI(ESClient *elastic.Client, index string, debug bool) *gin.Engine {
	if !debug {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(esInjector(ESClient, index))

	//API v0 endpoints
	v0 := router.Group("/v0")
	{
		v0.GET("/status", status)
		v0.POST("/indexDoc", indexDoc)
		v0.GET("/getDoc", getDoc)
	}
	return router
}
