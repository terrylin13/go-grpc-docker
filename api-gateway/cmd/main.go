package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/terrylin13/go-grpc-docker/api-gateway/pkg/auth"
	"github.com/terrylin13/go-grpc-docker/api-gateway/pkg/config"
	"github.com/terrylin13/go-grpc-docker/api-gateway/pkg/order"
	"github.com/terrylin13/go-grpc-docker/api-gateway/pkg/product"
)

func main() {
	c, err := config.LoadConfig()
	if err != nil {
		log.Fatalln("Failed at config", err)
	}

	r := gin.Default()

	authSvc := auth.RegisterRoutes(r, &c)
	order.RegisterRoutes(r, &c, authSvc)
	product.RegisterRoutes(r, &c, authSvc)
	r.Run(c.Port)
}
