package helper

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ErrorResponseBody struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
	Traceback  string `json:"traceback"`
}

type SuccessResponseBody struct {
	StatusCode int         `json:"status_code"`
	Data       interface{} `json:"data"`
}

func SuccessResponse(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusOK, SuccessResponseBody{
		StatusCode: http.StatusOK,
		Data:       data,
	})
}

func BadRequestResponse(ctx *gin.Context, err error) {
	ctx.JSON(http.StatusBadRequest, ErrorResponseBody{
		StatusCode: http.StatusBadRequest,
		Message:    err.Error(),
		Traceback:  fmt.Sprintf("%+v", err),
	})
}

func UnauthorizedResponse(ctx *gin.Context, err error) {
	ctx.JSON(http.StatusUnauthorized, ErrorResponseBody{
		StatusCode: http.StatusUnauthorized,
		Message:    err.Error(),
		Traceback:  fmt.Sprintf("%+v", err),
	})
}

func InternalServerErrorResponse(ctx *gin.Context, err error) {
	ctx.JSON(http.StatusInternalServerError, ErrorResponseBody{
		StatusCode: http.StatusInternalServerError,
		Message:    err.Error(),
		Traceback:  fmt.Sprintf("%+v", err),
	})
}
