package gapi

import (
	"context"

	db "github.com/leedrum/simplebank/db/sqlc"
	pb "github.com/leedrum/simplebank/pb"
	"github.com/leedrum/simplebank/util"
	"github.com/lib/pq"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	hashedPassword, err := util.HashPassword(req.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "hashing password failed: %v", err)
	}

	arg := db.CreateUserParams{
		Username:       req.GetUsername(),
		HashedPassword: hashedPassword,
		FullName:       req.GetFullName(),
		Email:          req.GetEmail(),
	}

	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code {
			case "unique_violation":
				return nil, status.Errorf(codes.AlreadyExists, "username or email already exists: %v", err)
			}
		}

		return nil, status.Errorf(codes.Internal, "cannot create user: %v", err)
	}

	rsp := &pb.CreateUserResponse{
		User: convertUserToResponse(user),
	}

	return rsp, nil
}
