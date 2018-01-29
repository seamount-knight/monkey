package inter

import (
	"monkey/infra/log"
	"github.com/gin-gonic/gin"
	"net/http"
)

type MonkeyController interface {
	ListMonkeys()
}

type MonkeyHandler struct {
	log log.Logger
	controller MonkeyController
}


func NewMonkeyHandler(controller MonkeyController, log log.Logger) *MonkeyHandler {
	return &MonkeyHandler{
		log: log,
		controller: controller,
	}
}

func (h *MonkeyHandler) AddRoutes(router *gin.RouterGroup) {

	// List monkeys
	router.GET("", h.ListMonkeys)
	router.GET("/", h.ListMonkeys)
}

func (h *MonkeyHandler) ListMonkeys(ctx *gin.Context) {
	ctx.String(http.StatusOK, "monkeys")
}