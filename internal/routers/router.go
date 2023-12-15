package routers

import (
	"gilangrizaltin/Backend_Golang/internal/middlewares"
	"gilangrizaltin/Backend_Golang/internal/repositories"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func New(db *sqlx.DB) *gin.Engine {
	router := gin.Default()
	router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "success",
		})
	})
	authRepo := repositories.InitializeAuthRepository(db)
	router.Use(middlewares.CORSMiddleware)
	RouterAuth(authRepo, router)
	return router
}
