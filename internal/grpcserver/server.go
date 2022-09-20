package grpcservera

import (
	"context"
	"fmt"
	u "github.com/AXlIS/go-shortener"
	"github.com/AXlIS/go-shortener/internal/config"
	"github.com/AXlIS/go-shortener/internal/service"
	pb "github.com/AXlIS/go-shortener/proto"
	"github.com/jackc/pgerrcode"
	"github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"log"
	"net"
)

type ShortenServer struct {
	service service.Service
	config  config.Config
	pb.UnimplementedShortenerServer
}

func NewGRPCServer(s service.Service) *ShortenServer {
	return &ShortenServer{
		service: s,
	}
}

func (s *ShortenServer) CreateShorten(ctx context.Context, in *pb.CreateShortenRequest) (*pb.CreateShortenResponse, error) {
	var response pb.CreateShortenResponse
	shortURL, err := s.service.AddURL(in.URL, in.Id)

	if err, ok := err.(*pq.Error); ok {
		if err.Code == pgerrcode.UniqueViolation {
			response.Url = fmt.Sprintf("%s/%s", s.config.BaseURL, shortURL)
			return &response, nil
		}
	}

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	response.Url = fmt.Sprintf("%s/%s", s.config.BaseURL, shortURL)
	return &response, nil
}

func (s *ShortenServer) GetShorten(ctx context.Context, in *pb.GetShortenRequest) (*pb.GetShortenResponse, error) {
	response := pb.GetShortenResponse{}

	url, err := s.service.GetURL(in.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if url == "" {
		return nil, status.Error(codes.Internal, err.Error())
	}

	header := metadata.Pairs("Location", url)
	_ = grpc.SendHeader(ctx, header)
	return &response, nil
}

func (s *ShortenServer) GetPing(ctx context.Context, in *pb.GetPingRequest) (*pb.GetPingResponse, error) {
	ping, err := s.service.Ping()
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	response := pb.GetPingResponse{Active: ping}
	return &response, nil
}

func (s *ShortenServer) CreateShortenBatch(ctx context.Context, in *pb.CreateShortenBatchRequest) (*pb.CreateShortenBatchResponse, error) {
	var (
		input    []*u.ShortenBatchInput
		response pb.CreateShortenBatchResponse
	)

	if len(in.Batches) == 0 {
		return nil, status.Error(codes.InvalidArgument, "the list is empty")
	}

	for _, item := range in.Batches {
		input = append(input, &u.ShortenBatchInput{
			UserID:        item.UserId,
			CorrelationID: item.CorrelationId,
			OriginalURL:   item.OriginalUrl,
			ShortenURL:    item.ShortUrl,
		})
	}

	urls, err := s.service.AddBatchURL(input, in.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	for _, item := range urls {

		response.Urls = append(response.Urls, &pb.BatchURL{ShortUrl: item.ShortURL, CorrelationId: item.CorrelationID})
	}

	return &response, err
}

func (s *ShortenServer) GetAllShortens(ctx context.Context, in *pb.GetAllShortensRequest) (*pb.GetAllShortensResponse, error) {
	var response pb.GetAllShortensResponse

	items, err := s.service.GetAllURLS(in.Id)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	for _, item := range items {
		response.Items = append(response.Items, &pb.URLItem{ShortUrl: item.ShortURL, OriginalUrl: item.OriginalURL, IsDeleted: item.IsDeleted})
	}

	return &response, nil
}

func (s *ShortenServer) DeleteShortens(ctx context.Context, in *pb.DeleteShortensRequest) (*pb.DeleteShortensResponse, error) {
	var response pb.DeleteShortensResponse
	s.service.DeleteURLS(in.Urls, in.Id)

	return &response, nil
}

func (g *ShortenServer) Start() {
	listen, err := net.Listen("tcp", ":3200")
	if err != nil {
		panic(err)
	}

	s := grpc.NewServer()
	pb.RegisterShortenerServer(s, &ShortenServer{})

	if err := s.Serve(listen); err != nil {
		log.Fatal(err)
	}
}
