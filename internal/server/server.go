package server

import (
    "context"

    pb "your_org/your_auth/api/auth/v1"
    "github.com/golang/protobuf/ptypes/empty"
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"

    "your_org/your_auth/internal/auth"
)

type AuthServer struct {
    pb.UnimplementedAuthServiceServer
    svc auth.Service
}

func NewAuthServer(svc auth.Service) *AuthServer {
    return &AuthServer{svc: svc}
}

func (s *AuthServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
    tok, exp, err := s.svc.Login(ctx, req.Email, req.Password)
    if err != nil {
        return nil, status.Errorf(codes.Unauthenticated, "login failed: %v", err)
    }
    return &pb.LoginResponse{
        Token:     tok,
        ExpiresAt: exp.Unix(),
    }, nil
}

func (s *AuthServer) ValidateToken(ctx context.Context, req *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error) {
    claims, err := s.svc.ValidateToken(ctx, req.Token)
    if err != nil {
        return nil, status.Errorf(codes.Unauthenticated, "invalid token: %v", err)
    }
    return &pb.ValidateTokenResponse{
        UserId:    claims.UserID,
        ExpiresAt: claims.ExpiresAt.Unix(),
    }, nil
}

func (s *AuthServer) Revoke(ctx context.Context, req *pb.RevokeRequest) (*empty.Empty, error) {
    if err := s.svc.Revoke(ctx, req.Token); err != nil {
        return nil, status.Errorf(codes.Internal, "revoke failed: %v", err)
    }
    return &empty.Empty{}, nil
}
