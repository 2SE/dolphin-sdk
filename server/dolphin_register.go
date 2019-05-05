package server

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

func registerMethod(version, resource, action string) {
	appInfo.registerMethod(version, resource, action)
}
