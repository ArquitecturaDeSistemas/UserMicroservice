package database

import (
	"log"

	"github.com/tam210/model"
	"gorm.io/gorm"
)

// EjecutarMigraciones realiza todas las migraciones necesarias en la base de datos.
func EjecutarMigraciones(db *gorm.DB) {

	db.AutoMigrate(&model.Usuario{})

	log.Println("Migraciones completadas")
}
