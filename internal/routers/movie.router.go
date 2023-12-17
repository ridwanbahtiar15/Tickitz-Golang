package routers

import (
	"gilangrizaltin/Backend_Golang/internal/handlers"
	"gilangrizaltin/Backend_Golang/internal/middlewares"
	"gilangrizaltin/Backend_Golang/internal/repositories"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func RouterMovie(authRepo *repositories.AuthRepository, g *gin.Engine, db *sqlx.DB) {
	route := g.Group("/movie")
	repository := repositories.InitializeMovieRepository(db)
	handler := handlers.InitializeMovieHandler(repository)
	route.GET("", middlewares.JWTGate(authRepo, "Admin", "Normal User"), handler.GetAllMovie)
	route.GET("/movie/:movie_id", middlewares.JWTGate(authRepo, "Admin", "Normal User"), handler.GetMovieDetails)
	route.GET("/schedule/:movie_id", middlewares.JWTGate(authRepo, "Admin", "Normal User"), handler.GetMovieSchedule)
	route.POST("", middlewares.JWTGate(authRepo, "Admin"), handler.AddMovie)
	route.PATCH("/:movie_id", middlewares.JWTGate(authRepo, "Admin"), handler.UpdateMovie)
	route.DELETE("/:movie_id", middlewares.JWTGate(authRepo, "Admin"), handler.DeleteMovie)
}
