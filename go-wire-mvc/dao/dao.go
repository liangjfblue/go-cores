/*
@Time : 2020/8/20 21:58
@Author : liangjiefan
*/
package dao

import "go-wire-mvc/models"

type Dao struct {
	DB *models.DB
}

func ProvideDao(db *models.DB) (*Dao, error) {
	return &Dao{
		DB: db,
	}, nil
}
