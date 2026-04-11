package response

import (
	"net/http"

	"github.com/labstack/echo/v4"

	//revisar este paquete log.Warnf("%s", e.Error())

	"github.com/marlonlyb/portfolioforge/model"
)

const (
	RecordCreated   = "record_created"
	RecordUpdate    = "record_updated"
	RecordDeleted   = "record_deleted"
	OK              = "ok"
	BindFailed      = "bind_failed"
	AuthError       = "authorizattion_error"
	UnexpectedError = "unexpected_error"
)

type API struct {
}

func New() *API {
	return &API{}
}

/* Estas funciones se utilizan para responder al cliente en ECHO
se ingresará la información en data interface{} y se le da una
información adicional. Retornamos los argumentos para c.JSON */

func (a *API) Created(data interface{}) (int, model.MessageResponse) {
	return http.StatusCreated, model.MessageResponse{
		Data:     data,
		Messages: model.Responses{{Code: RecordCreated, Message: "¡Listo!"}},
	}
}

func (a *API) Updated(data interface{}) (int, model.MessageResponse) {
	return http.StatusOK, model.MessageResponse{
		Data:     data,
		Messages: model.Responses{{Code: RecordUpdate, Message: "¡Listo!"}},
	}
}

func (a *API) Deleted(data interface{}) (int, model.MessageResponse) {
	return http.StatusOK, model.MessageResponse{
		Data:     data,
		Messages: model.Responses{{Code: RecordDeleted, Message: "¡Listo!"}},
	}
}

// El OK es para las funciones como GetAll() o GetByID(), etc.
func (a *API) OK(data interface{}) (int, model.MessageResponse) {
	return http.StatusOK, model.MessageResponse{
		Data:     data,
		Messages: model.Responses{{Code: OK, Message: "¡listo!"}},
	}
}

/* Estas funciones se utilizan para los errores generados en los procesos
del servidor ECHO, le damos forma de *model.Error y la función
HTTPErrorHandler(err error, c echo.Context) se encargará de entrega el CJSON */

// Si ocurre un error en cualquier endpoint, usamos este formato
func (a *API) Error(c echo.Context, who string, err error) *model.Error {
	e := model.NewError()
	e.Code = UnexpectedError
	e.Err = err
	e.Who = who
	//e.Data = ""
	e.StatusHTTP = http.StatusInternalServerError
	e.APIMessage = "response.Error, we are working to solve it..."
	//e.UserID = a.parseUserID(c)

	//log.Warnf("log.Warnf - %s", e.Error())

	return &e
}

/*
	// Usar esta función para responder me va a dar el siguiente resultado:
	{
	    "data": null, (err.Data)
	    "errors": [
	        {
	            "code": "unexpected_error", (err.Code)
	            "message": "response.Error, we are working to solve it..." (err.APIMessage)
	        }
	    ],
	    "messages": null
	}
	// Además de un StatusHTTP: 500 Internal Server Error
*/

/* Ahora: un c.Bind del mismo ECHO da un error del siguiente formato:
	code=400,
	message=unexpected EOF,
	internal=unexpected EOF,
Por lo tanto, el StatusHTTP debe ser cambiado a 400 Bad Request,
ademas el "code": "bind_failed" y no quiero mensaje al cliente */

func (a *API) BindFailed(c echo.Context, who string, err error) *model.Error {
	e := model.NewError()
	e.Code = BindFailed
	e.Err = err
	e.Who = who
	//e.Data = ""
	e.StatusHTTP = http.StatusBadRequest
	//e.APIMessage = ""
	//e.UserID = a.parseUserID(c)

	//log.Warnf("log.Warnf - %s", e.Error())

	return &e
}

/*Ahora: Que pasa si quiero darle un mensaje al cliente
en un error en particular como en el Login(),
ya tendríamos el WHO y como no quiere userID tampoco necesito
el c echo.Context, si quiero mensaje al cliente...*/

func (a *API) HashedPassword(c echo.Context, who string, err error) *model.Error {
	e := model.NewError()
	e.Code = AuthError
	e.Err = err
	e.Who = who
	//e.Data = ""
	e.StatusHTTP = http.StatusBadRequest
	e.APIMessage = "wrong user or password...."
	//e.UserID = a.parseUserID(c)

	//log.Warnf("log.Warnf - %s", e.Error())

	return &e
}

/* func (a *API) parseUserID(c echo.Context) (userID string) {

	userIDuuid, ok := c.Get("userID").(uuid.UUID)
	// Only to avoid the panic error
	if !ok {
		log.Errorf("cannot get/parse uuid from userID")
	}
	userID = userIDuuid.String()

	return userID
} */
