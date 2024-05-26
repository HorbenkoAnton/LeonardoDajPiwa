package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"google.golang.org/grpc"
	"log/slog"
	"net"
	"net/http"
	"os"
	el "profiles/lib/env"
	"profiles/lib/logger"
	pm "profiles/migrations"
	pb "profiles/proto"
	"strconv"
	"strings"
	"time"
)

var (
	pg *pgxpool.Pool
)

const Timeout = 10 * time.Second

type server struct {
	logger *slog.Logger
	pb.UnimplementedProfileServiceServer
}

type CityResponse struct {
	DisplayName string `json:"display_name"`
}

func GetCity(city string, logger *slog.Logger) ([]string, error) {
	url := fmt.Sprintf("https://geocode.maps.co/search?q=%v&api_key=%v",
		city,
		el.LoadEnvVar("GEO_API_KEY"),
	)

	response, err := http.Get(url)
	if err != nil {
		logger.Error("failed to request", slog.String("url", url), slog.String("error", err.Error()))
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		logger.Error("incorrect response status code", slog.Int("status_code", response.StatusCode))
		return nil, fmt.Errorf("incorrect status code: %d", response.StatusCode)
	}

	var cityResponse []*CityResponse
	if err := json.NewDecoder(response.Body).Decode(&cityResponse); err != nil {
		logger.Error("failed to decode body", slog.String("error", err.Error()))
		return nil, err
	}

	displayNames := strings.Split(cityResponse[0].DisplayName, ", ")

	if len(displayNames) <= 0 {
		logger.Error("display names length = 0")
		return nil, err
	}

	return displayNames, nil
}

func (s server) CreateProfile(_ context.Context, request *pb.ProfileRequest) (*pb.ErrorResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), Timeout)
	defer cancel()

	if request.Profile == nil {
		s.logger.Error("profile request is nil")
		return &pb.ErrorResponse{ErrorMessage: "profile request is nil"}, nil
	}

	responseCity, err := GetCity(request.Profile.Location, s.logger)
	if responseCity == nil {
		s.logger.Error("failed to parse city", slog.String("error", err.Error()))
		return &pb.ErrorResponse{ErrorMessage: "city retrieval error"}, err
	}

	if _, err := pg.Exec(ctx, "INSERT INTO profiles (id, name, age, description, pfp_id, user_location, location) VALUES ($1, $2, $3, $4, $5, $6, $7)",
		request.Profile.ID,
		request.Profile.Name,
		request.Profile.Age,
		request.Profile.Description,
		request.Profile.Pfp,
		request.Profile.Location,
		responseCity,
	); err != nil {
		s.logger.Error("failed to insert profile", slog.Int64("id", request.Profile.ID), slog.String("error", err.Error()))
		return &pb.ErrorResponse{ErrorMessage: err.Error()}, err
	}

	return &pb.ErrorResponse{ErrorMessage: "OK"}, nil
}

func (s server) ReadProfile(_ context.Context, request *pb.IdRequest) (*pb.Profile, error) {
	ctx, cancel := context.WithTimeout(context.Background(), Timeout)
	defer cancel()

	var profile pb.Profile
	if err := pg.QueryRow(ctx, "SELECT id, name, age, description, pfp_id, user_location FROM profiles WHERE id=$1", request.Id).Scan(
		&profile.ID,
		&profile.Name,
		&profile.Age,
		&profile.Description,
		&profile.Pfp,
		&profile.Location,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.logger.Error("profile not found", slog.Int64("id", request.GetId()), slog.String("error", err.Error()))
			return nil, errors.New("profile not found")
		}
		s.logger.Error("failed to select profile", slog.Int64("id", request.GetId()), slog.String("error", err.Error()))
		return nil, err
	}

	return &profile, nil
}

func (s server) UpdateProfile(_ context.Context, request *pb.ProfileRequest) (*pb.ErrorResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), Timeout)
	defer cancel()

	responseCity, err := GetCity(request.Profile.Location, s.logger)
	if responseCity == nil {
		s.logger.Error("failed to parse city", slog.String("error", err.Error()))
		return &pb.ErrorResponse{ErrorMessage: "city retrieval error"}, err
	}

	if _, err = pg.Exec(ctx, "UPDATE profiles SET name = $1, age = $2, description = $3, pfp_id=$4, user_location = $5, location = $6  WHERE id=$7",
		request.Profile.Name,
		request.Profile.Age,
		request.Profile.Description,
		request.Profile.Pfp,
		request.Profile.Location,
		responseCity,
		request.Profile.ID,
	); err != nil {
		s.logger.Error("failed to update profile", slog.String("error", err.Error()))
		return &pb.ErrorResponse{ErrorMessage: "error, profile wasn't updated"}, err
	}

	return &pb.ErrorResponse{ErrorMessage: "OK"}, nil
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
	db := stdlib.OpenDBFromPool(pg)

	reload, err := strconv.ParseBool(el.LoadEnvVar("DB_RELOAD"))
	if err != nil {
		setupLogger.Error("failed to parse bool env var", slog.Bool("reload_bool", reload), slog.String("error", err.Error()))
		os.Exit(1)
	}

	pm.Migrate(reload, db)

	srv := grpc.NewServer()
	pb.RegisterProfileServiceServer(srv, &server{logger: setupLogger})

	profilesPort := el.LoadEnvVar("PROFILES_PORT")

	lis, err := net.Listen("tcp", ":"+profilesPort)
	if err != nil {
		setupLogger.Error("failed to listen tcp", slog.String("profiles_port", profilesPort), slog.String("error", err.Error()))
		os.Exit(1)
	}

	setupLogger.Info("profiles server started")

	if err = srv.Serve(lis); err != nil {
		setupLogger.Error("failed to serve", slog.String("error", err.Error()))
		os.Exit(1)
	}
}
