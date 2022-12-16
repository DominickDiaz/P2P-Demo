package sec

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

func EncryptPass(password *string, secretKey *string) string {
	//   Encode the string and  secret key as bytes
	passBytes := []byte(*password)
	secretKeyBytes := []byte(*secretKey)

	// Compute the HMAC (Hash-based Message Authentication Code) of the string
	// using the secret key
	hmac := hmac.New(sha256.New, secretKeyBytes)
	hmac.Write(passBytes)
	hmacBytes := hmac.Sum(nil)

	// Return the hexadecimal representation of the HMAC as the encrypted string
	return hex.EncodeToString(hmacBytes)
}

// Encryption Method Tests below

//func GenerateUserKey(username string, password string) string {
//	// Combine username and password
//	combined := username + password
//	UUID := "8a6d8b184dde4a5a93eb08482b06f1d141f37c0aa390481ca37021624c660d220f40b7773fee416daf184e2c164256f79a33b35693b44b48894f3f53961fb28c"
//
//	return combined + UUID[:120-len(combined)]
//}

//
//func Encrypt(username string, password string) string {
//	secretKey := GenerateUserKey(username, password)
//	return EncryptPass(&password, &secretKey)
//}
