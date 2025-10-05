package container

import (
	"context"
	"kc-ewallet/internals/errors"
	"kc-ewallet/internals/helpers/array"
	"kc-ewallet/internals/helpers/pagination"
	"mime/multipart"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type Auth struct{}
type ExtraData struct{}

type UrlParams struct {
	Str map[string]string
	Arr map[string][]string
}

type RequestInf interface {
	SetPostParams(params interface{}) error
	ValidateFileType(file *multipart.FileHeader, allowTypes []string) error
}

type Request struct {
	Auth
	ExtraData
	Ctx         context.Context
	GinCtx      *gin.Context
	Pagination  pagination.Config
	UrlParams   UrlParams
	URIParams   interface{}
	QueryParams interface{}
	PostParams  interface{}
}

func (r *Request) SetURIParams(params interface{}) error {
	if params == nil {
		return nil
	}
	if err := r.GinCtx.ShouldBindUri(params); err != nil {
		_ = r.GinCtx.Error(err)
		return err
	}

	r.URIParams = params
	return nil
}

func (r *Request) SetQueryParams(params interface{}) error {
	if params == nil {
		return nil
	}
	err := r.GinCtx.ShouldBindQuery(params)
	if err != nil {
		_ = r.GinCtx.Error(err)
		return err
	}

	r.QueryParams = params
	return nil
}

func (r *Request) SetPostParams(params interface{}) error {
	if params == nil {
		return nil
	}
	b := binding.Default(r.GinCtx.Request.Method, r.GinCtx.ContentType())
	var i interface{} = b
	var err error
	bBody, ok := i.(binding.BindingBody)
	if ok {
		// Use ShouldBindBodyWith so we can reuse request body after we read it (so we can have multiple binding)
		err = r.GinCtx.ShouldBindBodyWith(params, bBody)
	} else {
		err = r.GinCtx.ShouldBind(params)
	}

	if err != nil {
		_ = r.GinCtx.Error(err)
		return err
	}

	r.PostParams = params
	return nil
}

func (r *Request) ValidateFileType(file *multipart.FileHeader, allowTypes []string) error {
	if file == nil {
		return nil
	}

	var mimes []string
	for _, allowType := range allowTypes {
		mimes = append(mimes, extractMimeFromType(allowType)...)
	}
	fileType := file.Header.Get("Content-Type")
	if exists := array.InArray(fileType, mimes); exists {
		return nil
	}

	return errors.Validation.New("File is only allow %s type", strings.Join(allowTypes, " / "))
}

func InitCoreRequest(c *gin.Context) *Request {
	var r Request
	r.GinCtx = c
	r.Ctx = context.Background()
	// r.DBM = NewCoreDBManager(database.GetDB()) // @FIXME
	initUrlParams(c, &r)
	initPagination(c, &r)
	return &r
}

func initUrlParams(c *gin.Context, r *Request) {
	r.UrlParams.Str = make(map[string]string)
	r.UrlParams.Arr = make(map[string][]string)

	if c != nil {
		params := c.Request.URL.Query()
		for param, value := range params {
			if len(value) <= 1 {
				r.UrlParams.Str[param] = value[0]
			}
			r.UrlParams.Arr[param] = value
		}
	}
}

func initPagination(c *gin.Context, r *Request) {
	var page, limit int
	var err error
	if c != nil {
		page, err = strconv.Atoi(c.Query("page"))
		if err != nil || page == 0 {
			page = 1
		}

		limit, err = strconv.Atoi(c.Query("limit"))
		if err != nil {
			limit = 25
		}
	}

	r.Pagination = pagination.GetPaginationConfig(page, limit)
}

func extractMimeFromType(typeStr string) []string {
	mimeType := getMimeTypeMap()
	if mimes, ok := mimeType[typeStr]; ok {
		return mimes
	}
	return []string{}
}

func getMimeTypeMap() map[string][]string {
	return map[string][]string{
		"pdf":   {"application/pdf"},
		"image": {"image/jpg", "image/jpeg", "image/png", "image/gif", "image/svg+xml", "image/heic"},
		"video": {"video/x-flv", "video/mp4", "application/x-mpegURL", "video/MP2T", "video/3gpp", "video/quicktime", "video/x-msvideo", "video/x-ms-wmv"},
		"docs":  {"application/vnd.openxmlformats-officedocument.wordprocessingml.document", "application/msword"},
		"rar":   {"application/vnd.rar", "application/x-rar-compressed", "application/octet-stream"},
		"zip":   {"application/zip", "application/octet-stream", "application/x-zip-compressed", "multipart/x-zip"},
	}
}
