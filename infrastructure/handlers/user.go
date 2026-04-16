package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
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

type updateProfileRequest struct {
	FullName string `json:"full_name"`
	Company  string `json:"company"`
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

func (h *User) GetAll(c echo.Context) error {
	users, err := h.service.AdminList()
	if err != nil {
		return response.ContractError(500, "unexpected_error", "Unable to load users")
	}

	return c.JSON(response.ContractOK(map[string]interface{}{"items": users}))
}

func (h *User) AdminGetByID(c echo.Context) error {
	userID, err := parseUserIDParam(c.Param("id"))
	if err != nil {
		return response.ContractError(400, "validation_error", "The user identifier is invalid")
	}

	userData, err := h.service.AdminGetByID(userID)
	if err != nil {
		if errors.Is(err, model.ErrInvalidID) || strings.Contains(err.Error(), "no rows in result set") {
			return response.ContractError(404, "not_found", "User not found")
		}
		return response.ContractError(500, "unexpected_error", "Unable to load the user")
	}

	return c.JSON(response.ContractOK(userData))
}

func (h *User) Me(c echo.Context) error {
	userID, ok := c.Get("userID").(uuid.UUID)
	if !ok || userID == uuid.Nil {
		return response.ContractError(401, "authentication_required", "You must sign in to continue")
	}

	userData, err := h.service.GetByID(userID)
	if err != nil {
		if errors.Is(err, model.ErrInvalidID) {
			return response.ContractError(404, "not_found", "User not found")
		}
		return response.ContractError(404, "not_found", "User not found")
	}

	return c.JSON(response.ContractOK(h.service.ToStoreUser(userData)))
}

func (h *User) UpdateProfile(c echo.Context) error {
	userID, ok := c.Get("userID").(uuid.UUID)
	if !ok || userID == uuid.Nil {
		return response.ContractError(401, "authentication_required", "You must sign in to continue")
	}

	var request updateProfileRequest
	if err := c.Bind(&request); err != nil {
		return response.ContractError(400, "validation_error", "Invalid profile payload")
	}

	updatedUser, err := h.service.UpdateProfile(userID, request.FullName, request.Company)
	if err != nil {
		if strings.Contains(err.Error(), "profile fields are required") {
			return response.ContractError(400, "validation_error", "Full name and company are required")
		}
		return response.ContractError(500, "unexpected_error", "Unable to update the profile")
	}

	return c.JSON(response.ContractOK(map[string]interface{}{
		"user": h.service.ToStoreUser(updatedUser),
	}))
}

func (h *User) AdminUpdate(c echo.Context) error {
	userID, err := parseUserIDParam(c.Param("id"))
	if err != nil {
		return response.ContractError(400, "validation_error", "The user identifier is invalid")
	}

	request, err := decodeAdminUserUpdateRequest(c.Request().Body)
	if err != nil {
		return response.ContractError(400, "validation_error", "Invalid admin user payload")
	}
	if request.IsAdmin == nil {
		return response.ContractError(400, "validation_error", "The is_admin field is required", model.APIErrorDetail{Field: "is_admin", Issue: "required"})
	}

	updatedUser, err := h.service.AdminUpdate(userID, request)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrAdminUserUpdateScope):
			return response.ContractError(400, "validation_error", "Only is_admin can be updated")
		case errors.Is(err, model.ErrAdminUserProtected):
			return response.ContractError(403, "forbidden", "Existing admins cannot be edited from this flow")
		case strings.Contains(err.Error(), "no rows in result set"):
			return response.ContractError(404, "not_found", "User not found")
		default:
			return response.ContractError(500, "unexpected_error", "Unable to update the user")
		}
	}

	return c.JSON(response.ContractOK(updatedUser))
}

func (h *User) AdminDelete(c echo.Context) error {
	actorID, ok := c.Get("userID").(uuid.UUID)
	if !ok || actorID == uuid.Nil {
		return response.ContractError(http.StatusUnauthorized, "authentication_required", "You must sign in to continue")
	}

	targetID, err := parseUserIDParam(c.Param("id"))
	if err != nil {
		return response.ContractError(400, "validation_error", "The user identifier is invalid")
	}

	err = h.service.AdminSoftDelete(actorID, targetID)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrAdminSelfDelete):
			return response.ContractError(403, "forbidden", "You cannot delete your own admin account")
		case errors.Is(err, model.ErrAdminUserProtected):
			return response.ContractError(403, "forbidden", "Admin accounts cannot be deleted from this flow")
		case strings.Contains(err.Error(), "no rows in result set"):
			return response.ContractError(404, "not_found", "User not found")
		default:
			return response.ContractError(500, "unexpected_error", "Unable to delete the user")
		}
	}

	return c.NoContent(http.StatusNoContent)
}

func toStoreUser(userData model.User) model.StoreUser {
	storeUser := model.StoreUser{
		ID:                     userData.ID,
		Email:                  userData.Email,
		IsAdmin:                userData.IsAdmin,
		AuthProvider:           userData.AuthProvider,
		EmailVerified:          userData.EmailVerified,
		FullName:               userData.FullName,
		Company:                userData.Company,
		ProfileCompleted:       userData.ProfileCompleted,
		AssistantEligible:      userData.AssistantEligible,
		CanUseProjectAssistant: userData.CanUseProjectAssistant,
		CreatedAt:              time.Unix(userData.CreatedAt, 0).UTC(),
	}

	if userData.LastLoginAt > 0 {
		storeUser.LastLoginAt = time.Unix(userData.LastLoginAt, 0).UTC()
	}

	return storeUser
}

func parseUserIDParam(value string) (uuid.UUID, error) {
	parsedID, err := uuid.Parse(strings.TrimSpace(value))
	if err != nil {
		return uuid.Nil, model.ErrInvalidID
	}
	return parsedID, nil
}

func decodeAdminUserUpdateRequest(body io.ReadCloser) (model.AdminUserUpdateRequest, error) {
	defer body.Close()
	decoder := json.NewDecoder(body)
	decoder.DisallowUnknownFields()

	var request model.AdminUserUpdateRequest
	if err := decoder.Decode(&request); err != nil {
		return model.AdminUserUpdateRequest{}, err
	}
	if decoder.More() {
		return model.AdminUserUpdateRequest{}, errors.New("unexpected trailing payload")
	}

	return request, nil
}
