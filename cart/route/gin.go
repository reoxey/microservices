package route

import (
	"github.com/gin-gonic/gin"

	"cart/api"
	"cart/core"
)

type Gin struct {
	*gin.Engine
}

func New(debug bool) *Gin {
	g := gin.New()
	if !debug {
		gin.SetMode(gin.ReleaseMode)
	}
	g.Use(gin.Logger())
	g.Use(gin.Recovery())
	return &Gin{
		g,
	}
}

func (g *Gin) Handle(cs core.CartService) {
	ch := api.NewHandler(cs)

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
