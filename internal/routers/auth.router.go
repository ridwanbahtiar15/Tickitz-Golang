package routers

import (
	"gilangrizaltin/Backend_Golang/internal/handlers"
	"gilangrizaltin/Backend_Golang/internal/repositories"

	"github.com/gin-gonic/gin"
)

func RouterAuth(authRepo *repositories.AuthRepository, g *gin.Engine) {
	route := g.Group("/auth")
	handler := handlers.InitializeAuthHandler(authRepo)
	route.POST("/register", handler.Register)
	route.POST("/activate", handler.ActivateAccount)
	route.POST("/login", handler.Login)
	route.DELETE("/logout", handler.Logout)
}
