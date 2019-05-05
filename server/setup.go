package server

func Run(c *Config) {
	base.run(c)
}
func RegisterService(s interface{}) error {
	return parseService(s)
}
