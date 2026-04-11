package services

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/marlonlyb/portfolioforge/domain/ports/user"
	"github.com/marlonlyb/portfolioforge/model"
)

// servicio, damain o usecase
type User struct {
	Repository user.Repository
}

// En go no hay constructores, asi que usamos esta función
func NewUser(ur user.Repository) *User {
	return &User{Repository: ur}
}

/*
Ojo aquí: aunque model.User es el dato de entrada, estamos
estableciendo que el usuario solo nos entregará los datos de
email, password y details, los demos los vamos a generar con
esta función Create...
*/

func (u *User) Create(m *model.User) error {

	//Creamos un ID
	ID, err := uuid.NewUUID()
	if err != nil {
		return fmt.Errorf("%s %w", "uuid.NewUUID()", err)
	}
	m.ID = ID

	if m.Email == "" {
		return fmt.Errorf("%s", "email is empty!")
	}

	if m.Password == "" {
		return fmt.Errorf("%s", "password is empty!")
	}

	//Encriptamos el password
	password, err := bcrypt.GenerateFromPassword([]byte(m.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("%s %w", "bcrypt.GenerateFromPassword()", err)
	}

	m.Password = string(password)

	//IsAdmin is False by default

	//recuerda que Details es tipo json
	if m.Details == nil {
		m.Details = []byte("{}")
	}

	m.CreatedAt = time.Now().Unix()

	err = u.Repository.Create(m)
	if err != nil {
		return fmt.Errorf("%s %w", "Repository.Create()", err)
	}

	//Limpiamos el password antes de entregar la información
	m.Password = ""
	return nil
}

func (u *User) GetByID(ID uuid.UUID) (model.User, error) {
	user, err := u.Repository.GetByID(ID)
	if err != nil {
		return model.User{}, fmt.Errorf("%s %w", "Repository.GetByID()", err)
	}

	user.Password = ""
	return user, nil
}

func (u *User) GetByEmail(email string) (model.User, error) {
	user, err := u.Repository.GetByEmail(email)
	if err != nil {
		return model.User{}, fmt.Errorf("%s %w", "Repository.GetByEmail()", err)
	}
	return user, nil
}

func (u *User) GetAll() (model.Users, error) {
	users, err := u.Repository.GetAll()
	if err != nil {
		return model.Users{}, fmt.Errorf("%s %w", "Repository.GetAll()", err)
	}
	return users, nil
}

func (u *User) Login(email, password string) (model.User, error) {
	m, err := u.GetByEmail(email)
	if err != nil {
		return model.User{}, fmt.Errorf("%s %w", "GetByEmail()", err)
	}

	//aquí comparo los passwords, pero no que sean iguales, comparo sus
	//comportamientos de cambio ya que estamos usando bcrypt
	err = bcrypt.CompareHashAndPassword([]byte(m.Password), []byte(password))
	if err != nil {
		return model.User{}, fmt.Errorf("%s %w", "CompareHashAndPassword()", err)
	}

	m.Password = ""

	return m, nil
}
