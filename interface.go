package gaudit

type Recorder interface {
	Record(content map[string]interface{})
}
