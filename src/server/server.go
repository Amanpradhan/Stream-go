package main

import (
	"io"
	"log"
	"math/rand"
	"net"
	"os"
	"os/signal"
	"strings"
	"time"

	pb "github.com/Amanpradhan/Stream-go/src/proto"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedAgentServer
}

const (
	interval_time = 2 * time.Second
)

func generateString() string {
	// generate random strings of 2 word from list of words
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

func (s server) Communicate(srv pb.Agent_CommunicateServer) error {
	log.Println("starting new server")
	ctx := srv.Context()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		req, err := srv.Recv()
		if err == io.EOF {
			log.Println("exit")
			return nil
		}
		if err != nil {
			log.Printf("got error %v", err)
			continue
		}
		recvStr := req.Message
		if strings.Contains(recvStr, "hello") {
			log.Println(recvStr)
		}
		ticker := time.NewTicker(interval_time)
		quit := make(chan struct{})
		go func() {
			for {
				select {
				case <-ticker.C:
					greetRes := pb.Response{Result: greet()}
					if err := srv.Send(&greetRes); err != nil {
						log.Fatalf("got error %v", err)
					}
					// uncomment below line to see all messages sent by server to client
					//fmt.Printf("message sent successfully %s to client\n", greetRes.Result)

				case <-quit:
					ticker.Stop()
					return
				}
			}
		}()
		<-time.After(time.Duration(interval_time))
		// generating random string below
		resp := pb.Response{Result: generateString()}
		if err := srv.Send(&resp); err != nil {
			log.Fatalf("got error %v", err)
		}
		// uncomment below line to see all messages sent by server to client
		//fmt.Printf("message sent successfully %s to client\n", resp.Result)
	}
}

func main() {
	// create a listener
	lis, err := net.Listen("tcp", ":50005")
	if err != nil {
		log.Fatalf("got error while trrying to listen: %v", err)
	}
	// create a grpc server
	s := grpc.NewServer()
	pb.RegisterAgentServer(s, &server{})

	// start the server...
	if err := s.Serve(lis); err != nil {
		log.Fatalf("got error while starting the server: %v", err)
	}
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	sig := <-quit
	log.Printf("Abrupt Server Shut Down. Reason: %v", sig)
	s.GracefulStop()
	log.Println("Server Stopped Gracefully")
}
