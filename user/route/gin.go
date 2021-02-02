package route

import (
	"github.com/gin-gonic/gin"

	"user/api"
	"user/logger"
	"user/profile"
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

func (g *Gin) Handle(us profile.UserService) {
	uh := api.NewHandler(us, g.log)

	loginRoute := g.Group("/api")
	loginRoute.POST("/login", uh.LoginUser)


	authRoute := g.Group("/api", uh.AuthorizeUser())

	// GET
	authRoute.GET("/users", uh.GetAllUsers)
	authRoute.GET("/users/:id", uh.GetUser)

	//POST
	authRoute.POST("/users", uh.AddUser)

	//PUT
	authRoute.PUT("/users", uh.EditUser)
}
