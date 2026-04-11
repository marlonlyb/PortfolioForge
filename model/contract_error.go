package model

import "fmt"

type ContractError struct {
	StatusHTTP int
	Response   APIErrorResponse
}

func (e *ContractError) Error() string {
	return fmt.Sprintf("status=%d code=%s message=%s", e.StatusHTTP, e.Response.Error.Code, e.Response.Error.Message)
}
