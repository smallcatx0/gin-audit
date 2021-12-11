# gin接口审计中间件



## 背景介绍



## 用法

要求:

- go >= 1.14
- mongo [可选]

安装

```GO
go get -u "github.com/smallcatx0/gaudit"
```



### 快速开始

```GO
func main() {
	r := gin.Default()
	// 添加些自定义的字段获取逻辑
	gaudit.AddFieldHandle("userphone", func(c *gin.Context, s ...string) string {
		token := c.GetHeader("token")
		log.Print(token)
		// ... token 换用户据信息
		return "110"
	})

	r.Use(gaudit.ReqLog("./conf.json"))
	routerHandler(r) 
	r.Run(":8090")
}

```



配置文件

```json
{
    "recorder":{
        "choose": "mongo",
        "mongo": {
            "dsn": "mongodb://root:example@mongo.serv:27017/?authSource=admin",
            "db": "rrlog",
            "collection": "demo",
            "part": "week"
        },
        "file": {
            "path": "/tmp/reqlog/"
        }
    },
    "rules":{ // 记录规则
        "get /v1/info": { // 路由地址 
            "extra_feilds": {  // 额外字段
                "type": "echo:课程", 
                "phone": "userphone"
            }
        }
    }
}
```





