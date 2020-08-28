//+build wireinject

/*
@Time : 2020/8/20 21:07
@Author : liangjiefan
*/
package server

import (
	"go-wire-mvc/config"
	"go-wire-mvc/controllers/coin"
	"go-wire-mvc/controllers/user"
	"go-wire-mvc/dao"
	"go-wire-mvc/models"
	"go-wire-mvc/router"
	"go-wire-mvc/service"
	coinSrv "go-wire-mvc/service/coin"
	userSrv "go-wire-mvc/service/user"

	"github.com/google/wire"
)

func NewAppServer(path string) (*AppServer, error) {
	panic(wire.Build(
		config.ProvideConfig,
		router.ProvideRouter,

		models.ProvideDB,
		dao.ProvideDao,

		//handle
		user.ProvideHandleUser,
		coin.ProvideHandleCoin,

		//service
		service.ProvideService,
		userSrv.ProvideSrvUser,
		coinSrv.ProvideSrvCoin,

		wire.Struct(new(AppServer), "*"),
	))
}
