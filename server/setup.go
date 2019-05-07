package server

//grpc server start
//address: dolphin address  http://www.xxx.com:1111
//services: business service
func Start(c *Config, services ...interface{}) {
	appInfo.setAppName(c.AppName)
	appInfo.setAddress(c.Address)
	md.setTitle(c.AppName)
	err := parseServices(services...)
	if err != nil {
		panic(err)
	}
	go base.run(c)
	err = registerServerOnDolpin(c.DolphinAddr)
	if err != nil {
		panic(err)
	}
	select {}
}

func StartGrpcOnly(c *Config, services ...interface{}) {
	appInfo.setAppName(c.AppName)
	appInfo.setAddress(c.Address)
	md.setTitle(c.AppName)
	err := parseServices(services...)
	if err != nil {
		panic(err)
	}
	base.run(c)
}
