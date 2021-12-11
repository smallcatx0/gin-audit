package gaudit

import (
	"bytes"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
)

type FieldHandle func(*gin.Context, ...string) string

func WatchReqBody(r *http.Request) ([]byte, error) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	r.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	return body, nil
}

// CustomRules 内置自定义规则
var CustomRules = map[string]FieldHandle{
	"echo": func(c *gin.Context, s ...string) string {
		return s[0]
	},
	"req_json_get": func(c *gin.Context, s ...string) string {
		// 获取request请求中的数据
		requestData, _ := c.GetRawData()
		// 将request.Body写回去
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(requestData))
		if len(s) == 0 {
			return string(requestData)
		}
		return gjson.GetBytes(requestData, s[0]).String()
	},
	"resp_json_get": func(c *gin.Context, s ...string) string {
		respData, _ := ioutil.ReadAll(c.Request.Response.Body)
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(respData))
		if len(s) == 0 {
			return string(respData)
		}
		return gjson.GetBytes(respData, s[0]).String()
	},
	"traceid": func(c *gin.Context, s ...string) string {
		return c.Request.Header.Get("x-b3-traceid")
	},
}

func AddFieldHandle(name string, h FieldHandle) {
	CustomRules[name] = h
}
