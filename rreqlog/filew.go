package rreqlog

// 接口实现检测
var _ Recorder = &FileRecord{}

type FileRecord struct {
	path string
	// part string
}

func NewFileRecord(path string) FileRecord {
	f := FileRecord{
		path: path,
	}
	return f
}

func (f *FileRecord) Record(data map[string]interface{}) {

}
