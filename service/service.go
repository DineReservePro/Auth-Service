package service

import (
	pb "auth-service/generated/auth_service"
	"auth-service/storage/postgres"
	"context"
)

type AuthService struct {
	pb.UnimplementedAuthServiceServer
	User *postgres.UserRepo
}

func NewAuthService(user *postgres.UserRepo) *AuthService {
	return &AuthService{User: user}
}

func (a *AuthService) GetUserProfile(ctx context.Context, in *pb.GetUserProfileRequest) (*pb.GetUserProfileResponse, error) {
	resp, err := a.User.GetUserProfile(in.Username)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (a *AuthService) UpdateUserProfile(ctx context.Context, in *pb.UpdateUserProfileRequest) (*pb.UpdateUserProfileResponse, error) {
	resp, err := a.User.UpdateUserProfile(in)

	if err != nil {
		return &pb.UpdateUserProfileResponse{
			Message: resp.Message,
		}, err
	}

	return &pb.UpdateUserProfileResponse{
		Message: resp.Message,
	}, nil
}

func (a *AuthService) LogoutUser(ctx context.Context, in *pb.LogoutRequest) (*pb.LogoutResponse, error) {
	resp, err := a.User.LogoutUser(in.UserId)

	if err != nil {
		return &pb.LogoutResponse{
			Message: resp.Message,
		}, err
	}

	return &pb.LogoutResponse{
		Message: resp.Message,
	}, nil
}
