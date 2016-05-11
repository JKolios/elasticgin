# elasticgin
A mini data warehouse in Golang using Elasticsearch and Gin.
Targets Elasticsearch v2.x.
Optionally supports consuming documents from AMQP queues.

##Usage
From the base directory:
```
go build
cp config/config_sample.json ./config.json
```
Edit config.json as needed.
```
./elasticgin
```

