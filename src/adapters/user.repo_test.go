package adapters

import (
	"testing"

	model "github.com/ArquitecturaDeSistemas/usermicroservice/dominio"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockDB struct {
	mock.Mock
}

func TestCrearUsuario(t *testing.T) {
	mockDB := new(MockDB)
	ur := &userRepository{db: mockDB}
	input := model.CrearUsuarioInput{
		Nombre:     "Test",
		Apellido:   "User",
		Correo:     "test@example.com",
		Contrasena: "password",
	}
	mockDB.On("CrearUsuario", input).Return(&model.Usuario{}, nil)
	_, err := ur.CrearUsuario(input)
	assert.NoError(t, err)
	mockDB.AssertExpectations(t)
}
