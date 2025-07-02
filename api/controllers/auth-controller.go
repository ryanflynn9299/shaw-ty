package controllers

import (
	"URLShortener/api/dto"
	"URLShortener/internal/middleware"
	"URLShortener/internal/services"
	"URLShortener/internal/storage/models"
	"URLShortener/internal/utils"
	"crypto/rand"
	"crypto/subtle"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/argon2"
	"net/http"
)

type AuthController struct {
	userService *services.UserService
	argonCfg    argonParams
}

type argonParams struct {
	Memory      uint32
	Iterations  uint32
	Parallelism uint8
	SaltLength  uint32
	KeyLength   uint32
}

func NewAuthController(userService *services.UserService) AuthController {
	return AuthController{
		userService: userService,
		argonCfg: argonParams{
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
	// check for user
	user, err := (*ac.userService).FindUserByUsername(c.Param("username"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// compare hashes
	if verifyPassword(c.Param("password"), user, ac.argonCfg) {
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
	// check for existing(dupe) user
	duplicateUser, err := (*ac.userService).FindUserByUsername(c.Param("username"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if duplicateUser != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User already exists, did you mean to Login?"})
		return
	}

	// generate salt
	salt, _ := getUserSalt(ac.argonCfg.SaltLength)

	// hash password
	pepper := utils.GetPasswordPepper()
	hashedPassword := getHashForPassword(c.Param("password"), salt, pepper, ac.argonCfg)

	// call user service
	user, err := (*ac.userService).RegisterUser(c.Param("firstname"), c.Param("lastname"), c.Param("email"), string(hashedPassword), salt)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "status": "failed", "msg": "Failed to create user, please try again later."})
		return
	}

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

	c.JSON(http.StatusOK, gin.H{
		"success": success,
		"message": "Successfully logged out",
	})
}

// getUserSalt: generate a random n-length byte array that represents a User's salt
func getUserSalt(n uint32) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

// getHashForPassword: encrypt a user's new password for persistence
func getHashForPassword(password string, salt []byte, pepper string, params argonParams) (hash []byte) {
	pepperedPassword := []byte(password + pepper)
	hash = argon2.IDKey(
		pepperedPassword,
		salt,
		params.Iterations,
		params.Memory,
		params.Parallelism,
		params.KeyLength)

	return hash
}

// verifyPassword checks an incoming password against the user's stored password hash
func verifyPassword(password string, user *models.User, params argonParams) bool {
	pepper := utils.GetPasswordPepper()
	hashedAttempt := getHashForPassword(password, user.Salt, pepper, params)

	// Perform a timed comparison to avoid a timing attack
	return subtle.ConstantTimeCompare(hashedAttempt, []byte(user.Password)) == 1
}
