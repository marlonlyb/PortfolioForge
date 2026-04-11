package model

import "fmt"

/* Esto se va a usar para errores de los procesos del servidor dándoles este formato,
trabajando con la función centralizada HTTPErrorHandler(err error, c echo.Context)
por lo tanto se usa Error en la definición de la función HTTPErrorHandler, ademas
de los error no directamente de ECHO pero si relacionado a sus procesos, por ejemplo:
BindFailed(err error) *domain.Error, y al que sí esta relacionado con ECHO:
Error(c echo.Context, who string, err error) *domain.Error */

type Error struct {
	Code       string
	Err        error
	Who        string
	StatusHTTP int
	Data       interface{} // any con generix en go
	APIMessage string
	UserID     string
}

func NewError() Error {
	return Error{}
}

//"Error struct" implementa el método "Error() string" por lo que cumple con el error type
func (e *Error) Error() string {
	return fmt.Sprintf("Code: %s, Err: %v, Who: %s, Status: %d, Data: %v, UserID: %s",
		e.Code,
		e.Err,
		e.Who,
		e.StatusHTTP,
		e.Data,
		e.UserID)
}

//Preguntamos si los campos están vacíos (Code)
func (e *Error) HasCode() bool {
	return e.Code != ""
}

//Preguntamos si los campos están vacíos (StatusHTTP)
func (e *Error) HasStatusHTTP() bool {
	return e.StatusHTTP > 0
}

//Preguntamos si los campos están vacíos (Data)
func (e *Error) HasData() bool {
	return e.Data != ""
}
