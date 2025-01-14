package server

import (
	"context"
	"strings"

	"github.com/chudik63/netevent/auth/internal/db/postgres/models"
	"github.com/chudik63/netevent/auth/internal/db/postgres/repository"
	"github.com/chudik63/netevent/auth/internal/token"
	pb "github.com/chudik63/netevent/auth/pkg/proto"
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
		return &pb.RegisterResponse{Id: 0}, status.Errorf(codes.Internal, err.Error())
	}

	err = sendToEvent(&models.Participant{
		UserId:    id,
		Name:      us.Name,
		Interests: us.Interests,
		Email:     us.Email,
	})
	if err != nil {
		return &pb.RegisterResponse{Id: 0}, status.Errorf(codes.Internal, err.Error())
	}
	return &pb.RegisterResponse{Id: id}, status.New(codes.OK, "Success").Err()
}

func (a *Auth) Authenticate(ctx context.Context, in *pb.AuthenticateRequest) (*pb.AuthenticateResponse, error) {
	name := in.GetName()
	pass := in.GetPassword()

	SmallToken, err := token.NewToken(name)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	LongToken, err := token.RefreshToken(name)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	tkn := &models.Token{
		AccessTkn:  SmallToken,
		AccessTtl:  int64(token.Small),
		RefreshTkn: LongToken,
		RefreshTtl: int64(token.Long),
	}
	err = a.repo.AuthUser(name, pass, tkn)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &pb.AuthenticateResponse{Tokens: &pb.Token{
		AccessToken:     SmallToken,
		AccessTokenTtl:  int64(token.Small),
		RefreshToken:    LongToken,
		RefreshTokenTtl: int64(token.Long),
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

	name, err := token.GetNameToken(auToken)
	if err != nil {
		return &pb.AuthoriseResponse{Role: ""}, status.Errorf(codes.Internal, err.Error())
	}

	role, err := a.repo.GetRole(name)
	if err != nil {
		return &pb.AuthoriseResponse{Role: ""}, status.Errorf(codes.Internal, err.Error())
	}
	return &pb.AuthoriseResponse{
		Role: role,
	}, status.New(codes.OK, "Success").Err()
}
