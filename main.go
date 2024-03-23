package main

import (
	"fmt"
	"github.com/dmitriibb/go-common/constants"
	"github.com/dmitriibb/go-common/db/mongo"
	"github.com/dmitriibb/go-common/logging"
	"github.com/dmitriibb/go-common/utils"
	"github.com/dmitriibb/go2-kitchen/manager"
	"github.com/dmitriibb/go2-kitchen/pkg/orders"
	"github.com/dmitriibb/go2-kitchen/recipes"
	"github.com/dmitriibb/go2-kitchen/storage"
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
		http.ListenAndServe(fmt.Sprintf(":%v", httpPort), nil)
		logger.Info("http started...")
	}()

	// grpc handle
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%s", grpcPort))
	if err != nil {
		panic(fmt.Sprintf("failed to listen: %v", err))
	}
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	orders.RegisterKitchenOrdersHandlerServer(grpcServer, orders.OrdersHandler)
	logger.Info("Kitchen service registered...")
	grpcServer.Serve(lis)
}
