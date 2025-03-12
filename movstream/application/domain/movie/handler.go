package movie

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"log"
	"net/http"
	"simple-go/pkg/response"
	"strconv"
)

type Service interface {
	GetMovies(ctx context.Context, customerId int, filter FilterMovie) ([]MoviesResponse, PaginateListing, response.ErrorResponse)

	WatchMovie(ctx context.Context, customerId, movieId int) response.ErrorResponse
	VoteMovie(ctx context.Context, customerId, movieId int) response.ErrorResponse
	UnVoteMovie(ctx context.Context, customerId, movieId int) response.ErrorResponse

	AdminCreateMovie(ctx context.Context, customerId int, req CreateMovieRequest) response.ErrorResponse
	AdminUpdateMovie(ctx context.Context, customerId, movieId int, req UpdateMovieRequest) response.ErrorResponse
	AdminDeleteMovie(ctx context.Context, customerId, movieId int) response.ErrorResponse
}

type handler struct {
	service Service
}

func NewHandler(svc Service) handler {
	return handler{
		service: svc,
	}
}

func (h handler) GetAllMovies(ctx *gin.Context) {
	customerId := ctx.GetInt("customerId")

	var req FilterMovie
	req.Search = ctx.Query("search")
	req.Page = ctx.Query("page")
	req.Limit = ctx.Query("limit")

	listing, paginate, err := h.service.GetMovies(ctx, customerId, req)
	if !err.IsNoError {
		resp := response.Error(err.Code).WithError(err.Message).WithStatusCode(err.StatusCode)
		ctx.AbortWithStatusJSON(resp.StatusCode, resp)
		return
	}

	resp := response.Success("22152").
		WithData(listing).
		WithTotalPage(int(paginate.TotalPage)).
		WithCount(int(paginate.TotalData))

	ctx.JSON(resp.StatusCode, resp)
}

func (h handler) MovieWatch(ctx *gin.Context) {
	customerId := ctx.GetInt("customerId")
	customerType := ctx.GetString("customerType")
	if customerType == "admin" {
		resp := response.Error("22101").WithStatusCode(http.StatusUnauthorized).WithError("unauthorized")
		ctx.AbortWithStatusJSON(resp.StatusCode, resp)
		return
	}

	ids := ctx.Param("movieId")
	if ids == "" {
		resp := response.Error("22102").WithStatusCode(http.StatusBadRequest)
		ctx.AbortWithStatusJSON(resp.StatusCode, resp)
		return
	}
	movieId, err := strconv.Atoi(ids)
	if err != nil {
		resp := response.Error("22102").WithStatusCode(http.StatusInternalServerError).WithError(err.Error())
		ctx.AbortWithStatusJSON(resp.StatusCode, resp)
		return
	}

	errResp := h.service.WatchMovie(ctx, customerId, movieId)
	if !errResp.IsNoError {
		resp := response.Error(errResp.Code).WithError(errResp.Message).WithStatusCode(errResp.StatusCode)
		ctx.AbortWithStatusJSON(resp.StatusCode, resp)
		return
	}

	resp := response.Success("22154")
	ctx.JSON(resp.StatusCode, resp)
}

func (h handler) MovieVote(ctx *gin.Context) {
	customerId := ctx.GetInt("customerId")
	customerType := ctx.GetString("customerType")
	if customerType == "admin" {
		resp := response.Error("22101").WithStatusCode(http.StatusUnauthorized).WithError("unauthorized")
		ctx.AbortWithStatusJSON(resp.StatusCode, resp)
		return
	}

	ids := ctx.Param("movieId")
	if ids == "" {
		resp := response.Error("22102").WithStatusCode(http.StatusBadRequest)
		ctx.AbortWithStatusJSON(resp.StatusCode, resp)
		return
	}
	movieId, err := strconv.Atoi(ids)
	if err != nil {
		resp := response.Error("22102").WithStatusCode(http.StatusInternalServerError).WithError(err.Error())
		ctx.AbortWithStatusJSON(resp.StatusCode, resp)
		return
	}

	errResp := h.service.VoteMovie(ctx, customerId, movieId)
	if !errResp.IsNoError {
		resp := response.Error(errResp.Code).WithError(errResp.Message).WithStatusCode(errResp.StatusCode)
		ctx.AbortWithStatusJSON(resp.StatusCode, resp)
		return
	}

	resp := response.Success("22154")
	ctx.JSON(resp.StatusCode, resp)
}

func (h handler) MovieUnVote(ctx *gin.Context) {
	customerId := ctx.GetInt("customerId")
	customerType := ctx.GetString("customerType")
	if customerType == "admin" {
		resp := response.Error("22101").WithStatusCode(http.StatusUnauthorized).WithError("unauthorized")
		ctx.AbortWithStatusJSON(resp.StatusCode, resp)
		return
	}

	ids := ctx.Param("movieId")
	if ids == "" {
		resp := response.Error("22102").WithStatusCode(http.StatusBadRequest)
		ctx.AbortWithStatusJSON(resp.StatusCode, resp)
		return
	}
	movieId, err := strconv.Atoi(ids)
	if err != nil {
		resp := response.Error("22102").WithStatusCode(http.StatusInternalServerError).WithError(err.Error())
		ctx.AbortWithStatusJSON(resp.StatusCode, resp)
		return
	}

	errResp := h.service.UnVoteMovie(ctx, customerId, movieId)
	if !errResp.IsNoError {
		resp := response.Error(errResp.Code).WithError(errResp.Message).WithStatusCode(errResp.StatusCode)
		ctx.AbortWithStatusJSON(resp.StatusCode, resp)
		return
	}

	resp := response.Success("22154")
	ctx.JSON(resp.StatusCode, resp)
}

func (h handler) AdminCreateMovie(ctx *gin.Context) {
	customerId := ctx.GetInt("customerId")
	customerType := ctx.GetString("customerType")
	if customerType != "admin" {
		resp := response.Error("22101").WithStatusCode(http.StatusUnauthorized).WithError("unauthorized")
		ctx.AbortWithStatusJSON(resp.StatusCode, resp)
		return
	}

	var req CreateMovieRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Println(err.Error())
		resp := response.Error("22102").WithStatusCode(http.StatusBadRequest)

		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			resp.WithArgsMessage(ve[0].Field(), ve[0].Tag())
		}

		ctx.AbortWithStatusJSON(resp.StatusCode, resp)
		return
	}

	err := h.service.AdminCreateMovie(ctx, customerId, req)
	if !err.IsNoError {
		resp := response.Error(err.Code).WithError(err.Message).WithStatusCode(err.StatusCode)
		ctx.AbortWithStatusJSON(resp.StatusCode, resp)
		return
	}

	resp := response.Success("22151")
	ctx.JSON(resp.StatusCode, resp)
}

func (h handler) AdminUpdateMovie(ctx *gin.Context) {
	customerId := ctx.GetInt("customerId")
	customerType := ctx.GetString("customerType")
	if customerType != "admin" {
		resp := response.Error("22101").WithStatusCode(http.StatusUnauthorized).WithError("unauthorized")
		ctx.AbortWithStatusJSON(resp.StatusCode, resp)
		return
	}

	ids := ctx.Param("movieId")
	if ids == "" {
		resp := response.Error("22102").WithStatusCode(http.StatusBadRequest)
		ctx.AbortWithStatusJSON(resp.StatusCode, resp)
		return
	}
	movieId, err := strconv.Atoi(ids)
	if err != nil {
		resp := response.Error("22102").WithStatusCode(http.StatusInternalServerError).WithError(err.Error())
		ctx.AbortWithStatusJSON(resp.StatusCode, resp)
		return
	}

	var req UpdateMovieRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Println(err.Error())
		resp := response.Error("22102").WithStatusCode(http.StatusBadRequest)

		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			resp.WithArgsMessage(ve[0].Field(), ve[0].Tag())
		}

		ctx.AbortWithStatusJSON(resp.StatusCode, resp)
		return
	}

	errResp := h.service.AdminUpdateMovie(ctx, customerId, movieId, req)
	if !errResp.IsNoError {
		resp := response.Error(errResp.Code).WithError(errResp.Message).WithStatusCode(errResp.StatusCode)
		ctx.AbortWithStatusJSON(resp.StatusCode, resp)
		return
	}

	resp := response.Success("22154")
	ctx.JSON(resp.StatusCode, resp)
}

func (h handler) AdminDeleteMovie(ctx *gin.Context) {
	customerId := ctx.GetInt("customerId")
	customerType := ctx.GetString("customerType")
	if customerType != "admin" {
		resp := response.Error("22101").WithStatusCode(http.StatusUnauthorized).WithError("unauthorized")
		ctx.AbortWithStatusJSON(resp.StatusCode, resp)
		return
	}

	ids := ctx.Param("movieId")
	if ids == "" {
		resp := response.Error("22102").WithStatusCode(http.StatusBadRequest)
		ctx.AbortWithStatusJSON(resp.StatusCode, resp)
		return
	}
	movieId, err := strconv.Atoi(ids)
	if err != nil {
		resp := response.Error("22102").WithStatusCode(http.StatusInternalServerError).WithError(err.Error())
		ctx.AbortWithStatusJSON(resp.StatusCode, resp)
		return
	}

	errResp := h.service.AdminDeleteMovie(ctx, customerId, movieId)
	if !errResp.IsNoError {
		resp := response.Error(errResp.Code).WithError(errResp.Message).WithStatusCode(errResp.StatusCode)
		ctx.AbortWithStatusJSON(resp.StatusCode, resp)
		return
	}

	resp := response.Success("22153")
	ctx.JSON(resp.StatusCode, resp)
}
