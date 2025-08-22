package auth

import (
	"URLShortener/internal/storage/models"
	"URLShortener/internal/utils"
	"crypto/rand"
	"crypto/subtle"
	"encoding/hex"
	"golang.org/x/crypto/argon2"
	"log"
)

type ArgonParams struct {
	Memory      uint32
	Iterations  uint32
	Parallelism uint8
	SaltLength  uint32
	KeyLength   uint32
}

// CreateUserSalt: generate a random n-length byte array that represents a User's salt
func CreateUserSalt(n uint32) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

// GetHashForPassword encrypt a user's new password for persistence
func GetHashForPassword(password string, salt []byte, pepper string, params ArgonParams) (hash []byte) {
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

// GetSaltAndHashedPassword: generate a salt and hash a user's new password for persistence used for migrations only
func GetSaltAndHashedPassword(password string) (string, string) {
	saltLength := 16
	salt, _ := CreateUserSalt(uint32(saltLength))
	pepper := utils.GetPasswordPepper()
	params := ArgonParams{Memory: 64 * 1024, Iterations: 10, Parallelism: 2, SaltLength: uint32(saltLength), KeyLength: 32}
	hashedPassword := GetHashForPassword(password, salt, pepper, params)

	// store as hex strings for direct persistence. This method should not be used in the application
	return hex.EncodeToString(salt), hex.EncodeToString(hashedPassword)
}

// VerifyPassword checks an incoming password against the user's stored password hash
func VerifyPassword(password string, user *models.User, params ArgonParams) bool {
	pepper := utils.GetPasswordPepper()
	decodeSalt, err := hex.DecodeString(user.Salt)
	if err != nil {
		log.Println("Could not parse salt from hex representation during authentication.")
		return false
	}
	hashedAttempt := GetHashForPassword(password, decodeSalt, pepper, params)
	decodePassword, err := hex.DecodeString(user.Password)
	if err != nil {
		log.Println("Could not parse password from hex representation during authentication.")
		return false
	}

	// Perform a timed comparison to avoid a timing attack
	return subtle.ConstantTimeCompare(hashedAttempt, decodePassword) == 1
}
