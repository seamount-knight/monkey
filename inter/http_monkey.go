package inter

import (
	"github.com/gin-gonic/gin"
	"monkey/domain"
	mhttp "monkey/infra/http"
	"monkey/infra/log"
	"net/http"
)

type MonkeyController interface {
	ListMonkeys() ([]*domain.Monkey, error)
	RetrieveMonkey(uuid string) (*domain.Monkey, error)
}

type MonkeyHandler struct {
	mhttp.Handler
	log        log.Logger
	controller MonkeyController
}

func NewMonkeyHandler(controller MonkeyController, log log.Logger) *MonkeyHandler {
	return &MonkeyHandler{
		log:        log,
		controller: controller,
	}
}

func (h *MonkeyHandler) AddRoutes(router *gin.RouterGroup) {

	// List monkeys
	router.GET("", h.ListMonkeys)
	router.GET("/", h.ListMonkeys)

	// Retrieve Monkey
	router.GET("/:uuid", h.RetrieveMonkey)
	router.GET("/:uuid/", h.RetrieveMonkey)
}

func (h *MonkeyHandler) ListMonkeys(ctx *gin.Context) {
	monkeys, err := h.controller.ListMonkeys()
	if err != nil {
		h.HandlerError(err, ctx, h.log)
		return
	}

	ctx.JSON(http.StatusOK, monkeys)
}

func (h *MonkeyHandler) RetrieveMonkey(ctx *gin.Context) {
	uuid := ctx.Param("uuid")
	monkey, err := h.controller.RetrieveMonkey(uuid)
	if err != nil {
		h.HandlerError(err, ctx, h.log)
		return
	}
	ctx.JSON(http.StatusOK, monkey)
}
