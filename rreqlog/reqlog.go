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
		fileConf := confParse.Get("recorder.file")
		dr = NewFileRecord(fileConf.Get("path").String())
	default:
		log.Fatal("[gin-reqlog-md] no such recoder")
	}
	reqlog := NewReqLog(dr, confParse.Get("rules"))
	return func(c *gin.Context) {
		st := time.Now()
		c.Next()
		rt := time.Now().UnixNano() - st.UnixNano()
		go func() {
			reqlog.ParseLog(c, int(rt))
		}()
	}
}

type extRule struct {
	field    string
	funName  string
	funParam []string
}

func (r extRule) String() string {
	s := r.field + " => " + r.funName + "("
	for _, param := range r.funParam {
		s += param + ", "
	}
	return s + ")"
}

type reqlog struct {
	c        *gin.Context
	recorder Recorder
	routers  map[string][]extRule
}

// parseRouteRules 解析路由规则
func parseRouteRules(routRules gjson.Result) map[string][]extRule {
	res := make(map[string][]extRule)
	// 解析路由
	routRules.ForEach(func(rout, rules gjson.Result) bool {
		rules = rules.Get("extra_feilds")
		res[rout.String()] = make([]extRule, 0, 3)
		// 解析额外字段规则
		rules.ForEach(func(field, rule gjson.Result) bool {
			funInfo := strings.SplitN(rule.String(), ":", 2)
			r := extRule{
				field:    field.String(),
				funName:  funInfo[0],
				funParam: make([]string, 0),
			}
			if len(funInfo) >= 2 {
				r.funParam = strings.Split(funInfo[1], ",")
			}
			res[rout.String()] = append(res[rout.String()], r)
			return true
		})
		return true
	})
	return res
}

func NewReqLog(recorder Recorder, rules gjson.Result) *reqlog {
	// 解析规则
	routRules := parseRouteRules(rules)
	l := &reqlog{
		recorder: recorder,
		routers:  routRules,
	}
	return l
}

func (r *reqlog) ParseLog(c *gin.Context, rt int) {
	r.c = c
	k := strings.ToLower(c.Request.Method) + " " + c.Request.RequestURI
	rules, ok := r.routers[k]
	if !ok {
		return
	}
	// 公共字段解析
	logData := map[string]interface{}{
		"method":         c.Request.Method,
		"url":            c.Request.RequestURI,
		"record_time":    time.Now(),
		"run_time":       rt / 10e6,
		"request_header": c.Request.Header,
	}
	// 自定义字段解析
	r.addExtRules(logData, rules)
	// 存储
	r.recorder.Record(logData)
}

func (r *reqlog) addExtRules(logData map[string]interface{}, rules []extRule) {
	for _, arule := range rules {
		h, ok := CustomRules[arule.funName]
		if !ok {
			continue
		}
		logData[arule.field] = h(r.c, arule.funParam...)
	}
}
