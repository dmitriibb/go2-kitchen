package model

type OrderItem struct {
	OrderId int             `bson:"orderId"`
	ItemId  int             `bson:"itemId"`
	Name    string          `bson:"name"`
	Comment string          `bson:"comment"`
	Status  OrderItemStatus `bson:"status"`
}

type OrderItemStatus string

const (
	OrderItemNew        OrderItemStatus = "New"
	OrderItemInProgress OrderItemStatus = "InProgress"
	OrderItemReady      OrderItemStatus = "Ready"
	OrderItemError      OrderItemStatus = "Error"
)
