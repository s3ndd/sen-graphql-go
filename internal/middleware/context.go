package middleware

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
)

type ctxKeyType struct{}

var ginCtxKey interface{} = ctxKeyType{}

// GinContextToContextMiddleware adds gin context to the context.Context
func GinContextToContextMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.WithValue(c.Request.Context(), ginCtxKey, c)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

// GinContextFromContext recovers the gin.Context from the context.Context struct
func GinContextFromContext(ctx context.Context) (*gin.Context, error) {
	ginContext := ctx.Value(ginCtxKey)
	if ginContext == nil {
		err := fmt.Errorf("could not retrieve gin.Context")
		return nil, err
	}

	gc, ok := ginContext.(*gin.Context)
	if !ok {
		err := fmt.Errorf("gin.Context has wrong type")
		return nil, err
	}
	return gc, nil
}

func GetAuthorizationHeader(ctx context.Context) string {
	if gc, err := GinContextFromContext(ctx); err == nil {
		return gc.GetHeader("Authorization")
	}
	return ""

}
