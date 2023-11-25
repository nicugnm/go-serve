package main

import (
	"flag"
	"fmt"
	"go-serve/nicolaemariusghergu/exercises"
	"io"
	"net"
	"sync"

	log "github.com/sirupsen/logrus"
	pb "go-serve/nicolaemariusghergu/proto-files"
	"google.golang.org/grpc"
)

var (
	port             = flag.Int("port", 50051, "The server port")
	MAX_SIZE_MESSAGE = 2952998688
)

const (
	SERVER_RECEIVED_MESSAGE                   = "Client %v has made a request with data: %v"
	RESPONSE_SERVER_RECEIVED_MESSAGE          = "Server has received the request. Data received: %v"
	SERVER_PROCESSING_DATA_MESSAGE            = "Server is processing the data"
	ERROR_WHILE_RECEIVING_MESSAGE             = "An error has occurred while receiving message. Err=%v"
	ERROR_SEND_SERVER_RECEIVED_MESSAGE        = "An error has occurred while sending SERVER_RECEIVED_MESSAGE message. Err=%v"
	ERROR_SEND_SERVER_PROCESSING_DATA_MESSAGE = "An error has occurred while sending SERVER_PROCESSING_DATA_MESSAGE message. Err=%v"
)

type server struct {
	pb.UnimplementedRouteGuideServer
}

func (s *server) RouteChat(stream pb.RouteGuide_RouteChatServer) error {
	var wg sync.WaitGroup
	defer wg.Wait()

	// Goroutine for sending responses
	wg.Add(1)
	go s.sendResponses(stream, &wg)

	// Goroutine for processing incoming messages
	wg.Add(1)
	go s.processIncomeMessages(stream, &wg)

	return nil
}

func (s *server) processIncomeMessages(stream pb.RouteGuide_RouteChatServer, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		request, errRecv := stream.Recv()
		if errRecv == io.EOF {
			// Client closed the connection, exit the loop
			break
		} else if errRecv != nil {
			log.Errorf(ERROR_WHILE_RECEIVING_MESSAGE, errRecv)
			break
		}

		log.Infof(SERVER_RECEIVED_MESSAGE, request.ClientInfo)
		sendResponseServerReceivedMessage(stream, request)
		sendProcessingDataMessage(stream, request)

		ExercisesServer := &exercises.ServerExercises{}

		switch request.ExerciseNumber {
		case 1:
			ExercisesServer.HandleExercise1(stream, request)
		case 2:
			ExercisesServer.HandleExercise2(stream, request)
		case 3:
			ExercisesServer.HandleExercise3(stream, request)
		case 4:
			ExercisesServer.HandleExercise4(stream, request)
		case 5:
			ExercisesServer.HandleExercise5(stream, request)
		case 6:
			ExercisesServer.HandleExercise6(stream, request)
		case 7:
			ExercisesServer.HandleExercise7(stream, request)
		case 8:
			ExercisesServer.HandleExercise8(stream, request)
		case 9:
			ExercisesServer.HandleExercise9(stream, request)
		case 10:
			ExercisesServer.HandleExercise10(stream, request)
		default:
			log.Errorf("Exercise not implemented: %v", request.ExerciseNumber)
		}
	}
}

func (s *server) sendResponses(stream pb.RouteGuide_RouteChatServer, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		if err := checkStream(stream); err != nil {
			log.Errorf("Error checking stream: %v", err)
			return
		}

		select {
		case <-stream.Context().Done():
			// Client closed the connection, exit the goroutine
			return
		default:
			// Continue processing
		}
	}
}

func sendResponseServerReceivedMessage(stream pb.RouteGuide_RouteChatServer, request *pb.Request) {
	errSend := stream.Send(&pb.Response{
		ClientInfo: request.ClientInfo,
		Response:   fmt.Sprintf(RESPONSE_SERVER_RECEIVED_MESSAGE, request.StringArray),
	})
	if errSend != nil {
		log.Errorf(ERROR_SEND_SERVER_RECEIVED_MESSAGE, errSend)
	}
}

func sendProcessingDataMessage(stream pb.RouteGuide_RouteChatServer, request *pb.Request) {
	errSend := stream.Send(&pb.Response{
		ClientInfo: request.ClientInfo,
		Response:   fmt.Sprintf(SERVER_PROCESSING_DATA_MESSAGE),
	})
	if errSend != nil {
		log.Errorf(ERROR_SEND_SERVER_PROCESSING_DATA_MESSAGE, errSend)
	}
}

// checkStream is a utility function to check if the stream is nil or closed.
func checkStream(stream pb.RouteGuide_RouteChatServer) error {
	if stream == nil {
		return fmt.Errorf("stream is nil")
	}
	if stream.Context() == nil || stream.Context().Err() != nil {
		return fmt.Errorf("stream context is nil or closed")
	}
	return nil
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("Couldn't start the server. Err= %v", err)
	}
	s := grpc.NewServer(
		grpc.MaxRecvMsgSize(MAX_SIZE_MESSAGE),
		grpc.MaxSendMsgSize(MAX_SIZE_MESSAGE))
	pb.RegisterRouteGuideServer(s, &server{})
	log.Printf("Server started at= %v", lis.Addr())

	// Defer stopping the server until the main function exits
	defer s.Stop()

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Ops! An error has occurred while serving. Err= %v", err)
	}
}
