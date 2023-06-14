package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"

	"github.com/ANMalko/grpc-server.git/config"
	"github.com/ANMalko/grpc-server.git/proto/users"
	usersserver "github.com/ANMalko/grpc-server.git/server/users"
	"github.com/ANMalko/grpc-server.git/utils"
)

var httpPortFlag = flag.String("http", "", "bind address")
var gRPCPortFlag = flag.String("grpc", "", "bind address")

func runRest(ctx context.Context, port int) error {
	usersServer := usersserver.NewServer()
	mux := runtime.NewServeMux()

	if err := users.RegisterUserServiceHandlerServer(ctx, mux, usersServer); err != nil {
		log.Fatal().Err(err)
	}

	address := fmt.Sprintf(":%d", port)
	server := &http.Server{Addr: address, Handler: mux}

	go func() {
		<-ctx.Done()
		log.Info().Msg("Shutting down the http gateway server")

		if err := server.Shutdown(ctx); err != nil {
			log.Error().Err(err).Msg("Failed to shutdown http gateway server")
		}
	}()

	log.Info().Msgf("server listening at %d", port)
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Error().Err(err).Msg("Failed to listen and serve")
		return err
	}

	return nil
}

func runGrpc(ctx context.Context, port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatal().Msgf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	userServer := usersserver.NewServer()
	users.RegisterUserServiceServer(grpcServer, userServer)

	go func() {
		<-ctx.Done()
		log.Info().Msg("Shutting down the gRPC server")
		grpcServer.GracefulStop()
	}()

	log.Info().Msgf("server listening at %v", lis.Addr())
	if err := grpcServer.Serve(lis); err != nil {
		return err
	}

	return nil
}

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

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		select {
		case <-ctx.Done():
			log.Info().Msg("terminating: context cancelled")
		case <-sigterm:
			log.Info().Msg("terminating: via signal")
		}
		cancel()
	}()

	go runRest(ctx, httpPort)
	runGrpc(ctx, gRPCPort)
}
