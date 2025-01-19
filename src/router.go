package src

import (
	"context"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/cbartram/hearthhub/src/handlers"
	"github.com/cbartram/hearthhub/src/handlers/cognito"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"log"
	"os"
)

func MakeRouter(ctx context.Context) *ginadapter.GinLambda {
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: false,
	})

	logLevel, err := logrus.ParseLevel(os.Getenv("LOG_LEVEL"))
	if err != nil {
		logLevel = logrus.InfoLevel
	}

	log.SetOutput(os.Stdout)
	logrus.SetLevel(logLevel)

	r := gin.New()

	gin.DefaultWriter = logger.Writer()
	gin.DefaultErrorWriter = logger.Writer()
	gin.SetMode(gin.ReleaseMode)

	r.Use(LogrusMiddleware(logger))

	apiGroup := r.Group("/api/v1")
	cognitoGroup := apiGroup.Group("/cognito")

	apiGroup.POST("/discord/oauth", func(c *gin.Context) {
		handler := handlers.DiscordRequestHandler{}
		handler.HandleRequest(c, ctx)
	})

	cognitoGroup.POST("/create-user", func(c *gin.Context) {
		handler := cognito.CognitoCreateUserRequestHandler{}
		handler.HandleRequest(c, ctx)
	})

	cognitoGroup.POST("/auth", func(c *gin.Context) {
		handler := cognito.CognitoAuthHandler{}
		handler.HandleRequest(c, ctx)
	})

	cognitoGroup.POST("/refresh-session", func(c *gin.Context) {
		handler := cognito.CognitoRefreshSessionHandler{}
		handler.HandleRequest(c, ctx)
	})

	cognitoGroup.GET("/get-user", func(c *gin.Context) {
		handler := cognito.CognitoGetUserHandler{}
		handler.HandleRequest(c, ctx)
	})

	return ginadapter.New(r)
}
