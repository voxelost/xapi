package xapi

func login(c *Client) error {
	type loginInput struct {
		UserId   int    `json:"userId"`
		Password string `json:"password"`
	}

	err := c.conn.WriteJSON(command[loginInput]{
		Command: "login",
		Arguments: loginInput{
			UserId:   c.userID,
			Password: c.password,
		},
	})
	if err != nil {
		return err
	}

	respBody := response[any]{}
	err = c.conn.ReadJSON(&respBody)
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

	return err
}
