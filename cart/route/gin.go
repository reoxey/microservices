package route

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"cart/api"
	"cart/core"
	log "cart/logger"
)

type routed struct {
	*gin.Engine
}

func New(cs core.CartService) http.Handler {
	g := gin.New()
	gin.SetMode(gin.ReleaseMode)
	g.Use(gin.Logger())
	g.Use(gin.Recovery())
	r := &routed{g}
	r.handle(cs)
	return r
}

func (r *routed) handle(cs core.CartService) {
	ch := api.NewHandler(cs)

	authRoute := r.Group("/api", ch.AuthorizeUser(), metrics())

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

func metrics() gin.HandlerFunc {
	return func(c *gin.Context) {

		start := time.Now()

		c.Next()

		msg := struct {
			Error string `json:"error"`
		}{}

		var status int

		errors := c.Errors.ByType(gin.ErrorTypeAny)
		if len(errors) > 0 {
			err := errors[0].Err
			httpDebug := c.GetHeader("X-Debug")
			switch et := err.(type) {
			case *log.GenericErr:
				status = et.Code
			case *log.ValidationErr:
				status = http.StatusBadRequest
			default:
				status = http.StatusInternalServerError
			}
			if httpDebug == "true" && err.Error() != "" {
				msg.Error = err.Error()
			}
			log.Error(err)
		}
		log.Msg("request time: " + time.Since(start).String())
		c.AbortWithStatusJSON(status, msg)
	}
}
