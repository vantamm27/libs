package influxdb

import (
	"github.com/influxdata/influxdb/client/v2"

	"time"
)

type InfluxDBConn struct {
	client.Client
}

type Result struct {
	client.Result
}

func NewConn(addr string, user string, passw string, timeout int64) (*InfluxDBConn, error) {
	cli, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:      addr,
		Username:  user,
		Password:  passw,
		UserAgent: "IotDevicesManager",
		Timeout:   time.Duration(timeout) * time.Second,
	})
	return &InfluxDBConn{cli}, err
}

func (inflx *InfluxDBConn) Query(db string, cmd string) ([]Result, error) {
	q := client.Query{
		Command:   cmd,
		Database:  db,
		Precision: "s",
	}
	var res []client.Result
	if response, err := inflx.Client.Query(q); err == nil {
		if response.Error() != nil {
			return []Result{}, response.Error()
		}
		res = response.Results
	} else {
		return []Result{}, err
	}

	var ret []Result
	for _, v := range res {
		ret = append(ret, Result{v})
	}
	return ret, nil

}
