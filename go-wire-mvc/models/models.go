/*
@Time : 2020/8/20 21:07
@Author : liangjiefan
*/
package models

import (
	"fmt"
	"go-wire-mvc/config"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type DB struct {
	GormDB *gorm.DB
}

func ProvideDB(conf *config.Config) (*DB, error) {
	var err error
	db := new(DB)

	str := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		conf.Mysql.User, conf.Mysql.Password, conf.Mysql.Addr, conf.Mysql.Db)
	db.GormDB, err = gorm.Open("mysql", str)
	if err != nil {
		return nil, err
	}

	db.GormDB.LogMode(true)
	db.GormDB.SingularTable(true)
	db.GormDB.DB().SetMaxIdleConns(conf.Mysql.MaxIdleCons)
	db.GormDB.DB().SetMaxOpenConns(conf.Mysql.MaxOpenCons)

	db.GormDB.AutoMigrate(new(TBUser), new(TBCoin))
	return db, err
}
