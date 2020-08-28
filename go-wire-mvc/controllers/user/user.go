/*
@Time : 2020/8/20 22:11
@Author : liangjiefan
*/
package user

import (
	"errors"
	"fmt"
	"go-wire-mvc/controllers"
	"go-wire-mvc/service"
	"go-wire-mvc/service/user"

	"github.com/gin-gonic/gin"
)

type HandleUser struct {
	srv *service.Service
}

func ProvideHandleUser(srv *service.Service) (*HandleUser, error) {
	return &HandleUser{
		srv: srv,
	}, nil
}

func (h *HandleUser) InitRouter(routerGroup *gin.RouterGroup) {
	routerGroup.GET("/info", h.Info)
	routerGroup.POST("/login", h.Login)
}

func (h *HandleUser) Login(c *gin.Context) {
	fmt.Println("HandleUser Login")

	var (
		err error
		req loginReq
	)

	if err = c.BindJSON(&req); err != nil {
		controllers.Json(c, nil, errors.New("param error"))
		return
	}

	dto := &user.LoginReqDto{
		Username: req.Username,
		Password: req.Password,
	}

	res, err := h.srv.SrvUser.Login(dto)
	controllers.Json(c, res, err)
}

func (h *HandleUser) Info(c *gin.Context) {
	fmt.Println("HandleUser Info")
	h.srv.SrvUser.Info()
}
