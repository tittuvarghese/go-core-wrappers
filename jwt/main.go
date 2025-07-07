package jwt

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/tittuvarghese/go-core-wrappers/constants"
	"strings"
	"time"
)

var jwtSecretKey = []byte(constants.JwtSecretKey)

type Claims struct {
	Data interface{} `json:"data"`
	jwt.RegisteredClaims
}

// Generate generates a new JWT token
func Generate(payload interface{}, issuer string, expiry time.Duration) (string, error) {
	// Set expiration time
	expirationTime := time.Now().Add(expiry * time.Hour)

	claims := Claims{
		Data: payload,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			Issuer:    issuer,
		},
	}

	// Create the token with the claims and sign it with the secret key
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token and return the signed string
	tokenString, err := token.SignedString(jwtSecretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateJWT validates the JWT token
func validateJWT(tokenString string) (*Claims, error) {
	// Parse the token and check the signature
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Ensure the token's method matches the expected signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		// Return the secret key used to sign the token
		return jwtSecretKey, nil
	})

	if err != nil {
		return nil, err
	}

	// Check if the token is valid and extract the claims
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, fmt.Errorf("invalid token")
	}
}

// Authorize to verify JWT token and extract claims
func Authorize() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get token from the Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(401, gin.H{"error": "Authorization header missing"})
			c.Abort()
			return
		}

		// Bearer <token>
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(401, gin.H{"error": "Invalid authorization format"})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// Parse and validate the token
		claims, err := validateJWT(tokenString)
		if err != nil {
			c.JSON(401, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Set the claims in the context
		c.Set("claims", claims)

		// Proceed with the next handler
		c.Next()
	}
}

func GetClaims(c *gin.Context) (*Claims, error) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return nil, fmt.Errorf("authorization header missing")
	}

	// Bearer <token>
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return nil, fmt.Errorf("invalid authorization format")
	}
	tokenString := parts[1]
	// Parse and validate the token
	claims, err := validateJWT(tokenString)

	if err != nil {
		return nil, err
	}

	return claims, err
}
