package dolregister

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

func NewRegisterManager() *RegisterManager {
	return &RegisterManager{
		AppInfoer: &AppInfo{},
		Doc: &markdown{
			content:  new(bytes.Buffer),
			resource: make(map[string]struct{}),
		},
	}
}

type RegisterManager struct {
	AppInfoer
	Doc
}

func (a *RegisterManager) Release() {
	a = nil
}
func (a *RegisterManager) RegisterServerOnDolpin(address string) error {
	err := a.AppInfoer.RegisterServerOnDolpin(address)
	if err != nil {
		return err
	}
	a = nil
	return nil
}

type AppInfoer interface {
	//设置本地服务名
	SetAppName(appName string)
	//设置本地服务地址
	SetAddress(address string)
	//添加方法到注册池等到注册到dolphin
	RegisterMethod(version, resource, action string)
	//将本地服务信息注册到dolphin
	RegisterServerOnDolpin(address string) error
}

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

func (a *AppInfo) SetAppName(appName string) {
	a.AppName = appName
}
func (a *AppInfo) SetAddress(address string) {
	a.Address = address
}
func (a *AppInfo) RegisterMethod(version, resource, action string) {
	for _, v := range a.Methods {
		if v.Action == version && v.Resource == resource && v.Action == action {
			return
		}
	}
	a.Methods = append(a.Methods, &MP{version, resource, action})
}

//address http://www.xxx.com:1111
//This method is called after the service starts successfully
func (a *AppInfo) RegisterServerOnDolpin(address string) error {
	appJson, err := json.Marshal(a)
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
