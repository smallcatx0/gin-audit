package rreqlog

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"time"
)

// FileExists 检查文件是否存在
func FileExists(file string) bool {
	_, err := os.Stat(file)
	if err == nil {
		return true //文件或者文件夹存在
	}
	if os.IsNotExist(err) {
		return false //不存在
	}
	return false //不存在，这里的err可以查到具体的错误信息
}

// TouchDir 创建文件夹
func TouchDir(path string) error {
	dir, _ := filepath.Split(path)
	if FileExists(dir) {
		return nil
	}
	err := os.MkdirAll(dir, 0666)
	return err
}

// 接口实现检测
var _ Recorder = &FileRecord{}

type FileRecord struct {
	path string
	// part string
	fp      *os.File
	recTime time.Time
}

func filenameR(t time.Time, path string) string {
	return path + t.Format("0601") + "/" + t.Format("060102") + ".log"
}

func NewFileRecord(path string) *FileRecord {
	t := time.Now()
	filename := filenameR(t, path)
	TouchDir(filename)
	fp, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("[fileRecord] open file err", err.Error())
	}
	f := FileRecord{
		path:    path,
		fp:      fp,
		recTime: t,
	}
	return &f
}

func (f *FileRecord) autoFile() {
	t := time.Now()
	if t.Month() == f.recTime.Month() && t.Day() == f.recTime.Day() {
		return
	}
	f.fp.Close()
	f.recTime = t
	filename := filenameR(t, f.path)
	TouchDir(filename)
	newfp, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("[fileRecord] open file err", err.Error())
	}
	f.fp = newfp
}

func (f *FileRecord) Record(data map[string]interface{}) {
	w, err := json.Marshal(data)
	if err != nil {
		log.Fatal("[json] marshal fail ", err.Error())
		return
	}
	f.autoFile()
	w = append(w, byte('\n'))
	_, err = f.fp.Write(w)
	if err != nil {
		log.Fatal("[fileRecord] write file fail", err.Error())
	}
}
