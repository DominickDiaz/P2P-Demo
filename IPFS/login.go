package main

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
)

// Function to generate an access token from a user's email and password
func generateAccessToken(email string, password string) string {
	// Create a map containing the user's email and password
	claims := map[string]interface{}{
		"email":    email,
		"password": password,
	}

	// Use the JSON Web Token (JWT) library to generate the access token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims(claims))
	accessToken, _ := token.SignedString([]byte("secret_key"))

	return accessToken
}

// Function to decrypt an access token
func decryptAccessToken(accessToken string) (map[string]interface{}, error) {
	// Use the JSON Web Token (JWT) library to decode the access token
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("secret_key"), nil
	})

	// Check for errors
	if err != nil {
		return nil, err
	}

	// Check if the access token is valid
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, fmt.Errorf("Invalid access token")
	}
}

func main() {
	// Generate an access token
	accessToken := generateAccessToken("Dom@testing.com", "Dom123")
	fmt.Println("Access token:", accessToken)

	// Decrypt the access token
	claims, err := decryptAccessToken(accessToken)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Claims:", claims)
	}
}
