/*
@Time : 2020/8/21 0:26
@Author : liangjiefan
*/
package user

type LoginReqDto struct {
	Username string
	Password string
}

type LoginRespDto struct {
	HeadImage string
	Nickname  string
}
