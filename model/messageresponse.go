package model

/* Esto es lo que entregamos siempre al cliente, ya sea mensaje o error
recuerda que los errores son tratados con la función:
HTTPErrorHandler(err error, c echo.Context) la cual ejecuta un:
c.JSON(int, domain.MessageResponse) */

type Response struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type Responses []Response

type MessageResponse struct {
	Data     interface{} `json:"data"`
	Errors   Responses   `json:"errors"`
	Messages Responses   `json:"messages"`
}

/*
{
	"data": {un objeto o un array de objetos},
	"errors": [
	{"code": "unexpected", "message": "algo pasó"},
	{"code": "not_found", "message": "no se encontró"}
	],
	"messages": [
	{"code": "ok", "message": "ok"},
	{"code": "record_created", "message": "registro creado"}
	]
}
*/
