package response

import (
	"fmt"
	"kc-ewallet/internals/errors"
	"kc-ewallet/internals/helpers/pagination"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type BaseResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type StandardResponse struct {
	BaseResponse
	Data interface{} `json:"data"`
}

type PaginatorResponse struct {
	StandardResponse
	Paginator *pagination.Paginator `json:"paginator"`
}

type ErrorDetail struct {
	IdMessage    string `json:"id_message"`
	EnMessage    string `json:"en_message"`
	RedirectCode int8   `json:"redirect_code"`
	Code         string `json:"code"`
}

type ErrorResponse struct {
	BaseResponse
	ErrorDetail ErrorDetail                     `json:"error_detail"`
	ErrorCode   string                          `json:"error_code"`
	Errors      []string                        `json:"errors"`
	Fields      map[string]errors.ErrorMessages `json:"fields"`
}

func (e ErrorResponse) WithDetail(errDetail ErrorDetail) ErrorResponse {
	e.ErrorDetail = errDetail
	return e
}

func BuildStandardResponse(status string, message string) BaseResponse {
	var response BaseResponse
	response.Status = status
	response.Message = message

	return response
}

func BuildSuccessResponse(message string, data interface{}) StandardResponse {
	var response StandardResponse
	response.BaseResponse = BuildStandardResponse("success", message)
	response.Data = data

	return response
}

func BuildErrorResponse(err error, errorCode string, message string) ErrorResponse {
	var response ErrorResponse
	var errMsgs []string

	response.BaseResponse = BuildStandardResponse("error", message)
	fields := errors.GetFields(err)
	for _, errorMessages := range fields {
		for _, errorMsg := range errorMessages {
			errMsgs = append(errMsgs, errorMsg)
		}
	}
	if response.Errors = errMsgs; errMsgs == nil {
		response.Errors = []string{message}
	}
	response.Fields = fields
	response.ErrorCode = errorCode

	return response
}

// RespondSuccess respond JSON with data
func RespondSuccess(c *gin.Context, data interface{}, message string) {
	c.JSON(http.StatusOK, BuildSuccessResponse(message, data))
}

// RespondSuccessWithPaginator respond with paginator
func RespondSuccessWithPaginator(c *gin.Context, data interface{}, paginator *pagination.Paginator, message string) {
	var response PaginatorResponse
	response.BaseResponse = BuildStandardResponse("success", message)
	response.Paginator = paginator
	response.Data = data

	c.JSON(http.StatusOK, response)
}

// RespondError respond error
func RespondError(c *gin.Context, err error) {
	var (
		defaultMessage, errorCode string
		errDetail                 ErrorDetail
	)

	errType := errors.GetType(err)
	switch errType {
	case errors.Unauthorized:
		defaultMessage = "Unauthorized."
	case errors.NotFound:
		defaultMessage = "Resource not found."
	case errors.UnprocessableEntity:
		defaultMessage = "Unprocessable entity error."
	case errors.BadRequest:
		defaultMessage = "Bad request."
	case errors.ExtendedError:
		extError, ok := err.(errors.ExtError)
		if ok {
			errDetail = ErrorDetail{
				IdMessage:    extError.GetIdMessage(),
				EnMessage:    extError.GetEnMessage(),
				Code:         extError.GetErrCode(),
				RedirectCode: extError.GetRedirectCode(),
			}
			errorCode = fmt.Sprintf("%d", extError.GetCode())
		}
	default:
		defaultMessage = "Internal server error."
	}

	message := err.Error()
	if message == "" {
		message = defaultMessage
	}
	if errorCode == "" {
		// Get error code from err type value
		errorCode = strconv.Itoa(int(errType))
	}
	statusCode := getStatusCode(errorCode)

	// Attach error to context for logging and send response
	c.Error(err)
	c.JSON(statusCode, BuildErrorResponse(err, errorCode, message).WithDetail(errDetail))
}

func getStatusCode(errorCode string) (statusCode int) {
	if len(errorCode) == 3 {
		statusCode, _ = strconv.Atoi(errorCode)
	} else if len(errorCode) > 3 {
		errorCodeStr := errorCode[:3]
		statusCode, _ = strconv.Atoi(errorCodeStr)
	}
	if statusCode == 0 {
		statusCode = 500
	}

	return statusCode
}
