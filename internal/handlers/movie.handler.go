package handlers

import (
	"encoding/json"
	"fmt"
	"gilangrizaltin/Backend_Golang/internal/helpers"
	"gilangrizaltin/Backend_Golang/internal/models"
	"gilangrizaltin/Backend_Golang/internal/repositories"
	"log"
	"net/http"
	"strconv"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
)

type HandlerMovie struct {
	repositories.IMovieRepository
}

func InitializeMovieHandler(r repositories.IMovieRepository) *HandlerMovie {
	return &HandlerMovie{r}
}

func (h *HandlerMovie) GetAllMovie(ctx *gin.Context) {
	var queryMovie models.QueryParamGetMovie
	if err := ctx.ShouldBind(&queryMovie); err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusBadRequest, helpers.NewResponse("Error in binding body movie", nil, nil))
		return
	}
	if _, err := govalidator.ValidateStruct(&queryMovie); err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusBadRequest, helpers.NewResponse("Input not valid", nil, nil))
		return
	}
	result, err := h.RepositoryGetAllMovie(&queryMovie)
	if err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, helpers.NewResponse("Internal Server Error", nil, nil))
		return
	}
	if len(result) == 0 {
		ctx.JSON(http.StatusNotFound, helpers.NewResponse("Movie not found", nil, nil))
		return
	}
	data, errCount := h.RepositoryCountAllMovie(&queryMovie)
	if errCount != nil {
		log.Println(errCount.Error())
		ctx.JSON(http.StatusInternalServerError, helpers.NewResponse("Internal Server Count Movie Error", nil, nil))
		return
	}
	if len(data) < 1 {
		ctx.JSON(http.StatusNotFound, helpers.NewResponse("Data not found", nil, nil))
		return
	}
	meta := helpers.GetPagination(ctx, data, queryMovie.Page)
	ctx.JSON(http.StatusOK, helpers.NewResponse("Successfully Get User", result, &meta))
}

func (h *HandlerMovie) GetMovieSchedule(ctx *gin.Context) {
	ID := ctx.Param("movie_id")
	movieID, err := strconv.Atoi(ID)
	if err != nil {
		log.Print(err.Error())
		ctx.JSON(http.StatusBadRequest, helpers.NewResponse("Failed to get id movie", nil, nil))
		return
	}
	result, err := h.RepositoryGetSchedule(movieID)
	if err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, helpers.NewResponse("Internal Server Error", nil, nil))
		return
	}
	if len(result) == 0 {
		ctx.JSON(http.StatusNotFound, helpers.NewResponse("Data not found", nil, nil))
		return
	}
	ctx.JSON(http.StatusOK, helpers.NewResponse("Successfully Get User", result, nil))
}

func (h *HandlerMovie) GetMovieDetails(ctx *gin.Context) {
	ID := ctx.Param("movie_id")
	movieID, err := strconv.Atoi(ID)
	if err != nil {
		log.Print(err.Error())
		ctx.JSON(http.StatusBadRequest, helpers.NewResponse("Failed to get id movie", nil, nil))
		return
	}
	resultMovie, errMovie := h.RepositoryGetMovie(movieID)
	if errMovie != nil {
		log.Print(err.Error())
		ctx.JSON(http.StatusInternalServerError, helpers.NewResponse("Internal Movie Server Error", nil, nil))
		return
	}
	resultSchedule, errSchedule := h.RepositoryGetSchedule(movieID)
	if errSchedule != nil {
		log.Print(errSchedule.Error())
		ctx.JSON(http.StatusInternalServerError, helpers.NewResponse("Internal Movie Schedule Server Error", nil, nil))
		return
	}
	if len(resultMovie) < 1 {
		ctx.JSON(http.StatusNotFound, helpers.NewResponse("Movie not found", nil, nil))
		return
	}
	movieInfo := make(map[string]interface{})
	movieInfo["movie"] = resultMovie
	movieInfo["schedule"] = resultSchedule
	ctx.JSON(http.StatusOK, helpers.NewResponse("Successfully get details movie", movieInfo, nil))
}

func (h *HandlerMovie) GetCinema(ctx *gin.Context) {
	var querySchedule models.QuerySchedule
	if err := ctx.ShouldBind(&querySchedule); err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusBadRequest, helpers.NewResponse("Error in binding body movie", nil, nil))
		return
	}
	if _, err := govalidator.ValidateStruct(&querySchedule); err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusBadRequest, helpers.NewResponse("Input not valid", nil, nil))
		return
	}
	result, err := h.RepositoryGetCinema(&querySchedule)
	if err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, helpers.NewResponse("Internal Server Error", nil, nil))
		return
	}
	if len(result) == 0 {
		ctx.JSON(http.StatusNotFound, helpers.NewResponse("Movie not found", nil, nil))
		return
	}
	ctx.JSON(http.StatusOK, helpers.NewResponse("Successfully Get User", result, nil))
}

func (h *HandlerMovie) AddMovie(ctx *gin.Context) {
	var bodyMovie models.NewMovieModel
	if err := ctx.ShouldBind(&bodyMovie); err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusBadRequest, helpers.NewResponse("Error in binding body movie", nil, nil))
		return
	}
	var schedule []models.NewMovieSchedule
	json.Unmarshal([]byte(bodyMovie.Schedules), &schedule)
	// userInfo := make(map[string]interface{})
	// userInfo["movie"] = bodyMovie
	// userInfo["schedule"] = schedule
	if _, err := govalidator.ValidateStruct(&bodyMovie); err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusBadRequest, helpers.NewResponse("Input not valid", nil, nil))
		return
	}
	tx, errTx := h.Begin()
	if errTx != nil {
		log.Println(errTx.Error())
		ctx.JSON(http.StatusInternalServerError, helpers.NewResponse("Error begin tx", nil, nil))
		return
	}
	cld, err := helpers.InitCloudinary()
	if err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, helpers.NewResponse("Error in initialize uploading image system", nil, nil))
		return
	}
	formFile, _ := ctx.FormFile("movie_cover")
	var dataUrl string
	if formFile != nil {
		file, err := formFile.Open()
		if err != nil {
			log.Println(err.Error())
			ctx.JSON(http.StatusInternalServerError, helpers.NewResponse("Error in opening file", nil, nil))
			return
		}
		defer file.Close()
		publicId := fmt.Sprintf("%s_%s", bodyMovie.Director_Name, bodyMovie.Movie_Name)
		folder := ""
		res, err := cld.Uploader(ctx, file, publicId, folder)
		if err != nil {
			log.Println(err.Error())
			ctx.JSON(http.StatusInternalServerError, helpers.NewResponse("Error in upload image", nil, nil))
			return
		}
		dataUrl = res.SecureURL
	}
	result, err := h.RepositoryAddMovie(&bodyMovie, dataUrl, tx)
	if err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, helpers.NewResponse("Error in insert movie", nil, nil))
		defer tx.Rollback()
		return
	}
	if result == "" {
		ctx.JSON(http.StatusInternalServerError, helpers.NewResponse("No movie id", nil, nil))
		defer tx.Rollback()
		return
	}
	if errInserMovieSchedule := h.RepositoryAddMovieSchedule(schedule, tx, result); errInserMovieSchedule != nil {
		log.Println(errInserMovieSchedule.Error())
		ctx.JSON(http.StatusInternalServerError, helpers.NewResponse("Error in insert order product", nil, nil))
		defer tx.Rollback()
		return
	}
	if err := tx.Commit(); err != nil {
		ctx.JSON(http.StatusInternalServerError, helpers.NewResponse("Error in comitting order", nil, nil))
		return
	}
	ctx.JSON(http.StatusCreated, helpers.NewResponse("Successfully insert movie", nil, nil))
}

func (h *HandlerMovie) UpdateMovie(ctx *gin.Context) {
	var bodyUpdateMovie models.UpdateMovieModel
	if err := ctx.ShouldBind(&bodyUpdateMovie); err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusBadRequest, helpers.NewResponse("Error in binding body movie", nil, nil))
		return
	}
	ID := ctx.Param("movie_id")
	movieID, err := strconv.Atoi(ID)
	if err != nil {
		log.Print(err.Error())
		ctx.JSON(http.StatusBadRequest, helpers.NewResponse("Failed to get id movie", nil, nil))
		return
	}
	if _, err := govalidator.ValidateStruct(&bodyUpdateMovie); err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusBadRequest, helpers.NewResponse("Input not valid", nil, nil))
		return
	}
	cld, err := helpers.InitCloudinary()
	if err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, helpers.NewResponse("Error in initialize uploading image system", nil, nil))
		return
	}
	formFile, _ := ctx.FormFile("movie_cover")
	var dataUrl string
	if formFile != nil {
		file, err := formFile.Open()
		if err != nil {
			log.Println(err.Error())
			ctx.JSON(http.StatusInternalServerError, helpers.NewResponse("Error in opening file", nil, nil))
			return
		}
		defer file.Close()
		publicId := fmt.Sprintf("%s_%s", "movie_cover", ID)
		folder := ""
		res, err := cld.Uploader(ctx, file, publicId, folder)
		if err != nil {
			log.Println(err.Error())
			ctx.JSON(http.StatusInternalServerError, helpers.NewResponse("Error in upload image", nil, nil))
			return
		}
		dataUrl = res.SecureURL
	}
	result, errEditMovie := h.RepositoryEditMovie(&bodyUpdateMovie, movieID, dataUrl)
	if errEditMovie != nil {
		log.Println(errEditMovie.Error())
		ctx.JSON(http.StatusInternalServerError, helpers.NewResponse("Interal Server Error", nil, nil))
		return
	}
	var successEdit int64 = 1
	if result < successEdit {
		ctx.JSON(http.StatusNotFound, helpers.NewResponse("Data not found", nil, nil))
		return
	}
	ctx.JSON(http.StatusOK, helpers.NewResponse("Successfully update movie", nil, nil))
}

func (h *HandlerMovie) DeleteMovie(ctx *gin.Context) {
	ID, _ := strconv.Atoi(ctx.Param("movie_id"))
	result, err := h.RepositoryDeleteMovie(ID)
	if err != nil {
		log.Print(err.Error())
		ctx.JSON(http.StatusInternalServerError, helpers.NewResponse("Internal Server Error", nil, nil))
		return
	}
	var dataExecuted int64 = 1
	if result < dataExecuted {
		ctx.JSON(http.StatusNotFound, helpers.NewResponse("Data not found", nil, nil))
		return
	}
	ctx.JSON(http.StatusOK, helpers.NewResponse("Successfully delete data", nil, nil))
}
