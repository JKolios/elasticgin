package es_requests

import (
	"code.google.com/p/go-uuid/uuid"
	"gopkg.in/olivere/elastic.v2"
)

func IndexDocMessage(client *elastic.Client, indexName, docType, body string) (*elastic.IndexResult, error) {
	resp, err := client.Index().Index(indexName).Type(docType).Id(uuid.New()).BodyJson(body).Do()
	return resp, err
}

func IndexDocMapping(client *elastic.Client, indexName, docType string, body map[string]string) (*elastic.IndexResult, error) {
	resp, err := client.Index().Index(indexName).Type(docType).Id(uuid.New()).BodyJson(body).Do()
	return resp, err
}

func GetDoc(client *elastic.Client, indexName, docType, id string) (*elastic.GetResult, error) {
	resp, err := client.Get().Index(indexName).Type(docType).Id(id).Do()
	return resp, err
}
