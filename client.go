package xapi

const (
	SYNC_PORT_REAL              = 5112
	SYNC_PORT_DEMO              = 5124
	SYNC_WEBSOCKET_ADDRESS_REAL = "wss://ws.xtb.com/real"
	SYNC_WEBSOCKET_ADDRESS_DEMO = "wss://ws.xtb.com/demo"
	API_ADDRESS_BASE            = "https://xapi.xtb.com"
)

type ClientMode string

var (
	ClientModeDemo ClientMode = "demo"
	ClientModeReal ClientMode = "real"
)

type ApiError struct {
	Code    string
	Message string
}

func (e ApiError) Error() string {
	return e.Message
}

type loginInput struct {
	UserId   int    `json:"userId"`
	Password string `json:"password"`
}

func getSync[T, R any](c *client, command string, data T) (R, error) {
	c.m.Lock()
	defer c.m.Unlock()

	var r R
	err := writeJSON(c, command, data)
	if err != nil {
		return r, err
	}

	err = getResponse(c, &r)
	return r, err
}

func writeJSON[T any](c *client, command string, data T) error {
	cmd := newCommand(command, data)
	return c.conn.WriteJSON(cmd)
}

func getResponse[T any](c *client, res *T) error {
	respBody := response[T]{}
	err := c.conn.ReadJSON(&respBody)
	if err != nil {
		return err
	}

	if !respBody.Status {
		var errorCode, errorDescription string
		if respBody.ErrorCode != nil {
			errorCode = *respBody.ErrorCode
		}
		if respBody.ErrorDescription != nil {
			errorDescription = *respBody.ErrorDescription
		}

		return ApiError{
			Code:    errorCode,
			Message: errorDescription,
		}
	}

	if respBody.ReturnData != nil {
		*res = *respBody.ReturnData
	}

	return nil
}
