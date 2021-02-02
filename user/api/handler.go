package api

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"user/logger"
	"user/profile"
	"user/repo/mysql"
)

type UserHandler interface {
	GetUser(c *gin.Context)
	GetAllUsers(c *gin.Context)
	AddUser(c *gin.Context)
	EditUser(c *gin.Context)
	LoginUser(c *gin.Context)
	AuthorizeUser() gin.HandlerFunc
}

type handler struct {
	service profile.UserService
	log *logger.Logger
}

const authVerify = "AUTHORIZE"

func (h handler) LoginUser(c *gin.Context) {
	var l *profile.Login

	if err := c.Bind(&l); err != nil {
		h.log.Println("WARNING:handler.Login", "invalid login")
		c.AbortWithStatusJSON(http.StatusBadRequest, nil)
		return
	}
	token, err := h.service.Login(context.Background(), l)
	if err != nil {
		h.log.Println("WARNING:handler.Login", "invalid login")
		c.AbortWithStatusJSON(http.StatusInternalServerError, nil)
		return
	}
	if token == "" {
		c.JSON(http.StatusUnauthorized, nil)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}

func (h handler) GetUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.log.Println("WARNING:handler.GetUser", "invalid id type")
		c.AbortWithStatusJSON(http.StatusBadRequest, nil)
		return
	}
	user, err := h.service.UserById(context.Background(), id)
	if err != nil {
		h.log.Println("ERROR:handler.GetUser", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, nil)
		return
	}
	c.JSON(http.StatusOK, user)
}

func (h handler) GetAllUsers(c *gin.Context) {

	val, ok := c.Get(authVerify)
	if !ok {
		h.log.Println("ERROR:handler.GetAllUsers", "interface assertion failed")
		c.AbortWithStatusJSON(http.StatusInternalServerError, nil)
		return
	}

	tokenMap := val.(map[string]interface{})

	if !tokenMap["is_admin"].(bool) {
		h.log.Println("WARNING:handler.GetAllUsers", "user denied "+tokenMap["email"].(string))
		c.AbortWithStatusJSON(http.StatusForbidden, nil)
		return
	}

	users, err := h.service.AllUsers(context.Background())
	if err != nil {
		h.log.Println("ERROR:handler.GetAllUsers", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, nil)
		return
	}

	c.JSON(http.StatusOK, users)
}

func (h handler) AddUser(c *gin.Context) {
	var u *profile.User

	if err := c.Bind(&u); err != nil {
		h.log.Println("ERROR:handler.AddUser", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, nil)
		return
	}

	id, err := h.service.AddUser(context.Background(), u)
	if err != nil {
		h.log.Println("ERROR:handler.AddUser", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, nil)
		return
	}
	c.Set("Location",
		fmt.Sprintf("/api/users/%d", id),
	)
	c.JSON(http.StatusCreated, nil)
}

func (h handler) EditUser(c *gin.Context) {
	var u *profile.User

	if err := c.Bind(&u); err != nil {
		h.log.Println("ERROR:handler.EditUser", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, nil)
		return
	}

	if err := h.service.EditUser(context.Background(), u); err != nil {
		if err == mysql.NoRowsAffected {
			h.log.Println("WARNING:handler.EditUser", err)
			c.AbortWithStatusJSON(http.StatusNoContent, nil)
			return
		}
		h.log.Println("ERROR:handler.EditUser", err)
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

func NewHandler(s profile.UserService, log *logger.Logger) UserHandler {
	return &handler{s, log}
}
