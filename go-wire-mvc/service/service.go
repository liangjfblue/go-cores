/*
@Time : 2020/8/20 21:08
@Author : liangjiefan
*/
package service

import (
	"go-wire-mvc/service/coin"
	"go-wire-mvc/service/user"
)

type Service struct {
	SrvUser *user.SrvUser
	SrvCoin *coin.SrvCoin
}

func ProvideService(
	srvUser *user.SrvUser,
	srvCoin *coin.SrvCoin,
) (*Service, error) {
	return &Service{
		SrvUser: srvUser,
		SrvCoin: srvCoin,
	}, nil
}
