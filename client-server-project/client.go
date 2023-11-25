package main

import (
	"context"
	"flag"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"sync"

	"time"

	pb "go-serve/nicolaemariusghergu/proto-files"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	defaultType = "Type Client"

	TIMEOUT_SECONDS = 5
)

var (
	addr = flag.String("addr", "localhost:50051", "the address to connect to")
)

type routeGuideClient struct {
}

func (s *routeGuideClient) RecordRoute(clientName string, stream pb.RouteGuide_RouteChatClient) error {
	var wg sync.WaitGroup
	defer wg.Wait()

	// Send a message to the server
	err := s.sendMessageToServer(clientName, stream)
	if err != nil {
		return err
	}

	// Receive and handle responses from the server
	wg.Add(1)
	go s.receiveServerMessage(stream, &wg)

	return nil
}

func (s *routeGuideClient) sendMessageToServer(clientName string, stream pb.RouteGuide_RouteChatClient) error {
	errSend := stream.Send(&pb.Request{ClientInfo: &pb.ClientInfo{Name: clientName}, Type: defaultType})
	if errSend != nil {
		log.Errorf("An error has occurred while sending CLIENT_REQUEST_MESSAGE message. Err=%v", errSend)
		return errSend
	}
	return nil
}

func (s *routeGuideClient) receiveServerMessage(stream pb.RouteGuide_RouteChatClient, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		response, errResp := stream.Recv()
		if errResp == io.EOF {
			// Server closed the stream
			break
		} else if errResp != nil {
			log.Errorf("Error receiving message from server: %v", errResp)
			break
		}

		log.Infof("Received message from server: %v", response)
	}
}

func main() {
	flag.Parse()
	// Set up a connection to the server.
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewRouteGuideClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT_SECONDS*time.Second)
	defer cancel()

	// Call the RouteChat method to get the stream
	stream, err := c.RouteChat(ctx)
	if err != nil {
		log.Fatalf("Error opening stream: %v", err)
	}

	var wg sync.WaitGroup

	// Multiple clients simulation through goroutines = 10 clients
	for i := 0; i < 10; i++ {
		i := i
		wg.Add(1)
		go func() {
			defer wg.Done()

			client := routeGuideClient{}

			// Call RecordRoute to handle the streaming
			if err := client.RecordRoute(fmt.Sprintf("Client %v", i), stream); err != nil {
				log.Fatalf("Error during streaming: %v", err)
			}
		}()
	}

	// Wait for all goroutines to finish
	wg.Wait()
}
