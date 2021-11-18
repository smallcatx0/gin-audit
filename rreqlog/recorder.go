package rreqlog

import (
	"io/ioutil"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func ReqLog(confpath string) gin.HandlerFunc {
	mongo := NewMongoWtDef("mongodb://root:example@localhost:27017/?authSource=admin", "reqLog", "demo")
	record := NewReqLog(confpath, mongo)
	return func(c *gin.Context) {
		st := time.Now()
		c.Next()
		rt := time.Now().Sub(st)
		record.ParseLog(c, rt)
	}
}

type extRule struct {
	field string
	funcs []struct {
		funcName string
		funParam []string
	}
}

// func ParseExtRule() []extRule {

// }

type reqlog struct {
	rulec    []byte
	c        *gin.Context
	rules    extRule
	recorder Recorder
}

func NewReqLog(confPath string, recorder Recorder) *reqlog {
	rulec, err := ioutil.ReadFile(confPath)
	if err != nil {
		log.Fatal("read config file faile ! err=", err.Error())
	}
	l := &reqlog{
		recorder: recorder,
		rulec:    rulec,
	}
	return l
}

func (r *reqlog) ParseLog(c *gin.Context, rt time.Duration) {
	r.c = c
	// todo:匹配路由
	// k := strings.ToLower(c.Request.Method) + c.Request.RequestURI
	// todo:解析请求以及响应参数

	// todo:调用存储

}
