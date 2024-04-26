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

func (s *server) GetNextProfile(_ context.Context, in *pb.IdReqResp) (*pb.IdReqResp, error) {
	id, err := cache.GetNext(pg, in.GetID())
	if err != nil {
		if errors.Is(err, cache.ErrNotFound) {
			return &pb.IdReqResp{ID: -1}, nil
		}
		log.Printf("Error getting next profile: %v\n", err)
		return nil, err
	}

	return &pb.IdReqResp{ID: int64(id)}, nil
}

func main() {
	connStr := fmt.Sprintf("postgres://%v:%v@%v:%v/%v",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
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
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Println("matching server started")
	go cache.InvalidateCache()

	if err := srv.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
