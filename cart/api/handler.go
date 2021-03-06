package api

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"cart/core"
	"cart/logger"
	"cart/repo/mysql"
)

type CartHandler interface {
	CreateCart(c *gin.Context)
	GetCart(c *gin.Context)
	AddToCart(c *gin.Context)
	UpdateQty(c *gin.Context)
	RemoveItems(c *gin.Context)
	Checkout(c *gin.Context)
	AuthorizeUser() gin.HandlerFunc
}

type handler struct {
	service core.CartService
	log     *logger.Logger
}

const authVerify = "AUTHORIZE"

func (h handler) CreateCart(c *gin.Context) {
	buyer, err := verifyAuth(c)
	if err != nil {
		h.log.Println("ERROR:handler.CreateCart", err)
		return
	}

	id, err := h.service.New(c, buyer)
	if err != nil {
		h.log.Println("ERROR:handler.CreateCart", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, nil)
		return
	}
	c.Header("Location",
		fmt.Sprintf("/api/cart/%d", id),
	)
	c.JSON(http.StatusCreated, nil)
}

func (h handler) GetCart(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.log.Println("WARNING:handler.GetCart", "invalid id type")
		c.AbortWithStatusJSON(http.StatusBadRequest, nil)
		return
	}
	cart, err := h.service.Show(c, id)
	if err != nil {
		h.log.Println("ERROR:handler.GetCart", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, nil)
		return
	}
	c.JSON(http.StatusOK, cart)
}

func (h handler) AddToCart(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.log.Println("WARNING:handler.AddToCart", "invalid id type")
		c.AbortWithStatusJSON(http.StatusBadRequest, nil)
		return
	}

	var item *core.Item
	if err := c.Bind(&item); err != nil {
		h.log.Println("ERROR:handler.AddToCart", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, nil)
		return
	}
	if err = h.service.AddToCart(c, id, item); err != nil {
		h.log.Println("ERROR:handler.AddToCart", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, nil)
		return
	}
	c.JSON(http.StatusOK, nil)
}

func (h handler) UpdateQty(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.log.Println("WARNING:handler.UpdateQty", "invalid id type")
		c.AbortWithStatusJSON(http.StatusBadRequest, nil)
		return
	}

	var item *core.Item
	if err := c.Bind(&item); err != nil {
		h.log.Println("ERROR:handler.UpdateQty", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, nil)
		return
	}
	if err = h.service.UpdateQty(c, id, item); err != nil {
		if err == mysql.NoRowsAffected {
			h.log.Println("WARNING:handler.UpdateQty", err)
			c.AbortWithStatusJSON(http.StatusOK, nil)
			return
		}
		h.log.Println("ERROR:handler.UpdateQty", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, nil)
		return
	}
	c.JSON(http.StatusOK, nil)
}

func (h handler) RemoveItems(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.log.Println("WARNING:handler.RemoveItems", "invalid id type")
		c.AbortWithStatusJSON(http.StatusBadRequest, nil)
		return
	}
	itemId, err := strconv.Atoi(c.Param("item_id"))
	if err != nil {
		h.log.Println("WARNING:handler.RemoveItems", "invalid item_id type")
		c.AbortWithStatusJSON(http.StatusBadRequest, nil)
		return
	}
	if err = h.service.DeleteItems(c, id, itemId); err != nil {
		if err == mysql.NoRowsAffected {
			h.log.Println("WARNING:handler.RemoveItems", err)
			c.AbortWithStatusJSON(http.StatusOK, nil)
			return
		}
		h.log.Println("ERROR:handler.RemoveItems", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, nil)
		return
	}
	c.JSON(http.StatusOK, nil)
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

func (h handler) Checkout(c *gin.Context) {

	cartId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.log.Println("WARNING:handler.Checkout", "invalid id type")
		c.AbortWithStatusJSON(http.StatusBadRequest, nil)
		return
	}

	var checkout *core.Checkout
	if err = c.Bind(&checkout); err != nil {
		h.log.Println("ERROR:handler.Checkout", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, nil)
		return
	}

	if err = h.service.Checkout(c, checkout, cartId); err != nil {
		h.log.Println("ERROR:handler.Checkout", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, nil)
		return
	}
	c.JSON(http.StatusOK, nil)
}

func NewHandler(s core.CartService, log *logger.Logger) CartHandler {
	return &handler{s, log}
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
