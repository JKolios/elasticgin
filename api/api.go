package api

import (
	"github.com/JKolios/elasticgin/config"
	"github.com/gin-gonic/gin"
	"gopkg.in/olivere/elastic.v2"
	"github.com/streadway/amqp"
)

func contextInjector(ESClient *elastic.Client, AMQPChannel *amqp.Channel, config *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("ESClient", ESClient)
		c.Set("AMQPChannel", AMQPChannel)
		c.Set("Config", config)
		
		c.Next()
	}
}

func SetupAPI(ESClient *elastic.Client, AMQPChannel *amqp.Channel, config *config.Config) *gin.Engine {
	if !config.GinDebug {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(contextInjector(ESClient, AMQPChannel, config))

	//API v0 endpoints
	v0 := router.Group("/v0")
	{
		v0.GET("/status", status)
		v0.POST("/indexDoc", indexDoc)
		v0.GET("/getDoc", getDoc)
	}
	return router
}
