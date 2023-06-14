package users

import (
	"context"
	"fmt"

	"google.golang.org/protobuf/types/known/emptypb"

	usersproto "github.com/ANMalko/grpc-server.git/proto/users"
)

type Server struct {
	usersproto.UnimplementedUserServiceServer
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) GetUser(ctx context.Context, req *usersproto.UserId) (*usersproto.User, error) {
	fmt.Println("GetUser")
	fmt.Printf("%v", req)

	response := usersproto.User{
		Id: req.Id,
		Name: "User Name",
		Email: "user_name@mail.ru",
		PhoneNumber: "+7-952-201-86-24",
    }

	return &response, nil
}

func (s *Server) CreateUser(ctx context.Context, req *usersproto.UserCreate) (*usersproto.User, error) {
	fmt.Println("CreateUser")
	fmt.Printf("%v", req)

	response := usersproto.User{
		Id: 123,
		Name: req.Name,
		Email: req.Email,
		PhoneNumber: req.PhoneNumber,
    }

	return &response, nil
}

func (s *Server) UpdateUser(ctx context.Context, req *usersproto.User) (*usersproto.User, error) {
	fmt.Println("UpdateUser")
	fmt.Printf("%v", req)

	response := usersproto.User{
		Id: req.Id,
		Name: req.Name,
		Email: req.Email,
		PhoneNumber: req.PhoneNumber,
    }

	return &response, nil
}

func (s *Server) DeleteUser(ctx context.Context, _ *emptypb.Empty) (*usersproto.User, error) {
	fmt.Println("DeleteUser")
	response := usersproto.User{
		Id: 333,
		Name: "Deleted user name",
		Email: "Deleted user email",
		PhoneNumber: "Deleted user phone",
    }

	return &response, nil
}
