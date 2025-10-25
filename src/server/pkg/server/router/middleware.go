package router

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"speak-sphere/pkg/server/conf"
	"speak-sphere/pkg/server/util"
)

func CorsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
			return
		}
		c.Next()
	}
}

func shouldPassThrough(req *http.Request) bool {
	path := req.URL.Path
	if path == "/" || path == "/user/register" || path == "/user/login" {
		return true
	}
	return false
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if shouldPassThrough(c.Request) {
			c.Next()
			return
		}
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			util.JsonHttpResponse(c, 2, "Authorization header is missing", nil)
			c.Abort()
			return
		}

		if tokenString[:7] != "Bearer " {
			util.JsonHttpResponse(c, 2, "Invalid token", nil)
			c.Abort()
			return
		}

		tokenString = tokenString[len("Bearer "):]

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(conf.Cfg.JWT.Secret), nil
		})

		if err != nil {
			util.JsonHttpResponse(c, 2, "Invalid token", nil)
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			util.JsonHttpResponse(c, 2, "Invalid token", nil)
			c.Abort()
			return
		}
		if userID, ok := claims["user_id"].(float64); ok {
			c.Set("user_id", uint(userID))
		} else {
			util.JsonHttpResponse(c, 2, "user_id is not an int", nil)
			c.Abort()
			return
		}
		c.Next()
	}
}
