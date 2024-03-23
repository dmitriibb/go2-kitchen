package orders

import (
	"context"
	"github.com/dmitriibb/go-common/logging"
)

var loggerService = logging.NewLogger("KitchenOrders")

type ordersHandler struct {
	NewOrders chan *PutNewOrderRequest
}

func (ko *ordersHandler) mustEmbedUnimplementedKitchenOrdersHandlerServer() {
	panic("Not implemented")
}

var OrdersHandler = &ordersHandler{NewOrders: make(chan *PutNewOrderRequest, 100)}

func (ko *ordersHandler) PutNewOrder(ctx context.Context, in *PutNewOrderRequest) (*PutNewOrderResponse, error) {
	loggerService.Debug("Received new order %s", in)

	if in.Items[0].Comment == "kitchen error" {
		loggerService.Warn("Fake error on receiving new order. Cancel ctx")
		return &PutNewOrderResponse{Status: "error"}, nil
	}

	go func() {
		ko.NewOrders <- in
	}()
	return &PutNewOrderResponse{Status: "Received"}, nil
}
