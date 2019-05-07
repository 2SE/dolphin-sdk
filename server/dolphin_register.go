package server

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

var appInfo = new(AppInfo)

type AppInfo struct {
	PeerName string
	AppName  string
	Address  string
	Methods  []*MP
}
type MP struct {
	Reversion string
	Resource  string
	Action    string
}

func (a *AppInfo) setAppName(appName string) {
	a.AppName = appName
}
func (a *AppInfo) setAddress(address string) {
	a.Address = address
}
func (a *AppInfo) registerMethod(version, resource, action string) {
	for _, v := range a.Methods {
		if v.Action == version && v.Resource == resource && v.Action == action {
			return
		}
	}
	a.Methods = append(a.Methods, &MP{version, resource, action})
}

//address http://www.xxx.com:1111
//This method is called after the service starts successfully
func registerServerOnDolpin(address string) error {
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
		logrus.Info("The service registered successfully with the dolphin.")
		return nil
	} else {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return errors.New(string(body))
	}
}
