package config

import (
	"io/ioutil"
	"log"
	"encoding/json"
	"github.com/Jkolios/elasticgin/utils"
	)

type Config struct {
	ApiURL       string `json:"apiURL"`
	ElasticURL   string `json:"elasticURL"`
	DefaultIndex string `json:"defaultIndex"`
	SniffCluster bool   `json:"sniffCluster, omitempty"`
	AmqpURL      string `json:"amqpURL"`
	AmqpQueue    string `json:"amqpQueue"`
	GinDebug     bool   `json:"ginDebug, omitempty"`
}

func GetConfFromJSONFile(filename string) *Config {

	confContent, err := ioutil.ReadFile(filename)
	utils.CheckFatalError(err)
	var config *Config = new(Config)
	err = json.Unmarshal(confContent, config)
	utils.CheckFatalError(err)
	log.Println("Configuration loaded")
	log.Printf("Configuration: %+v\n", config)
	return config
}