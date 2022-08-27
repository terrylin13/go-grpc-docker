package order

import (
	"github.com/gin-gonic/gin"
	"github.com/terrylin13/go-grpc-docker/api-gateway/pkg/auth"
	"github.com/terrylin13/go-grpc-docker/api-gateway/pkg/config"
	"github.com/terrylin13/go-grpc-docker/api-gateway/pkg/order/routes"
)

func RegisterRoutes(r *gin.Engine, c *config.Config, authSvc *auth.ServiceClient) {
	a := auth.InitAuthMiddleware(authSvc)
	svc := &ServiceClient{
		Client: InitServiceClient(c),
	}

	routes := r.Group("/order")
	routes.Use(a.AuthRequired)
	routes.POST("/", svc.CreatOrder)
}

func (svc *ServiceClient) CreatOrder(ctx *gin.Context) {
	routes.CreateOrder(ctx, svc.Client)
}
