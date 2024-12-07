package xapi

func getSync[T, R any](c *Client, command string, data T) (R, error) {
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

func writeJSON[T any](c *Client, command string, data T) error {
	cmd := newCommand(command, data)
	return c.conn.WriteJSON(cmd)
}

func getResponse[T any](c *Client, res *T) error {
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
