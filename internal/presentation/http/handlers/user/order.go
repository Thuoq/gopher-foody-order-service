package user

import (
	"gopher-order-service/internal/core/domain"
	"gopher-order-service/internal/core/ports"
	"gopher-order-service/internal/presentation/http/handlers/user/dto/request"
	"gopher-order-service/internal/presentation/http/handlers/user/dto/response"
	"gopher-order-service/pkg/app_response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	createUC ports.ICreateOrderUseCase
	//detailUC ports.IGetOrderDetailUseCase
}

func NewOrderHandler(
	createUC ports.ICreateOrderUseCase,
// detailUC ports.IGetOrderDetailUseCase,
) *OrderHandler {
	return &OrderHandler{
		createUC: createUC,
		//detailUC: detailUC,
	}
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var req request.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fieldErrors := app_response.ParseValidationErrors(err)
		app_response.ValidationError(c, http.StatusBadRequest, "invalid request body", fieldErrors)
		return
	}

	userID := c.GetString("public_user_id")
	if userID == "" {
		app_response.Error(c, http.StatusUnauthorized, "missing user identity")
		return
	}

	items := make([]ports.CreateOrderItemInput, len(req.Items))
	for i, item := range req.Items {
		items[i] = ports.CreateOrderItemInput{
			FoodID:   item.FoodID,
			Quantity: item.Quantity,
		}
	}

	input := ports.CreateOrderInput{
		UserID:          userID,
		RestaurantID:    req.RestaurantID,
		ShippingAddress: req.ShippingAddress,
		Note:            req.Note,
		Items:           items,
	}

	order, err := h.createUC.Execute(c.Request.Context(), input)
	if err != nil {
		app_response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	app_response.Success(c, http.StatusCreated, mapToOrderResponse(order))
}

func mapToOrderResponse(o *domain.Order) response.OrderResponse {
	items := make([]response.OrderItemResponse, len(o.Items))
	for i, item := range o.Items {
		items[i] = response.OrderItemResponse{
			FoodID:   item.FoodID,
			Name:     item.Name,
			Price:    item.Price,
			Quantity: item.Quantity,
			Subtotal: item.Subtotal,
		}
	}

	return response.OrderResponse{
		ID:              o.PublicID,
		UserID:          o.UserID,
		RestaurantID:    o.RestaurantID,
		Status:          string(o.Status),
		TotalPrice:      o.TotalPrice,
		ShippingAddress: o.ShippingAddress,
		PaymentStatus:   o.PaymentStatus,
		PaymentMethod:   o.PaymentMethod,
		Note:            o.Note,
		Items:           items,
		CreatedAt:       o.CreatedAt,
	}
}
