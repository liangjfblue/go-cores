/*
@Time : 2020/8/20 23:20
@Author : liangjiefan
*/
package user

type loginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResp struct {
	HeadImage string `json:"headImage"`
	Nickname  string `json:"nickname"`
}
