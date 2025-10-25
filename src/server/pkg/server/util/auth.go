package util

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"time"
	"speak-sphere/pkg/server/conf"
)

func GenerateSalt() string {
	salt := make([]byte, 32)
	_, _ = rand.Read(salt)
	return hex.EncodeToString(salt)
}

func HashPassword(password, salt string) string {
	hash := sha256.New()
	hash.Write([]byte(password + salt))
	return hex.EncodeToString(hash.Sum(nil))
}

func GenerateJWT(userID int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24 * 365).Unix(),
	})
	return token.SignedString([]byte(conf.Cfg.JWT.Secret))
}

func GetUserID(c *gin.Context) (int, error) {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0, errors.New("invalid user id")
	} else {
		res, _ := userID.(uint)
		return int(res), nil
	}
}
