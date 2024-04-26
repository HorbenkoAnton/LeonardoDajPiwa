package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"google.golang.org/grpc"
	lm "likes/migrations"
	pb "likes/proto"
	"log"
	"net"
	"os"
	"strconv"
	"time"
)

var (
	pg *pgxpool.Pool
)

const Timeout = 60 * time.Second

type server struct {
	pb.UnsafeLikeServiceServer
}

func (s *server) Like(_ context.Context, in *pb.TargetRequest) (*pb.ErrorResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), Timeout)
	defer cancel()

	if in.GetTgtId() == 0 || in.GetId() == 0 {
		return &pb.ErrorResponse{ErrorMessage: "empty target id"}, nil
	}

	result, err := pg.Exec(ctx, "INSERT INTO likes (liker, liked) VALUES ($1, $2) ON CONFLICT (liker, liked) DO NOTHING", in.GetId(), in.GetTgtId())
	if err != nil {
		return &pb.ErrorResponse{ErrorMessage: err.Error()}, nil
	}

	if rowsAffected := result.RowsAffected(); rowsAffected == 0 {
		return &pb.ErrorResponse{ErrorMessage: "like already exists"}, nil
	}

	return &pb.ErrorResponse{ErrorMessage: "OK"}, nil
}

func (s *server) GetLikes(_ context.Context, in *pb.IdRequest) (*pb.LikesResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), Timeout)
	defer cancel()

	if in.GetId() == 0 {
		return &pb.LikesResponse{}, errors.New("empty id")
	}

	rows, err := pg.Query(ctx, "SELECT * FROM likes WHERE liked = $1", in.GetId())
	if err != nil {
		return nil, err
	}

	var likes []*pb.Profile
	for rows.Next() {
		var liker, liked int64
		if err = rows.Scan(&liker, &liked); err != nil {
			return nil, err
		}

		profile := &pb.Profile{
			ID: liker,
		}

		likes = append(likes, profile)
	}

	if _, err = pg.Exec(ctx, "DELETE FROM likes WHERE liked = $1", in.GetId()); err != nil {
		return nil, err
	}

	return &pb.LikesResponse{Likes: likes}, nil
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
	db := stdlib.OpenDBFromPool(pg)

	reload, err := strconv.ParseBool(os.Getenv("DB_RELOAD"))
	if err != nil {
		log.Fatal(err)
	}

	lm.Migrate(reload, db)

	srv := grpc.NewServer()
	pb.RegisterLikeServiceServer(srv, &server{})

	port := os.Getenv("LIKES_PORT")
	if port == "" {
		log.Fatalf("Error: port not provided, add LIKES_PORT env var")
	}

	lis, err2 := net.Listen("tcp", ":"+port)
	if err2 != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	fmt.Println("likes server started")

	if err = srv.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
