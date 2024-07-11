package handler

import (
	"auth-service/auth/token"
	pb "auth-service/generated/auth_service"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// RegisterHandler handles user registration
// @Summary Register a new user
// @Description Register a new user with email and password
// @Tags Auth
// @Accept json
// @Produce json
// @Param Register body auth_service.RegisterRequest true "User Registration"
// @Success 201 {object} models.Success
// @Failure 400 {object} models.Errors
// @Router /auth/register [post]
func (h *Handler) RegisterHandler(ctx *gin.Context) {
	user := pb.RegisterRequest{}

	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"Error": err.Error(),
		})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"Error": err.Error(),
		})
		return
	}

	user.Password = string(hashedPassword)

	resp, err := h.UserRepo.CreateUser(&user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"Error": err.Error(),
		})
		return
	}
	
	ctx.JSON(http.StatusCreated, gin.H{
		"Message": resp.Message,
	})
}

// LoginHandler handles user login
// @Summary Login a user
// @Description Login a user with email and password
// @Tags Auth
// @Accept json
// @Produce json
// @Param Login body auth_service.LoginRequest true "User Login"
// @Success 200 {object} models.Token
// @Failure 400 {object} models.Errors
// @Failure 404 {object} models.Errors
// @Failure 500 {object} models.Errors
// @Router /auth/login [post]
func (h *Handler) LoginHandler(ctx *gin.Context) {
	user := pb.LoginRequest{}

	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"Error": err.Error(),
		})
		return
	}

	storedUser, err := h.UserRepo.GetByEmail(user.Email)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"Error": err.Error(),
		})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(user.Password)); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"Error": err.Error(),
		})
		return
	}

	accessToken, err := token.GenerateAccessJWT(storedUser)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate access token"})
		return
	}

	refreshToken, err := token.GenerateRefreshJWT(storedUser)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate refresh token"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

// RefreshToken generates a new access token using a refresh token
// @Summary Refresh access token
// @Description Refresh the access token using the refresh token
// @Tags Auth
// @Accept json
// @Produce json
// @Param Authorization header string true "Refresh token"
// @Success 200 {object} models.Request
// @Failure 400 {object} models.Errors
// @Failure 401 {object} models.Errors
// @Failure 500 {object} models.Errors
// @Security ApiKeyAuth
// @Router /auth/refresh_token [get]
func (h *Handler) RefreshToken(c *gin.Context) {
	refreshToken := c.GetHeader("Authorization")
	if refreshToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Refresh token is required"})
		return
	}

	claims, err := token.ExtractClaim(refreshToken)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid token"})
		return
	}

	if claims.ExpiresAt < time.Now().Unix() {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token has expired"})
		return
	}

	newAccessToken, err := token.GenerateAccessJWT(&pb.LoginResponse{
		UserId:   claims.UserId,
		Username: claims.Username,
		Email:    claims.Email,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate new access token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"access_token": newAccessToken})
}
