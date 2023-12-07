package repository

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/ArquitecturaDeSistemas/usermicroservice/database"
	"github.com/ArquitecturaDeSistemas/usermicroservice/model"
	"github.com/dgrijalva/jwt-go"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserRepository interface {
	CrearUsuario(input model.CrearUsuarioInput) (*model.Usuario, error)
	// Usuario(id string) (string, error)
	// ActualizarUsuario(id string, input *model.ActualizarUsuarioInput) (string, error)
	// EliminarUsuario(id string) (string, error)
	Usuarios() ([]*model.Usuario, error)
	// ExistePorCorreo(correo string) (bool, error)
	// Retrieve(correo string, contrasena string) (*model.Usuario, error)
	// Login(input model.LoginInput) (string, error)
	// Logout(id string) (string, error)
}

type userRepository struct {
	db             *database.DB
	activeSessions map[string]string
}

func NewUserRepository(db *database.DB) UserRepository {
	return &userRepository{
		db:             db,
		activeSessions: make(map[string]string),
	}
}

func ToJSON(obj interface{}) (string, error) {
	jsonData, err := json.Marshal(obj)
	if err != nil {
		return "", err
	}
	return string(jsonData), err
}

// ExistePorCorreo verifica si existe un usuario con el correo proporcionado.
func (ur *userRepository) ExistePorCorreo(correo string) (bool, error) {
	var usuarioGORM model.UsuarioGORM
	result := ur.db.GetConn().Where("correo = ?", correo).First(&usuarioGORM)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return false, nil
		}
		log.Printf("Error al buscar el usuario con correo %s: %v", correo, result.Error)
		return false, result.Error
	}

	return true, result.Error
}

// Retrieve obtiene un usuario por su correo y contraseña.
// Retorna nil si no se encuentra el usuario.
func (ur *userRepository) Retrieve(correo string, contrasena string) (*model.Usuario, error) {
	var usuarioGORM model.UsuarioGORM
	fmt.Printf("correo: %s\n", correo)

	if err := ur.db.GetConn().Where("correo = ?", correo).First(&usuarioGORM).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("Usuario con correo %s no encontrado", correo)
		}
		return nil, fmt.Errorf("Error al buscar usuario: %v", err)
	}

	// Verificar la contraseña con bcrypt
	if err := bcrypt.CompareHashAndPassword([]byte(usuarioGORM.Contrasena), []byte(contrasena)); err != nil {
		// Contraseña incorrecta
		return nil, fmt.Errorf("Credenciales incorrectas")
	}
	return usuarioGORM.ToGQL()
}

// ObtenerTrabajo obtiene un trabajo por su ID.
func (ur *userRepository) Usuario(id string) (string, error) {
	var usuarioGORM model.UsuarioGORM
	//result := ur.db.GetConn().First(&usuarioGORM, id)
	result := ur.db.GetConn().Preload("Direcciones.Ciudad").First(&usuarioGORM, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return "No encontrado", result.Error
		}
		log.Printf("Error al obtener el usuario con ID %s: %v", id, result.Error)
		return "Error al obtener el usuario con ID", result.Error
	}

	response, _ := usuarioGORM.ToGQL()
	return ToJSON(response)
}

// Usuarios obtiene todos los usuarios de la base de datos.
func (ur *userRepository) Usuarios() ([]*model.Usuario, error) {
	var usuariosGORM []model.UsuarioGORM
	result := ur.db.GetConn().Preload("Direcciones.Ciudad").Find(&usuariosGORM)

	if result.Error != nil {
		log.Printf("Error al obtener los usuarios: %v", result.Error)
		return nil, result.Error
	}

	var usuarios []*model.Usuario
	for _, usuarioGORM := range usuariosGORM {
		usuario, _ := usuarioGORM.ToGQL()
		usuarios = append(usuarios, usuario)
	}

	// usuariosJSON, err := json.Marshal(usuarios)
	// if err != nil {
	// 	log.Printf("Error al convertir usuarios a JSON: %v", err)
	// 	return "[]", err
	// }
	// return ToJSON(usuarios)
	return usuarios, nil
}
func (ur *userRepository) CrearUsuario(input model.CrearUsuarioInput) (*model.Usuario, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Contrasena), bcrypt.DefaultCost)
	log.Printf("Hashed password: %s", string(hashedPassword))

	if err != nil {
		log.Printf("Error al crear el hash de la contraseña: %v", err)
		return nil, err
	}

	usuarioGORM :=
		&model.UsuarioGORM{
			Nombre:     input.Nombre,
			Apellido:   input.Apellido,
			Correo:     input.Correo,
			Contrasena: string(hashedPassword),
		}
	result := ur.db.GetConn().Create(&usuarioGORM)
	if result.Error != nil {
		log.Printf("Error al crear el usuario: %v", result.Error)
		return nil, result.Error
	}

	response, err := usuarioGORM.ToGQL()
	return response, err
}
func (ur *userRepository) ActualizarUsuario(id string, input *model.ActualizarUsuarioInput) (string, error) {
	// Buscar el usuario existente por su ID
	var usuarioGORM model.UsuarioGORM
	result := ur.db.GetConn().First(&usuarioGORM, id)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return "No se pudo actualizar el usuario", fmt.Errorf("Usuario con ID %s no encontrado", id)
		}
		log.Printf("Error al buscar el usuario con ID %s: %v", id, result.Error)
		return "No se pudo actualizar el usuario", result.Error
	}

	// Actualizar los campos proporcionados en el input
	if input.Nombre != nil {
		usuarioGORM.Nombre = *input.Nombre
	}
	if input.Apellido != nil {
		usuarioGORM.Apellido = *input.Apellido
	}
	if input.Correo != nil {
		usuarioGORM.Correo = *input.Correo
	}

	// Actualizar el registro en la base de datos
	result = ur.db.GetConn().Save(&usuarioGORM)
	if result.Error != nil {
		log.Printf("Error al actualizar el usuario con ID %s: %v", id, result.Error)
		return "No se puede actualizar el usuario con ID dado", result.Error
	}

	// Devolver el usuario actualizado
	response, _ := usuarioGORM.ToGQL()
	return ToJSON(response)
}

// EliminarUsuario elimina un usuario de la base de datos por su ID.
func (ur *userRepository) EliminarUsuario(id string) (string, error) {
	// Intenta buscar el usuario por su ID
	var usuarioGORM model.UsuarioGORM
	result := ur.db.GetConn().First(&usuarioGORM, id)

	if result.Error != nil {
		// Manejo de errores
		if result.Error == gorm.ErrRecordNotFound {
			// El usuario no se encontró en la base de datos
			response, _ := &model.RespuestaEliminacion{
				Mensaje:          "El usuario no existe",
				CodigoEstadoHTTP: http.StatusNotFound,
			}, result.Error
			return ToJSON(response)

		}
		log.Printf("Error al buscar el usuario con ID %s: %v", id, result.Error)
		response, _ := &model.RespuestaEliminacion{
			Mensaje:          "Error al buscar el usuario",
			CodigoEstadoHTTP: http.StatusInternalServerError,
		}, result.Error
		return ToJSON(response)
	}

	// Elimina el usuario de la base de datos
	result = ur.db.GetConn().Delete(&usuarioGORM, id)

	if result.Error != nil {
		log.Printf("Error al eliminar el usuario con ID %s: %v", id, result.Error)
		response, _ := &model.RespuestaEliminacion{
			Mensaje:          "Error al eliminar el usuario",
			CodigoEstadoHTTP: http.StatusInternalServerError,
		}, result.Error
		return ToJSON(response)
	}

	// Éxito al eliminar el usuario
	response, _ := &model.RespuestaEliminacion{
		Mensaje:          "Usuario eliminado con éxito",
		CodigoEstadoHTTP: http.StatusOK,
	}, result.Error
	return ToJSON(response)

}

func (ur *userRepository) Login(input model.LoginInput) (string, error) {
	// Verificar las credenciales del usuario (correo y contraseña)
	if input.Correo == "" || input.Contrasena == "" {
		return "", errors.New("Correo y contraseña son requeridos")
	}
	if len(input.Contrasena) < 6 || len(input.Contrasena) > 50 {
		return "", errors.New("La contraseña debe tener al menos 6 caracteres")
	}
	if len(input.Correo) < 3 || len(input.Correo) > 50 {
		return "", errors.New("El correo debe tener al menos 3 caracteres")
	}

	usuario, err := ur.Retrieve(input.Correo, input.Contrasena)
	if err != nil {
		fmt.Printf("Error al verificar las credenciales: %v", err)
		return "", errors.New("Credenciales inválidas")
	}

	// Comprueba si el usuario ya tiene una sesión activa (esto podría ser a través de una base de datos)
	if ur.isSessionActive(usuario.ID) {
		return "", errors.New("Ya existe una sesión activa")
	}
	// Generar un token de autenticación para el usuario
	token, err := CreateToken(usuario)
	if err != nil {
		fmt.Printf("Error al generar el token de autenticación: %v", err)
		return "", fmt.Errorf("Error al generar el token: %v", err)
	}

	ur.registerSession(usuario.ID, token)

	// Crear el objeto AuthPayload con el token y los datos del usuario
	authPayload := &model.AuthPayload{
		Token:   token,
		Usuario: usuario,
	}
	log.Printf("Usuario autenticado: %v", usuario.ID)
	return ToJSON(authPayload)

}

// func (ur *userRepository) Logout(userID string) (string, error) {

// 	if userID == "" {
// 		return "false", errors.New("El ID de usuario es requerido")
// 	}

// 	// Verifica si la sesión del usuario existe
// 	if !ur.isSessionActive(userID) {
// 		return "false", errors.New("No hay una sesión activa para este usuario")
// 	}

// 	// Elimina la sesión del usuario
// 	delete(activeSessions, userID)

// 	log.Printf("Sesión cerrada para el usuario: %v", userID)
// 	return "true", nil
// }

func (ur *userRepository) Logout(userID string) (string, error) {
	var respuesta model.RespuestaEliminacion
	if userID == "" {
		var re, _ = ToJSON(respuesta)
		return re, errors.New("El ID de usuario es requerido")
	}
	if !ur.isSessionActive(userID) {
		var re, _ = ToJSON(respuesta)
		return re, errors.New("No hay una sesión activa para este usuario")
	}
	delete(activeSessions, userID)
	log.Printf("Sesión cerrada para el usuario: %v", userID)
	respuesta = model.RespuestaEliminacion{
		Mensaje:          "Sesión cerrada exitosamente",
		CodigoEstadoHTTP: 200, // Código HTTP para éxito
	}
	var re, _ = ToJSON(respuesta)
	return re, nil
}

// Clave secreta que no se expone! es una clvve
// del servidor
var jwtKey = []byte("clave_secreta")

// Estructura del token
type Claims struct {
	UserID string `json:"user_id"`
	//Role   string `json:"role"`
	jwt.StandardClaims
}

func CreateToken(user *model.Usuario) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)

	claims := &Claims{
		UserID: user.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

// ExtraerInfoToken es una función que decodifica un token JWT y extrae los claims (afirmaciones) del mismo.
func ExtraerInfoToken(tokenStr string) (*Claims, error) {
	// jwt.ParseWithClaims intenta analizar el token JWT.
	// Se le pasa el token como string, una instancia de Claims para mapear los datos del token,
	// y una función de callback para validar el algoritmo de firma del token.
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Comprueba que el algoritmo de codificación del token sea el esperado.
		// En este caso, se espera que el algoritmo sea HMAC (jwt.SigningMethodHMAC).
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// Si el algoritmo es el esperado, se devuelve la clave secreta utilizada para firmar el token.
		return jwtKey, nil
	})

	// Si no hay errores y el token es válido, extrae los claims y los devuelve.
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	} else {
		// Si hay un error o el token no es válido, devuelve el error.
		return nil, err
	}
}

var activeSessions = make(map[string]string) // Mapa de ID de usuario a token

func (ur *userRepository) isSessionActive(userID string) bool {
	_, active := activeSessions[userID]
	return active
}
func (ur *userRepository) registerSession(userID, token string) {
	activeSessions[userID] = token
}

func (ur *userRepository) endSession(userID string) {
	delete(activeSessions, userID)
}
