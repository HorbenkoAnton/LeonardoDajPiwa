package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"log"
	"matching/cache"
	"matching/env"
	pb "matching/proto"
	"net"
)

type server struct {
	pb.UnimplementedProfileServiceServer
}

var pg *pgxpool.Pool

func (s *server) GetNextProfile(_ context.Context, in *pb.IdRequest) (*pb.Profile, error) {
	id, err := cache.GetNext(pg, in.GetID())
	if err != nil {
		if errors.Is(err, cache.ErrNotFound) {
			return &pb.Profile{ID: -1}, nil
		}
		log.Printf("Error getting next profile: %v\n", err)
		return nil, err
	}

	return &pb.Profile{ID: int64(id)}, nil
}

func main() {
	connStr := fmt.Sprintf("postgres://%v:%v@%v:%v/%v",
		env.LoadEnvVar("DB_USER"),
		env.LoadEnvVar("DB_PASS"),
		env.LoadEnvVar("DB_HOST"),
		env.LoadEnvVar("DB_PORT"),
		env.LoadEnvVar("DB_NAME"),
	)

	pgconn, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	pg = pgconn

	srv := grpc.NewServer()
	pb.RegisterProfileServiceServer(srv, &server{})

	port := env.LoadEnvVar("MATCHING_PORT")

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Println("matching server started")
	go cache.InvalidateCache()

	if err := srv.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
