package workers

import (
	"context"
	commonNongo "github.com/dmitriibb/go-common/db/mongo"
	"github.com/dmitriibb/go-common/logging"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	orderItemStatusesCollection = "order_item_statuses"
)

var loggerRepo = logging.NewLogger("WorkersRepository")

func saveOrderItemWrapper(wrapper *OrderItemWrapper) {
	ctx := context.TODO()
	filter := bson.D{
		{"orderItem.orderId", wrapper.OrderItem.OrderId},
		{"orderItem.itemId", wrapper.OrderItem.ItemId},
	}
	update := bson.D{{"$set", wrapper}}
	f := func(client *mongo.Client) any {
		collection := client.Database(commonNongo.GetDbName()).Collection(orderItemStatusesCollection)
		collection.UpdateOne(ctx, filter, update, options.Update().SetUpsert(true))
		return 1
	}
	commonNongo.UseClient(ctx, f)
}
