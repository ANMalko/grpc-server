package users

import (
	"context"
	"fmt"

	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	usersproto "github.com/ANMalko/grpc-server.git/proto/users"
	"github.com/ANMalko/grpc-server.git/db/model"
	"github.com/ANMalko/grpc-server.git/db/error"
)

type UsersDAO interface {
	GetUser(ctx context.Context, userID uint32) *model.User
	CreateUser(ctx context.Context, user *model.User) (*model.User, error)
	DeleteUser(ctx context.Context, userID uint32)
	UpdateUser(ctx context.Context, user *model.User) error
}

type Server struct {
	dao UsersDAO
	usersproto.UnimplementedUserServiceServer
}

func NewServer(dao UsersDAO) *Server {
	return &Server{dao: dao}
}

func (s *Server) GetUser(ctx context.Context, req *usersproto.UserId) (*usersproto.User, error) {
	user := s.dao.GetUser(ctx, req.Id)

	if user == nil {
		return nil, status.Error(codes.NotFound, fmt.Sprint(dberror.ENOTFOUND))
	}

	response := usersproto.User{
		Id: user.Id,
		Name: user.Name,
		Email: user.Email,
		PhoneNumber: user.PhoneNumber,
    }

	return &response, nil
}

func (s *Server) CreateUser(ctx context.Context, req *usersproto.User) (*usersproto.User, error) {
	user := model.User{
		Id: req.Id,
		Name: req.Name,
		Email: req.Email,
		PhoneNumber: req.PhoneNumber,
	}


	new_user, err := s.dao.CreateUser(ctx, &user,)

	if err != nil && err.Error() == fmt.Sprint(dberror.EALREADYEXISTS) {
		return nil, status.Error(codes.AlreadyExists, fmt.Sprint(dberror.EALREADYEXISTS))
	}

	response := usersproto.User{
		Id: new_user.Id,
		Name: new_user.Name,
		Email: new_user.Email,
		PhoneNumber: new_user.PhoneNumber,
    }

	return &response, nil
}

func (s *Server) UpdateUser(ctx context.Context, req *usersproto.User) (*emptypb.Empty, error) {
	user := model.User{
		Id: req.Id,
		Name: req.Name,
		Email: req.Email,
		PhoneNumber: req.PhoneNumber,
	}

	err := s.dao.UpdateUser(ctx, &user)

	if err != nil && err.Error() == fmt.Sprint(dberror.ENOTFOUND) {
		return nil, status.Error(codes.AlreadyExists, fmt.Sprint(dberror.ENOTFOUND))
	}

	return nil, nil
}

func (s *Server) DeleteUser(ctx context.Context, req *usersproto.UserId) (*emptypb.Empty, error) {
	s.dao.DeleteUser(ctx, req.Id)
	return nil, nil
}
