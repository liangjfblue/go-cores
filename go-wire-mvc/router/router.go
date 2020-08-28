/*
@Time : 2020/8/20 21:07
@Author : liangjiefan
*/
package router

import (
	"go-wire-mvc/config"
	"go-wire-mvc/controllers/coin"
	"go-wire-mvc/controllers/user"
	"go-wire-mvc/router/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Router struct {
	G          *gin.Engine
	HandleUser *user.HandleUser
	HandleCoin *coin.HandleCoin
}

func ProvideRouter(
	conf *config.Config,
	handleUser *user.HandleUser,
	handleCoin *coin.HandleCoin,
) (*Router, error) {

	gin.SetMode(conf.Server.RunMode)
	r := &Router{
		G:          gin.Default(),
		HandleUser: handleUser,
		HandleCoin: handleCoin,
	}

	return r, nil
}

func (r *Router) init() {
	r.G.Use(gin.Recovery())
	r.G.NoRoute(func(c *gin.Context) {
		c.String(http.StatusNotFound, "The incorrect API route")
	})

	r.initRouter()
}

func (r *Router) initRouter() {
	u := r.G.Group("/v1/user")
	u.Use(middleware.LogMiddleware())
	{
		r.HandleUser.InitRouter(u)
	}

	c := r.G.Group("/v1/coin")
	c.Use(middleware.LogMiddleware())
	{
		r.HandleCoin.InitRouter(c)
	}
}
