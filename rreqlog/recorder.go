package rreqlog

import (
	"io/ioutil"
	"log"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
)

// gin-requestlog 中间件
func ReqLog(confpath string) gin.HandlerFunc {
	confC, err := ioutil.ReadFile(confpath)
	if err != nil {
		log.Fatal("read config file faile ! err=", err.Error())
	}
	confParse := gjson.ParseBytes(confC)
	channel := confParse.Get("recorder.choose").String()
	// 记录者配置
	var dr Recorder
	switch channel {
	case "mongo":
		mongoConf := confParse.Get("recorder.mongo")
		dsn := mongoConf.Get("dsn").String()
		db := mongoConf.Get("db").String()
		collName := mongoConf.Get("collection").String()
		dr = NewMongoWtDef(dsn, db, collName)
	case "file":
		// fileConf := confParse.Get("recorder.file")

	default:
		log.Fatal("[gin-reqlog-md] no such recoder")
	}
	reqlog := NewReqLog(dr, confParse.Get("rules"))
	return func(c *gin.Context) {
		st := time.Now()
		c.Next()
		rt := time.Now().Sub(st)
		reqlog.ParseLog(c, rt)
	}
}

type extRule struct {
	field string
	funcs []struct {
		funcName string
		funParam []string
	}
}
type reqlog struct {
	c        *gin.Context
	recorder Recorder
	routers  map[string][]extRule
}

func NewReqLog(recorder Recorder, rules gjson.Result) *reqlog {
	l := &reqlog{
		recorder: recorder,
	}
	return l
}

func ParseRouteRules(routeRules gjson.Result) {

}

func (r *reqlog) ParseLog(c *gin.Context, rt time.Duration) {
	r.c = c
	k := strings.ToLower(c.Request.Method) + c.Request.RequestURI
	_, ok := r.routers[k]
	if ok {
		return
	}
	// todo:解析请求以及响应参数

	// todo:调用存储

}
