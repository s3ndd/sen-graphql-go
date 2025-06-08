package router

import (
	"github.com/gin-gonic/gin"

	"github.com/s3ndd/sen-graphql-go/graph/dataloader"
	"github.com/s3ndd/sen-graphql-go/handler"
	"github.com/s3ndd/sen-graphql-go/internal/middleware"
)

// InitRouter initialize routing information
func InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery()).Use(middleware.GinContextToContextMiddleware()).Use(dataloader.LoaderMiddleware())

	apiV1 := r.Group("/api/v1")
	apiV1.GET("/healthcheck", handler.HealthCheck())
	aipGraphql := r.Group("/graphql")
	aipGraphql.POST("/query", handler.GraphqlHandler())
	aipGraphql.GET("/graphiql", handler.PlaygroundHandler())

	return r
}
