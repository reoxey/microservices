package api

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"order/core"
	"order/logger"
)

type OrderHandler interface {
	GetOrder(c *gin.Context)
	AllOrders(c *gin.Context)
	AuthorizeUser() gin.HandlerFunc
}

type handler struct {
	service  core.OrderService
	log     *logger.Logger
}

const authVerify = "AUTHORIZE"

func NewHandler(s core.OrderService, log *logger.Logger) OrderHandler {
	return &handler{s, log}
}

func (h handler) GetOrder(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.log.Println("WARNING:handler.GetOrder", "invalid id type")
		c.AbortWithStatusJSON(http.StatusBadRequest, nil)
		return
	}
	order, err := h.service.GetOrder(c, id)
	if err != nil {
		h.log.Println("ERROR:handler.GetOrder", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, nil)
		return
	}
	c.JSON(http.StatusOK, order)
}

func (h handler) AllOrders(c *gin.Context) {

	buyer, err := verifyAuth(c)
	if err != nil {
		h.log.Println("ERROR:handler.CreateCart", err)
		return
	}

	orders, err := h.service.AllOrders(c, buyer)
	if err != nil {
		h.log.Println("ERROR:handler.GetOrder", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, nil)
		return
	}
	c.JSON(http.StatusOK, orders)
}

func (h handler) AuthorizeUser() gin.HandlerFunc {
	return func(c *gin.Context) {

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			h.log.Println("WARNING:handler.AuthorizeUser", "Not Authorize")
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		tokenString := authHeader[len("Bearer "):]

		auth, err := h.service.Authorize(tokenString)
		if err != nil {
			h.log.Println("ERROR:handler.AuthorizeUser", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, nil)
			return
		}

		if auth == nil {
			h.log.Println("WARNING:handler.AuthorizeUser", "Not Authorize")
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Set(authVerify, auth)
	}
}

func verifyAuth(c *gin.Context) (id int, err error) {

	val, ok := c.Get(authVerify)
	if !ok {
		c.AbortWithStatusJSON(http.StatusInternalServerError, nil)
		return id, errors.New("interface assertion failed")
	}

	tokenMap := val.(map[string]interface{})
	id = int(tokenMap["id"].(float64))

	if !tokenMap["is_admin"].(bool) {
		c.AbortWithStatusJSON(http.StatusForbidden, nil)
		return id, errors.New("user denied "+tokenMap["email"].(string))
	}

	return
}
