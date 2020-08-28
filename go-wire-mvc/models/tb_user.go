/*
@Time : 2020/8/20 21:35
@Author : liangjiefan
*/
package models

import "github.com/jinzhu/gorm"

type TBUser struct {
	gorm.Model
	Username  string `gorm:"column:username;type:varchar(255)" description:"用户名"`
	Password  string `gorm:"column:password;type:varchar(128)" description:"密码"`
	Nickname  string `gorm:"column:nickname;type:varchar(255)" description:"昵称"`
	Address   string `gorm:"column:address;type:varchar(100)" description:"地址"`
	HeadImage string `gorm:"column:head_image;type:varchar(255)" description:"头像图片"`
	Phone     string `gorm:"column:phone;type:varchar(100)" description:"手机号码"`
}

func (d *TBUser) TableName() string {
	return "tb_user"
}
