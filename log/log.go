package log

import (
	influx "github.com/influxdata/influxdb1-client/v2"
	"github.com/sirupsen/logrus"

	"time"
)

type Config struct {
	//是否使用db
	WithDB bool
	//DB配置
	DBConf *influx.HTTPConfig
	//批量写入长度
	BatchSize int
	//批量写入间隔
	BatchInterval time.Duration
	//log存储的数据库名
	DataBase string
	//log存储的表名
	Measurement string
	//log落库索引键名
	Tags []string
}

func NewLog(c *Config) (*logrus.Logger, error) {
	log := logrus.New()
	if c.WithDB {
		hook, err := NewInfluxHook(c)
		if err != nil {
			return nil, err
		}
		log.Hooks.Add(hook)
	}
	return log, nil
}
