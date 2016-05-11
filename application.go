package main

import (
	"github.com/JKolios/elasticgin/api"
	"github.com/JKolios/elasticgin/rabbitmq"
	"github.com/JKolios/elasticgin/config"
	"github.com/JKolios/elasticgin/utils"
	"github.com/streadway/amqp"
	"gopkg.in/olivere/elastic.v3"
	"log"
)

func initESClient(url string, indices []string, doSniff bool) *elastic.Client {

	log.Printf("Connecting to ES on: %v", url)
	elasticClient, err := elastic.NewClient(elastic.SetURL(url), elastic.SetSniff(doSniff))
	utils.CheckFatalError(err)

	log.Println("Connected to ES")
	
	for _, index := range(indices){
		
		log.Printf("Initializing Index: %s", index)
	
		indexExists, err := elasticClient.IndexExists(index).Do()
		utils.CheckFatalError(err)
		if !indexExists {
			resp, err := elasticClient.CreateIndex(index).Do()
			utils.CheckFatalError(err)
			if !resp.Acknowledged {
				log.Fatalf("Cannot create index: %s on ES", index)
			}
			log.Printf("Created index: %s on ES", index)

		} else {
			log.Printf("Index: %s already exists on ES", index)
		}

		_, err = elasticClient.OpenIndex(index).Do()
		utils.CheckFatalError(err)
		
		mapping, err := elasticClient.GetMapping().Index(index).Do()
		if err != nil {
			log.Printf("Cannot get mapping for index: %s", index)
		}
		log.Printf("Mapping for index %s: %s", index,  mapping)
	}
		
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
	config := config.GetConfFromJSONFile("config.json")

	//ES init
	esClient := initESClient(config.ElasticURL, config.Indices, config.SniffCluster)	
	defer esClient.Stop()

	//Rabbitmq init
	var amqpChannel *amqp.Channel
	
	if config.UseAMQP{
		
		amqpConnection, amqpChannel := initAMQPClient(config)
		defer amqpConnection.Close()
		defer amqpChannel.Close()

		rabbitmq.StartSubscribers(amqpChannel, esClient, config)
	}else{
		amqpChannel = nil
	}

	api := api.SetupAPI(esClient, amqpChannel, config)
	api.Run(config.ApiURL)
}
