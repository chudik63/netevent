package main

import (
	"context"
	"fmt"
	"log"
	"time"

	sq "github.com/Masterminds/squirrel"
	pb "gitlab.crja72.ru/gospec/go9/netevent/auth/pkg/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var addr = "localhost:5100"

func errFatal(function string, err error) {
	if err != nil {
		log.Fatalf("error fatal in %s: %s", function, err.Error())
	}
}

func main() {
	fmt.Println(time.Now().Unix())

	sql, args, err := sq.Select("id").
		From("tuser").Where(sq.Eq{"name": "anton"}).
		PlaceholderFormat(sq.Dollar).PlaceholderFormat(sq.Dollar).ToSql()
	fmt.Println(sql, args, err)

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %s", err.Error())
	}
	defer conn.Close()

	//создаётся один раз
	//важно что б ай ди шли по порядку иначе в бд ошибки будут
	//NewUser(conn)
	Authenticate(conn)
	Authorise(conn)

}

func NewUser(conn *grpc.ClientConn) {
	c := pb.NewAuthServiceClient(conn)
	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := c.Register(ctx, &pb.RegisterRequest{
		User: &pb.User{
			Id:        int64(2),
			Name:      "Anton04",
			Email:     "Nekt@gmail.ru",
			Password:  "adf123",
			Role:      "admin",
			Interests: []string{"adsf", "asdffff"},
		}})

	errFatal("Register", err)
	log.Printf("Register: %s", res.Message)

}

func Authenticate(conn *grpc.ClientConn) {
	c := pb.NewAuthServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := c.Authenticate(ctx, &pb.AuthenticateRequest{
		Name:     "Anton04",
		Password: "adf123",
	})
	errFatal("Authenticate", err)
	log.Printf("Authenticate: %s", res.Tokens)
}

func Authorise(conn *grpc.ClientConn) {
	c := pb.NewAuthServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	//строку с токеном надо заменить
	//мб ошибка если просто копировать
	//попробовать копировать как html
	res, err := c.Authorise(ctx, &pb.AuthoriseRequest{
		Token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MzQ0NTE0NDksImlhdCI6MTczNDQ1MTA0OSwic3ViIjoiQW50b24wNCJ9.Z180v2GcZuozERWsP8noP2otBzyJLOSp8pG2seEN1vY",
	})
	if err != nil {
		fmt.Println(err.Error())
	}
	errFatal("Authorise", err)

	log.Printf("Authorise: %s", res.Message, res.Role)
}

// func GetInterests(conn *grpc.ClientConn) {
// 	c := pb.NewAuthServiceClient(conn)
// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
// 	defer cancel()

// 	res, err := c.GetInterests(ctx, &pb.GetInterestsRequest{
// 		UserId: int64(4),
// 	})
// 	errFatal("GetInterests", err)
// 	log.Printf("GetInterests: %s", res.Interests)
// }
