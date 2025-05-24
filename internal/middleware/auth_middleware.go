package middleware

import (
	"DH52111659-api-quan-ly-suc-khoe/common"
	"DH52111659-api-quan-ly-suc-khoe/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func JWTAuthMiddleware(tokenService utils.TokenService, requireRoles ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.JSON(http.StatusUnauthorized, common.NewResponseError("Authorization header is required"))
			ctx.Abort()
			return
		}

		if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
			ctx.JSON(http.StatusUnauthorized, common.NewResponseError("Token must be in Bearer format"))
			return
		}

		tokenString := authHeader[len("Bearer "):]
		claims, err := tokenService.VerifyToken(tokenString)
		if err != nil {
			ctx.JSON(401, common.NewResponseError("Invalid token"))
			ctx.Abort()
			return
		}

		// Check if the user has the required role
		roleIsValid := false
		for _, role := range requireRoles {
			if claims.Role == role {
				roleIsValid = true
				break
			}
		}

		if !roleIsValid && len(requireRoles) > 0 {
			ctx.JSON(http.StatusForbidden, common.NewResponseError("You do not have permission to access this resource"))
			ctx.Abort()
			return
		}

		// Store the user ID in the context for later use
		ctx.Set("userID", claims.UserID)
		ctx.Set("role", claims.Role)
		ctx.Next()
	}
}
