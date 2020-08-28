/*
@Time : 2020/8/20 22:02
@Author : liangjiefan
*/
package coin

import (
	"fmt"
	"go-wire-mvc/dao"
)

type SrvCoin struct {
	dao *dao.Dao
}

func ProvideSrvCoin(dao *dao.Dao) (*SrvCoin, error) {
	return &SrvCoin{dao: dao}, nil
}

func (s *SrvCoin) Add() {
	fmt.Println("SrvCoin Add")
}

func (s *SrvCoin) Get() {
	fmt.Println("SrvCoin Get")
}
