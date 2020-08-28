/*
@Time : 2020/8/20 21:35
@Author : liangjiefan
*/
package models

import (
	"github.com/jinzhu/gorm"
)

type TBCoin struct {
	gorm.Model
	UserId int  `gorm:"column:user_id;" description:"关联用户id"`
	Coin   uint `gorm:"column:coin;" description:"金币"`
}
