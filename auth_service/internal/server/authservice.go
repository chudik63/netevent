package server

import (
	"context"
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
	repo *repository.UserRepository
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
		return &pb.RegisterResponse{}, status.Errorf(codes.Internal, err.Error())
	}

	err = sendToEvent(&models.Participant{
		UserId:    id,
		Name:      us.Name,
		Interests: us.Interests,
		Email:     us.Email,
	})
	if err != nil {
		return &pb.RegisterResponse{}, status.Errorf(codes.Internal, err.Error())
	}

	return &pb.RegisterResponse{}, status.New(codes.OK, "Success").Err()
}

func (a *Auth) Authenticate(ctx context.Context, in *pb.AuthenticateRequest) (*pb.AuthenticateResponse, error) {
	name := in.GetName()
	pass := in.GetPassword()

	user, err := a.repo.AuthUser(name, pass)
	if err != nil {
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
	}}, status.New(codes.OK, "Success").Err()
}

func (a *Auth) Authorise(ctx context.Context, in *pb.AuthoriseRequest) (*pb.AuthoriseResponse, error) {
	auToken := in.GetToken()

	flag, err := token.ValidToken(auToken)
	if err != nil {
		return &pb.AuthoriseResponse{Role: ""}, status.Errorf(codes.Internal, err.Error())
	}
	if !flag {
		return &pb.AuthoriseResponse{Role: ""}, nil
	}

	id, err := token.GetIdToken(auToken)
	if err != nil {
		return &pb.AuthoriseResponse{Role: ""}, status.Errorf(codes.Internal, err.Error())
	}

	role, err := token.GetRoleToken(auToken)
	if err != nil {
		return &pb.AuthoriseResponse{Role: ""}, status.Errorf(codes.Internal, err.Error())
	}

	return &pb.AuthoriseResponse{
		Id:   id,
		Role: role,
	}, status.New(codes.OK, "Success").Err()
}
