/*
@Time : 2020/8/20 22:08
@Author : liangjiefan
*/
package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Result struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"code": 1,
		"msg":  "ok",
		"data": data,
	})
}

func Error(c *gin.Context, msg string) {
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  msg,
	})
	c.Abort()
}

func Json(c *gin.Context, data interface{}, err error) {
	if err == nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "ok",
			"data": data,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  "error",
			"desc": err.Error(),
		})
		c.Abort()
	}
}
