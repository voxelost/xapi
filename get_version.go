package xapi

type GetVersionResponse struct {
	Version string `json:"version"`
}

func (c *client) GetVersion() (string, error) {
	versionResponse, err := getSync[any, GetVersionResponse](c, "getVersion", nil)
	if err != nil {
		return "", err
	}
	return versionResponse.Version, nil
}
