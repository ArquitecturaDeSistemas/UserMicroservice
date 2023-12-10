package ports

import (
	model "github.com/ArquitecturaDeSistemas/usermicroservice/dominio"
)

// puerto de salida
type UserRepository interface {
	CrearUsuario(input model.CrearUsuarioInput) (*model.Usuario, error)
	Usuario(id string) (*model.Usuario, error)
	ActualizarUsuario(id string, input *model.ActualizarUsuarioInput) (*model.Usuario, error)
	EliminarUsuario(id string) (*model.RespuestaEliminacion, error)
	Usuarios() ([]*model.Usuario, error)
	ExistePorCorreo(correo string) (bool, error)
	Retrieve(correo string, contrasena string) (*model.Usuario, error)
	Login(input model.LoginInput) (*model.AuthPayload, error)
	Logout(id string) (model.RespuestaEliminacion, error)
}
