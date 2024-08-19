package gapi

import (
	"context"
	"errors"

	db "github.com/leedrum/simplebank/db/sqlc"
	"github.com/leedrum/simplebank/pb"
	"github.com/leedrum/simplebank/util"
	"github.com/leedrum/simplebank/val"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *Server) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	violations := validateLoginUserRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	user, err := server.store.GetUser(ctx, req.GetUsername())
	if err != nil {
		if errors.Is(err, db.ErrorRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "user not found")
		}

		return nil, status.Errorf(codes.Internal, "cannot get user: %v", err)
	}

	if err := util.CheckPassword(req.GetPassword(), user.HashedPassword); err != nil {
		return nil, status.Errorf(codes.NotFound, "invalid password")
	}

	accessToken, accessPayload, err := server.tokenMaker.CreateToken(user.Username, server.config.AccessTokenDuration)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot create access token: %v", err)
	}

	refreshToken, refreshPayload, err := server.tokenMaker.CreateToken(
		user.Username,
		server.config.RefreshTokenDuration,
	)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot create refresh token: %v", err)
	}
	metaData := server.extractMetadata(ctx)
	session, err := server.store.CreateSession(ctx, db.CreateSessionParams{
		ID:           refreshPayload.ID,
		Username:     user.Username,
		RefreshToken: refreshToken,
		ExpiresAt:    refreshPayload.ExpiredAt,
		ClientIp:     metaData.ClientIP,
		UserAgent:    metaData.UserAgent,
		IsBlocked:    false,
	})

	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot create session: %v", err)
	}

	rsp := &pb.LoginUserResponse{
		AccessToken:           accessToken,
		RefreshToken:          refreshToken,
		User:                  convertUserToResponse(user),
		SessionId:             session.ID.String(),
		AccessTokenExpiresAt:  timestamppb.New(accessPayload.ExpiredAt),
		RefreshTokenExpiresAt: timestamppb.New(refreshPayload.ExpiredAt),
	}

	return rsp, nil
}

func validateLoginUserRequest(req *pb.LoginUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := val.ValidateUsername(req.GetUsername()); err != nil {
		violations = append(violations, fieldViolation("username", err))
	}

	if err := val.ValidatePassword(req.GetPassword()); err != nil {
		violations = append(violations, fieldViolation("password", err))
	}

	return violations
}
