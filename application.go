package main

import (
	"github.com/Jkolios/elasticgin/api"
	"github.com/Jkolios/elasticgin/rabbitmq"
	"github.com/Jkolios/elasticgin/config"
	"github.com/Jkolios/elasticgin/utils"
	"github.com/streadway/amqp"
	"gopkg.in/olivere/elastic.v2"
	"log"
)

func initESClient(config *config.Config) *elastic.Client {

	log.Printf("Connecting to ES on: %v", config.ElasticURL)
	elasticClient, err := elastic.NewClient(elastic.SetURL(config.ElasticURL), elastic.SetSniff(config.SniffCluster))
	utils.CheckFatalError(err)

	log.Println("Connected to ES")
	indexExists, err := elasticClient.IndexExists(config.DefaultIndex).Do()
	utils.CheckFatalError(err)
	if !indexExists {
		resp, err := elasticClient.CreateIndex(config.DefaultIndex).Do()
		utils.CheckFatalError(err)
		if !resp.Acknowledged {
			log.Fatal("Cannot create index on ES")
		}
		log.Println("Created index on ES")

	} else {
		log.Println("Index already exists on ES")
	}

	_, err = elasticClient.OpenIndex(config.DefaultIndex).Do()
	utils.CheckFatalError(err)
	return elasticClient
}

func initAMQPClient(config *config.Config) (*amqp.Connection, *amqp.Channel) {

	log.Printf("Connecting to RabbitMQ on: %v", config.AmqpURL)
	conn, err := amqp.Dial(config.AmqpURL)
	utils.CheckFatalError(err)
	ch, err := conn.Channel()
	utils.CheckFatalError(err)
	log.Println("Connected to RabbitMQ.")
	for _, queue := range(config.AmqpQueues){
		log.Printf("Declaring Queue: %v", queue)
		_, err = ch.QueueDeclare(
			queue,
			false, // durable
			false, // delete when unused
			false, // exclusive
			false, // no-wait
			nil,   // arguments
		)
		utils.CheckFatalError(err)
		log.Println("Queue Declared")
}
	return conn, ch
}

func main() {

	log.Println("Starting elasticgin")

	//Config fetch
	config := config.GetConfFromJSONFile("config/config.json")

	//ES init
	esClient := initESClient(config)

	//Rabbitmq init
	amqpConnection, amqpChannel := initAMQPClient(config)
	defer amqpConnection.Close()
	defer amqpChannel.Close()

	rabbitmq.StartSubscribers(amqpChannel, esClient, config)

	api := api.SetupAPI(esClient, amqpChannel, config)
	api.Run(config.ApiURL)
}
