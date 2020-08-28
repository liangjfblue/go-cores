/*
@Time : 2020/8/20 22:02
@Author : liangjiefan
*/
package user

import (
	"fmt"
	"go-wire-mvc/dao"
)

type SrvUser struct {
	dao *dao.Dao
}

func ProvideSrvUser(dao *dao.Dao) (*SrvUser, error) {
	return &SrvUser{dao: dao}, nil
}

func (s *SrvUser) Login(req *LoginReqDto) (*LoginRespDto, error) {
	fmt.Println("SrvUser Login")

	user, err := s.dao.GetTBUser("username = ?", req.Username)
	if err != nil {
		return nil, err
	}

	if user.Password != req.Password {
		return nil, err
	}

	resp := &LoginRespDto{
		HeadImage: user.HeadImage,
		Nickname:  user.Nickname,
	}

	return resp, nil
}

func (s *SrvUser) Info() {
	fmt.Println("SrvUser Info")
}
