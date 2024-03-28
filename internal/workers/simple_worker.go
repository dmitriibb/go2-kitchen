package workers

import (
	"fmt"
	"github.com/dmitriibb/go-common/logging"
	"github.com/dmitriibb/go2-kitchen/internal/buffers"
	"github.com/dmitriibb/go2-kitchen/internal/model"
	"github.com/dmitriibb/go2-kitchen/internal/recipes"
	"github.com/dmitriibb/go2-kitchen/internal/storage"
	"time"
)

type simpleWorker struct {
	id       string
	stopChan chan string
	logger   logging.Logger
}

func New(name string) Worker {
	return &simpleWorker{
		id:       name,
		stopChan: make(chan string),
		logger:   logging.NewLogger(fmt.Sprintf("worker-%v", name)),
	}
}

type Worker interface {
	Start()
	Stop()
}

func (worker *simpleWorker) Start() {
	worker.logger.Debug("Init working")
	go func() {
		for {
			select {
			case newOrderItem := <-buffers.NewOrderItems:
				worker.logger.Info("take new order item %v", newOrderItem)
				worker.processOrderItem(newOrderItem)
			case orderItemWrapper := <-conveyorItems:
				if orderItemWrapper.RecipeStage != nil {
					worker.logger.Info("take item from conveyor %v", orderItemWrapper)
					worker.cook(orderItemWrapper)
				} else {
					worker.logger.Info("take item without recipe from conveyor %v", orderItemWrapper)
					worker.processOrderItem(orderItemWrapper.OrderItem)
				}
			case stop := <-worker.stopChan:
				worker.logger.Debug("Stop because &v", stop)
				return
			}
		}
	}()
}

func (worker *simpleWorker) Stop() {
	worker.logger.Debug("Finish")
	worker.stopChan <- "Stop() called"
}

func (worker *simpleWorker) processOrderItem(item *model.OrderItem) {
	item.Status = model.OrderItemInProgress
	itemWrapper := &OrderItemWrapper{
		OrderItem: item,
	}
	recipe, err := recipes.GetRecipe(item.Name)
	if err != nil {
		comment := fmt.Sprintf("can't get recipe for '%v' because - '%v', Retry again after 30 sec", item.Name, err.Error())
		worker.logger.Error(comment)
		item.Status = model.OrderItemError
		item.Comment = comment
		startConveyorTimer(itemWrapper, 30)
		return
	}

	itemWrapper.RecipeStage = &recipe

	worker.cook(itemWrapper)
}

func (worker *simpleWorker) cook(itemWrapper *OrderItemWrapper) {
	ready := worker.cookRecipeStage(itemWrapper, itemWrapper.RecipeStage)

	if ready && itemWrapper.RecipeStage.Status == recipes.RecipeStageStatusFinished {
		worker.logger.Info("finished to cook %v", itemWrapper.OrderItem)
		saveOrderItemWrapper(itemWrapper)
		itemWrapper.OrderItem.Status = model.OrderItemReady
		buffers.ReadyOrderItems <- itemWrapper.OrderItem
	}
}

// return true if can proceed to the next stage
func (worker *simpleWorker) cookRecipeStage(itemWrapper *OrderItemWrapper, currentStage *recipes.RecipeStage) bool {

	switch currentStage.Status {
	case recipes.RecipeStageStatusFinished:
		return true
	case recipes.RecipeStageStatusInProgress:
		worker.logger.Info("received %s  stage: '%s' with status in progress. Check if it is ready", itemWrapper, currentStage.Name)
		currentTime := time.Now().Unix()
		timePassed := currentTime - currentStage.TimeStarted
		if currentStage.TimeStarted > 0 && timePassed >= currentStage.TimeToWaitSec {
			worker.logger.Info("%s stage '%s' is ready after %s sec. Continue sub stages", itemWrapper, currentStage.Name, timePassed)
			ready := worker.cookSubStages(itemWrapper, currentStage)
			if ready {
				worker.logger.Info("%s stage '%s' is ready after %v sec", itemWrapper, currentStage.Name, time.Now().Unix()-currentStage.TimeStarted)
				currentStage.Status = recipes.RecipeStageStatusFinished
			}
			return ready
		} else {
			return worker.cookCurrentStages(itemWrapper, currentStage)
		}
	case recipes.RecipeStageStatusEmpty, recipes.RecipeStageStatusError:
		// TODO maybe need to check all ingredients at the root stage
		responseChan := make(chan *storage.IngredientsResponse)
		storage.RequireIngredients(currentStage.Ingredients, responseChan)

		response := <-responseChan
		if !response.Success {
			currentStage.Status = recipes.RecipeStageStatusError
			currentStage.Comment = fmt.Sprintf("can't get ingredients because %v. Will try again after 30 sec", response.Comment)
			saveOrderItemWrapper(itemWrapper)
			startConveyorTimer(itemWrapper, 30)
			return false
		}
		return worker.cookCurrentStages(itemWrapper, currentStage)
	}
	panic(fmt.Sprintf("Unexpected for %s, %v", itemWrapper, currentStage))
}

func (worker *simpleWorker) cookCurrentStages(itemWrapper *OrderItemWrapper, currentStage *recipes.RecipeStage) bool {
	currentStage.Status = recipes.RecipeStageStatusInProgress
	currentStage.TimeStarted = time.Now().Unix()
	if currentStage.TimeToWaitSec > 5 {
		worker.logger.Info("keep stage '%v' cooking for %v sec and continue it later", currentStage.Name, currentStage.TimeToWaitSec)
		saveOrderItemWrapper(itemWrapper)
		startConveyorTimer(itemWrapper, currentStage.TimeToWaitSec)
		return false
	} else {
		worker.logger.Info("cooking - %v", currentStage.Name)
		time.Sleep(time.Duration(currentStage.TimeToWaitSec) * time.Second)
		ready := worker.cookSubStages(itemWrapper, currentStage)
		if ready {
			worker.logger.Info("%s stage '%s' is ready after %v sec", itemWrapper, currentStage.Name, time.Now().Unix()-currentStage.TimeStarted)
		}
		return ready
	}
}

func (worker *simpleWorker) cookSubStages(itemWrapper *OrderItemWrapper, currentStage *recipes.RecipeStage) bool {
	if currentStage.SubStages != nil {
		for _, subStage := range currentStage.SubStages {
			if !worker.cookRecipeStage(itemWrapper, subStage) {
				return false
			}
		}
	}
	currentStage.Status = recipes.RecipeStageStatusFinished
	return true
}
