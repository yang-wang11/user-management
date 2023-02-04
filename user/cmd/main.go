package main

import (
	"fmt"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"user/config"
	"user/discovery"
	"user/pkg/handler"
	"user/pkg/pb"
	"user/pkg/repository"
	"user/pkg/util"
)

func main() {

	// init logger
	logger := util.GetLogger()

	// init configuration
	config.InitConfig()

	// init DB
	db, err := repository.InitDB(logger)
	if err != nil {
		log.Fatalln(err)
	}

	// register service on etcd
	grpcAddress := fmt.Sprintf("%s:%s", viper.GetString("server.address"), viper.GetString("server.port"))
	server := discovery.Server{
		Name:    viper.GetString("server.name"),
		Addr:    grpcAddress,
		Version: viper.GetString("server.version"),
	}

	etcdRegister := discovery.NewRegister([]string{viper.GetString("etcd.address")}, logger)
	if err = etcdRegister.Register(&server, 10); err != nil {
		log.Fatalln(err)
	}

	exitCh := make(chan os.Signal, 1)
	signal.Notify(exitCh, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL, os.Interrupt)

	go func() {
		<-exitCh
		etcdRegister.Stop()
	}()

	// start the gcp server
	logger.Infof("start Gateway on :%s\n", viper.GetString("server.port"))
	grpcServer := grpc.NewServer()
	defer grpcServer.Stop()
	pb.RegisterUserServiceServer(grpcServer, handler.NewUserService(db))
	listener, err := net.Listen("tcp", grpcAddress)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("server run on %s\n", grpcAddress)

	if err = grpcServer.Serve(listener); err != nil {
		log.Fatalln(err)
	}
}
