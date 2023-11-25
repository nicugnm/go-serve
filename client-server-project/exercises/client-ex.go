package exercises

import (
	log "github.com/sirupsen/logrus"
	pb "go-serve/nicolaemariusghergu/proto-files"
)

type RouteGuideClientStrategy interface {
	RecordRoute(clientName string, stream pb.RouteGuide_RouteChatClient) error
}

type ExerciseStrategy struct {
	ExerciseNumber int32
	Factory        ExerciseFactory
}

func (s ExerciseStrategy) SendExercise(clientName string, stream pb.RouteGuide_RouteChatClient) error {
	inputArray := s.Factory.CreateInputArray(s.ExerciseNumber)
	errSend := stream.Send(&pb.Request{
		ClientInfo:     &pb.ClientInfo{Name: clientName, ExerciseNumber: s.ExerciseNumber},
		ExerciseNumber: s.ExerciseNumber,
		StringArray:    inputArray,
	})
	if errSend != nil {
		log.Errorf("An error has occurred while sending CLIENT_REQUEST_MESSAGE message. Err=%v", errSend)
		return errSend
	}
	return nil
}

type ExerciseFactory struct{}

func (f *ExerciseFactory) CreateInputArray(exerciseNumber int32) []string {
	switch exerciseNumber {
	case 1:
		return []string{"casa", "masa", "trei", "tanc", "4321"}
	case 2:
		return []string{"abd4g5", "1sdf6fd", "fd2fdsf5"}
	case 3:
		return []string{"12", "13", "14"}
	case 4:
		return []string{"2", "10", "5", "11", "39", "32", "80", "84"}
	case 5:
		return []string{"2dasdas", "12", "dasdas", "1010", "101"}
	case 6:
		return []string{"abcdef"}
	case 7:
		return []string{"1G11o1L"}
	case 8:
		return []string{"23", "17", "15", "3", "18"}
	case 9:
		return []string{"mama", "iris", "bunica", "ala"}
	case 10:
		return []string{"24", "16", "8", "aaa", "bbb"}
	default:
		log.Errorf("An error has been occurred while establishing inputArray. No input array for the exercise.")
		return []string{}
	}
}
