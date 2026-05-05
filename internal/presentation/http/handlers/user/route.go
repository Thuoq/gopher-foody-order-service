package user

import (
	"github.com/gin-gonic/gin"
)

type Router struct {
	orderHander *OrderHandler
}

func NewRouter(orderHander *OrderHandler) *Router {
	return &Router{
		orderHander: orderHander,
	}
}

func (r *Router) Register(api *gin.RouterGroup) {
	userGroup := api.Group("/sso")
	{
		userGroup.POST("/orders", r.orderHander.CreateOrder)
	}
}
