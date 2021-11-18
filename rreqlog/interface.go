package rreqlog

type Recorder interface {
	Record(content map[string]interface{})
}
