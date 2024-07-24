package gapi

import (
	db "github.com/leedrum/simplebank/db/sqlc"
	pb "github.com/leedrum/simplebank/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func convertUserToResponse(user db.User) *pb.User {
	return &pb.User{
		Username:          user.Username,
		FullName:          user.FullName,
		Email:             user.Email,
		PasswordChangedAt: timestamppb.New(user.PasswordChangedAt),
		CreatedAt:         timestamppb.New(user.CreatedAt),
	}
}
