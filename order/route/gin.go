package route

import (
	"github.com/gin-gonic/gin"

	"order/api"
	"order/core"
	"order/logger"
)

type Gin struct {
	*gin.Engine
	log *logger.Logger
}

func New(l *logger.Logger, debug bool) *Gin {
	g := gin.New()
	if !debug {
		gin.SetMode(gin.ReleaseMode)
	}
	g.Use(gin.Logger())
	g.Use(gin.Recovery())
	return &Gin{
		g,
		l,
	}
}

func (g *Gin) Handle(cs core.OrderService) {
	ch := api.NewHandler(cs, g.log)


	authRoute := g.Group("/api", ch.AuthorizeUser())

	// GET
	authRoute.GET("/orders/:id", ch.GetOrder)
	authRoute.GET("/orders", ch.AllOrders)
}
