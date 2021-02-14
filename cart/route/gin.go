package route

import (
	"github.com/gin-gonic/gin"

	"cart/api"
	"cart/core"
	"cart/logger"
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

func (g *Gin) Handle(cs core.CartService) {
	ch := api.NewHandler(cs, g.log)


	authRoute := g.Group("/api", ch.AuthorizeUser())

	// GET
	authRoute.GET("/cart/:id", ch.GetCart)

	//POST
	authRoute.POST("/cart", ch.CreateCart)
	authRoute.POST("/cart/:id", ch.AddToCart)
	authRoute.POST("/cart/:id/checkout", ch.Checkout)

	//PUT
	authRoute.PUT("/cart/:id", ch.UpdateQty)

	//DELETE
	authRoute.DELETE("/cart/:id/:item_id", ch.RemoveItems)
}
