package main

import (
	"context"
	pb "github.com/Amanpradhan/Stream-go/src/proto"
	"google.golang.org/grpc"
	"io"
	"log"
	"math/rand"
	"strings"
	"time"
)

const (
	interval_time = 2 * time.Second
	num_request   = 100
)

func generateString() string {
	var dummy = []string{"hello", "sun", "world", "space", "moon", "crypto", "sky", "ocean", "universe", "human"}
	rnd1 := int32(rand.Intn(len(dummy)))
	rnd2 := int32(rand.Intn(len(dummy)))
	return dummy[rnd1] + " " + dummy[rnd2]
}

func greet() string {
	t := time.Now()
	var st string
	switch {
	case t.Hour() < 12:
		st = "Good morning!"
	case t.Hour() < 17:
		st = "Good afternoon!"
	case t.Hour() < 19:
		st = "Good evening!"
	default:
		st = "Good night!"
	}
	return st
}

func main() {
	rand.Seed(time.Now().Unix())
	// start connection by dialing through grpc
	conn, err := grpc.Dial(":50005", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("got error while trying to connect: %v", err)
	}
	// create stream
	client := pb.NewAgentClient(conn)
	stream, err := client.Communicate(context.Background())
	if err != nil {
		log.Fatalf("got error while opening the stream %v", err)
	}

	ctx := stream.Context()
	done := make(chan bool)

	// started a goroutine which will send 100 requests to client
	go func() {
		for i := 1; i <= num_request; i++ { // sending the request for specified times
			ticker := time.NewTicker(interval_time)
			quit := make(chan struct{})
			go func() {
				for {
					select {
					case <-ticker.C:
						m := greet()
						req := pb.Request{Message: m}
						if err := stream.Send(&req); err != nil {
							log.Fatalf("got error while sending %v", err)
						}
						// uncomment below line to see all messages sent by server to client
						//fmt.Printf("message sent successfully %s to server\n", req.Message)

					case <-quit:
						ticker.Stop()
						return
					}
				}
			}()
			<-time.After(time.Duration(interval_time))
			m := generateString()
			req := pb.Request{Message: m}
			if err := stream.Send(&req); err != nil {
				log.Fatalf("got error while sending %v", err)
			}
			// uncomment line below to check the message that client sent
			//log.Printf("random message sent to server %s", req.Message)
		}
		if err := stream.CloseSend(); err != nil {
			log.Println(err)
		}
	}()

	// second goroutine receives data from stream
	// if stream is finished it closes done channel
	go func() {
		for {
			resp, err := stream.Recv()
			if err == io.EOF {
				close(done)
				return
			}
			if err != nil {
				log.Fatalf("can not receive %v", err)
			}
			recvMsg := resp.Result
			if strings.Contains(recvMsg, "hello") {
				log.Println("Message containing string hello: ", recvMsg)
			}
		}
	}()

	// third goroutine closes done channel
	// if context is done
	go func() {
		<-ctx.Done()
		if err := ctx.Err(); err != nil {
			log.Println(err)
		}
		close(done)
	}()

	<-done
}
