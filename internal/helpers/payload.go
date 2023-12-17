package helpers

import (
	"fmt"
	"gilangrizaltin/Backend_Golang/pkg"

	"github.com/gin-gonic/gin"
)

func GetPayload(ctx *gin.Context) (id int, role string) {
	payload, exists := ctx.Get("Payload")
	if !exists {
		// ctx.JSON(http.StatusUnauthorized, gin.H{
		// 	"message": "dont have token",
		// })
		fmt.Println("dont have token")
		return
	}
	data := payload.(*pkg.Claims)
	return data.Id, data.Role
}
