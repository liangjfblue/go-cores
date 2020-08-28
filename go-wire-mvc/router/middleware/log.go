/*
@Time : 2020/8/28 22:11
@Author : liangjiefan
*/
package middleware

import (
	"log"

	"github.com/gin-gonic/gin"
)

func LogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println(c.Request.URL.Path)

		c.Next()
	}
}
