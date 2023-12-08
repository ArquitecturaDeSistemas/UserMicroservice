package model

type ActualizarUsuarioInput struct {
	Nombre     *string `json:"nombre,omitempty"`
	Apellido   *string `json:"apellido,omitempty"`
	Correo     *string `json:"correo,omitempty"`
	Contrasena *string `json:"contrasena,omitempty"`
}

type AuthPayload struct {
	Token   string   `json:"token"`
	Usuario *Usuario `json:"usuario"`
}

type CrearUsuarioInput struct {
	Nombre     string `json:"nombre"`
	Apellido   string `json:"apellido"`
	Correo     string `json:"correo"`
	Contrasena string `json:"contrasena"`
}

type LoginInput struct {
	Correo     string `json:"correo"`
	Contrasena string `json:"contrasena"`
}

type RespuestaEliminacion struct {
	Mensaje string `json:"mensaje"`
}

type Usuario struct {
	ID         string `json:"id"`
	Nombre     string `json:"nombre"`
	Apellido   string `json:"apellido"`
	Correo     string `json:"correo"`
	Contrasena string `json:"contrasena"`
}

func (Usuario) IsEntity() {}
