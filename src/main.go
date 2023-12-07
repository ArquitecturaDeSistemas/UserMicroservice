package main

import (
	"context"
	"log"

	"github.com/tam210/model"
	"github.com/tam210/repository"
	pb "github.com/tam210/usermicroservice/proto"
)

type server struct {
	pb.UnimplementedUserMicroServiceServer
	repo repository.UserRepository
}

func (s *server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	nuevoUsuario := model.CrearUsuarioInput{
		Nombre:     req.GetNombre(),
		Apellido:   req.GetApellido(),
		Correo:     req.GetCorreo(),
		Contrasena: req.GetContrasena(),
	}
	u, err := s.repo.CrearUsuario(nuevoUsuario)
	if err != nil {
		log.Printf("Error al crear usuario", err)
		return nil, err
	}
	response := &pb.CreateUserResponse{
		Id: u.ID,
	}
	return response, nil
}
