package handlers

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/marlonlyb/portfolioforge/domain/ports/user"
	"github.com/marlonlyb/portfolioforge/infrastructure/handlers/response"
	"github.com/marlonlyb/portfolioforge/model"
)

/* date cuenta que a diferencia de model
aqui los tipos y funciones son privadas
porque serán consultadas desde la ruta */

type User struct {
	service   user.Service
	responser response.API
}

func NewUser(us user.Service) *User {
	return &User{service: us}
}

func (h *User) Create(c echo.Context) error {
	m := model.User{}

	//vinculamos (bind) la información del cuerpo de la solicitud
	err := c.Bind(&m)
	if err != nil {
		return h.responser.BindFailed(c, "handlers-User-Create-c.Bind(&m)", err)
		//return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	err = h.service.Create(&m)
	if err != nil {
		return h.responser.Error(c, "handlers-User-Create-h.service.Create((&m))", err)
		//return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(h.responser.Created(m))
	//return c.JSON(http.StatusOK, m)
}

func (h *User) Register(c echo.Context) error {
	m := model.User{}
	if err := c.Bind(&m); err != nil {
		return response.ContractError(400, "validation_error", "Los datos enviados no son válidos")
	}

	if err := h.service.Create(&m); err != nil {
		return response.ContractError(400, "validation_error", "No fue posible registrar el usuario")
	}

	return c.JSON(response.ContractCreated(map[string]interface{}{
		"user": toStoreUser(m),
	}))
}

func (h *User) GetAll(c echo.Context) error {
	users, err := h.service.GetAll()
	if err != nil {
		return h.responser.Error(c, "handlers-User-Create-h.service.GetAll()", err)
	}

	return c.JSON(h.responser.OK(users))
}

func (h *User) Me(c echo.Context) error {
	userID, ok := c.Get("userID").(uuid.UUID)
	if !ok || userID == uuid.Nil {
		return response.ContractError(401, "authentication_required", "Debes iniciar sesión para continuar")
	}

	userData, err := h.service.GetByID(userID)
	if err != nil {
		if errors.Is(err, model.ErrInvalidID) {
			return response.ContractError(404, "not_found", "Usuario no encontrado")
		}
		return response.ContractError(404, "not_found", "Usuario no encontrado")
	}

	return c.JSON(response.ContractOK(toStoreUser(userData)))
}

func toStoreUser(userData model.User) model.StoreUser {
	return model.StoreUser{
		ID:        userData.ID,
		Email:     userData.Email,
		IsAdmin:   userData.IsAdmin,
		CreatedAt: time.Unix(userData.CreatedAt, 0).UTC(),
	}
}
