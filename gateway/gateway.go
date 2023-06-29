package gateway

import (
	"context"
	"fmt"
	"net"
	"net/http"
	// "sync"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"

	"github.com/ANMalko/grpc-server.git/proto/users"
	usersserver "github.com/ANMalko/grpc-server.git/server/users"
	usersfiledao "github.com/ANMalko/grpc-server.git/db/dao/users/filedb"

)

func Run(ctx context.Context, usersFileDAO *usersfiledao.DAO, httpPort int, gRPCPort int) {
	go runRest(ctx, usersFileDAO, httpPort)
	runGrpc(ctx, usersFileDAO, gRPCPort)
}

func runRest(ctx context.Context, usersFileDAO *usersfiledao.DAO, port int) error {
	usersServer := usersserver.NewServer(usersFileDAO)
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

func runGrpc(ctx context.Context, usersFileDAO *usersfiledao.DAO, port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatal().Msgf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	userServer := usersserver.NewServer(usersFileDAO)
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
