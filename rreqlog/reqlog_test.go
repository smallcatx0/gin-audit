package rreqlog

import (
	"testing"

	"github.com/tidwall/gjson"
)

var routRules = `
{
	"get /v1/info": {
		"extra_feilds": {
			"id": "respget:id",
			"type": "echo:课程",
			"uid": "authuid",
		}
	}
}
`

func TestRoutRuleParse(t *testing.T) {
	parseRouteRules(gjson.Parse(routRules))
}

func BenchmarkRoutRulePase(b *testing.B) {
	con := gjson.Parse(routRules)
	for i := 0; i < b.N; i++ {
		parseRouteRules(con)
	}

}
