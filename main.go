package main

import "github.com/gin-gonic/gin"

func main() {
	r := gin.Default()
	r.Use()
	routerHandler(r)
	r.Run(":8090")
}

func routerHandler(r *gin.Engine) {
	v1 := r.Group("/v1")
	v1.GET("/info", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"code": 1,
			"msg":  "操作成功",
			"data": gin.H{
				"id": 1, "name": "tkkk",
			},
		})
	})
	v1.POST("/info", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"code": 1,
			"msg":  "操作成功",
			"data": gin.H{
				"autoId": 5,
			},
		})
	})
}
