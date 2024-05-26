package main

import (
	"context"
	"gateway/lib/env"
	"gateway/lib/logger"
	pb "gateway/protos/gatewaypb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log/slog"
	"net"
	"os"
	"time"
)

var (
	profileStub  pb.ProfileServiceClient
	matchingStub pb.ProfileServiceClient
	likesStub    pb.ProfileServiceClient
)

type server struct {
	logger *slog.Logger
	pb.UnimplementedProfileServiceServer
}

func (s *server) CreateProfile(ctx context.Context, in *pb.ProfileRequest) (*pb.ErrorResponse, error) {
	result, err := profileStub.CreateProfile(ctx, in)
	if err != nil {
		s.logger.Error("CreateProfile error, ", err)
		return nil, err
	}

	return result, nil
}

func (s *server) ReadProfile(ctx context.Context, in *pb.IdRequest) (*pb.Profile, error) {
	profile, err := profileStub.ReadProfile(ctx, in)
	if err != nil {
		s.logger.Error("ReadProfile error:", err)
		return nil, err
	}

	return profile, nil
}

func (s *server) UpdateProfile(ctx context.Context, in *pb.ProfileRequest) (*pb.ErrorResponse, error) {
	result, err := profileStub.UpdateProfile(ctx, in)
	if err != nil {
		s.logger.Error("UpdateProfile error:", err)
		return nil, err
	}

	return result, nil
}

func (s *server) GetNextProfile(ctx context.Context, in *pb.IdRequest) (*pb.Profile, error) {
	nextProfileId, err := matchingStub.GetNextProfile(ctx, in)
	if err != nil {
		s.logger.Error("GetNextProfile error, ", err)
		return nil, err
	}

	if nextProfileId.GetID() == -1 {
		return &pb.Profile{ID: -1}, nil
	}

	profile, err := profileStub.ReadProfile(ctx, &pb.IdRequest{Id: nextProfileId.GetID()})
	if err != nil {
		s.logger.Error("ReadProfile (getnext) error, ", err)
		return nil, err
	}

	return profile, nil
}

func (s *server) Like(ctx context.Context, in *pb.TargetRequest) (*pb.ErrorResponse, error) {
	arg := &pb.TargetRequest{
		Id:    in.Id,
		TgtId: in.TgtId,
	}
	result, err := likesStub.Like(ctx, arg)
	if err != nil {
		s.logger.Error("Like error: ", err)
		return nil, err
	}

	resultGW := &pb.ErrorResponse{ErrorMessage: result.ErrorMessage}

	return resultGW, nil
}

func (s *server) GetLikes(ctx context.Context, in *pb.IdRequest) (*pb.LikesResponse, error) {
	arg := &pb.IdRequest{Id: in.Id}

	result, err := likesStub.GetLikes(ctx, arg)

	if err != nil {
		s.logger.Error("GetLikes error, ", err)
		return nil, err
	}

	return result, nil
}

func main() {
	setupLogger := logger.SetupLogger(env.LoadEnvVar("LOG_LEVEL"))

	// sleep to give services time to start
	time.Sleep(5 * time.Second)

	opts := grpc.WithTransportCredentials(insecure.NewCredentials())
	matchingConn, err := grpc.Dial(
		net.JoinHostPort("matching", env.LoadEnvVar("MATCHING_PORT")),
		opts)
	if err != nil {
		setupLogger.Error("Cannot connect to matching, ", err)
		os.Exit(1)
	}
	matchingStub = pb.NewProfileServiceClient(matchingConn)

	profileConn, err := grpc.Dial(
		net.JoinHostPort("profiles", env.LoadEnvVar("PROFILES_PORT")),
		opts)
	if err != nil {
		setupLogger.Error("Cannot connect to profile, ", err)
		os.Exit(1)
	}
	profileStub = pb.NewProfileServiceClient(profileConn)

	likeConn, err := grpc.Dial(
		net.JoinHostPort("likes", env.LoadEnvVar("LIKES_PORT")),
		opts)
	if err != nil {
		setupLogger.Error("Cannot connect to likes, ", err)
		os.Exit(1)
	}
	likesStub = pb.NewProfileServiceClient(likeConn)

	srv := grpc.NewServer()
	pb.RegisterProfileServiceServer(srv, &server{logger: setupLogger})

	port := env.LoadEnvVar("GATEWAY_PORT")

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		setupLogger.Error("failed to listen: %v", err)
		os.Exit(1)
	}

	setupLogger.Info("Server started")

	if err := srv.Serve(lis); err != nil {
		setupLogger.Error("failed to serve: %v", err)
		os.Exit(1)
	}
}
