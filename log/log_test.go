package log

import (
	influx "github.com/influxdata/influxdb1-client/v2"
	"github.com/sirupsen/logrus"
	"os"
	"testing"
	"time"
)

func TestLog(t *testing.T) {
	log, err := NewLog(&Config{
		//是否使用db
		WithDB: false,
		//DB配置
		DBConf: &influx.HTTPConfig{
			Addr:     "http://localhost:8086",
			Username: "rennbon",
			Password: "111111",
			Timeout:  time.Second * 5,
		},
		//批量写入长度
		BatchSize: 20,
		//批量写入间隔
		BatchInterval: time.Second,
		//log存储的数据库名
		DataBase: "hashhash",
		//log存储的表名
		Measurement: "tracelog",
		//log落库索引键名
		Tags: []string{"level", "spanId", "appName", "resource"},
	})

	if err != nil {
		t.Error(err)
		return
	}
	//设置格式
	log.SetFormatter(&logrus.TextFormatter{ForceColors: true, FullTimestamp: true})
	//设置控制台输出
	log.SetOutput(os.Stdout)
	//设置落库等级
	log.SetLevel(logrus.TraceLevel)
	log.WithFields(logrus.Fields{
		"spanId":   "spanId-123456",
		"appName":  "user-node1",
		"resource": "user",
	}).Trace("This is a trace id: 12345")

	log.WithFields(logrus.Fields{
		"spanId":   "spanId-123457",
		"appName":  "user-node1",
		"resource": "user",
	}).Error("This is an error")
	time.Sleep(time.Second * 5)
}
