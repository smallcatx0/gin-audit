package rreqlog_test

import (
	"gin-reqlog-md/rreqlog"
	"testing"
)

func TestMongoRecoder(t *testing.T) {
	r := rreqlog.NewMongoWtDef(
		"mongodb://root:example@localhost:27017/?authSource=admin",
		"reqLog", "demo",
	)
	content := map[string]interface{}{"name": "小李", "age": 25}
	r.Record(content)
	r.Record(content)
}
