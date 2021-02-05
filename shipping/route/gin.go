package route

import (
	"github.com/gin-gonic/gin"

	"shipping/api"
	"shipping/logger"
	"shipping/core"
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

func (g *Gin) Handle(s core.ShippingService) {
	h := api.NewHandler(s, g.log)


	authRoute := g.Group("/api/shipping", h.AuthoriseUser())

	// GET
	authRoute.GET("/addresses", h.GetAllAddresses)
	authRoute.GET("/addresses/:id", h.GetAddressById)

	//POST
	authRoute.POST("/addresses", h.AddAddress)

	//PUT
	authRoute.PUT("/addresses", h.UpdateAddress)

	authRoute.PUT("/status/:id/:status", h.UpdateShippingStatus)
	authRoute.PUT("/payment/:id/:status", h.UpdatePaymentStatus)
}
