package gaudit

import (
	"bytes"
	"io/ioutil"
	"log"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
)

// gin-requestlog 中间件
func ApiAudit(confpath string) gin.HandlerFunc {
	confC, err := ioutil.ReadFile(confpath)
	h := func(c *gin.Context) {}
	if err != nil {
		log.Println("read config file faile ! err=", err.Error())
		return h
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
		log.Fatal("[gin-audit-md] no such recoder")
	}
	reqlog := NewReqLog(dr, confParse.Get("rules"))
	h = func(c *gin.Context) {
		if !reqlog.shouldRecord(c) {
			return
		}
		reqBody, _ := c.GetRawData()
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(reqBody))
		rw := ResponseWriterWrapper{
			Body:           bytes.NewBufferString(""),
			ResponseWriter: c.Writer,
		}
		c.Writer = rw
		st := time.Now()
		c.Next()
		rt := time.Now().UnixNano() - st.UnixNano()
		resBody := rw.Body.Bytes()
		go func() {
			reqlog.ParseAndRecorde(rt, c, reqBody, resBody)
		}()
	}
	return h
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
	return s[:len(s)-2] + ")"
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

type reqlog struct {
	recorder Recorder
	routers  map[string][]extRule
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

func (r *reqlog) ruleKey(c *gin.Context) string {
	return strings.ToLower(c.Request.Method) + " " + c.Request.URL.Path
}

func (r *reqlog) shouldRecord(c *gin.Context) bool {
	_, ok := r.routers[r.ruleKey(c)]
	return ok
}

func (r *reqlog) ParseAndRecorde(
	rt int64,
	c *gin.Context,
	req, res []byte,
) {
	k := r.ruleKey(c)
	rules, ok := r.routers[k]
	if !ok {
		return
	}
	// 公共字段解析
	logData := map[string]interface{}{
		"method":      c.Request.Method,
		"url":         c.Request.URL.Path,
		"query":       c.Request.URL.Query().Encode(),
		"form":        c.Request.Form,
		"record_time": time.Now(),
		"run_time":    rt / 10e6,
		"header":      c.Request.Header,
		"body":        string(req),
		"status":      c.Writer.Status(),
		"size":        c.Writer.Size(),
		"resp_h":      c.Writer.Header(),
		"resp":        string(res),
	}
	// 自定义字段解析
	r.addExtRules(logData, rules, c, req, res)
	// 存储
	r.recorder.Record(logData)
}

func (r *reqlog) addExtRules(
	data map[string]interface{},
	rules []extRule,
	c *gin.Context,
	req, res []byte,
) {
	for _, arule := range rules {
		h, ok := CustomRules[arule.funName]
		if !ok {
			continue
		}
		data[arule.field] = h(c, arule.funParam...)
	}
}

type ResponseWriterWrapper struct {
	gin.ResponseWriter
	Body *bytes.Buffer // 缓存
}

func (w ResponseWriterWrapper) Write(b []byte) (int, error) {
	w.Body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w ResponseWriterWrapper) WriteString(s string) (int, error) {
	w.Body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}
