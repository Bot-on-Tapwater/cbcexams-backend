package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var JWT_SECRET = os.Getenv("JWT_SECRET")

// GenerateJWT generates a JSON Web Token (JWT) for the given user ID.
// The token is signed using the HS256 signing method and includes the
// user ID as a claim, along with an expiration time set to 72 hours
// (3 days) from the time of generation.
//
// Parameters:
//   - userID: A UUID representing the unique identifier of the user.
//
// Returns:
//   - A string containing the signed JWT.
//   - An error if the token signing process fails.
func GenerateJWT(userID uuid.UUID) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp": time.Now().Add(time.Hour * 72).Unix(), /* 3 days expiry */
	})

	return token.SignedString([]byte(JWT_SECRET))
}

// ValidateJWT validates a given JWT token string and returns the parsed token
// if it is valid, or an error if the validation fails.
// 
// Parameters:
//   - tokenString: The JWT token string to be validated.
//
// Returns:
//   - *jwt.Token: The parsed JWT token if validation is successful.
//   - error: An error if the token is invalid or if there is an issue during parsing.
//
// The function checks if the token's signing method is HMAC and uses the
// predefined JWT_SECRET to validate the token's signature.
func ValidateJWT(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface {}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(JWT_SECRET), nil
	})
}