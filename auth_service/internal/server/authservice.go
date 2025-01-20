package server

import (
	"context"
	"errors"
	"strings"

	"github.com/chudik63/netevent/auth_service/internal/db/postgres/models"
	"github.com/chudik63/netevent/auth_service/internal/db/postgres/repository"
	"github.com/chudik63/netevent/auth_service/internal/token"
	pb "github.com/chudik63/netevent/auth_service/pkg/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Auth struct {
	pb.UnimplementedAuthServiceServer
	repo        *repository.UserRepository
	eventAdress string
}

func (a *Auth) Register(ctx context.Context, in *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	us := in.GetUser()

	mod := &models.User{
		Name:      us.Name,
		Password:  us.Password,
		Email:     us.Email,
		Role:      us.Role,
		Interests: strings.Join(us.Interests, " "),
	}

	id, err := a.repo.NewUser(mod)
	if err != nil {

		if errors.Is(err, models.ErrUserAlreadyExists) {
			return nil, status.Errorf(codes.InvalidArgument, err.Error())
		}

		return nil, status.Errorf(codes.Internal, err.Error())
	}

	err = sendToEvent(&models.Participant{
		UserId:    id,
		Name:      us.Name,
		Interests: us.Interests,
		Email:     us.Email,
	}, a.eventAdress)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &pb.RegisterResponse{}, nil
}

func (a *Auth) Authenticate(ctx context.Context, in *pb.AuthenticateRequest) (*pb.AuthenticateResponse, error) {
	name := in.GetName()
	pass := in.GetPassword()

	user, err := a.repo.AuthUser(name, pass)
	if err != nil {
		if errors.Is(err, models.ErrUserNotFound) {
			return nil, status.Errorf(codes.InvalidArgument, err.Error())
		}

		return nil, status.Errorf(codes.Internal, err.Error())
	}

	tokens, err := token.NewTokens(user.Id, user.Role)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &pb.AuthenticateResponse{Tokens: &pb.Token{
		AccessToken:     tokens.AccessTkn,
		AccessTokenTtl:  tokens.AccessTtl,
		RefreshToken:    tokens.RefreshTkn,
		RefreshTokenTtl: tokens.RefreshTtl,
	}}, nil
}

func (a *Auth) Authorise(ctx context.Context, in *pb.AuthoriseRequest) (*pb.AuthoriseResponse, error) {
	auToken := in.GetToken()

	valid, err := token.ValidToken(auToken)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	if !valid {
		return nil, status.Error(codes.PermissionDenied, "invalid token")
	}

	id, err := token.GetIdToken(auToken)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	role, err := token.GetRoleToken(auToken)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &pb.AuthoriseResponse{
		Id:   id,
		Role: role,
	}, nil
}
