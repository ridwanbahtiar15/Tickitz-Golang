package handlers

import (
	"fmt"
	"gilangrizaltin/Backend_Golang/internal/helpers"
	"gilangrizaltin/Backend_Golang/internal/models"
	"gilangrizaltin/Backend_Golang/internal/repositories"
	"gilangrizaltin/Backend_Golang/pkg"
	"log"
	"net/http"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
)

type HandlerUser struct {
	repositories.IUserRepository
}

func InitializeUserHandler(r repositories.IUserRepository) *HandlerUser {
	return &HandlerUser{r}
}

func (h *HandlerUser) GetUserProfile(ctx *gin.Context) {
	id, _ := helpers.GetPayload(ctx)
	result, err := h.RepositoryGetUserProfile(id)
	if err != nil {
		log.Print(err.Error())
		ctx.JSON(http.StatusInternalServerError, helpers.NewResponse("Internal Server Error", nil, nil))
		return
	}
	if len(result) == 0 {
		ctx.JSON(http.StatusNotFound, helpers.NewResponse("Data user not found", nil, nil))
		return
	}
	ctx.JSON(http.StatusOK, helpers.NewResponse("Successfully Get Profile user", result, nil))
}

func (h *HandlerUser) UpdateProfileUser(ctx *gin.Context) {
	var body models.UserProfileUpdateModel
	err := ctx.ShouldBind(&body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, helpers.NewResponse("Error in binding body update", nil, nil))
		return
	}
	ID, _ := helpers.GetPayload(ctx)
	if _, err := govalidator.ValidateStruct(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, helpers.NewResponse("Input not valid", nil, nil))
		return
	}
	var newPassword string
	if body.New_Password != "" {
		if body.Last_Password == "" {
			ctx.JSON(http.StatusInternalServerError, helpers.NewResponse("Please input last password to verify", nil, nil))
			return
		}
		result, err := h.RepositorySensitiveData(ID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, helpers.NewResponse("Error in get sensitive data", nil, nil))
			return
		}
		hs := pkg.HashConfig{
			Time:    3,
			Memory:  64 * 1024,
			Threads: 2,
			KeyLen:  32,
			SaltLen: 16,
		}
		isValid, _ := hs.ComparePasswordAndHash(body.Last_Password, result[0].Password)
		if !isValid {
			ctx.JSON(http.StatusUnauthorized, helpers.NewResponse("Last password doesnt match", nil, nil))
			return
		}
		hashedPassword, errHash := hs.GenHashedPassword(body.New_Password)
		if errHash != nil {
			log.Println(errHash.Error())
			ctx.JSON(http.StatusInternalServerError, helpers.NewResponse("Error hashing password", nil, nil))
			return
		}
		newPassword = hashedPassword
	}
	cld, err := helpers.InitCloudinary()
	if err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, helpers.NewResponse("Error in initialize uploading image system", nil, nil))
		return
	}
	formFile, _ := ctx.FormFile("user_photo")
	var dataUrl string
	if formFile != nil {
		file, err := formFile.Open()
		if err != nil {
			log.Println(err.Error())
			ctx.JSON(http.StatusInternalServerError, helpers.NewResponse("Error in opening file", nil, nil))
			return
		}
		defer file.Close()
		publicId := fmt.Sprintf("%s_%s-%d", "users", "photo_profile", ID)
		folder := "Tickitz"
		res, err := cld.Uploader(ctx, file, publicId, folder)
		if err != nil {
			log.Println(err.Error())
			ctx.JSON(http.StatusInternalServerError, helpers.NewResponse("Error in upload image", nil, nil))
			return
		}
		dataUrl = res.SecureURL
	}
	result, errUpdate := h.RepositoryUpdateUser(ID, &body, dataUrl, newPassword)
	if errUpdate != nil {
		log.Println(errUpdate.Error())
		ctx.JSON(http.StatusInternalServerError, helpers.NewResponse("Internal Server Error", nil, nil))
		return
	}
	var data int64 = 1
	if result < data {
		ctx.JSON(http.StatusNotFound, helpers.NewResponse("User not found", nil, nil))
		return
	}
	ctx.JSON(http.StatusOK, helpers.NewResponse("Successfully update user", &body, nil))
}
