package influxdb

import (
	influx "github.com/influxdata/influxdb1-client/v2"
)

type Config struct {
	Addr     string
	UserName string
	Password string
}

//TODO udpClient  等client-go稳定后添加或替换influxdb1-client/v2,目前出来时间太短，暂时不用
//httpClient
func NewInfluxdb(c *influx.HTTPConfig) (influx.Client, error) {
	cli, err := influx.NewHTTPClient(*c)
	if err != nil {
		return nil, err
	}
	_, _, err = cli.Ping(c.Timeout)
	if err != nil {
		return nil, err
	}
	return cli, nil
}
