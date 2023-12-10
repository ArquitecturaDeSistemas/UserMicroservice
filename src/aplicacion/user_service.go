package service

import (
	"context"
	"fmt"

	model "github.com/ArquitecturaDeSistemas/usermicroservice/dominio"
	repository "github.com/ArquitecturaDeSistemas/usermicroservice/ports"
	pb "github.com/ArquitecturaDeSistemas/usermicroservice/proto"
)

// este servicio implementa la interfaz UserServiceServer
// que se genera a partir del archivo proto
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
		Id:       u.ID,
		Nombre:   u.Nombre,
		Apellido: u.Apellido,
		Correo:   u.Correo,
	}
	fmt.Printf("Usuario creado: %v", response)
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

func (s *UserService) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	nombre := req.GetNombre()
	apellido := req.GetApellido()
	correo := req.GetCorreo()
	fmt.Printf("Nombre: %v", nombre)
	actualizarUsuarioInput := &model.ActualizarUsuarioInput{
		Nombre:   &nombre,
		Apellido: &apellido,
		Correo:   &correo,
	}
	fmt.Printf("Usuario actualizado input: %v", actualizarUsuarioInput)
	u, err := s.repo.ActualizarUsuario(req.GetId(), actualizarUsuarioInput)
	if err != nil {
		return nil, err
	}
	fmt.Printf("Usuario actualizado: %v", u)
	response := &pb.UpdateUserResponse{
		Id:       u.ID,
		Nombre:   u.Nombre,
		Apellido: u.Apellido,
		Correo:   u.Correo,
	}
	return response, nil
}

func (s *UserService) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	respuesta, err := s.repo.EliminarUsuario(req.GetId())
	if err != nil {
		return nil, err
	}
	response := &pb.DeleteUserResponse{
		Mensaje: respuesta.Mensaje,
	}
	return response, nil
}

func (s *UserService) LoginUser(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	loginInput := model.LoginInput{
		Correo:     req.GetCorreo(),
		Contrasena: req.GetContrasena(),
	}
	authPayload, err := s.repo.Login(loginInput)
	if err != nil {
		return nil, err
	}
	response := &pb.LoginResponse{
		Token: authPayload.Token,
		User: &pb.User{
			Id:       authPayload.Usuario.ID,
			Nombre:   authPayload.Usuario.Nombre,
			Apellido: authPayload.Usuario.Apellido,
			Correo:   authPayload.Usuario.Correo,
		},
	}
	return response, nil
}

func (s *UserService) LogoutUser(ctx context.Context, req *pb.LogoutRequest) (*pb.LogoutResponse, error) {
	respuesta, err := s.repo.Logout(req.GetUserID())
	if err != nil {
		return nil, err
	}
	response := &pb.LogoutResponse{
		Mensaje: respuesta.Mensaje,
	}
	return response, nil
}
