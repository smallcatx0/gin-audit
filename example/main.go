package main

import (
	"gaudit"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	// 添加些自定义的字段获取逻辑
	gaudit.AddFieldHandle("userphone", func(c *gin.Context, s ...string) string {
		token := c.GetHeader("token")
		log.Print(token)
		// ... token 换用户据信息
		return "110"
	})

	r.Use(gaudit.ApiAudit("./conf.json"))
	routerHandler(r)
	r.Run(":8090")
}

func routerHandler(r *gin.Engine) {
	v1 := r.Group("/v1")
	v1.GET("/info", func(c *gin.Context) {
		time.Sleep(time.Millisecond * 20)
		c.JSON(200, gin.H{
			"code": 1,
			"msg":  "操作成功",
			"data": gin.H{
				"id": 1, "name": "tkkk",
			},
		})
	})
	v1.POST("/info", func(c *gin.Context) {
		time.Sleep(time.Millisecond * 20)
		c.JSON(200, gin.H{
			"code": 1,
			"msg":  "操作成功",
			"data": gin.H{
				"autoId": 5,
			},
		})
	})
}
