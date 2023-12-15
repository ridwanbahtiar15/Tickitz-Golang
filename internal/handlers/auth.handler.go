package handlers

import (
	"fmt"
	"gilangrizaltin/Backend_Golang/internal/helpers"
	"gilangrizaltin/Backend_Golang/internal/models"
	"gilangrizaltin/Backend_Golang/internal/repositories"
	"gilangrizaltin/Backend_Golang/pkg"
	"log"
	"math/rand"
	"net/http"
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
)

type HandlerAuth struct {
	repositories.IAuthRepository
}

func InitializeAuthHandler(r repositories.IAuthRepository) *HandlerAuth {
	return &HandlerAuth{r}
}

func (h *HandlerAuth) Register(ctx *gin.Context) {
	body := &models.AuthUserModel{}
	if err := ctx.ShouldBind(body); err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, helpers.NewResponse("Error binding body register", nil, nil))
		return
	}
	if _, err := govalidator.ValidateStruct(body); err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusBadRequest, helpers.NewResponse("Input not valid", nil, nil))
		return
	}
	hs := pkg.HashConfig{
		Time:    3,
		Memory:  64 * 1024,
		Threads: 2,
		KeyLen:  32,
		SaltLen: 16,
	}
	hashedPassword, err := hs.GenHashedPassword(body.Password)
	if err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, helpers.NewResponse("Error in hashing password", nil, nil))
		return
	}
	otp := 100000 + rand.Intn(900000)
	err = h.RepositoryRegisterUser(body, hashedPassword, otp)
	if err != nil {
		if strings.Contains(err.Error(), "users_email_key") {
			ctx.JSON(http.StatusBadRequest, helpers.NewResponse("Email already registered", nil, nil))
			return
		}
		log.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, helpers.NewResponse("Internal Server Error", nil, nil))
		return
	}
	dataRegister := make(map[string]interface{})
	dataRegister["email"] = body.Email
	ctx.JSON(http.StatusCreated, helpers.NewResponse("Successfully register", dataRegister, nil))
}

func (h *HandlerAuth) ActivateAccount(ctx *gin.Context) {
	var query models.AuthUserActivateModel
	if err := ctx.ShouldBindQuery(&query); err != nil {
		ctx.JSON(http.StatusBadRequest, helpers.NewResponse("Error binding query user", nil, nil))
		log.Println(err.Error())
		return
	}
	if _, err := govalidator.ValidateStruct(&query); err != nil {
		ctx.JSON(http.StatusBadRequest, helpers.NewResponse("Input not valid", nil, nil))
		log.Println(err.Error())
		return
	}
	result, err := h.RepositorySelectPrivateData(query.Email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helpers.NewResponse("Internal Server Error in Private data", nil, nil))
		log.Println(err.Error())
		return
	}
	if len(result) == 0 {
		ctx.JSON(http.StatusNotFound, helpers.NewResponse("Account not found", nil, nil))
		return
	}
	if result[0].Otp != query.Otp {
		fmt.Println(result[0].Otp)
		fmt.Println(query.Otp)
		ctx.JSON(http.StatusBadRequest, helpers.NewResponse("Your OTP is wrong", nil, nil))
		return
	}
	errActivate := h.RepositoryActivateAccount(query.Email)
	if errActivate != nil {
		log.Println(errActivate.Error())
		ctx.JSON(http.StatusInternalServerError, helpers.NewResponse("Internal Server Error in activating account", nil, nil))
		return
	}
	ctx.JSON(http.StatusOK, helpers.NewResponse("Activate account success", nil, nil))
}

func (h *HandlerAuth) Login(ctx *gin.Context) {
	body := &models.AuthUserModel{}
	if err := ctx.ShouldBind(body); err != nil {
		ctx.JSON(http.StatusInternalServerError, helpers.NewResponse("Error in binding body login", nil, nil))
		log.Println(err.Error())
		return
	}
	if _, err := govalidator.ValidateStruct(body); err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusBadRequest, helpers.NewResponse("Input not valid", nil, nil))
		return
	}
	result, err := h.RepositorySelectPrivateData(body.Email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helpers.NewResponse("Internal Server Error in Private data", nil, nil))
		log.Println(err.Error())
		return
	}
	if len(result) == 0 {
		ctx.JSON(http.StatusNotFound, helpers.NewResponse("Account not found", nil, nil))
		return
	}
	if !result[0].Activate {
		ctx.JSON(http.StatusNotFound, helpers.NewResponse("Please activate your account first", nil, nil))
		return
	}
	hs := pkg.HashConfig{}
	isValid, err := hs.ComparePasswordAndHash(body.Password, result[0].Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helpers.NewResponse("Error during verification data", nil, nil))
		log.Println(err.Error())
		return
	}
	if !isValid {
		ctx.JSON(http.StatusUnauthorized, helpers.NewResponse("Email or password is wrong", nil, nil))
		return
	}
	payload := pkg.NewPayload(result[0].Id, result[0].Role)
	token, err := payload.GenerateToken()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helpers.NewResponse("Error Generating token", nil, nil))
		return
	}
	userInfo := make(map[string]interface{})
	userInfo["token"] = token
	userInfo["email"] = body.Email
	// userInfo["fullname"] = fmt.Sprintf("%p %p", result[0].First_Name, result[0].Last_Name)
	userInfo["role"] = result[0].Role
	userInfo["photo_profile"] = result[0].Photo_Profile
	ctx.JSON(http.StatusOK, helpers.NewResponse("Successfully login", userInfo, nil))
}

func (h *HandlerAuth) Logout(ctx *gin.Context) {
	bearerToken := ctx.GetHeader("Authorization")
	token := strings.Replace(bearerToken, "Bearer ", "", -1)
	if err := h.RepositoryLogOut(token); err != nil {
		ctx.JSON(http.StatusInternalServerError, helpers.NewResponse("Internal Server Error", nil, nil))
		return
	}
	ctx.JSON(http.StatusOK, helpers.NewResponse("Successfully logout. Thank you", nil, nil))
}
