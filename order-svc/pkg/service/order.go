package service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/terrylin13/go-grpc-docker/order-svc/pkg/client"
	"github.com/terrylin13/go-grpc-docker/order-svc/pkg/db"
	"github.com/terrylin13/go-grpc-docker/order-svc/pkg/models"
	"github.com/terrylin13/go-grpc-docker/order-svc/pkg/pb"
)

type Server struct {
	H          db.Handler
	ProductSvc client.ProductServiceClient
	pb.UnimplementedOrderServiceServer
}

func (s *Server) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (resp *pb.CreateOrderResponse, err error) {
	product, err := s.ProductSvc.FindOne(req.ProductId)
	fmt.Println(product.Data.Stock)
	if err != nil {
		resp = &pb.CreateOrderResponse{Status: http.StatusBadRequest, Error: err.Error()}
		return
	} else if product.Status >= http.StatusNotFound {
		resp = &pb.CreateOrderResponse{Status: product.Status, Error: product.Error}
		return
	} else if product.Data.Stock < req.Quantity {
		resp = &pb.CreateOrderResponse{Status: http.StatusConflict, Error: "Stock too less"}
		return
	}

	order := models.Order{
		Price:     product.Data.Price,
		ProductId: product.Data.Id,
		UserId:    req.UserId,
	}

	s.H.DB.Create(&order)

	res, err := s.ProductSvc.DecreaseStock(req.ProductId, order.Id)

	if err != nil {
		resp = &pb.CreateOrderResponse{Status: http.StatusBadRequest, Error: err.Error()}
		return
	} else if res.Status == http.StatusConflict {
		s.H.DB.Delete(&models.Order{}, order.Id)
		resp = &pb.CreateOrderResponse{Status: http.StatusConflict, Error: res.Error}
		return
	}

	resp = &pb.CreateOrderResponse{
		Status: http.StatusCreated,
		Id:     order.Id,
	}
	return
}
