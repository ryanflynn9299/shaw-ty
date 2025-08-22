package controllers

import (
	"URLShortener/api/dto"
	"URLShortener/internal/auth"
	"URLShortener/internal/middleware"
	"URLShortener/internal/services"
	"URLShortener/internal/utils"
	"database/sql"
	"encoding/hex"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

// TODO: add input sanitization and validation

type AuthController struct {
	userService *services.UserService
	argonCfg    auth.ArgonParams
}

func NewAuthController(userService *services.UserService) AuthController {
	return AuthController{
		userService: userService,
		argonCfg: auth.ArgonParams{
			Memory:      64 * 1024,
			Iterations:  10,
			Parallelism: 2,
			SaltLength:  16,
			KeyLength:   32,
		},
	}
}

// Login log in an existing user, providing a session token for further interaction with the API
func (ac *AuthController) Login(c *gin.Context) {
	loginRequest := new(dto.LoginRequest)
	err := c.ShouldBindBodyWithJSON(loginRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body, please check your request and try again."})
		return
	}

	// check for user
	user, err := (*ac.userService).FindUserByUsername(loginRequest.Username)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// compare hashes
	if auth.VerifyPassword(loginRequest.Password, user, ac.argonCfg) {
		// login successful, provide JWT token
		token, err := middleware.GenerateToken(uint(user.UUID))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "msg": "Failed to generate token, Login Failed, please try again."})
			return
		}
		resp := &dto.LoginResponse{
			Status: "success",
			Msg:    "Successfully logged in",
			Token:  token,
		}
		c.JSON(http.StatusOK, resp)

	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Incorrect password, Login Failed", "status": "failed"})
	}
}

// Register a new user, persisting the new client data and sending back a session token
func (ac *AuthController) Register(c *gin.Context) {
	registerRequest := new(dto.RegistryRequest)
	err := c.ShouldBindBodyWithJSON(registerRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body, please check your request and try again."})
		return
	}

	// check for existing(dupe) user
	duplicateUser, err := (*ac.userService).FindUserByUsername(registerRequest.Username)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if duplicateUser != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "error": "User already exists, did you mean to Login?"})
		return
	}

	// TODO: review auth workflow and add RBAP/OWASP compliance

	// generate salt
	salt, _ := auth.CreateUserSalt(ac.argonCfg.SaltLength)
	hexSalt := hex.EncodeToString(salt)

	// hash password
	pepper := utils.GetPasswordPepper()
	hashedPassword := auth.GetHashForPassword(registerRequest.Password, salt, pepper, ac.argonCfg)
	hexedHashedPassword := hex.EncodeToString(hashedPassword)

	// call user service
	user, err := (*ac.userService).RegisterUser(
		registerRequest.FirstName, registerRequest.LastName, registerRequest.Email,
		registerRequest.Username, hexedHashedPassword, hexSalt)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "status": "failed", "msg": "Failed to create user, please try again later."})
		return
	}

	// login successful, provide JWT token
	token, err := middleware.GenerateToken(uint(user.UUID))
	if err != nil {
		// TODO: generify this error message and enable RBAP error messaging (dev permissions)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "msg": "Failed to generate token, Login Failed, please try again."})
		return
	}
	resp := &dto.LoginResponse{
		Status: "success",
		Msg:    "Successfully logged in",
		Token:  token,
	}
	c.JSON(http.StatusOK, resp)
}

// Logout expires active token(s) and sends a status message
func (ac *AuthController) Logout(c *gin.Context) {
	// Extract token from request
	token, err := middleware.ExtractTokenFromRequest(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	// Expire the token
	success, err := middleware.ExpireToken(token)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Failed to invalidate token: " + err.Error(),
		})
		return
	}

	// TODO: improve this logic/messaging
	c.JSON(http.StatusOK, gin.H{
		"success": success,
		"message": "Successfully logged out",
	})
}
