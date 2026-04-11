package response

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/marlonlyb/portfolioforge/model"
)

/* Esta función la utilizaremos para el tratamiento de los
errores producidos en los procesos del servidor*/

/* HTTPErrorHandler is a centralized HTTP error handler.
type HTTPErrorHandler func(err error, c Context) */

func HTTPErrorHandler(err error, c echo.Context) {
	contractErr, ok := err.(*model.ContractError)
	if ok {
		_ = c.JSON(contractErr.StatusHTTP, contractErr.Response)
		return
	}

	/* verificamos si es del tipo *model.Error */
	e, ok := err.(*model.Error)
	if ok {
		_ = c.JSON(getResponseError(e))
		return
	}

	/* Check echo error: HTTPError represents an error that occurred
	while handling a request.
	type HTTPError struct {
		Code     int         `json:"-"`
		Message  interface{} `json:"message"`
		Internal error       `json:"-"` // Stores the error returned
										// by an external dependency
	} */

	echoErr, ok := err.(*echo.HTTPError)
	if ok {
		msg, ok := echoErr.Message.(string)
		if !ok {
			msg = "echo.HTTP Error, we are working to solve it..."
		}

		_ = c.JSON(echoErr.Code, model.MessageResponse{
			Errors: model.Responses{
				{Code: UnexpectedError, Message: "echo.HTTP Error, " + msg},
			},
		})
		return
	}

	/* if the handler not return a "model.Error" or "*echo.HTTPError"
	the url return a generic error JSON response */
	_ = c.JSON(http.StatusInternalServerError, model.MessageResponse{
		Errors: model.Responses{
			{Code: UnexpectedError, Message: "Internal Error, we are working to solve it..."},
		},
	})
}

/* Esta firma cumple con el mensaje que se entrega al cliente
en ECHO c.JSON(...)*/

func getResponseError(err *model.Error) (int, model.MessageResponse) {

	outputStatus := 0
	outputResponse := model.MessageResponse{}

	/* Si yo uso response.error entonces ya tiene Code y StatusHTTP
	ademas de no poder modificar el APIMessage, mientras que si
	quiero un APIMessage particular genero un model.error*/

	if !err.HasCode() {
		err.Code = UnexpectedError
	}

	if !err.HasStatusHTTP() {
		err.StatusHTTP = http.StatusInternalServerError
	}

	if err.HasData() {
		outputResponse.Data = err.Data
	}

	outputStatus = err.StatusHTTP
	outputResponse.Errors = model.Responses{
		{Code: err.Code, Message: err.APIMessage},
	}

	return outputStatus, outputResponse
}
