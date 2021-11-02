package model

import (
	"log"

	elastic "github.com/olivere/elastic/v7"
)

var DB *elastic.Client

func ElasticSearchClient() {
	var err error
	DB, err = elastic.NewClient(elastic.SetURL("http://localhost:9200"),
		elastic.SetSniff(false),
		elastic.SetHealthcheck(false))

	if err != nil {
		log.Fatalln("Error occured while initializing Elastic Client")

	}

}
