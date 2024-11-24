package xapi

func (c *client) Ping() error {
	_, err := getSync[any, any](c, "ping", nil)
	return err
}
