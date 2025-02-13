package model

import (
	"context"
	"github.com/gin-gonic/gin"
)

type RequestHandler interface {
	HandleRequest(c *gin.Context, ctx context.Context)
}

// DiscordTokenResponse represents Discord's token response
type DiscordTokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
}

type CognitoCreateUserRequest struct {
	DiscordID       string `json:"discord_id"`
	DiscordUsername string `json:"discord_username"`
	DiscordEmail    string `json:"discord_email"`
	AvatarId        string `json:"avatar_id"`
}

type CognitoUserStatusRequest struct {
	AccountEnabled bool   `json:"accountEnabled"`
	DiscordID      string `json:"discordId"`
}

type CognitoAuthRequest struct {
	RefreshToken string `json:"refreshToken"`
	DiscordID    string `json:"discordId"`
}

type SimpleS3Object struct {
	Key  string `json:"key"`
	Size int64  `json:"fileSize"`
}

type CognitoCredentials struct {
	RefreshToken    string `json:"refresh_token,omitempty"`
	TokenExpiration int32  `json:"token_expiration_seconds,omitempty"`
	AccessToken     string `json:"access_token,omitempty"`
	IdToken         string `json:"id_token,omitempty"`
}

type CognitoUser struct {
	CognitoID        string             `json:"cognitoId,omitempty"`
	DiscordUsername  string             `json:"discordUsername,omitempty"`
	Email            string             `json:"email,omitempty"`
	AvatarId         string             `json:"avatarId"`
	DiscordID        string             `json:"discordId,omitempty"`
	InstalledMods    map[string]bool    `json:"installedMods"`
	InstalledBackups map[string]bool    `json:"installedBackups"`
	AccountEnabled   bool               `json:"accountEnabled,omitempty"`
	Credentials      CognitoCredentials `json:"credentials,omitempty"`
}
