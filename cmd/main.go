package main

import (
	"fmt"
	"context"
	"flag"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/rs/zerolog/log"

	"github.com/ANMalko/grpc-server.git/config"
	"github.com/ANMalko/grpc-server.git/loger"
	"github.com/ANMalko/grpc-server.git/gateway"
	usersfiledao"github.com/ANMalko/grpc-server.git/db/dao/users/filedb"
)

var httpPortFlag = flag.String("http", "", "bind address")
var gRPCPortFlag = flag.String("grpc", "", "bind address")

func main() {
	cfg := config.GetAPIConfig()
	ctx, cancel := context.WithCancel(context.Background())
	zlog.InitZerolog()
	flag.Parse()

	var httpPort int
	var gRPCPort int

	if *httpPortFlag != "" {
		parsedPort, err := strconv.Atoi(*httpPortFlag)
		if err != nil {

			log.Fatal().Msg("Invalid http port")
		}
		httpPort = parsedPort

	} else {
		httpPort = cfg.HTTPPort
	}

	if *gRPCPortFlag != "" {
		parsedPort, err := strconv.Atoi(*gRPCPortFlag)
		if err != nil {

			log.Fatal().Msg("Invalid gRPC port")
		}
		gRPCPort = parsedPort

	} else {
		gRPCPort = cfg.GRPCPort
	}

	usersDAO := usersfiledao.NewDAO(ctx, cfg.UserDbFilename)

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		select {
		case <-ctx.Done():
			log.Info().Msg("terminating: context cancelled")
		case <-sigterm:
			log.Info().Msg("terminating: via signal")
		}
		fmt.Println("-------1-------")
		usersDAO.DB().DumpDB()
		fmt.Println("-------2-------")
		cancel()
	}()

	gateway.Run(ctx, usersDAO, httpPort, gRPCPort)
}
