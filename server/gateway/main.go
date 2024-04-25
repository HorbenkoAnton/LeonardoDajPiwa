package main

import (
	"context"
	"fmt"
	pb "gateway/protos/gatewaypb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
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
	pb.UnimplementedProfileServiceServer
}

func (s *server) CreateProfile(ctx context.Context, in *pb.ProfileRequest) (*pb.ErrorResponse, error) {
	result, err := profileStub.CreateProfile(ctx, in)
	if err != nil {
		fmt.Println("CreateProfile error, ", err)
		return nil, err
	}

	return result, nil
}

func (s *server) ReadProfile(ctx context.Context, in *pb.IdRequest) (*pb.Profile, error) {
	profile, err := profileStub.ReadProfile(ctx, in)
	if err != nil {
		fmt.Println("ReadProfile error:", err)
		return nil, err
	}

	return profile, nil
}

func (s *server) UpdateProfile(ctx context.Context, in *pb.ProfileRequest) (*pb.ErrorResponse, error) {
	result, err := profileStub.UpdateProfile(ctx, in)
	if err != nil {
		fmt.Println("UpdateProfile error:", err)
		return nil, err
	}

	return result, nil
}

func (s *server) GetNextProfile(ctx context.Context, in *pb.IdRequest) (*pb.Profile, error) {
	nextProfileId, err := matchingStub.GetNextProfile(ctx, in)
	if err != nil {
		fmt.Println("GetNextProfile error, ", err)
		return nil, err
	}

	if nextProfileId.GetID() == -1 {
		return &pb.Profile{ID: -1}, nil
	}

	profile, err := profileStub.ReadProfile(ctx, &pb.IdRequest{Id: nextProfileId.GetID()})
	if err != nil {
		fmt.Println("ReadProfile (getnext) error, ", err)
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
		fmt.Println("Like error: ", err)
		return nil, err
	}

	resultGW := &pb.ErrorResponse{ErrorMessage: result.ErrorMessage}

	return resultGW, nil
}

func (s *server) GetLikes(ctx context.Context, in *pb.IdRequest) (*pb.LikesResponse, error) {
	arg := &pb.IdRequest{Id: in.Id}

	result, err := likesStub.GetLikes(ctx, arg)

	if err != nil {
		fmt.Println("GetLikes error, ", err)
		return nil, err
	}

	return result, nil
}

func main() {
	// sleep to give services time to start
	time.Sleep(5 * time.Second)
	opts := grpc.WithTransportCredentials(insecure.NewCredentials())
	matchingConn, err := grpc.Dial(
		net.JoinHostPort("matching", os.Getenv("MATCHING_PORT")),
		opts)
	if err != nil {
		log.Fatalln("Cannot connect to matching, ", err)
	}
	matchingStub = pb.NewProfileServiceClient(matchingConn)

	profileConn, err := grpc.Dial(
		net.JoinHostPort("profiles", os.Getenv("PROFILES_PORT")),
		opts)
	if err != nil {
		log.Fatalln("Cannot connect to profile, ", err)
	}
	profileStub = pb.NewProfileServiceClient(profileConn)

	likeConn, err := grpc.Dial(
		net.JoinHostPort("likes", os.Getenv("LIKES_PORT")),
		opts)
	if err != nil {
		log.Fatalln("Cannot connect to likes, ", err)
	}
	likesStub = pb.NewProfileServiceClient(likeConn)

	srv := grpc.NewServer()
	pb.RegisterProfileServiceServer(srv, &server{})

	port := os.Getenv("GATEWAY_PORT")
	if port == "" {
		log.Fatalf("Error: port not provided, add GATEWAY env var")
	}

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Println("Server started")

	if err := srv.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
