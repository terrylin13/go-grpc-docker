package services

import (
	"context"
	"net/http"

	"github.com/terrylin13/go-grpc-docker/product-svc/pkg/db"
	"github.com/terrylin13/go-grpc-docker/product-svc/pkg/models"
	"github.com/terrylin13/go-grpc-docker/product-svc/pkg/pb"
)

type Server struct {
	H db.Handler
	pb.UnimplementedProductServiceServer
}

func (s *Server) CreateProduct(ctx context.Context, req *pb.CreateProductRequest) (resp *pb.CreateProductResponse, err error) {
	product := &models.Product{
		Name:  req.Name,
		Price: req.Price,
		Stock: req.Stock,
	}

	if res := s.H.DB.Create(&product); res.Error != nil {
		return &pb.CreateProductResponse{
			Status: http.StatusConflict,
			Error:  res.Error.Error(),
		}, nil
	}

	resp = &pb.CreateProductResponse{
		Status: http.StatusCreated,
		Id:     product.Id,
	}
	return
}

func (s *Server) FindOne(ctx context.Context, req *pb.FindOneRequest) (resp *pb.FindOneResponse, err error) {
	var product models.Product

	if res := s.H.DB.First(&product, req.Id); res.Error != nil {
		return &pb.FindOneResponse{
			Status: http.StatusNotFound,
			Error:  res.Error.Error(),
		}, nil
	}

	data := &pb.FindOneData{
		Id:    product.Id,
		Name:  product.Name,
		Stock: product.Stock,
		Price: product.Price,
	}

	resp = &pb.FindOneResponse{
		Status: http.StatusOK,
		Data:   data,
	}

	return
}

func (s *Server) DecreaseStock(ctx context.Context, req *pb.DecreaseStockRequest) (resp *pb.DecreaseStockResponse, err error) {
	var product models.Product

	if res := s.H.DB.First(&product, req.Id); res.Error != nil {
		resp = &pb.DecreaseStockResponse{
			Status: http.StatusNotFound,
			Error:  res.Error.Error(),
		}
		return
	}

	if product.Stock <= 0 {
		resp = &pb.DecreaseStockResponse{
			Status: http.StatusConflict,
			Error:  "Stock too low",
		}
		return
	}

	var log models.StockDecreaseLog

	if res := s.H.DB.Where(&models.StockDecreaseLog{OrderId: req.OrderId}).First(&log); res.Error == nil {
		resp = &pb.DecreaseStockResponse{
			Status: http.StatusConflict,
			Error:  "Stock already decreased",
		}
		return
	}

	product.Stock = product.Stock - 1

	s.H.DB.Save(&product)

	log.OrderId = req.OrderId
	log.ProductRefer = product.Id

	s.H.DB.Create(&log)

	resp = &pb.DecreaseStockResponse{
		Status: http.StatusOK,
	}

	return
}
