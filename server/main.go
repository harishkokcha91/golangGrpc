package main

import (
	"context"
	"log"
	"math/rand"
	"net"
	"strconv"

	"google.golang.org/grpc"
	pb "moviesapp.com/grpc/protos"
)

const (
	port = ":50051"
)

var movies []*pb.MovieInfo

type movieServer struct {
	pb.UnimplementedMovieServer
}

func main() {
	initMovies()
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen : %v", err)
	}

	s := grpc.NewServer()

	pb.RegisterMovieServer(s, &movieServer{})

	log.Printf("Server listening at %v", lis.Addr())

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

}

func initMovies() {
	movie1 := &pb.MovieInfo{Id: "1", Title: "The batman", Director: &pb.Director{Firstname: "Harish", Lastname: "Kokcha"}}
	movie2 := &pb.MovieInfo{Id: "2", Title: "The team", Director: &pb.Director{Firstname: "Kapil", Lastname: "Kokcha"}}

	movies = append(movies, movie1, movie2)
}

func (s *movieServer) GetMovies(in *pb.Empty, stream pb.Movie_GetMoviesServer) error {
	log.Printf("Received: %v ", in)

	for _, movie := range movies {
		if err := stream.Send(movie); err != nil {
			return err
		}
	}
	return nil
}

func (s *movieServer) GetMovie(ctx context.Context, in *pb.Id) (*pb.MovieInfo, error) {
	log.Printf("Received: %v", in)

	res := &pb.MovieInfo{}

	for _, movie := range movies {
		if movie.GetId() == in.GetValue() {
			res = movie
			break
		}
	}

	return res, nil
}

func (s *movieServer) CreatMovie(ctx context.Context, in *pb.MovieInfo) (*pb.Id, error) {
	log.Printf("Received: %v", in)
	res := pb.Id{}
	res.Value = strconv.Itoa(rand.Intn(1000000000))
	in.Id = res.GetValue()
	movies = append(movies, in)

	return &res, nil
}

func (s *movieServer) UpdateMoview(ctx context.Context, in *pb.MovieInfo) (*pb.Status, error) {
	log.Printf("Received: %v", in)
	res := pb.Status{}

	for index, movie := range movies {
		if movie.GetId() == in.GetId() {
			movies = append(movies[:index], movies[index+1:]...)
			in.Id = movie.GetId()
			movies = append(movies, in)
			res.Value = 1
			break
		}
	}
	return &res, nil
}

func (s *movieServer) DeleteMovie(ctx context.Context, in *pb.Id) (*pb.Status, error) {
	log.Printf("Received: %v", in)
	res := pb.Status{}
	for index, movie := range movies {
		if movie.GetId() == in.GetValue() {
			movies = append(movies[:index], movies[index+1:]...)
			res.Value = 1
			break
		}
	}
	return &res, nil
}
