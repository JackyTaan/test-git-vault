package main

import (
	"context"
	"io"
	"log"
	"net"

	"github.com/jackytaan/go-grpc-examples/stream/bi-directional-streaming/feeds/feedpb"

	"google.golang.org/grpc"

	"github.com/redis/go-redis/v9"
)

type server struct{}

func main() {

	lis, err := net.Listen("tcp", "localhost:50051")
	if err != nil {
		log.Fatalf("could not listen: %v", err)
	} else {
		log.Println("Server is listening on", "localhost:50051")
	}

	s := grpc.NewServer()
	feedpb.RegisterFeedsServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("could not start the server: %v", err)
	}
}

// Broadcast reads client stream and broadcasts recieved feeds
func (*server) Broadcast(stream feedpb.Feeds_BroadcastServer) error {
	////init db redis
	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	_ = rdb.FlushDB(ctx).Err()
	pipe := rdb.Pipeline()

	//////////////////////////
	//////////////////////////
	for {
		//Receiving process
		msg, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatalf("could not recieve from stream : %v", err)
			return err
		}
		feed := msg.GetFeed()
		///set redis
		pipe.Set(ctx, "TK:"+feed, feed, 0)

		_, err := pipe.Exec(ctx)
		if err != nil {
			panic(err)
		}
		/////////////////
		//Sending process
		stream.Send(&feedpb.FeedResponse{Feed: feed})

	}

}
