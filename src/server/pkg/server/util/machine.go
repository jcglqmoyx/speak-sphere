package util

import (
	"crypto/sha512"
	"fmt"
)

func generateSecretKey(s string) string {
	secretKey := ""
	for i := 0; i < 10; i++ {
		hash := sha512.New()
		hash.Write([]byte(s))
		hashedBytes := hash.Sum(nil)
		s = fmt.Sprintf("%x", hashedBytes)
		secretKey = s
	}
	return secretKey
}
