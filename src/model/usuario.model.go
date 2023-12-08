package model

import (
	"strconv"
)

// UsuarioGORM es el modelo de usuario para GORM de Usuario
type UsuarioGORM struct {
	ID         uint   `gorm:"primaryKey:autoIncrement" json:"id"`
	Nombre     string `gorm:"type:varchar(255);not null"`
	Apellido   string `gorm:"type:varchar(255);not null"`
	Correo     string `gorm:"type:varchar(255);not null;unique"`
	Contrasena string `gorm:"type:varchar(255);not null"`
}

// TableName especifica el nombre de la tabla para UsuarioGORM
func (UsuarioGORM) TableName() string {
	return "usuarios"
}

func (usuarioGORM *UsuarioGORM) ToGQL() (*Usuario, error) {

	return &Usuario{
		ID:         strconv.Itoa(int(usuarioGORM.ID)),
		Nombre:     usuarioGORM.Nombre,
		Apellido:   usuarioGORM.Apellido,
		Correo:     usuarioGORM.Correo,
		Contrasena: usuarioGORM.Contrasena,
	}, nil
}
