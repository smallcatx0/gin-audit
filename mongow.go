package gaudit

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func IsSameWeek(t1, t2 time.Time) bool {
	return ((t1.Unix() - 288000) / 604800) == ((t2.Unix() - 288000) / 604800)
}

func WeekStime(t time.Time) time.Time {
	if (t.Unix()-316800)/86400%7 == 0 {
		return time.Unix(t.Unix()/86400*86400, 0)
	}
	return time.Unix((t.Unix())/604800*604800+316800, 0)
}

// 接口实现检测
var _ Recorder = &MongoRecord{}

type MongoRecord struct {
	dsn          string
	db, collName string

	recordTime time.Time
	partMod    string
	cli        *mongo.Client
	coll       *mongo.Collection
}

func NewMongoWtDef(dsn, db, collName string) *MongoRecord {
	m := &MongoRecord{
		dsn:      dsn,
		db:       db,
		collName: collName,
		partMod:  "week",
	}

	err := m.conn()
	if err != nil {
		log.Fatal("mongo connect failed err=", err.Error())
	}
	now := time.Now()
	switch m.partMod {
	case "week":
		m.recordTime = WeekStime(now)
	default:
		m.recordTime = time.Date(now.Year(), now.Month(), 0, 0, 0, 0, 0, now.Location())
	}
	m.chooseColl()
	return m
}

func (m *MongoRecord) conn() error {
	cliOpt := options.Client().ApplyURI(m.dsn)
	cliOpt.SetMaxPoolSize(10)
	cliOpt.SetConnectTimeout(time.Millisecond * 2000)
	cliOpt.SetConnectTimeout(time.Millisecond * 10000)

	c, _ := context.WithTimeout(context.Background(), 5*time.Second)
	var err error
	m.cli, err = mongo.Connect(c, cliOpt)
	if err != nil {
		return err
	}
	// ping 2s超时
	timeout, _ := context.WithTimeout(context.Background(), 2*time.Second)
	err = m.cli.Ping(timeout, nil)
	if err != nil {
		return err
	}
	return nil
}

func (m *MongoRecord) chooseColl() {
	collName := fmt.Sprintf("%s_%d%d%d",
		m.collName,
		m.recordTime.Year()-2000,
		m.recordTime.Month(),
		m.recordTime.Day(),
	)
	m.coll = m.cli.Database(m.db).Collection(collName)
}

func (m *MongoRecord) autoCollect() {
	now := time.Now()
	switch m.partMod {
	case "week":
		// 判断同周
		if IsSameWeek(now, m.recordTime) {
			return
		}
		m.recordTime = now
		m.chooseColl()
	default:
		// 判断同月
		if (int(m.recordTime.Month()) + 100*m.recordTime.Year()) == (int(now.Month()) + 100*now.Year()) {
			return
		}
		m.recordTime = now
		m.chooseColl()
	}
}

func (m *MongoRecord) Record(content map[string]interface{}) {
	m.autoCollect()
	_, err := m.coll.InsertOne(context.TODO(), content)
	if err != nil {
		log.Println("[reqlog] save to mongo faile err=", err.Error())
	}
}
