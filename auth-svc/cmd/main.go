package main

import (
	"fmt"
	"log"
	"net"

	"github.com/terrylin13/go-grpc-docker/auth-svc/pkg/config"
	"github.com/terrylin13/go-grpc-docker/auth-svc/pkg/db"
	"github.com/terrylin13/go-grpc-docker/auth-svc/pkg/pb"
	"github.com/terrylin13/go-grpc-docker/auth-svc/pkg/services"
	"github.com/terrylin13/go-grpc-docker/auth-svc/pkg/utils"
	"google.golang.org/grpc"
)

func main() {
	c, err := config.LoadConfig()
	if err != nil {
		log.Fatalln("Failed at config", err)
	}

	h := db.Init(c.DBUrl)

	jwt := utils.JwtWrapper{
		SecretKey:       c.JWTSecretKey,
		Issuer:          "go-grpc-auth-svc",
		ExpirationHours: 24 * 365,
	}

	lis, err := net.Listen("tcp", c.Port)
	if err != nil {
		log.Fatalln("Failed to listing:", err)
	}
	fmt.Println("Auth Svc on", c.Port)

	s := &services.Server{
		H:   h,
		Jwt: jwt,
	}

	gSvc := grpc.NewServer()
	pb.RegisterAuthServiceServer(gSvc, s)

	if err := gSvc.Serve(lis); err != nil {
		log.Fatalln(err)
	}

}
