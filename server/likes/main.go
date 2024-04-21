package main

import (
	"context"
	pb "dp_mock/backend/likes/proto"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"log"
	"net"
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
		return &pb.ErrorResponse{ErrorMessage: "the like already exists"}, nil
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
	pgconn, err := pgxpool.New(context.Background(), "postgres://postgres:password@localhost:5432/db_likes")
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	pg = pgconn

	srv := grpc.NewServer()
	pb.RegisterLikeServiceServer(srv, &server{})

	lis, err2 := net.Listen("tcp", ":50051")
	if err2 != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	fmt.Println("Server is running on port :50051")

	if err = srv.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
