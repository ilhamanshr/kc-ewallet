package middleware

import (
	"fmt"
	"kc-ewallet/internals/errors"
	"kc-ewallet/protocols/http/response"
	"log"
	"net"
	"os"
	"strings"
	"syscall"
	"unicode"

	"github.com/davecgh/go-spew/spew"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

func isAllUpper(s string) bool {
	for _, v := range s {
		if !unicode.IsUpper(v) {
			return false
		}
	}
	return true
}

func lowerCaseFirst(s string) string {
	if len(s) < 2 {
		return strings.ToLower(s)
	}
	for i, v := range s {
		return string(unicode.ToLower(v)) + s[i+1:]
	}
	return ""
}

func fieldErrorToText(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", e.Field())
	case "required_without":
		fields := strings.Split(strings.ToLower(e.Param()), " ")
		return fmt.Sprintf("%s is required if %s is not present", e.Field(), strings.Join(fields, " or "))
	case "required_without_all":
		fields := strings.Split(strings.ToLower(e.Param()), " ")
		return fmt.Sprintf("%s is required if %s is not present", e.Field(), strings.Join(fields, " and "))
	case "required_with":
		fields := strings.Split(strings.ToLower(e.Param()), " ")
		return fmt.Sprintf("%s is required if %s is present", e.Field(), strings.Join(fields, " or "))
	case "required_with_all":
		fields := strings.Split(strings.ToLower(e.Param()), " ")
		return fmt.Sprintf("%s is required if %s is present", e.Field(), strings.Join(fields, " and "))
	case "max":
		return fmt.Sprintf("%s cannot be longer than or equal to %s", e.Field(), e.Param())
	case "min":
		return fmt.Sprintf("%s must be longer than or equal to %s", e.Field(), e.Param())
	case "email":
		return "Invalid email format"
	case "len":
		return fmt.Sprintf("%s must be %s characters long", e.Field(), e.Param())
	case "uuid", "uuid3", "uuid4", "uuid5":
		return fmt.Sprintf("%s must be in UUID format", e.Field())
	case "datetime":
		return fmt.Sprintf("%s must be in a valid datetime format (%s)", e.Field(), e.Param())
	case "startswith":
		return fmt.Sprintf("%s should starts with %s", e.Field(), e.Param())
	case "endswith":
		return fmt.Sprintf("%s should ends with %s", e.Field(), e.Param())
	}
	return fmt.Sprintf("%s is not valid", e.Field())
}

func handleValidationError(e *gin.Error, c *gin.Context) {
	validationErrs := e.Err.(validator.ValidationErrors)
	err := errors.Validation.New("Validation error")
	var errMessage string
	for _, validationErr := range validationErrs {
		errMessage = fieldErrorToText(validationErr)

		err = errors.AddFieldError(err, validationErr.Field(), errMessage)
	}
	err = errors.Msg(err, errMessage) // replace err message to the latest error message

	response.RespondError(c, err)
}

// HandleError Middleware for handling error
func HandleError() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// If there are errors and a response is already written before, do nothing
		if len(c.Errors) != 0 && c.Writer.Written() {
			return
		}

		// If there are errors and nothing is in the response, make a default error response
		for _, ginErr := range c.Errors {
			switch e := ginErr.Err.(type) {
			case validator.ValidationErrors:
				handleValidationError(ginErr, c)
			case *net.OpError:
				if se, ok := e.Err.(*os.SyscallError); ok {
					switch se.Err {
					case syscall.EPIPE:
						response.RespondError(c, errors.NewAndDontReport("Broken Pipe"))
						log.Printf("Error: Broken Pipe | %+v", ginErr)
					case syscall.ECONNRESET:
						response.RespondError(c, errors.NewAndDontReport("Connection Reset"))
						log.Printf("Error: Connection Reset | %+v", ginErr)
					}
				}
			default:
				err := errors.New("Unknown error. Error: %s", spew.Sdump(e))
				response.RespondError(c, errors.Msg(err, "Unknown error"))
				log.Printf("Error: Unknown error | %v", ginErr)
			}
		}
	}
}

func PanicRecoveryHandler() gin.HandlerFunc {
	// Add more panic recovery functionalities here
	return ginzap.RecoveryWithZap(zap.L(), false)
}
