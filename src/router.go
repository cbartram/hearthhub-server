package src

import (
	"context"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/cbartram/hearthhub/src/handlers"
	"github.com/cbartram/hearthhub/src/handlers/cognito"
	"github.com/cbartram/hearthhub/src/service"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"log"
	"os"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		logrus.Infof("setting CORS response headers")
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

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

	s3, err := service.MakeS3Service("us-east-1")
	if err != nil {
		logrus.Errorf("failed to create s3 client: %v", err)
	}

	apiGroup := r.Group("/api/v1", CORSMiddleware())
	cognitoGroup := apiGroup.Group("/cognito", CORSMiddleware())

	apiGroup.POST("/discord/oauth", func(c *gin.Context) {
		handler := handlers.DiscordRequestHandler{}
		handler.HandleRequest(c, ctx)
	})

	apiGroup.GET("/file", func(c *gin.Context) {
		handler := handlers.FileHandler{}
		handler.HandleRequest(c, s3)
	})

	apiGroup.POST("/file/upload", func(c *gin.Context) {
		handler := handlers.UploadFileHandler{}
		handler.HandleRequest(c, s3)
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
