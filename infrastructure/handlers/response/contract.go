package response

import (
	"net/http"

	"github.com/marlonlyb/portfolioforge/model"
)

func ContractData(status int, data interface{}) (int, map[string]interface{}) {
	return status, map[string]interface{}{"data": data}
}

func ContractOK(data interface{}) (int, map[string]interface{}) {
	return ContractData(http.StatusOK, data)
}

func ContractCreated(data interface{}) (int, map[string]interface{}) {
	return ContractData(http.StatusCreated, data)
}

func ContractError(status int, code, message string, details ...model.APIErrorDetail) *model.ContractError {
	payload := model.APIErrorPayload{
		Code:    code,
		Message: message,
	}

	if len(details) > 0 {
		payload.Details = details
	}

	return &model.ContractError{
		StatusHTTP: status,
		Response: model.APIErrorResponse{
			Error: payload,
		},
	}
}
