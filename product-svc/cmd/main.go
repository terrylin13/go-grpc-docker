package main

import (
	"fmt"
	"log"
	"net"

	"github.com/terrylin13/go-grpc-docker/product-svc/pkg/config"
	"github.com/terrylin13/go-grpc-docker/product-svc/pkg/db"
	"github.com/terrylin13/go-grpc-docker/product-svc/pkg/pb"
	"github.com/terrylin13/go-grpc-docker/product-svc/pkg/services"
	"google.golang.org/grpc"
)

func main() {
	c, err := config.LoadConfig()
	if err != nil {
		log.Fatalln("Failed at config", err)
	}

	h := db.Init(c.DBUrl)

	lis, err := net.Listen("tcp", c.Port)
	if err != nil {
		log.Fatalln("Failed to listing:", err)
	}

	fmt.Println("Product Svc on", c.Port)

	s := &services.Server{
		H: h,
	}

	gSvc := grpc.NewServer()

	pb.RegisterProductServiceServer(gSvc, s)

	if err := gSvc.Serve(lis); err != nil {
		log.Fatalln("Failed to serve:", err)
	}
}
