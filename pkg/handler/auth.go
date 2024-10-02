package handler

import (
	"authBack/pkg/service"
	"authBack/pkg/storage"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func getToken(c *gin.Context) {
	guid := c.Query("user_id")
	clientIp := c.ClientIP()

	if guid == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id required"})
		return
	}

	accessToken, err := service.GenerateAccessToken(guid, clientIp)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate access token"})
		log.Print(err)
		return
	}

	refToken := service.GenerateRefreshToken()
	log.Print(refToken)
	hashedRefToken, err := service.HashRefreshToken(refToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate refresh token"})
		log.Print(err)

		return
	}

	err = storage.SaveRefreshToken(guid, string(hashedRefToken), clientIp)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save refresh token"})
		log.Print(err)

        return
    }

	c.JSON(http.StatusOK, gin.H{
        "access_token":  accessToken,
        "refresh_token": refToken,
    })
}

func refreshToken(c *gin.Context) {
	var request struct {
        AccessToken  string `json:"access_token"`
        RefreshToken string `json:"refresh_token"`
    }

	if err := c.ShouldBindJSON(&request); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		log.Print(err)
        return
    }

	claims, err := service.ValidateAccessToken(request.AccessToken)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid access token"})
		log.Print(err)
        return
    }

	userId := claims.UserID
    clientIP := c.ClientIP()

	storedTokenInfo, err := storage.GetRefreshTokenInfo(userId)
    if err != nil || storedTokenInfo == nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh token"})
		log.Print(err)
        return
    }

	if clientIP != storedTokenInfo.ClientIp {
        // Отправляем предупреждение по email
        service.SendWarningEmail(storedTokenInfo.Email, "IP address changed during token refresh.")
		
    }

	err = service.CompareRefreshToken(storedTokenInfo.HashedToken, request.RefreshToken)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh token"})
		log.Print(err)
        return
    }

	newAccessToken, err := service.GenerateAccessToken(userId, clientIP)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate new access token"})
		log.Print(err)
        return
    }


	c.JSON(http.StatusOK, gin.H{
        "access_token": newAccessToken,
    })
}
