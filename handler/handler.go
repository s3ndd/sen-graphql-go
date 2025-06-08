package handler

import (
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"

	"github.com/s3ndd/sen-graphql-go/graph/generated"
	"github.com/s3ndd/sen-graphql-go/graph/resolver"
)

// GraphqlHandler defines the Graphql handler
func GraphqlHandler() gin.HandlerFunc {
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &resolver.Resolver{}}))

	return func(ctx *gin.Context) {

		srv.ServeHTTP(ctx.Writer, ctx.Request)
	}
}

// PlaygroundHandler defines the Playground handler
func PlaygroundHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		playground.Handler("GraphQL playground", "/graphql/query").ServeHTTP(ctx.Writer, ctx.Request)
	}
}

// HealthCheck returns 200 and success
func HealthCheck() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, "success")
	}
}
