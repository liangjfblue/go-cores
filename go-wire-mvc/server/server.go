/*
@Time : 2020/8/20 21:07
@Author : liangjiefan
*/
package server

import (
	"fmt"
	"go-wire-mvc/config"
	router2 "go-wire-mvc/router"
	"log"
	"net/http"
)

type IAppServer interface {
	Init()
	Run()
	Stop()
}

type AppServer struct {
	Config *config.Config
	router *router2.Router
}

func (app *AppServer) Init() {
	app.router.Init()
}

func (app *AppServer) Run() {
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", app.Config.Server.Port), app.router.G).Error())
}

func (app *AppServer) Stop() {
	//TODO clean something
}
