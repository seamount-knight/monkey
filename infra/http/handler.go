package http

import (
	"github.com/gin-gonic/gin"
	"monkey/infra/errors"
	"monkey/infra/log"
	"net/http"
)

type Handler struct {
}

func (Handler) HandlerError(err error, ctx *gin.Context, log log.Logger) {
	status := getErrorStatusCode(err)
	log.Debugf("Error: %v - returning status: %d", err, status)
	ctx.JSON(status, NewHTTPError(err))
}

func (Handler) HandleErrors(errs []error, ctx *gin.Context, log log.Logger) bool {
	status := getErrorsStatusCode(errs)
	if status == 0 {
		return false
	}
	log.Debugf("Error: %v - returning status: %d", errs, status)
	ctx.JSON(
		status,
		NewHTTPError(errs...),
	)
	return true
}

func getErrorsStatusCode(errs []error) int {
	if errs == nil || len(errs) == 0 {
		return http.StatusInternalServerError
	}
	for _, e := range errs {
		if e != nil {
			return getErrorStatusCode(e)
		}
	}
	return http.StatusInternalServerError
}

func getErrorStatusCode(err error) int {
	if e, ok := err.(*errors.Error); ok {
		if e.StatusCode != 0 {
			return e.StatusCode
		}
	}
	return http.StatusInternalServerError
}

type HTTPError struct {
	Errors []error `json:"errors"`
}

func NewHTTPError(err ...error) *HTTPError {
	return &HTTPError{
		Errors: err,
	}
}
