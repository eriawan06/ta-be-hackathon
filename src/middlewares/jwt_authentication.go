package middlewares

import (
	um "be-sagara-hackathon/src/modules/user"
	"be-sagara-hackathon/src/utils"
	"be-sagara-hackathon/src/utils/common"
	"fmt"
	"math"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// JwtAuthMiddleware Token Authentication
func JwtAuthMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {

		tokenString := utils.ExtractToken(context)

		// Check token string
		if tokenString == "" {
			common.SendError(context, http.StatusUnauthorized, "Unauthorized", []string{"Authentication Token Required"})
			context.Abort()
			return
		}

		key := os.Getenv("API_JWT_SECRET")
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Check sign in method token
			if jwt.GetSigningMethod("HS256") != token.Method {
				// When sign in method not same
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			// return key
			return []byte(key), nil
		})
		if err != nil {
			common.SendError(context, http.StatusUnauthorized, "Unauthorized", []string{err.Error()})
			context.Abort()
			return
		}

		if !token.Valid {
			common.SendError(context, http.StatusUnauthorized, "Unauthorized", []string{"Invalid token"})
			context.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			common.SendError(context, http.StatusUnauthorized, "Unauthorized", []string{"Invalid token"})
			context.Abort()
			return
		}

		// Check if token expired
		if time.Now().Unix() > int64(math.Round(claims["expired"].(float64))) {
			common.SendError(context, http.StatusUnauthorized, "Unauthorized", []string{"Token expired"})
			context.Abort()
			return
		}

		email := fmt.Sprintf("%v", claims["email"])
		user, err := um.GetUserRepository().FindByEmail(email)
		if err != nil {
			common.SendError(context, http.StatusUnauthorized, "Unauthorized", []string{"User not found"})
			context.Abort()
			return
		}

		// Next
		context.Set("userID", user.ID)
		context.Set("user", user)
		context.Set("token", tokenString)
		context.Next()
	}
}
