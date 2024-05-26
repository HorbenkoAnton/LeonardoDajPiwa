package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	el "likes/lib/env"
	"likes/lib/logger"
	pb "likes/proto"
	"log/slog"
	"net"
	"os"
	"time"
)

var (
	pg *pgxpool.Pool
)

const Timeout = 60 * time.Second

type server struct {
	logger *slog.Logger
	pb.UnimplementedProfileServiceServer
}

func (s *server) Like(_ context.Context, in *pb.TargetRequest) (*pb.ErrorResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), Timeout)
	defer cancel()

	if in.GetTgtId() == 0 || in.GetId() == 0 {
		s.logger.Error("empty id or target id")
		return &pb.ErrorResponse{ErrorMessage: "empty target id"}, nil
	}

	result, err := pg.Exec(ctx, "INSERT INTO likes (liker, liked) VALUES ($1, $2) ON CONFLICT (liker, liked) DO NOTHING", in.GetId(), in.GetTgtId())
	if err != nil {
		s.logger.Error("failed to insert into likes", slog.Int64("id", in.GetId()), slog.Int64("target_id", in.GetTgtId()), slog.String("error", err.Error()))
		return &pb.ErrorResponse{ErrorMessage: err.Error()}, nil
	}

	if rowsAffected := result.RowsAffected(); rowsAffected == 0 {
		s.logger.Info("like already exists")
	}

	return &pb.ErrorResponse{ErrorMessage: "OK"}, nil
}

func (s *server) GetLikes(_ context.Context, in *pb.IdRequest) (*pb.LikesResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), Timeout)
	defer cancel()

	rows, err := pg.Query(ctx, "SELECT * FROM likes WHERE liked = $1", in.GetId())
	if err != nil {
		s.logger.Error("failed to select from likes", slog.Int64("id", in.GetId()), slog.String("error", err.Error()))
		return nil, err
	}

	var likes []*pb.Profile
	for rows.Next() {
		var liker, liked int64
		if err = rows.Scan(&liker, &liked); err != nil {
			s.logger.Error("failed to scan rows in likes", slog.String("error", err.Error()))
			return nil, err
		}

		profile := &pb.Profile{
			ID: liker,
		}

		likes = append(likes, profile)
	}

	if _, err = pg.Exec(ctx, "DELETE FROM likes WHERE liked = $1", in.GetId()); err != nil {
		s.logger.Error("failed to delete from likes", slog.Int64("id", in.GetId()), slog.String("error: ", err.Error()))
		return nil, err
	}

	return &pb.LikesResponse{Likes: likes}, nil
}

func main() {
	setupLogger := logger.SetupLogger(el.LoadEnvVar("LOG_LEVEL"))

	connStr := fmt.Sprintf("postgres://%v:%v@%v:%v/%v",
		el.LoadEnvVar("DB_USER"),
		el.LoadEnvVar("DB_PASS"),
		el.LoadEnvVar("DB_HOST"),
		el.LoadEnvVar("DB_PORT"),
		el.LoadEnvVar("DB_NAME"),
	)

	pgconn, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		setupLogger.Error("unable to connect to database", slog.String("error", err.Error()))
		os.Exit(1)
	}
	pg = pgconn

	srv := grpc.NewServer()
	pb.RegisterProfileServiceServer(srv, &server{logger: setupLogger})

	likesPort := el.LoadEnvVar("LIKES_PORT")

	lis, err := net.Listen("tcp", ":"+likesPort)
	if err != nil {
		setupLogger.Error("failed to listen tcp", slog.String("likes_port", likesPort), slog.String("error", err.Error()))
		os.Exit(1)
	}

	setupLogger.Info("likes server started")

	if err = srv.Serve(lis); err != nil {
		setupLogger.Error("failed to serve", slog.String("error", err.Error()))
		os.Exit(1)
	}
}
