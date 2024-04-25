package main

import (
	"context"
	"fmt"
	pb "gateway/protos/gatewaypb"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
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
	matchingConn, err := grpc.Dial(net.JoinHostPort("localhost", os.Getenv("MATCHING_PORT")))
	if err != nil {
		log.Fatalln("Cannot connect to matching, ", err)
	}
	matchingStub = pb.NewProfileServiceClient(matchingConn)

	profileConn, err := grpc.Dial(net.JoinHostPort("localhost", os.Getenv("PROFILES_PORT")))
	if err != nil {
		log.Fatalln("Cannot connect to profile, ", err)
	}
	profileStub = pb.NewProfileServiceClient(profileConn)

	likeConn, err := grpc.Dial(net.JoinHostPort("localhost", os.Getenv("LIKES_PORT")))
	if err != nil {
		log.Fatalln("Cannot connect to likes, ", err)
	}
	likesStub = pb.NewProfileServiceClient(likeConn)

	lis, err := net.Listen("tcp", net.JoinHostPort("localhost", os.Getenv("GATEWAY_PORT")))
	if err != nil {
		fmt.Println("Failed to listen:", err)
		return
	}

	s := grpc.NewServer()

	pb.RegisterProfileServiceServer(s, &server{})
	fmt.Println("Server started at port", os.Getenv("GATEWAY_PORT"))

	if err := s.Serve(lis); err != nil {
		fmt.Println("Failed to serve:", err)
		return
	}
}
