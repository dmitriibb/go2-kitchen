package storage

import (
	"github.com/dmitriibb/go-common/logging"
	"time"
)

var logger = logging.NewLogger("StorageRepository")
var ingredientsRequests = make(chan *IngredientRequest, 100)
var initialized = false

func Init(closeChan chan string) {
	if initialized {
		logger.Warn("Already initialized")
		return
	}

	go func() {
		for {
			select {
			case request := <-ingredientsRequests:
				processRequest(request)
			case closeMessage := <-closeChan:
				logger.Debug("Stop because %v", closeMessage)
				return
			}
		}
	}()
	initialized = true
	logger.Debug("initialized")
}

func RequireIngredients(ingredients []string, responseChan chan *IngredientsResponse) {
	ingredientsRequests <- &IngredientRequest{ingredients, responseChan}
}

func processRequest(request *IngredientRequest) {
	// TODO use db
	time.Sleep(200 * time.Millisecond)
	request.ResponseChan <- &IngredientsResponse{Success: true}
}

type IngredientRequest struct {
	Ingredients  []string
	ResponseChan chan *IngredientsResponse
}

type IngredientsResponse struct {
	Success bool
	Comment string
}
