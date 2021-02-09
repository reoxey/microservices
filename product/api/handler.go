package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"product/core"
	"product/logger"
	"product/repo/mysql"
)

type UserHandler interface {
	GetProduct(c *gin.Context)
	GetAllProducts(c *gin.Context)
	AddProduct(c *gin.Context)
	EditProduct(c *gin.Context)
	AuthorizeUser() gin.HandlerFunc
}

type handler struct {
	service core.ProductService
	log     *logger.Logger
}

const authVerify = "AUTHORIZE"

func (h handler) GetProduct(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.log.Println("WARNING:handler.GetProduct", "invalid id type")
		c.AbortWithStatusJSON(http.StatusBadRequest, nil)
		return
	}
	user, err := h.service.ProductById(c, id)
	if err != nil {
		h.log.Println("ERROR:handler.GetProduct", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, nil)
		return
	}
	c.JSON(http.StatusOK, user)
}

func (h handler) GetAllProducts(c *gin.Context) {

	val, ok := c.Get(authVerify)
	if !ok {
		h.log.Println("ERROR:handler.GetAllProducts", "interface assertion failed")
		c.AbortWithStatusJSON(http.StatusInternalServerError, nil)
		return
	}

	tokenMap := val.(map[string]interface{})

	if !tokenMap["is_admin"].(bool) {
		h.log.Println("WARNING:handler.GetAllProducts", "user denied "+tokenMap["email"].(string))
		c.AbortWithStatusJSON(http.StatusForbidden, nil)
		return
	}

	users, err := h.service.AllProducts(c)
	if err != nil {
		h.log.Println("ERROR:handler.GetAllProducts", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, nil)
		return
	}

	c.JSON(http.StatusOK, users)
}

func (h handler) AddProduct(c *gin.Context) {
	var p *core.Product

	if err := c.Bind(&p); err != nil {
		h.log.Println("ERROR:handler.AddProduct", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, nil)
		return
	}

	id, err := h.service.AddProduct(c, p)
	if err != nil {
		h.log.Println("ERROR:handler.AddProduct", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, nil)
		return
	}
	c.Set("Location",
		fmt.Sprintf("/api/products/%d", id),
	)
	c.JSON(http.StatusCreated, nil)
}

func (h handler) EditProduct(c *gin.Context) {
	var p *core.Product

	if err := c.Bind(&p); err != nil {
		h.log.Println("ERROR:handler.EditProduct", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, nil)
		return
	}

	err := h.service.EditProduct(c, p)

	if err != nil {
		if err == mysql.NoRowsAffected {
			h.log.Println("WARNING:handler.EditProduct", err)
			c.AbortWithStatusJSON(http.StatusOK, nil)
			return
		}
		h.log.Println("ERROR:handler.EditProduct", err)
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

func NewHandler(s core.ProductService, log *logger.Logger) UserHandler {
	return &handler{s, log}
}
