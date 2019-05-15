package influxdb

import (
	client "github.com/influxdata/influxdb1-client/v2"
	"testing"
	"time"
)

func TestNewInfluxdb(t *testing.T) {

	cli, err := NewInfluxdb(&client.HTTPConfig{
		Addr:     "http://localhost:8086",
		Username: "rennbon",
		Password: "111111",
		Timeout:  time.Second * 5,
	})
	if err != nil {
		t.Error(err)
		return
	}
	q := client.NewQuery("select * from tracelog", "hashhash", "ns")
	if response, err := cli.Query(q); err == nil && response.Error() == nil {
		t.Log(response.Results)
	} else {
		t.Error(err)
	}
}
