package server

import (
	"context"
	"fmt"
	"strings"

	"github.com/chudik63/netevent/auth/internal/db/postgres/models"
	"github.com/chudik63/netevent/auth/internal/db/postgres/repository"
	"github.com/chudik63/netevent/auth/internal/token"
	pb "github.com/chudik63/netevent/auth/pkg/proto"
)

type Auth struct {
	pb.UnimplementedAuthServiceServer
	repo *repository.UserRepository
}

func (a *Auth) Register(ctx context.Context, in *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	us := in.GetUser()
	mod := &models.User{
		Id:        us.Id,
		Name:      us.Name,
		Password:  us.Password,
		Email:     us.Email,
		Role:      us.Role,
		Interests: strings.Join(us.Interests, " "),
	}
	err := a.repo.NewUser(mod)
	if err != nil {
		return &pb.RegisterResponse{Message: "err in db: " + err.Error()}, err
	}

	err = sendToEvent(&models.Participant{
		UserId:    us.Id,
		Name:      us.Name,
		Interests: us.Interests,
		Email:     us.Email,
	})
	if err != nil {
		return &pb.RegisterResponse{Message: "err to send event service: " + err.Error()}, err
	}
	return &pb.RegisterResponse{Message: "OK"}, nil
}

func (a *Auth) Authenticate(ctx context.Context, in *pb.AuthenticateRequest) (*pb.AuthenticateResponse, error) {
	name := in.GetName()
	pass := in.GetPassword()

	SmallToken, err := token.NewToken(name)
	if err != nil {
		return nil, err
	}
	LongToken, err := token.RefreshToken(name)
	if err != nil {
		return nil, err
	}

	tkn := &models.Token{
		AccessTkn:  SmallToken,
		AccessTtl:  int64(token.Small),
		RefreshTkn: LongToken,
		RefreshTtl: int64(token.Long),
	}
	err = a.repo.AuthUser(name, pass, tkn)
	if err != nil {
		return nil, err
	}

	return &pb.AuthenticateResponse{Tokens: &pb.Token{
		AccessToken:     SmallToken,
		AccessTokenTtl:  int64(token.Small),
		RefreshToken:    LongToken,
		RefreshTokenTtl: int64(token.Long),
	}}, nil
}

func (a *Auth) Authorise(ctx context.Context, in *pb.AuthoriseRequest) (*pb.AuthoriseResponse, error) {
	auToken := in.GetToken()
	flag, err := token.ValidToken(auToken)
	if err != nil {
		return &pb.AuthoriseResponse{Role: "", Message: string(fmt.Sprintf("%v", err))}, err
	}
	if !flag {
		return &pb.AuthoriseResponse{Role: "", Message: "error token not valid"}, nil
	}

	name, err := token.GetNameToken(auToken)
	if err != nil {
		return &pb.AuthoriseResponse{Role: "", Message: "error token not valid"}, nil
	}

	role, err := a.repo.GetRole(name)
	if err != nil {
		return &pb.AuthoriseResponse{Role: "", Message: "error not found in db on name: " + name + " " + err.Error()}, nil
	}
	return &pb.AuthoriseResponse{
		Role:    role,
		Message: "OK",
	}, nil
}
