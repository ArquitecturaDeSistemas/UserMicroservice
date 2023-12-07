package user

import (
	"context"

	pb "github.com/ArquitecturaDeSistemas/usermicroservice/proto"
	"github.com/tam210/model"
	"github.com/tam210/repository"
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
	return &pb.CreateUserResponse{User: u}, nil
}

func (s *UserService) Usuarios(ctx context.Context, req *pb.UsuariosRequest) (*pb.UsuariosResponse, error) {
	usuarios, err := s.repo.Usuarios()
	if err != nil {
		return nil, err
	}
	return &pb.UsuariosResponse{Users: usuarios}, nil
}

// ... otras funciones gRPC ...
