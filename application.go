package main

import (
	"encoding/json"
	"github.com/Jkolios/elasticgin/api"
	"gopkg.in/olivere/elastic.v2"
	"io/ioutil"
	"log"
)

type Config struct {
	ApiURL       string `json:"apiURL, omitempty"`
	ElasticURL   string `json:"elasticURL, omitempty"`
	DefaultIndex string `json:"defaultIndex, omitempty"`
	SniffCluster bool   `json:"sniffCluster, omitempty"`
	GinDebug     bool   `json:"ginDebug, omitempty"`
}

func checkFatalError(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}

func checkPassableError(err error) {
	if err != nil {
		log.Println(err.Error())
	}
}

func getConfFromFile(filename string) *Config {

	confContent, err := ioutil.ReadFile(filename)
	checkFatalError(err)
	var config *Config = new(Config)
	err = json.Unmarshal(confContent, config)
	checkFatalError(err)
	log.Println("Configuration loaded")
	log.Printf("Configuration: %+v\n", config)
	return config
}

func initESClient(config *Config) *elastic.Client {

	log.Printf("Connecting to ES on: %v", config.ElasticURL)
	elasticClient, err := elastic.NewClient(elastic.SetURL(config.ElasticURL), elastic.SetSniff(config.SniffCluster))
	checkFatalError(err)

	log.Println("Connected to ES")
	indexExists, err := elasticClient.IndexExists(config.DefaultIndex).Do()
	checkFatalError(err)
	if !indexExists {
		resp, err := elasticClient.CreateIndex(config.DefaultIndex).Do()
		checkFatalError(err)
		if !resp.Acknowledged {
			log.Fatal("Cannot create index on ES")
		}
		log.Println("Created index on ES")

	} else {
		log.Println("Index already exists on ES")
	}

	_, err = elasticClient.OpenIndex(config.DefaultIndex).Do()
	checkFatalError(err)
	return elasticClient
}

func main() {

	log.Println("Starting elasticgin")

	config := getConfFromFile("config.json")
	client := initESClient(config)
	defer client.CloseIndex(config.DefaultIndex).Do()
	api := api.SetupAPI(client, config.DefaultIndex, config.GinDebug)
	api.Run(config.ApiURL)
}
