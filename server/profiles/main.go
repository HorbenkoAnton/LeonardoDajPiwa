package main

import (
	"context"
	"database/sql"
	pb "dp_mock/backend/profiles/proto"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
	"os"
	"time"
)

var (
	pg *pgxpool.Pool
)

const Timeout = 10 * time.Second

type server struct {
	pb.UnimplementedProfileServiceServer
}

type CityResponse struct {
	DisplayName string `json:"display_name"`
}

func GetCity(city string) ([]string, error) {
	if err := godotenv.Load(".env"); err != nil {
		return nil, err
	}
	url := os.Getenv("GEO_API_URL")
	url = fmt.Sprintf(url, city)

	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("incorrect status code: %s", response.Status)
	}

	var cityResponse []CityResponse
	if err := json.NewDecoder(response.Body).Decode(&cityResponse); err != nil {
		return nil, err
	}

	displayNames := make([]string, len(cityResponse))
	for i, city := range cityResponse {
		displayNames[i] = city.DisplayName
	}

	if len(displayNames) <= 0 {
		return nil, err
	}

	return displayNames, nil
}

func (s server) CreateProfile(_ context.Context, request *pb.ProfileRequest) (*pb.ErrorResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), Timeout)
	defer cancel()

	if request.Profile == nil {
		return &pb.ErrorResponse{ErrorMessage: "profile request is nil"}, nil
	}

	responseCity, err := GetCity(request.Profile.Location)
	if responseCity == nil {
		return &pb.ErrorResponse{ErrorMessage: "error getting city"}, err
	}

	if _, err := pg.Exec(ctx, "INSERT INTO profiles (id, name, age, description, user_location, location) VALUES ($1, $2, $3, $4, $5, $6)",
		request.Profile.ID,
		request.Profile.Name,
		request.Profile.Age,
		request.Profile.Description,
		request.Profile.Location,
		responseCity,
	); err != nil {
		return &pb.ErrorResponse{ErrorMessage: err.Error()}, err
	}

	return &pb.ErrorResponse{ErrorMessage: "OK"}, nil
}

func (s server) ReadProfile(_ context.Context, request *pb.IdRequest) (*pb.Profile, error) {
	ctx, cancel := context.WithTimeout(context.Background(), Timeout)
	defer cancel()

	var profile pb.Profile
	if err := pg.QueryRow(ctx, "SELECT id, name, age, description, user_location FROM profiles WHERE id=$1", request.Id).Scan(
		&profile.ID,
		&profile.Name,
		&profile.Age,
		&profile.Description,
		&profile.Location,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("profile not found")
		}
		return nil, err
	}

	return &profile, nil
}

func (s server) UpdateProfile(_ context.Context, request *pb.ProfileRequest) (*pb.ErrorResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), Timeout)
	defer cancel()

	responseCity, err := GetCity(request.Profile.Location)
	if responseCity == nil {
		return &pb.ErrorResponse{ErrorMessage: "error in parsing city"}, err
	}

	if _, err = pg.Exec(ctx, "UPDATE profiles SET name = $1, age = $2, description = $3, user_location = $4, location = $5  WHERE id=$6",
		request.Profile.Name,
		request.Profile.Age,
		request.Profile.Description,
		request.Profile.Location,
		responseCity,
		request.Profile.ID,
	); err != nil {
		return &pb.ErrorResponse{ErrorMessage: "error, profile wasn't updated"}, err
	}

	return &pb.ErrorResponse{ErrorMessage: "OK"}, nil
}

func main() {
	pgconn, err := pgxpool.New(context.Background(), "postgres://postgres:password@localhost:5432/db_profiles")
	if err != nil {
		log.Fatal(err)
	}
	pg = pgconn

	srv := grpc.NewServer()
	pb.RegisterProfileServiceServer(srv, &server{})

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Server is running on port :50051")

	if err = srv.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
