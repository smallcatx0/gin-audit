package gaudit

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var timeFomat = "2006-01-02 15:04:05"

func parsetime(t string) time.Time {
	t1, _ := time.ParseInLocation(timeFomat, t, time.Local)
	return t1
}

func TestFilenamer(t *testing.T) {
	ass := assert.New(t)
	ass.Equal(
		"/home/rlog/2109/210916.log",
		filenameR(parsetime("2021-09-16 08:15:09"), "/home/rlog/"),
	)
	ass.Equal(
		"/home/rlog/2111/211111.log",
		filenameR(parsetime("2021-11-11 09:15:09"), "/home/rlog/"),
	)
	ass.Equal(
		"/home/rlog/0809/080915.log",
		filenameR(parsetime("2008-09-15 08:15:09"), "/home/rlog/"),
	)
}

func TestFileW() {

}
