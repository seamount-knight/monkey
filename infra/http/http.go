package http

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"monkey/infra/log"
	"net/http"
	"time"
)

type Router interface {
	AddRoutes(router *gin.RouterGroup)
}

type Config struct {
	Host           string
	Port           string
	AddHealthCheck bool
	Component      string
}

type Server struct {
	config   Config
	start    time.Time
	log      log.Logger
	app      *gin.Engine
	versions map[int]*gin.RouterGroup
}

func NewServer(config Config, log log.Logger) *Server {
	return &Server{
		config:   config,
		log:      log,
		app:      gin.New(),
		versions: map[int]*gin.RouterGroup{},
	}
}

func (s *Server) Init() *Server {
	gin.SetMode(viper.GetString("GIN_MODE"))
	if s.config.AddHealthCheck {
		s.app.Any("/", s.HealthCheck)
		s.app.Any("/_ping", s.HealthCheck)
	}
	return s
}

func (s *Server) AddVersion(version int) *Server {
	if _, ok := s.versions[version]; !ok {
		s.versions[version] = s.app.Group(fmt.Sprintf("/v%d", version))
	}
	return s
}

func (s *Server) AddEndpoint(relativePath string, handler Router) *Server {
	router := s.app.Group(relativePath)
	handler.AddRoutes(router)
	return s
}

func (s *Server) AddVersionEndpoint(version int, relativePath string, handler Router) *Server {
	s.AddVersion(version)
	handler.AddRoutes(s.versions[version].Group(relativePath))

	return s
}

func (s *Server) Start() {
	s.start = time.Now()
	s.app.Run(fmt.Sprintf("%s:%s", s.config.Host, s.config.Port))
}

func (s *Server) HealthCheck(ctx *gin.Context) {
	ctx.String(http.StatusOK, "%s OK", s.config.Component)
}
