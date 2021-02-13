package api

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"shipping/logger"
	"shipping/repo/mysql"
	"shipping/core"
)

type ShippingHandler interface {
	AddAddress(c *gin.Context)
	GetAddressById(c *gin.Context)
	GetAllAddresses(c *gin.Context)
	UpdateAddress(c *gin.Context)
	UpdateShippingStatus(c *gin.Context)
	UpdatePaymentStatus(c *gin.Context)
	AuthoriseUser() gin.HandlerFunc
}

type handler struct {
	service core.ShippingService
	log     *logger.Logger
}

func (h handler) AddAddress(c *gin.Context) {
	id, err := verifyAuth(c)
	if err != nil {
		h.log.Println("ERROR:handler.AddAddress", err)
		return
	}
	var addr *core.Address
	if err := c.Bind(&addr); err != nil {
		h.log.Println("ERROR:handler.AddAddress", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, nil)
		return
	}
	addr.User = id

	id, err = h.service.AddAddress(c, addr)
	if err != nil {
		h.log.Println("ERROR:handler.AddAddress", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, nil)
		return
	}
	c.Header("Location",
		fmt.Sprintf("/api/shipping/addresses/%d", id),
	)
	c.JSON(http.StatusCreated, nil)
}

func (h handler) GetAddressById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.log.Println("WARNING:handler.GetAddressById", "invalid id type")
		c.AbortWithStatusJSON(http.StatusBadRequest, nil)
		return
	}
	addr, err := h.service.AddressById(c, id)
	if err != nil {
		h.log.Println("ERROR:handler.GetAddressById", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, nil)
		return
	}
	c.JSON(http.StatusOK, addr)
}

func (h handler) GetAllAddresses(c *gin.Context) {

	userId, err := verifyAuth(c)
	if err != nil {
		h.log.Println("ERROR:handler.AddAddress", err)
		return
	}

	addrs, err := h.service.AllAddresses(c, userId)
	if err != nil {
		h.log.Println("ERROR:handler.GetAllAddresses", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, nil)
		return
	}

	c.JSON(http.StatusOK, addrs)
}

func (h handler) UpdateAddress(c *gin.Context) {
	var addr *core.Address
	if err := c.Bind(&addr); err != nil {
		h.log.Println("ERROR:handler.UpdateAddress", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, nil)
		return
	}

	if err := h.service.EditAddress(c, addr); err != nil {
		if err == mysql.NoRowsAffected {
			h.log.Println("WARNING:handler.UpdateAddress", err)
			c.AbortWithStatusJSON(http.StatusOK, nil)
			return
		}
		h.log.Println("ERROR:handler.UpdateAddress", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, nil)
		return
	}

	c.JSON(http.StatusOK, nil)
}

func (h handler) UpdateShippingStatus(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.log.Println("WARNING:handler.UpdateShippingStatus", "invalid id type")
		c.AbortWithStatusJSON(http.StatusBadRequest, nil)
		return
	}
	status, err := strconv.Atoi(c.Param("status"))
	if err != nil {
		h.log.Println("WARNING:handler.UpdateShippingStatus", "invalid status type")
		c.AbortWithStatusJSON(http.StatusBadRequest, nil)
		return
	}

	if err = h.service.UpdateStatus(c, id, core.ShippingStatus(status)); err != nil {
		h.log.Println("WARNING:handler.UpdateShippingStatus", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, nil)
		return
	}
	c.JSON(http.StatusOK, nil)
}

func (h handler) UpdatePaymentStatus(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.log.Println("WARNING:handler.UpdatePaymentStatus", "invalid id type")
		c.AbortWithStatusJSON(http.StatusBadRequest, nil)
		return
	}
	status, err := strconv.Atoi(c.Param("status"))
	if err != nil {
		h.log.Println("WARNING:handler.UpdatePaymentStatus", "invalid status type")
		c.AbortWithStatusJSON(http.StatusBadRequest, nil)
		return
	}

	if err = h.service.UpdatePayment(c, id, status); err != nil {
		h.log.Println("WARNING:handler.UpdatePaymentStatus", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, nil)
		return
	}
	c.JSON(http.StatusOK, nil)
}

func (h handler) AuthoriseUser() gin.HandlerFunc {
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

const authVerify = "AUTHORIZE"

func NewHandler(s core.ShippingService, log *logger.Logger) ShippingHandler {
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

	// if !tokenMap["is_admin"].(bool) {
	// 	c.AbortWithStatusJSON(http.StatusForbidden, nil)
	// 	return id, errors.New("user denied "+tokenMap["email"].(string))
	// }

	return
}
