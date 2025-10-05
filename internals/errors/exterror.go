package errors

type ExtError interface {
	Error() string

	GetIdMessage() string
	GetEnMessage() string
	GetCode() int
	GetType() ErrorType
	GetRedirectCode() int8
	SetRedirectCode(rc int8)
	GetErrCode() string
}

var _ = ExtError(&CustomError{})

type CustomError struct {
	code         int
	errCode      string
	redirectCode int8
	errType      ErrorType
	err          string
	idMessage    string
	enMessage    string
}

type ExtErrorArg struct {
	Code         int
	ErrCode      string
	RedirectCode int8
	Err          string
	IdMessage    string
	EnMessage    string
}

func (c *CustomError) WithOriginalError(err error) *CustomError {
	c.err = err.Error()
	return c
}

func (c *CustomError) Error() string {
	return c.err
}

func (c *CustomError) GetIdMessage() string {
	return c.idMessage
}

func (c *CustomError) GetEnMessage() string {
	return c.enMessage
}

func (c *CustomError) GetCode() int {
	return c.code
}

func (c *CustomError) GetType() ErrorType {
	return c.errType
}

func (c *CustomError) GetRedirectCode() int8 {
	return c.redirectCode
}

func (c *CustomError) SetRedirectCode(rc int8) {
	c.redirectCode = rc
}

func (c *CustomError) GetErrCode() string {
	return c.errCode
}

func NewExtError(arg ExtErrorArg) *CustomError {
	return &CustomError{
		code:         arg.Code,
		errCode:      arg.ErrCode,
		redirectCode: arg.RedirectCode,
		err:          arg.Err,
		errType:      ExtendedError,
		idMessage:    arg.IdMessage,
		enMessage:    arg.EnMessage,
	}
}
