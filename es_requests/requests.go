package es_requests

import (
	"code.google.com/p/go-uuid/uuid"
	"gopkg.in/olivere/elastic.v3"
)

func IndexDocMessage(client *elastic.Client, indexName, docType, body string) (*elastic.IndexResponse, error) {
	resp, err := client.Index().Index(indexName).Type(docType).Id(uuid.New()).BodyJson(body).Do()
	return resp, err
}

func IndexDocJSONBytes(client *elastic.Client, indexName, docType string, body string) (*elastic.IndexResponse, error) {
	resp, err := client.Index().Index(indexName).Type(docType).Id(uuid.New()).BodyString(body).Do()
	return resp, err
}

func IndexDocMapping(client *elastic.Client, indexName, docType string, body map[string]interface{}) (*elastic.IndexResponse, error) {
	resp, err := client.Index().Index(indexName).Type(docType).Id(uuid.New()).BodyJson(body).Do()
	return resp, err
}

func GetDoc(client *elastic.Client, indexName, docType, id string) (*elastic.GetResult, error) {
	resp, err := client.Get().Index(indexName).Type(docType).Id(id).Do()
	return resp, err
}
