package log

import (
	"fmt"
	"github.com/2se/dolphin-sdk/database/influxdb"
	tw "github.com/RussellLuo/timingwheel"
	influx "github.com/influxdata/influxdb1-client/v2"
	"github.com/sirupsen/logrus"
	"strconv"
	"sync"
	"time"
)

const (
	empty  = ""
	failed = false
	msg    = "message"
)

type InfluxHook struct {
	sync.Mutex
	ticker        *tw.TimingWheel
	client        influx.Client
	tags          []string
	database      string
	measurement   string
	batchSize     int
	batchInterval time.Duration
	pointChan     chan *influx.Point
	bp            influx.BatchPoints
}

func NewInfluxHook(c *Config) (*InfluxHook, error) {

	cli, err := influxdb.NewInfluxdb(c.DBConf)
	if err != nil {
		return nil, err
	}
	bp, err := influx.NewBatchPoints(influx.BatchPointsConfig{
		Precision: "ms",
		Database:  c.DataBase,
	})
	if err != nil {
		return nil, err
	}
	hk := &InfluxHook{
		ticker:        tw.NewTimingWheel(time.Millisecond*100, 600),
		client:        cli,
		database:      c.DataBase,
		measurement:   c.Measurement,
		batchSize:     c.BatchSize,
		batchInterval: c.BatchInterval,
		tags:          c.Tags,
		pointChan:     make(chan *influx.Point, 5000),
		bp:            bp,
	}
	go hk.writePoint()
	go hk.timeTrigger()
	return hk, nil
}
func (h *InfluxHook) writePoint() {
	for ch := range h.pointChan {
		h.bp.AddPoint(ch)
		if len(h.bp.Points()) >= h.batchSize {
			h.save()
		}
	}
}

func (h *InfluxHook) timeTrigger() {
	h.ticker.Start()
	defer h.ticker.Stop()
	ch := make(chan struct{})
	for {
		h.ticker.AfterFunc(h.batchInterval, func() {
			h.save()
			ch <- struct{}{}
		})
		<-ch
	}
}

func (h *InfluxHook) save() {
	h.Lock()
	defer h.Unlock()
	if len(h.bp.Points()) > 0 {
		h.client.Write(h.bp)
	}
}

// Levels implementation allows for level logging.
func (h *InfluxHook) Levels() []logrus.Level {
	return logrus.AllLevels
}
func (h *InfluxHook) Fire(entry *logrus.Entry) error {
	tags := map[string]string{
		"level": entry.Level.String(),
	}
	for _, tag := range h.tags {
		if tagValue, ok := getTag(entry.Data, tag); ok {
			tags[tag] = tagValue
		}
	}
	fields := make(map[string]interface{}, len(entry.Data)+1)
	fields[msg] = entry.Message
	for k, v := range entry.Data {
		fields[k] = v
	}
	for _, tag := range h.tags {
		delete(fields, tag)
	}
	pt, err := influx.NewPoint(h.measurement, tags, fields, entry.Time)
	if err != nil {
		return fmt.Errorf("Could not create new InfluxDB point: %v", err)
	}
	h.pointChan <- pt
	return nil
}

func getTag(fields logrus.Fields, tag string) (string, bool) {
	value, ok := fields[tag]
	if !ok {
		return empty, ok
	}
	switch vs := value.(type) {
	case fmt.Stringer:
		return vs.String(), ok
	case string:
		return vs, ok
	case byte:
		return string(vs), ok
	case int:
		return strconv.FormatInt(int64(vs), 10), ok
	case int32:
		return strconv.FormatInt(int64(vs), 10), ok
	case int64:
		return strconv.FormatInt(vs, 10), ok
	case uint:
		return strconv.FormatUint(uint64(vs), 10), ok
	case uint32:
		return strconv.FormatUint(uint64(vs), 10), ok
	case uint64:
		return strconv.FormatUint(vs, 10), ok
	default:
		return empty, failed
	}
	return empty, failed
}
