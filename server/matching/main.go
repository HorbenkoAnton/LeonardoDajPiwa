package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"log"
	"log/slog"
	"matching/cache"
	"matching/lib/env"
	"matching/lib/logger"
	pb "matching/proto"
	"net"
	"os"
)

type server struct {
	logger *slog.Logger
	pb.UnimplementedProfileServiceServer
}

var pg *pgxpool.Pool

func (s *server) GetNextProfile(_ context.Context, in *pb.IdRequest) (*pb.Profile, error) {
	id, err := cache.GetNext(pg, in.GetID())
	if err != nil {
		if errors.Is(err, cache.ErrNotFound) {
			return &pb.Profile{ID: -1}, nil
		}
		s.logger.Error("Error getting next profile: %v\n", err)
		return nil, err
	}

	return &pb.Profile{ID: id}, nil
}

func main() {
	setupLogger := logger.SetupLogger(env.LoadEnvVar("LOG_LEVEL"))

	connStr := fmt.Sprintf("postgres://%v:%v@%v:%v/%v",
		env.LoadEnvVar("DB_USER"),
		env.LoadEnvVar("DB_PASS"),
		env.LoadEnvVar("DB_HOST"),
		env.LoadEnvVar("DB_PORT"),
		env.LoadEnvVar("DB_NAME"),
	)

	pgconn, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		setupLogger.Error("Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	pg = pgconn

	srv := grpc.NewServer()
	pb.RegisterProfileServiceServer(srv, &server{logger: setupLogger})

	port := env.LoadEnvVar("MATCHING_PORT")

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		setupLogger.Error("Failed to listen: %v", err)
		os.Exit(1)
	}

	setupLogger.Info("matching server started")
	go cache.InvalidateCache()

	if err := srv.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
