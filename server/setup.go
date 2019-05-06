package server

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

func Run(c *Config) {
	base.run(c)
}
func Stop() {
	base.stop()
}
func RegisterService(s interface{}) error {
	return parseService(s)
}

//address http://www.xxx.com:1111
func RegisterServerOnDolpin(address string) error {
	appJson, err := json.Marshal(appInfo)
	if err != nil {
		return err
	}
	resp, err := http.Post(address, "application/json; charset=utf-8", bytes.NewReader(appJson))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		return nil
	} else {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return errors.New(string(body))
	}
}
