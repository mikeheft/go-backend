package gapi

import (
	"context"
	"errors"
	"fmt"

	db "github.com/mikeheft/go-backend/db/sqlc"
	"github.com/mikeheft/go-backend/pb"
	"github.com/mikeheft/go-backend/util"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *Server) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	user, err := server.store.GetUser(ctx, req.GetUsername())
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "user with username (%s) not found: %s", user, err)
		}
		return nil, status.Errorf(codes.Internal, "failed to find user: %s", err)
	}

	err = util.CheckPassword(req.Password, user.HashedPassword)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "incorrect password: %s", err)
	}

	accessToken, accessPayload, err := server.tokenMaker.CreateToken(
		req.Username,
		server.config.AccessTokenDuration,
	)
	if err != nil {
		fmt.Println("Error in access create")
		return nil, status.Errorf(codes.Internal, "failed to create access token: %s", err)
	}

	refreshToken, refreshPayload, err := server.tokenMaker.CreateToken(
		user.Username,
		server.config.RefreshTokenDuration,
	)
	if err != nil {
		fmt.Println("Error in refresh create")
		return nil, status.Errorf(codes.Internal, "failed to create refresh token: %s", err)
	}

	metadata := server.extractMetadata(ctx)

	session, err := server.store.CreateSession(ctx, db.CreateSessionParams{
		ID:           refreshPayload.ID,
		Username:     user.Username,
		RefreshToken: refreshToken,
		ExpiresAt:    refreshPayload.ExpiredAt,
		IsBlocked:    false,
		UserAgent:    metadata.UserAgent,
		ClientIp:     metadata.ClientIP,
	})
	if err != nil {
		fmt.Println("Error in session create")
		return nil, status.Errorf(codes.Internal, "failed to create session: %s", err)
	}

	rsp := &pb.LoginUserResponse{
		SessionId:             session.ID.String(),
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  timestamppb.New(accessPayload.ExpiredAt),
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: timestamppb.New(refreshPayload.ExpiredAt),
		User:                  convertUser(user),
	}

	return rsp, nil
}
