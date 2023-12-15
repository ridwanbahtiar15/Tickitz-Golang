package routers

import (
	"gilangrizaltin/Backend_Golang/internal/handlers"
	"gilangrizaltin/Backend_Golang/internal/middlewares"
	"gilangrizaltin/Backend_Golang/internal/repositories"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func RouterUser(authRepo *repositories.AuthRepository, g *gin.Engine, db *sqlx.DB) {
	route := g.Group("/user")
	repository := repositories.InitializeUserRepository(db)
	handler := handlers.InitializeUserHandler(repository)
	route.GET("/profile", middlewares.JWTGate(authRepo, "Admin", "Normal User"), handler.GetUserProfile)
	route.PATCH("", middlewares.JWTGate(authRepo, "Admin", "Normal User"), handler.UpdateProfileUser)
}
