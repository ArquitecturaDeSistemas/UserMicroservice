package service

import (
	"context"
	"fmt"

	"github.com/ArquitecturaDeSistemas/usermicroservice/model"
	pb "github.com/ArquitecturaDeSistemas/usermicroservice/proto"
	"github.com/ArquitecturaDeSistemas/usermicroservice/repository"
)

type UserService struct {
	pb.UnimplementedUserServiceServer
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {

	crearUsuarioInput := model.CrearUsuarioInput{
		Nombre:     req.GetNombre(),
		Apellido:   req.GetApellido(),
		Correo:     req.GetCorreo(),
		Contrasena: req.GetContrasena(),
	}
	u, err := s.repo.CrearUsuario(crearUsuarioInput)
	if err != nil {
		return nil, err
	}
	fmt.Printf("Usuario creado: %v", u)
	response := &pb.CreateUserResponse{
		Id: u.ID,
	}
	return response, nil
}

// func (s *UserService) Usuarios(ctx context.Context, req *pb.UsuariosRequest) (*pb.UsuariosResponse, error) {
// 	usuarios, err := s.repo.Usuarios()
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &pb.UsuariosResponse{Users: usuarios}, nil
// }

// ... otras funciones gRPC ...
