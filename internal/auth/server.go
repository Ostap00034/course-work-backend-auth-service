package auth

import (
	"context"

	pb "github.com/Ostap00034/course-work-backend-api-specs/gen/go/auth/v1"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthServer struct {
	pb.UnimplementedAuthServiceServer
	svc Service
}

func NewAuthServer(svc Service) *AuthServer {
	return &AuthServer{svc: svc}
}

func (s *AuthServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	tok, exp, err := s.svc.Login(ctx, req.Email, req.Password)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	return &pb.LoginResponse{
		Token:     tok,
		ExpiresAt: exp.Unix(),
	}, nil
}

func (s *AuthServer) ValidateToken(ctx context.Context, req *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error) {
	claims, user, err := s.svc.ValidateToken(ctx, req.Token)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, err.Error())
	}

	return &pb.ValidateTokenResponse{
		UserId:    claims.UserID,
		User:      user,
		ExpiresAt: claims.ExpiresAt.Unix(),
	}, nil
}

func (s *AuthServer) Revoke(ctx context.Context, req *pb.RevokeRequest) (*empty.Empty, error) {
	if err := s.svc.Revoke(ctx, req.Token); err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &empty.Empty{}, nil
}
