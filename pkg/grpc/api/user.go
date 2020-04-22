package api

import (
	"context"
	"errors"
	"time"

	"github.com/douglaszuqueto/golang-grpc/pkg/storage"
	"github.com/douglaszuqueto/golang-grpc/proto"
	"github.com/golang/protobuf/ptypes"

	"google.golang.org/grpc"
)

// UserService UserService
type UserService struct{}

// NewUserService NewUserService
func NewUserService(s *grpc.Server) *UserService {
	server := &UserService{}

	if s != nil {
		proto.RegisterUserServiceServer(s, server)
	}

	return server
}

// Get Get
func (s UserService) Get(ctx context.Context, req *proto.GetUserRequest) (*proto.GetUserResponse, error) {
	user, err := storage.GetUser(req.Id)
	if err != nil {
		return nil, err
	}

	protoUser, _ := userToProtoStruct(user)

	resp := &proto.GetUserResponse{
		User: &protoUser,
	}

	return resp, nil
}

// List List
func (s UserService) List(ctx context.Context, req *proto.ListUserRequest) (*proto.ListUserResponse, error) {
	l := []*proto.User{}

	users, err := storage.ListUser()
	if err != nil {
		return nil, err
	}

	for _, u := range users {
		user, _ := userToProtoStruct(u)

		l = append(l, &user)
	}

	resp := &proto.ListUserResponse{
		User: l,
	}

	return resp, nil
}

// Create Create
func (s UserService) Create(ctx context.Context, req *proto.CreateUserRequest) (*proto.CreateUserResponse, error) {
	if _, err := storage.GetUser(req.User.Id); err == nil {
		return nil, errors.New("user already exists")
	}

	user := storage.User{
		ID:        req.User.Id,
		Username:  req.User.Username,
		Email:     req.User.Email,
		State:     req.User.State,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := storage.CreateUser(user); err != nil {
		return nil, err
	}

	resp := &proto.CreateUserResponse{
		Result: "ok",
	}

	return resp, nil
}

// Update Update
func (s UserService) Update(ctx context.Context, req *proto.UpdateUserRequest) (*proto.UpdateUserResponse, error) {
	userOld, err := storage.GetUser(req.User.Id)
	if err != nil {
		return nil, err
	}

	user := storage.User{
		ID:        req.User.Id,
		Username:  req.User.Username,
		Email:     req.User.Email,
		State:     req.User.State,
		CreatedAt: userOld.CreatedAt,
		UpdatedAt: time.Now(),
	}

	if err := storage.UpdateUser(user); err != nil {
		return nil, err
	}

	resp := &proto.UpdateUserResponse{
		Result: "ok",
	}

	return resp, nil
}

// Delete Delete
func (s UserService) Delete(ctx context.Context, req *proto.DeleteUserRequest) (*proto.DeleteUserResponse, error) {
	if _, err := storage.GetUser(req.Id); err != nil {
		return nil, err
	}

	if err := storage.DeleteUser(req.Id); err != nil {
		return nil, err
	}

	resp := &proto.DeleteUserResponse{
		Result: "ok",
	}

	return resp, nil
}

func userToProtoStruct(u storage.User) (proto.User, error) {
	user := proto.User{
		Id:       u.ID,
		Username: u.Username,
		Email:    u.Email,
		State:    u.State,
	}

	var err error

	user.CreatedAt, err = ptypes.TimestampProto(u.CreatedAt)
	if err != nil {
		return user, err
	}

	user.UpdatedAt, err = ptypes.TimestampProto(u.UpdatedAt)
	if err != nil {
		return user, err
	}

	return user, nil
}
