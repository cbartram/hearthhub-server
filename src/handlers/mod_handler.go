package handlers

import (
	"context"
	"github.com/cbartram/hearthhub/src/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ModHandler struct{}

// HandleRequest Handles the request for listing files under a given prefix. Since this route is deployed
// to a lambda function and backed by the Cognito Authorizer only authorized users can invoke this.
func (m *ModHandler) HandleRequest(c *gin.Context, s3 *service.S3Service) {
	discordId := c.Query("discordId")

	if discordId == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "discordId query parameter is required",
		})
		return
	}

	s3.ListObjects(context.Background())
}
