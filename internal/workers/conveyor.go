package workers

import (
	"fmt"
	"github.com/dmitriibb/go2-kitchen/internal/model"
	"github.com/dmitriibb/go2-kitchen/internal/recipes"
	"time"
)

var conveyorItems = make(chan *OrderItemWrapper, 100)

type OrderItemWrapper struct {
	OrderItem   *model.OrderItem     `bson:"orderItem"`
	RecipeStage *recipes.RecipeStage `bson:"recipeStage"`
	Comment     string               `bson:"comment"`
}

func (wrapper *OrderItemWrapper) String() string {
	return fmt.Sprintf("{order: %v, item: %v, dishName: %v}", wrapper.OrderItem.OrderId, wrapper.OrderItem.ItemId, wrapper.RecipeStage.Name)
}

// TODO create a list or db with timers and shut them down properly when application stops
func startConveyorTimer(orderItemWrapper *OrderItemWrapper, timeDelaySec int64) {
	go func() {
		time.Sleep(time.Duration(timeDelaySec) * time.Second)
		conveyorItems <- orderItemWrapper
	}()
}
