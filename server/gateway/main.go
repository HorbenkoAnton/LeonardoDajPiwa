package main

import (
	"context"
	"fmt"
	"gateway/protos/gatewaypb"
	"gateway/protos/likespb"
	"gateway/protos/matchingpb"
	profilespb "gateway/protos/profilepb"
	"google.golang.org/grpc"
	"log"
	"net"
)

var matchingStub matchingpb.MatchingServiceClient
var profileStub profilespb.ProfileServiceClient
var likeStub likespb.LikeServiceClient

type server struct {
	gatewaypb.UnimplementedProfileServiceServer
}

func (s *server) CreateProfile(ctx context.Context, in *gatewaypb.ProfileRequest) (*gatewaypb.ErrorResponse, error) {
	arg := &profilespb.ProfileRequest{}
	arg.Profile = &profilespb.Profile{
		ID:          in.Profile.ID,
		Name:        in.Profile.Name,
		Age:         in.Profile.Age,
		Description: in.Profile.Description,
		Location:    in.Profile.Location,
	}

	result, err := profileStub.CreateProfile(ctx, arg)
	if err != nil {
		fmt.Println("CreateProfile error, ", err)
		return nil, err
	}

	resultGW := &gatewaypb.ErrorResponse{
		ErrorMessage: result.ErrorMessage,
	}

	return resultGW, nil
}

func (s *server) ReadProfile(ctx context.Context, in *gatewaypb.IdRequest) (*gatewaypb.Profile, error) {
	arg := &profilespb.IdRequest{Id: in.Id}
	profile, err := profileStub.ReadProfile(ctx, arg)
	if err != nil {
		fmt.Println("ReadProfile error:", err)
		return nil, err
	}

	result := &gatewaypb.Profile{
		ID:          profile.ID,
		Name:        profile.Name,
		Age:         profile.Age,
		Description: profile.Description,
		Location:    profile.Location,
	}

	return result, nil
}

func (s *server) UpdateProfile(ctx context.Context, in *gatewaypb.ProfileRequest) (*gatewaypb.ErrorResponse, error) {
	arg := &profilespb.ProfileRequest{}
	arg.Profile = &profilespb.Profile{
		ID:          in.Profile.ID,
		Name:        in.Profile.Name,
		Age:         in.Profile.Age,
		Description: in.Profile.Description,
		Location:    in.Profile.Location,
	}
	result, err := profileStub.UpdateProfile(ctx, arg)
	if err != nil {
		fmt.Println("UpdateProfile error:", err)
		return nil, err
	}

	resultGW := &gatewaypb.ErrorResponse{
		ErrorMessage: result.ErrorMessage,
	}

	return resultGW, nil
}

func (s *server) GetNextProfile(ctx context.Context, in *gatewaypb.IdRequest) (*gatewaypb.Profile, error) {
	selfID := in.GetId()
	fmt.Println(ctx.Deadline())

	nextID, err := matchingStub.GetNextProfile(ctx, &matchingpb.IdReqResp{ID: selfID})
	if err != nil {
		fmt.Println("GetNextProfile error, ", err)
		return nil, err
	}

	profile, err := profileStub.ReadProfile(ctx, &profilespb.IdRequest{
		Id: nextID.GetID(),
	})
	if err != nil {
		fmt.Println("ReadProfile error, ", err)
		return nil, err
	}
	profileGW := &gatewaypb.Profile{
		ID:          profile.ID,
		Name:        profile.Name,
		Age:         profile.Age,
		Description: profile.Description,
		Location:    profile.Location,
	}

	return profileGW, nil
}

func (s *server) Like(ctx context.Context, in *gatewaypb.TargetRequest) (*gatewaypb.ErrorResponse, error) {
	arg := &likespb.TargetRequest{
		Id:    in.Id,
		TgtId: in.TgtId,
	}
	result, err := likeStub.Like(ctx, arg)
	if err != nil {
		fmt.Println("Like error: ", err)
		return nil, err
	}

	resultGW := &gatewaypb.ErrorResponse{ErrorMessage: result.ErrorMessage}

	return resultGW, nil
}

func (s *server) GetLikes(ctx context.Context, in *gatewaypb.IdRequest) (*gatewaypb.LikesResponse, error) {
	arg := &likespb.IdRequest{Id: in.Id}

	result, err := likeStub.GetLikes(ctx, arg)

	if err != nil {
		fmt.Println("GetLikes error, ", err)
		return nil, err
	}

	likesGW := make([]*gatewaypb.Profile, 0)

	for _, like := range result.Likes {
		profile := &gatewaypb.Profile{
			ID:          like.ID,
			Name:        like.Name,
			Age:         like.Age,
			Description: like.Description,
			Location:    like.Location,
		}
		likesGW = append(likesGW, profile)
	}

	resultGW := &gatewaypb.LikesResponse{Likes: likesGW}
	return resultGW, nil
}

func main() {
	matchingConn, err := grpc.Dial("localhost:50051")
	if err != nil {
		log.Fatalln("Cannot connect to matching, ", err)
	}
	matchingStub = matchingpb.NewMatchingServiceClient(matchingConn)

	profileConn, err := grpc.Dial("localhost:50052")
	if err != nil {
		log.Fatalln("Cannot connect to profile, ", err)
	}
	profileStub = profilespb.NewProfileServiceClient(profileConn)

	likeConn, err := grpc.Dial("localhost:50053")
	if err != nil {
		log.Fatalln("Cannot connect to likes, ", err)
	}
	likeStub = likespb.NewLikeServiceClient(likeConn)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 50050))
	if err != nil {
		fmt.Println("Failed to listen:", err)
		return
	}

	s := grpc.NewServer()

	gatewaypb.RegisterProfileServiceServer(s, &server{})
	fmt.Println("Server started at port", 50050)

	if err := s.Serve(lis); err != nil {
		fmt.Println("Failed to serve:", err)
		return
	}
}
