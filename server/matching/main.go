package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"log"
	"matching/cache"
	pb "matching/proto"
	"net"
	"os"
)

type server struct {
	pb.UnimplementedMatchingServiceServer
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
		os.Getenv("PG_USER"),
		os.Getenv("PG_PASS"),
		os.Getenv("PG_HOST"),
		os.Getenv("PG_PORT"),
		os.Getenv("PG_DB"),
	)

	pgconn, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	pg = pgconn

	srv := grpc.NewServer()
	pb.RegisterMatchingServiceServer(srv, &server{})

	port := os.Getenv("MATCHING_PORT")
	if port == "" {
		log.Fatalf("Error: port not provided, add MATCHING_PORT env var")
	}

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Println("Server started")
	go cache.InvalidateCache()

	if err := srv.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
