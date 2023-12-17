package middlewares

import (
	"gilangrizaltin/Backend_Golang/internal/helpers"
	"gilangrizaltin/Backend_Golang/internal/repositories"
	"gilangrizaltin/Backend_Golang/pkg"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func JWTGate(authRepo *repositories.AuthRepository, allowedRole ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		bearerToken := ctx.GetHeader("Authorization")
		if bearerToken == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, helpers.NewResponse("Please login first", nil, nil))
			return
		}
		if !strings.Contains(bearerToken, "Bearer") {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, helpers.NewResponse("Please login again", nil, nil))
			return
		}

		token := strings.Replace(bearerToken, "Bearer ", "", -1)
		result, err := authRepo.RepositoryIsTokenBlacklisted(token)
		if err != nil {
			log.Println(err.Error())
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, helpers.NewResponse("Internal Server Token Error", nil, nil))
			return
		}
		if result {
			ctx.AbortWithStatusJSON(http.StatusForbidden, helpers.NewResponse("You have logout. Please login again", nil, nil))
			return
		}
		payload, err := pkg.VerifyToken(token)
		if err != nil {
			if strings.Contains(err.Error(), "expired") {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, helpers.NewResponse("Please login again", nil, nil))
				return
			}
			log.Println(err.Error())
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, helpers.NewResponse("Error verifying token", nil, nil))
			return
		}

		var allowed = false
		for _, role := range allowedRole {
			if payload.Role == role {
				allowed = true
				break
			}
		}
		if !allowed {
			ctx.AbortWithStatusJSON(http.StatusForbidden, helpers.NewResponse("Access denied", nil, nil))
			return
		}
		ctx.Set("Payload", payload)
		ctx.Next()
	}
}
