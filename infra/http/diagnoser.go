package http

import (
	"github.com/gin-gonic/gin"
	"monkey/infra/diagnose"
	"net/http"
)

// DiagnoseRouter struct to add a route
type DiagnoseRouter struct {
	*diagnose.HealthChecker
}

func NewDiagnoser(checker *diagnose.HealthChecker) *DiagnoseRouter {
	return &DiagnoseRouter{
		HealthChecker: checker,
	}
}

func (h *DiagnoseRouter) AddRoutes(router *gin.RouterGroup) {
	router.Any("", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, h.Check())
	})
}
