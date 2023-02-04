package main

import (
	"api-gateway/config"
	"api-gateway/pkg/common"
	"api-gateway/pkg/discovery"
	"api-gateway/pkg/pb"
	"api-gateway/pkg/routers"
	"fmt"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	// init logger
	logger := common.GetLogger()

	// init configuration
	config.InitConfig()

	etcdResolver := discovery.NewResolver([]string{viper.GetString("etcd.address")}, logger)
	resolver.Register(etcdResolver)

	logger.Infof("start Gateway on :%s\n", viper.GetString("server.port"))
	go startServer()

	exitCh := make(chan os.Signal)
	signal.Notify(exitCh, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL, os.Interrupt)

	go func() {
		<-exitCh
		stopServer()
	}()
}

func startServer() {

	conn, _ := grpc.Dial(viper.GetString("server.grpcAddress"), []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}...)

	// grpc client
	userService := pb.NewUserServiceClient(conn)

	// init the router
	ginRouter := routers.NewRouter(userService)

	server := http.Server{
		Addr:           fmt.Sprintf(":%s", viper.GetString("server.port")),
		Handler:        ginRouter,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("failed to start the server: %s\n", err.Error())
	} else {
		fmt.Println("start the server successfully")
	}
}

func stopServer() {
	os.Exit(0)
}
