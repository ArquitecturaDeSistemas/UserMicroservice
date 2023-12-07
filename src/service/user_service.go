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

func (s *UserService) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	u, err := s.repo.Usuario(req.GetId())
	if err != nil {
		return nil, err
	}
	response := &pb.GetUserResponse{
		Id:       u.ID,
		Nombre:   u.Nombre,
		Apellido: u.Apellido,
		Correo:   u.Correo,
	}
	return response, nil
}
func (s *UserService) ListUsers(ctx context.Context, req *pb.ListUsersRequest) (*pb.ListUsersResponse, error) {
	users, err := s.repo.Usuarios()
	if err != nil {
		return nil, err
	}
	var response []*pb.User
	for _, u := range users {
		user := &pb.User{
			Id:       u.ID,
			Nombre:   u.Nombre,
			Apellido: u.Apellido,
			Correo:   u.Correo,
		}
		response = append(response, user)
	}

	return &pb.ListUsersResponse{Users: response}, nil
}
