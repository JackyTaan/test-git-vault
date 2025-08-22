package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"strconv"
	"time"

	"github.com/jackytaan/go-grpc-examples/stream/bi-directional-streaming/feeds/feedpb"

	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}

	c := feedpb.NewFeedsClient(conn)
	//  get client stream
	stream, err := c.Broadcast(context.Background())
	if err != nil {
		log.Fatalf("failed to call Broadcast: %v", err)
	}

	//record time begin sending to streaming server
	startTime := time.Now()
	fmt.Println("start at:", startTime)
	// make blocking channel
	waitc := make(chan struct{})

	// send feeds to the stream ( go routine )
	go func() {
		for i := 1; i <= 50000; i++ {
			feed := strconv.Itoa(i)
			// fmt.Println("Client send: ", feed)
			if err := stream.Send(&feedpb.FeedRequest{Feed: feed}); err != nil {
				log.Fatalf("error while sending feed: %v", err)
			}
			//time.Sleep(time.Second)
		}
		if err := stream.CloseSend(); err != nil {
			log.Fatalf("failed to close stream: %v", err)
		}
	}()

	// recieve feeds frrom the stream ( go routine )
	go func() {
		for {
			//msg, err := stream.Recv()
			_, err := stream.Recv()
			if err == io.EOF {
				close(waitc)
				return
			}
			if err != nil {
				log.Fatalf("failed to recieve: %v", err)
				close(waitc)
				return
			}

			//fmt.Println("Client recieved: ", msg.GetFeed())
		}

	}()

	//
	<-waitc
	//
	endTime := time.Now()
	fmt.Println("end at:", endTime)
	duration := time.Now().Sub(startTime).Seconds()
	fmt.Println("total:", duration)
	//
}
