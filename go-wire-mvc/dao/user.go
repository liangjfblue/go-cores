/*
@Time : 2020/8/21 0:12
@Author : liangjiefan
*/
package dao

import "go-wire-mvc/models"

func (d *Dao) CreateTBUser(t *models.TBUser) (err error) {
	err = d.DB.GormDB.Create(t).Error
	return
}

func (d *Dao) GetTBUser(query interface{}, args ...interface{}) (t models.TBUser, err error) {
	err = d.DB.GormDB.Where(query, args).First(&t).Error
	return
}

func (d *Dao) DeleteTBUser(t *models.TBUser) (err error) {
	err = d.DB.GormDB.Where("id", t.ID).Delete(t).Error
	return
}

func (d *Dao) UpdateTBUser(t *models.TBUser) (err error) {
	err = d.DB.GormDB.Where("id", t.ID).Update(t).Error
	return
}
