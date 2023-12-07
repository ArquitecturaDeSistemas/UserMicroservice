package middleware

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/tam210/model"
)

// Clave secreta que no se expone! es una clvve
// del servidor
var jwtKey = []byte("clave_secreta")

// Estructura del token
type Claims struct {
	UserID string `json:"user_id"`
	jwt.StandardClaims
}

// Crea un token con el ID del usuario
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

func VerifyToken(tokenStr string) (string, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		return "ERROR", err
	}

	if !token.Valid {
		return "ERROR", errors.New("Token inválido")
	}

	return claims.UserID, nil
}
