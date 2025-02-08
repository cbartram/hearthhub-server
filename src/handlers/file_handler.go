package handlers

import (
	"context"
	"fmt"
	"github.com/cbartram/hearthhub/src/service"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

type FileHandler struct{}

// HandleRequest Handles the request for listing files under a given prefix. Since this route is deployed
// to a lambda function and backed by the Cognito Authorizer only authorized users can invoke this.
func (f *FileHandler) HandleRequest(c *gin.Context, s3Client *service.S3Service) {
	discordId := c.Query("discordId")
	refreshToken := c.Query("refreshToken")
	prefix := c.Query("prefix")

	if discordId == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "discordId query parameter is required",
		})
		return
	}

	if refreshToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "refreshToken query parameter is required",
		})
		return
	}

	// Anyone can list mods, only users can list their own backups and configuration
	// Therefore, we need to verify that the given discord id belongs to the refresh token that was given
	// in the Authorization header
	authManager := service.MakeCognitoService()
	log.Infof("authenticating user with discord id: %s", discordId)
	_, user := authManager.AuthUser(context.Background(), &refreshToken, &discordId)

	if user.DiscordID != discordId {
		log.Errorf("authenticated user id: %s does not match given discord id: %s", user.DiscordID, discordId)
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized: authenticated user id does not match given discord id",
		})
		return
	}

	// valid prefixes are stored in file_upload_handler.go and essentially are just:
	// config, backups, mods to direct the s3 operation at where to list or put user files
	_, ok := ValidPrefixes[prefix]
	if !ok {
		log.Errorf("invalid prefix: %s", prefix)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("invalid prefix: %s", prefix),
		})
		return
	}

	var sanitizedPrefix = prefix
	if strings.HasSuffix("/", prefix) {
		sanitizedPrefix = prefix[0 : len(prefix)-1]
	}

	path := fmt.Sprintf("%s/%s/", sanitizedPrefix, discordId)
	log.Infof("prefix is sanitized and valid: %s, listing objects for path: %s", sanitizedPrefix, path)

	objs, err := s3Client.ListObjects(context.Background(), path)
	if err != nil {
		log.Errorf("failed to list objects: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("failed to list objects: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"files": objs,
	})
}
