package main

import (
	"context"
	"database/sql"
	"net"
	"net/http"
	"os"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	migrate "github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	_ "github.com/lib/pq"
	"github.com/mikeheft/go-backend/api"
	db "github.com/mikeheft/go-backend/db/sqlc"
	"github.com/mikeheft/go-backend/gapi"
	"github.com/mikeheft/go-backend/pb"
	"github.com/mikeheft/go-backend/util"
	"github.com/mikeheft/go-backend/worker"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
)

const (
	serverCreateErrorMsg   = "cannot create server"
	listenerCreateErrorMsg = "cannot create listener"
)

func main() {
	config, err := util.LoadConfig(".")

	if config.Environment == "development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	if err != nil {
		log.Fatal().Err(err).Msg("cannot load config")
	}
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot connect to db:")
	}

	runDBMigration(config.MigrationUrl, config.DBSource)

	store := db.NewStore(conn)

	redisOpt := asynq.RedisClientOpt{
		Addr: config.RedisAddress,
	}
	taskDistributor := worker.NewRedisTaskDistributor(redisOpt)

	go runTaskProcessor(redisOpt, store)
	go runGatewayServer(config, store, taskDistributor)

	runGrpcServer(config, store, taskDistributor)
}

func runGrpcServer(config util.Config, store db.Store, taskDistributor worker.TaskDistributor) {
	server, err := gapi.NewServer(config, store, taskDistributor)
	if err != nil {
		log.Fatal().Err(err).Msg(serverCreateErrorMsg)
	}

	grpcLogger := grpc.UnaryInterceptor(gapi.GrpcLogger)

	grpcServer := grpc.NewServer(grpcLogger)
	pb.RegisterSimpleBankServer(grpcServer, server)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GrpcServerAddress)
	if err != nil {
		log.Fatal().Err(err).Msg(listenerCreateErrorMsg)
	}

	log.Info().Msgf("start gRPC server at: %s", listener.Addr().String())

	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal().Err(err).Msg(serverCreateErrorMsg)
	}
}

func runDBMigration(url string, dbSource string) {
	migration, err := migrate.New(url, dbSource)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create migration instance")
	}

	if err := migration.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal().Err(err).Msg("cannot run migrations")
	}

	log.Info().Msg("DB migrated successfully")
}

func runGatewayServer(config util.Config, store db.Store, taskDistributor worker.TaskDistributor) {
	server, err := gapi.NewServer(config, store, taskDistributor)
	if err != nil {
		log.Fatal().Err(err).Msg(serverCreateErrorMsg)
	}

	grpcMux := runtime.NewServeMux(
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
			MarshalOptions: protojson.MarshalOptions{
				UseProtoNames: true,
			},
			UnmarshalOptions: protojson.UnmarshalOptions{
				DiscardUnknown: true,
			},
		}),
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	err = pb.RegisterSimpleBankHandlerServer(ctx, grpcMux, server)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create server handler")
	}

	httpMux := http.NewServeMux()
	httpMux.Handle("/", grpcMux)

	listener, err := net.Listen("tcp", config.HttpServerAddress)
	if err != nil {
		log.Fatal().Err(err).Msg(listenerCreateErrorMsg)
	}

	log.Info().Msgf("start HTTP Gateway server at: %s", listener.Addr().String())

	handler := gapi.HttpLogger(httpMux)
	err = http.Serve(listener, handler)
	if err != nil {
		log.Fatal().Err(err).Msg(serverCreateErrorMsg)
	}
}

func runGinServer(config util.Config, store db.Store) {
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal().Err(err).Msg(serverCreateErrorMsg)
	}

	err = server.Start(config.HttpServerAddress)
	if err != nil {
		log.Fatal().Err(err).Msg(serverCreateErrorMsg)
	}
}

func runTaskProcessor(redisOpt asynq.RedisClientOpt, store db.Store) {
	taskProcessor := worker.NewRedisTaskProcessor(redisOpt, store)
	log.Info().Msg("Start task processor")
	err := taskProcessor.Start()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to start task processor")
	}
}
