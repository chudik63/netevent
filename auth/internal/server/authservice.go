package server

import (
	"context"
	"fmt"
	"strings"

	"gitlab.crja72.ru/gospec/go9/netevent/auth/internal/db/postgres/models"
	"gitlab.crja72.ru/gospec/go9/netevent/auth/internal/db/postgres/repository"
	"gitlab.crja72.ru/gospec/go9/netevent/auth/internal/token"
	pb "gitlab.crja72.ru/gospec/go9/netevent/auth/pkg/proto"
)

type Auth struct {
	pb.UnimplementedAuthServiceServer
	repo *repository.UserRepository
}

/*message RegisterRequest {
  User user = 1;
}
message RegisterResponse {
  string message = 1; // "OK" or error message
}*/

func (a *Auth) Register(ctx context.Context, in *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	us := in.GetUser()
	mod := &models.User{
		Id:        us.Id,
		Name:      us.Name,
		Password:  us.Password,
		Email:     us.Email,
		Interests: strings.Join(us.Interests, " "),
	}
	err := a.repo.NewUser(mod)
	if err != nil {
		return nil, err
	}
	return &pb.RegisterResponse{Message: "OK"}, nil
}

/*
	message AuthenticateRequest {
	  string name = 1;
	  string password = 2;
	}

	message AuthenticateResponse {
	  Token tokens = 1;
	}
*/
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
	flag, err := token.ValidTocken(auToken)
	if err != nil {
		return &pb.AuthoriseResponse{UserId: int64(-1), Message: string(fmt.Sprintf("%v", err))}, err
	}
	if !flag {
		return &pb.AuthoriseResponse{UserId: int64(-1), Message: "error token not valid"}, nil
	}

	name, err := token.GetNameToken(auToken)
	if err != nil {
		return &pb.AuthoriseResponse{UserId: int64(-1), Message: "error token not valid"}, nil
	}

	id, err := a.repo.GetId(name)
	if err != nil {
		return &pb.AuthoriseResponse{UserId: int64(-1), Message: "error token not valid"}, nil
	}
	return &pb.AuthoriseResponse{
		UserId:  int64(id),
		Message: "OK",
	}, nil
}

func (a *Auth) GetInterests(ctx context.Context, in *pb.GetInterestsRequest) (*pb.GetInterestsResponse, error) {
	interests, err := a.repo.GetInterests(int(in.GetUserId()))
	if err != nil {
		return &pb.GetInterestsResponse{Interests: []string{""}}, err
	}

	return &pb.GetInterestsResponse{Interests: strings.Split(interests, " ")}, nil
}
