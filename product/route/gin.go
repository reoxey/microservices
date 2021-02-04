package route

import (
	"github.com/gin-gonic/gin"

	"product/api"
	"product/core"
	"product/logger"
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

func (g *Gin) Handle(ps core.ProductService) {
	ph := api.NewHandler(ps, g.log)


	authRoute := g.Group("/api", ph.AuthorizeUser())

	// GET
	authRoute.GET("/products", ph.GetAllProducts)
	authRoute.GET("/products/:id", ph.GetProduct)

	//POST
	authRoute.POST("/products", ph.AddProduct)

	//PUT
	authRoute.PUT("/products", ph.EditProduct)
}
