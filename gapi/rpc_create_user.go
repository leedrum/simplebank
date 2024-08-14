package gapi

import (
	"context"
	"time"

	"github.com/hibiken/asynq"
	db "github.com/leedrum/simplebank/db/sqlc"
	pb "github.com/leedrum/simplebank/pb"
	"github.com/leedrum/simplebank/util"
	"github.com/leedrum/simplebank/val"
	"github.com/leedrum/simplebank/worker"
	"github.com/lib/pq"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	violations := validateCreateUserRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	hashedPassword, err := util.HashPassword(req.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "hashing password failed: %v", err)
	}

	arg := db.CreateUserTxParams{
		CreateUserParams: db.CreateUserParams{
			Username:       req.GetUsername(),
			HashedPassword: hashedPassword,
			FullName:       req.GetFullName(),
			Email:          req.GetEmail(),
		},
		AfterCreate: func(user db.User) error {
			taskPayload := worker.PayloadSendVerifyEmail{
				Username: user.Username,
			}
			taskOpts := []asynq.Option{
				asynq.ProcessIn(10 * time.Second),
				asynq.MaxRetry(10),
				asynq.Queue(worker.QueueCritical),
			}
			return server.taskDistributor.DistributeTaskSendVerifyEmail(ctx, taskPayload, taskOpts...)
		},
	}

	txResult, err := server.store.CreateUserTx(ctx, arg)
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
		User: convertUserToResponse(txResult.User),
	}

	return rsp, nil
}

func validateCreateUserRequest(req *pb.CreateUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := val.ValidateUsername(req.GetUsername()); err != nil {
		violations = append(violations, fieldViolation("username", err))
	}

	if err := val.ValidateFullName(req.GetFullName()); err != nil {
		violations = append(violations, fieldViolation("full_name", err))
	}

	if err := val.ValidatePassword(req.GetPassword()); err != nil {
		violations = append(violations, fieldViolation("password", err))
	}

	if err := val.ValidateEmail(req.GetEmail()); err != nil {
		violations = append(violations, fieldViolation("email", err))
	}

	return violations
}
