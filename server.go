package main

import (
	"context"
	"flag"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"net"

	pb "go-serve/nicolaemariusghergu/proto-files"
	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

const (
	SERVER_RECEIVED_MESSAGE          = "Server has received the request."
	RESPONSE_SERVER_RECEIVED_MESSAGE = "Data received: %v"
)

type server struct {
	pb.UnimplementedRouteGuideServer
}

func (s *server) RouteChat(stream pb.RouteGuide_RouteChatServer) error {
	for {
		request, errRecv := stream.Recv()

		if errRecv == io.EOF {
			// Client closed the connection, exit the loop
			break
		} else if errRecv != nil {
			log.Errorf("An error has occurred while receiving message. Err=%v", errRecv)
			break
		}

		errSend := stream.Send(&pb.Response{ClientInfo: request.ClientInfo, Response: fmt.Sprintf(RESPONSE_SERVER_RECEIVED_MESSAGE, request.Type)})
		if errSend != nil {
			log.Errorf("An error has occurred while sending SERVER_RECEIVED_MESSAGE message. Err=%v", errSend)
			break
		}
	}

	return nil
}

func (s *server) GetFeature(ctx context.Context, in *pb.Request) (*pb.Response, error) {
	return nil, nil
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("Couldn't start the server. Err= %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterRouteGuideServer(s, &server{})
	log.Printf("Server started at= %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Ops! An error has occured while serve. Err= %v", err)
	}
}
