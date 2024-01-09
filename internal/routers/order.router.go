package routers

import (
	"gilangrizaltin/Backend_Golang/internal/handlers"
	"gilangrizaltin/Backend_Golang/internal/middlewares"
	"gilangrizaltin/Backend_Golang/internal/repositories"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func RouterOrder(authRepo *repositories.AuthRepository, g *gin.Engine, db *sqlx.DB) {
	route := g.Group("/order")
	repository := repositories.InitializeOrderRepository(db)
	handler := handlers.InitializeOrderHandler(repository)
	route.GET("", middlewares.JWTGate(authRepo, "Admin", "Normal User"), handler.GetOrder)
	route.GET("/stat", middlewares.JWTGate(authRepo, "Admin"), handler.GetOrderStatistic)
	route.GET("/:schedule_id", middlewares.JWTGate(authRepo, "Admin", "Normal User"), handler.GetDetailSchedule)
	route.POST("", middlewares.JWTGate(authRepo, "Admin", "Normal User"), handler.CreateTransaction)
	route.POST("/success", handler.SubmitPayment)
	route.POST("/failed", handler.FailedPayment)
}
