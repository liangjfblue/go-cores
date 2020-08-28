/*
@Time : 2020/8/21 0:15
@Author : liangjiefan
*/
package dao

import (
	"go-wire-mvc/models"
)

func (d *Dao) CreateTBCoin(t *models.TBCoin) (err error) {
	err = d.DB.GormDB.Create(t).Error
	return
}

func (d *Dao) GetTBCoin(query interface{}, args ...interface{}) (t *models.TBCoin, err error) {
	err = d.DB.GormDB.Where(query, args).First(&t).Error
	return
}

func (d *Dao) DeleteTBCoin(t *models.TBCoin) (err error) {
	err = d.DB.GormDB.Where("id", t.ID).Delete(t).Error
	return
}

func (d *Dao) UpdateTBCoin(t *models.TBCoin) (err error) {
	err = d.DB.GormDB.Where("id", t.ID).Update(t).Error
	return
}
