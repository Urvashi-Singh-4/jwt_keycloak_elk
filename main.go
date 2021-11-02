package main

import (
	"token_based_auth/logger"
	"token_based_auth/model"
	"token_based_auth/routes"
)

func main() {
	logger.InitLogger()
	model.ElasticSearchClient()

	app := routes.NewRoute()
	defer logger.CloseStash()

	app.Run("localhost:9000")
}
