package server

import (
	"context"
	"fmt"

	"gitlab.crja72.ru/gospec/go9/netevent/auth_service/internal/token"
	pb "gitlab.crja72.ru/gospec/go9/netevent/auth_service/pkg/proto"
)

type Auth struct {
	pb.UnimplementedAuthServiceServer
}

/*
message RegisterRequest {
  User user = 1;
}
message RegisterResponse {
  string message = 1; // "OK" or error message
}
*/

func (a *Auth) Register(ctx context.Context, in *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	//conn db
	//err := db.Register(id, name ...)

	return nil, nil
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
	//conn db
	name := in.GetName()
	pass := in.GetPassword()
	err := a.repo.AuthUser(name, pass)
	if err != nil {
		return nil, err
	}
	SmallToken, err := token.NewToken(name)
	if err != nil {
		return nil, err
	}
	LongToken, err := token.RefreshToken(name)
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
	//conn db
	auToken := in.GetToken()
	flag, err := token.ValidTocken(auToken)
	if err != nil {
		return &pb.AuthoriseResponse{UserId: int64(-1), Message: string(fmt.Sprintf("%v", err))}, err
	}
	if !flag {
		return &pb.AuthoriseResponse{UserId: int64(-1), Message: "error token not valid"}, nil
	}

	_, err = token.GetNameToken(auToken)
	if err != nil {
		return &pb.AuthoriseResponse{UserId: int64(-1), Message: "error token not valid"}, nil
	}

	//db get id on name

	return &pb.AuthoriseResponse{
		UserId:  int64(1),
		Message: "OK",
	}, nil
}

func (a *Auth) GetInterests(ctx context.Context, in *pb.GetInterestsRequest) (*pb.GetInterestsResponse, error) {
	//conn db
	//db GetInterests
	return nil, nil
}
