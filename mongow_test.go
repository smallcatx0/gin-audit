package gaudit

import (
	"testing"
)

func TestMongoRecoder(t *testing.T) {
	r := NewMongoWtDef(
		"mongodb://root:example@localhost:27017/?authSource=admin",
		"reqLog", "demo",
	)
	content := map[string]interface{}{"name": "小李", "age": 25}
	r.Record(content)
	r.Record(content)
}
