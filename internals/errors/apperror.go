package errors

import (
	"fmt"
	log_color "kc-ewallet/internals/helpers/color"
	"net/http"
	"runtime"
	"strconv"

	"github.com/getsentry/sentry-go"
	"github.com/pkg/errors"
)

// ErrorType is the type of an error
type ErrorType uint

// NoType error
const (
	NoType              ErrorType = http.StatusInternalServerError
	InternalServer      ErrorType = http.StatusInternalServerError
	ServiceUnavailable  ErrorType = http.StatusServiceUnavailable
	BadRequest          ErrorType = http.StatusBadRequest
	Unauthorized        ErrorType = http.StatusUnauthorized
	Forbidden           ErrorType = http.StatusForbidden
	NotFound            ErrorType = http.StatusNotFound
	Validation          ErrorType = http.StatusUnprocessableEntity
	UnprocessableEntity ErrorType = http.StatusUnprocessableEntity
	TooManyRequests     ErrorType = http.StatusTooManyRequests

	// from customer service
	DefaultAppError ErrorType = 42201

	// need to find the convention for this error outside http status error
	ExtendedError ErrorType = 50001
)

type ErrorMessages []string

type AppError struct {
	errorType     ErrorType
	originalError error
	fields        map[string]ErrorMessages
	stackTrace    []string
}

// New creates a new AppError with formatted message
func (errorType ErrorType) New(msg string, args ...interface{}) error {
	log_color.PrintRed(fmt.Sprintf(msg, args...))
	shouldReport := shouldReport(errorType)
	appError := AppError{errorType: errorType, originalError: fmt.Errorf(msg, args...), stackTrace: []string{msg}}
	if shouldReport {
		appError.Report()
	}

	return appError
}

// NewWithUserMsg creates a new error with custom formatted user message
func (errorType ErrorType) NewWithUserMsg(err error, userMsg string, args ...interface{}) error {
	appError := errorType.New("Error: %v", err)

	return Msg(appError, userMsg, args...)
}

// NewAndReport creates a new AppError with formatted message and report
func (errorType ErrorType) NewAndReport(msg string, args ...interface{}) error {
	log_color.PrintRed(fmt.Sprintf(msg, args...))
	appError := AppError{errorType: errorType, originalError: fmt.Errorf(msg, args...), stackTrace: []string{msg}}
	appError.Report()

	return appError
}

func shouldReport(errorType ErrorType) bool {
	statusCodeStr := strconv.Itoa(int(errorType))[:3] // Get first 3 digits
	statusCode, err := strconv.Atoi(statusCodeStr)
	if err != nil {
		return true
	}

	switch ErrorType(statusCode) {
	case Unauthorized, Forbidden, NotFound, UnprocessableEntity, TooManyRequests:
		return false
	default:
		return true
	}
}

// Error returns the mssage of a AppError
func (error AppError) Error() string {
	return error.originalError.Error()
}

func (error AppError) Report() {
	log_color.PrintRed("========== Error Stack Strace ==========")
	var stackTrace []string
	for i := 4; i < 9; i++ { // Skip 4 function, Get last 5 error trace
		file, line, fnName := TraceCaller(i)
		traceMsg := fmt.Sprintf("%s:%d@%s", file, line, fnName)
		log_color.PrintYellow(traceMsg)
		stackTrace = append(stackTrace, traceMsg)
	}

	sentry.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetExtra("stack_trace", stackTrace)
	})
	sentry.CaptureException(error)
}

// New creates a no type error with formatted message and report to sentry
func New(msg string, args ...interface{}) error {
	log_color.PrintRed(fmt.Sprintf(msg, args...))
	err := AppError{errorType: NoType, originalError: errors.New(fmt.Sprintf(msg, args...)), stackTrace: []string{msg}}

	err.Report()

	return err
}

// NewAndDontReport creates a new AppError and don't report it
func NewAndDontReport(msg string, args ...interface{}) error {
	log_color.PrintRed(fmt.Sprintf(msg, args...))
	err := AppError{errorType: NoType, originalError: errors.New(fmt.Sprintf(msg, args...)), stackTrace: []string{msg}}

	return err
}

func Msg(err error, msg string, args ...interface{}) error {
	fileName, line, fnName := TraceCaller(3)
	errorMsg := fmt.Sprintf("%s:%d@%s()", fileName, line, fnName)
	formattedMsg := fmt.Sprintf(msg, args...)
	if appError, ok := err.(AppError); ok {
		return AppError{
			errorType:     appError.errorType,
			originalError: errors.New(formattedMsg),
			fields:        appError.fields,
			stackTrace:    append([]string{errorMsg}, appError.stackTrace...),
		}
	}

	return AppError{errorType: NoType, originalError: errors.New(formattedMsg), stackTrace: append([]string{errorMsg}, getOriginalErrorStackTrace(err)...)}
}

// Wrap an error with formatted message
func Wrap(err error, msg string, args ...interface{}) error {
	wrappedError := errors.Wrapf(err, msg, args...)
	if appError, ok := err.(AppError); ok {
		return AppError{
			errorType:     appError.errorType,
			originalError: wrappedError,
			fields:        appError.fields,
			stackTrace:    appError.stackTrace,
		}
	}

	return AppError{errorType: NoType, originalError: wrappedError, stackTrace: append([]string{fmt.Sprintf(msg, args...)}, getOriginalErrorStackTrace(err)...)}
}

// Cause gives the original error
func Cause(err error) error {
	return errors.Cause(err)
}

// AddStackTrace an error with format string
func AddStackTrace(err error, msg string) error {
	if appError, ok := err.(AppError); ok {
		stackTrace := append([]string{msg}, appError.stackTrace...)
		return AppError{errorType: appError.errorType, originalError: appError.originalError, fields: appError.fields, stackTrace: stackTrace}
	}

	stackTrace := append([]string{msg}, getOriginalErrorStackTrace(err)...)
	return AppError{errorType: NoType, originalError: err, stackTrace: stackTrace}
}

// GetStackTrace returns the error stack trace
func GetStackTrace(err error) []string {
	if appError, ok := err.(AppError); ok {
		return appError.stackTrace
	}
	return []string{}
}

// AddFieldError Append an error message to a field
func AddFieldError(err error, fieldName string, errorMessage string) error {
	appError, isAppError := err.(AppError)
	if !isAppError {
		fields := map[string]ErrorMessages{
			fieldName: {errorMessage},
		}
		return AppError{errorType: NoType, originalError: err, fields: fields, stackTrace: getOriginalErrorStackTrace(err)}
	}

	if appError.fields == nil {
		appError.fields = make(map[string]ErrorMessages)
	}
	appError.fields[fieldName] = append(appError.fields[fieldName], errorMessage)

	return appError
}

// SetFieldErrors adds an error messages to a field
func SetFieldErrors(err error, fieldName string, errorMessages []string) error {
	if appError, ok := err.(AppError); ok {
		appError.fields[fieldName] = errorMessages
		return appError
	}

	fields := map[string]ErrorMessages{
		fieldName: ErrorMessages(errorMessages),
	}
	return AppError{errorType: NoType, originalError: err, fields: fields, stackTrace: getOriginalErrorStackTrace(err)}
}

// GetFields returns the error fields
func GetFields(err error) map[string]ErrorMessages {
	if appError, ok := err.(AppError); ok {
		return appError.fields
	}
	return make(map[string]ErrorMessages)
}

// GetType returns the error type
func GetType(err error) ErrorType {
	switch v := err.(type) {
	case AppError:
		return v.errorType
	case ExtError:
		return v.GetType()
	default:
		return NoType
	}
}

// Is Check if error is the specified error type
func Is(err error, errorType ErrorType) bool {
	if err == nil {
		return false
	}

	errType := NoType
	if appError, ok := err.(AppError); ok {
		errType = appError.errorType
	}

	return errType == errorType
}

// IsNotFound Check if error is NotFound error type
func IsNotFound(err error) bool {
	return Is(err, NotFound)
}

func SetExtra(key string, value interface{}) {
	// @TODO: Need to export to an interface!
	sentry.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetExtra(key, value)
	})
}

func TraceCaller(skip int) (file string, line int, fnName string) {
	pc := make([]uintptr, 10) // at least 1 entry needed
	runtime.Callers(skip, pc)

	if pc[0] == uintptr(0) {
		return
	}

	f := runtime.FuncForPC(pc[0])
	file, line = f.FileLine(pc[0])
	fnName = f.Name()

	return
}
