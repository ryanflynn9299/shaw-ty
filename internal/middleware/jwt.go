package middleware

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

var (
	// Default secret used only if environment variable is not set
	defaultSecret = []byte("your-default-secret-only-for-development")
)

// GetJWTSecret returns the JWT secret from environment or falls back to default
func GetJWTSecret() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		// Log a warning that default secret is being used
		// logger.Warn("Using default JWT secret. Set JWT_SECRET environment variable in production.")
		return defaultSecret
	}
	return []byte(secret)
}

// TokenBlacklist stores invalidated tokens and their expiry times
var tokenBlacklist = struct {
	tokens map[string]int64
	mutex  sync.RWMutex
}{
	tokens: make(map[string]int64),
}

// Claims represents the JWT claims structure
type Claims struct {
	UserID uint `json:"user_id"`
	jwt.StandardClaims
}

// GenerateToken creates a new JWT token for a user
func GenerateToken(userID uint) (string, error) {
	// Set expiration time - 24 hours from now
	expirationTime := time.Now().Add(24 * time.Hour)

	// Create claims with user ID and expiration time
	claims := &Claims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "url-shortener-api",
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with our secret
	tokenString, err := token.SignedString(GetJWTSecret())
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ParseToken validates and parses the JWT token
func ParseToken(tokenString string) (*Claims, error) {
	// Parse the token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return GetJWTSecret(), nil
	})

	if err != nil {
		return nil, err
	}

	// Extract claims
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// ExpireToken invalidates a token by adding it to the blacklist
// It accepts a token string and returns a success status and error if any
func ExpireToken(tokenString string) (bool, error) {
	// Validate the token first
	claims, err := ParseToken(tokenString)
	if err != nil {
		return false, err
	}

	// Add token to blacklist with its expiry time
	tokenBlacklist.mutex.Lock()
	tokenBlacklist.tokens[tokenString] = claims.ExpiresAt
	tokenBlacklist.mutex.Unlock()

	// Start a cleanup routine if needed
	go cleanupBlacklist()

	return true, nil
}

// AuthMiddleware validates JWT tokens in requests
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		// Check if the header has the Bearer prefix
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header must be in format: Bearer {token}"})
			c.Abort()
			return
		}

		// Extract the token
		tokenString := parts[1]

		// Parse and validate the token
		claims, err := ParseToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Set the user ID in the context for use in protected routes
		c.Set("userID", claims.UserID)

		// Continue to the next middleware/handler
		c.Next()
	}
}

// GetUserID extracts the user ID from the context after authentication
func GetUserID(c *gin.Context) (uint, bool) {
	userID, exists := c.Get("userID")
	if !exists {
		return 0, false
	}

	id, ok := userID.(uint)
	return id, ok
}

// ExtractTokenFromRequest extracts the JWT token from the Authorization header
func ExtractTokenFromRequest(c *gin.Context) (string, error) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return "", errors.New("authorization header is required")
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", errors.New("authorization header must be in format: Bearer {token}")
	}

	return parts[1], nil
}

// cleanupBlacklist periodically removes expired tokens from the blacklist
func cleanupBlacklist() {
	// Only run cleanup occasionally to avoid excessive locking
	tokenBlacklist.mutex.Lock()
	defer tokenBlacklist.mutex.Unlock()

	now := time.Now().Unix()
	for token, expiry := range tokenBlacklist.tokens {
		if expiry < now {
			delete(tokenBlacklist.tokens, token)
		}
	}
}
