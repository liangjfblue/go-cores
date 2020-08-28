/*
@Time : 2020/8/20 21:14
@Author : liangjiefan
*/
package coin

import (
	"fmt"
	"go-wire-mvc/service"

	"github.com/gin-gonic/gin"
)

type HandleCoin struct {
	srv *service.Service
}

func ProvideHandleCoin(srv *service.Service) (*HandleCoin, error) {
	return &HandleCoin{
		srv: srv,
	}, nil
}

func (h *HandleCoin) InitRouter(routerGroup *gin.RouterGroup) {
	routerGroup.GET("/add", h.Add)
	routerGroup.POST("/get", h.Get)
}

func (h *HandleCoin) Add(c *gin.Context) {
	fmt.Println("HandleCoin Add")
	h.srv.SrvCoin.Add()
}

func (h *HandleCoin) Get(c *gin.Context) {
	fmt.Println("HandleCoin Get")
	h.srv.SrvCoin.Get()
}
