package main

import (
	"context"
	"io"
	"log"
	"time"

	"google.golang.org/grpc"
	pb "moviesapp.com/grpc/protos"
)

const (
	address = "localhost:50051"
)

func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())

	if err != nil {
		log.Fatalf("did not connect : %v ", err)
	}

	defer conn.Close()

	client := pb.NewMovieClient(conn)

	// runGetMovies(client)
	// runGetMovie(client, "1")
	runCreateMovie(client, "928374903", "the a team", "Kajal", "Verma")
	// runUpdateMovie(client, "2", "928374903", "The A Team", "Kajal", "Verma")
	// runDeleteMovie(client, "2")
}

func runGetMovies(client pb.MovieClient) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req := &pb.Empty{}

	stream, err := client.GetMovies(ctx, req)
	if err != nil {
		log.Fatalf("%v.GetMovies() =_,%v", client, err)

	}

	for {
		row, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("%v.GetMovies() =_,%v", client, err)
		}

		log.Printf("MovieInfo : %v", row)
	}
}

func runGetMovie(client pb.MovieClient, movieid string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)

	defer cancel()
	req := &pb.Id{Value: movieid}
	res, err := client.GetMovie(ctx, req)
	if err != nil {
		log.Fatalf("%v.GetMovies() =_,%v", client, err)
	}
	log.Printf("MovieInfo : %v", res)
}

func runCreateMovie(client pb.MovieClient, isbn string, title string, firstname string, lastname string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req := &pb.MovieInfo{Isbn: isbn, Title: title, Director: &pb.Director{Firstname: firstname, Lastname: lastname}}

	res, err := client.CreatMovie(ctx, req)

	if err != nil {
		log.Fatalf("%v.GetMovies() =_,%v", client, err)
	}

	if res.GetValue() != "" {
		log.Printf("CreateMovie Id %v", res)
	} else {
		log.Printf("CreateMovie failed")
	}

}

func runUpdateMovie(client pb.MovieClient, movieid string, isbn string, title string, firstname string, lastname string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req := &pb.MovieInfo{Id: movieid, Isbn: isbn, Title: title, Director: &pb.Director{Firstname: firstname, Lastname: lastname}}

	res, err := client.UpdateMoview(ctx, req)

	if err != nil {
		log.Fatalf("%v.UpdateMovies() =_,%v", client, err)
	}

	if int(res.GetValue()) == 1 {
		log.Printf("UpdateMovie success")
	} else {
		log.Printf("UpdateMovie failed")
	}

}

func runDeleteMovie(client pb.MovieClient, movieid string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)

	defer cancel()

	req := &pb.Id{Value: movieid}

	res, err := client.DeleteMovie(ctx, req)

	if err != nil {
		log.Fatalf("%v.DeleteMovies() =_,%v", client, err)
	}

	if int(res.GetValue()) == 1 {
		log.Printf("DMovie success")
	} else {
		log.Printf("DMovie failed")
	}
}
