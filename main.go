package main

import (
	"fmt"
	"github.com/dmitriibb/go-common/constants"
	"github.com/dmitriibb/go-common/db/mongo"
	"github.com/dmitriibb/go-common/logging"
	"github.com/dmitriibb/go-common/utils"
	"github.com/dmitriibb/go2-kitchen/internal/manager"
	"github.com/dmitriibb/go2-kitchen/internal/recipes"
	"github.com/dmitriibb/go2-kitchen/internal/storage"
	"github.com/dmitriibb/go2-kitchen/pkg/orders"
	"github.com/dmitriibb/go3"
	"google.golang.org/grpc"
	"net"
	"net/http"
)

var logger = logging.NewLogger("KitchenMain")

func main() {
	go3.TestHelloWorld()
	logger.Info("start")
	httpPort := utils.GetEnvProperty(constants.HttpPortEnv)
	grpcHost := utils.GetEnvProperty("GRPC_HOST", "localhost")
	grpcPort := utils.GetEnvProperty(constants.GrpcPortEnv)

	// init
	mongo.Init()
	recipes.Init()
	closeManagerChan := make(chan string)
	manager.Init(orders.OrdersHandler.NewOrders, closeManagerChan)
	closeStorageChan := make(chan string)
	storage.Init(closeStorageChan)

	// http handle
	go func() {
		http.HandleFunc("/recipes/menu", recipes.GetAllRecipesAsMenu)
		http.ListenAndServe(fmt.Sprintf(":%v", httpPort), nil)
	}()
	logger.Info("http started on %s", fmt.Sprintf(":%v", httpPort))

	// grpc handle
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%s", grpcHost, grpcPort))
	if err != nil {
		panic(fmt.Sprintf("failed to listen: %v", err))
	}
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	orders.RegisterKitchenOrdersHandlerServer(grpcServer, orders.OrdersHandler)
	logger.Info("Kitchen service registered...")

	go func() {
		grpcServer.Serve(lis)
	}()
	logger.Info("grpcServer.Serving on %s:%s...", grpcHost, grpcPort)

	forever := make(chan int)
	<-forever
}
