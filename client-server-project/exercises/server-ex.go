package exercises

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"math"
	"strconv"
	"strings"
	"unicode"

	pb "go-serve/nicolaemariusghergu/proto-files"
)

const (
	SERVER_PROCESSED_DATA_MESSAGE                           = "Server has procesed the data. Response: %v"
	ERROR_WHILE_SERVER_SENDING_PROCESED_DATA                = "An error has occurred while sending SERVER_PROCESSED_DATA_MESSAGE. Err=%v"
	ERROR_WHILE_PROCESSING_DATA_MESSAGE                     = "An error has occurred while processing data. Err=%v"
	ERROR_WHILE_SENDING_ERROR_WHILE_PROCESSING_DATA_MESSAGE = "An error has occurred while sending ERROR_WHILE_PROCESSING_DATA_MESSAGE message. Err=%v"
)

type ServerExercises struct {
}

func (s *ServerExercises) HandleExercise1(stream pb.RouteGuide_RouteChatServer, request *pb.Request) {
	inputStrings := request.GetStringArray()

	if len(inputStrings) == 0 {
		SendErrorProcessingDataMessage(stream, request, "Empty input array")
		return
	}

	var outputString string

	for i := range inputStrings[0] {
		for j := range inputStrings {
			outputString += string(inputStrings[j][i])
		}
	}

	SendProcessedDataMessage(stream, request, outputString)
}

func (s *ServerExercises) HandleExercise2(stream pb.RouteGuide_RouteChatServer, request *pb.Request) {
	inputStrings := request.GetStringArray()

	if len(inputStrings) == 0 {
		SendErrorProcessingDataMessage(stream, request, "Empty input array")
		return
	}

	countPerfectSquares := 0

	for _, str := range inputStrings {
		num, err := strconv.Atoi(str)
		if err == nil {
			squareRoot := int(math.Sqrt(float64(num)))
			if squareRoot*squareRoot == num {
				countPerfectSquares++
			}
		}
	}

	SendProcessedDataMessage(stream, request, fmt.Sprintf("Number of perfect squares: %d", countPerfectSquares))
}

func (s *ServerExercises) HandleExercise3(stream pb.RouteGuide_RouteChatServer, request *pb.Request) {
	inputStrings := request.GetStringArray()

	inputNumbers, err := convertStringArrayToIntArray(inputStrings)
	if err != nil {
		SendErrorProcessingDataMessage(stream, request, "Invalid input format for exercise 3")
		return
	}

	sum := calculateSumOfReversedNumbers(inputNumbers)

	SendProcessedDataMessage(stream, request, fmt.Sprintf("Sum of reversed numbers: %d", sum))
}

func (s *ServerExercises) HandleExercise4(stream pb.RouteGuide_RouteChatServer, request *pb.Request) {
	inputStrings := request.GetStringArray()

	inputNumbers, err := convertStringArrayToIntArray(inputStrings)
	if err != nil {
		SendErrorProcessingDataMessage(stream, request, "Invalid input format for exercise 4")
		return
	}

	if len(inputNumbers) < 3 {
		SendErrorProcessingDataMessage(stream, request, "Insufficient input for exercise 4")
		return
	}

	lowerBound := inputNumbers[0]
	upperBound := inputNumbers[1]
	numbers := inputNumbers[2:]

	average, err := calculateAverageOfNumbersInRange(numbers, lowerBound, upperBound)
	if err != nil {
		SendErrorProcessingDataMessage(stream, request, err.Error())
		return
	}

	SendProcessedDataMessage(stream, request, fmt.Sprintf("Average of numbers with sum of digits in range: %.2f", average))
}

func (s *ServerExercises) HandleExercise5(stream pb.RouteGuide_RouteChatServer, request *pb.Request) {
	inputStrings := request.GetStringArray()

	var binaryNumbers []string

	for _, str := range inputStrings {
		if isBinary(str) {
			decimalValue, _ := strconv.ParseInt(str, 2, 64)
			binaryNumbers = append(binaryNumbers, fmt.Sprintf("%d", decimalValue))
		}
	}

	SendProcessedDataMessage(stream, request, fmt.Sprintf("Converted binary numbers: %v", binaryNumbers))
}

func (s *ServerExercises) HandleExercise6(stream pb.RouteGuide_RouteChatServer, request *pb.Request) {
	inputStrings := request.GetStringArray()

	var encryptedStrings []string

	for _, str := range inputStrings {
		encryptedStrings = append(encryptedStrings, caesarCipher(str, 3)) // Using LEFT shift with 3 characters
	}

	SendProcessedDataMessage(stream, request, fmt.Sprintf("Encrypted strings: %v", encryptedStrings))
}

func (s *ServerExercises) HandleExercise7(stream pb.RouteGuide_RouteChatServer, request *pb.Request) {
	encodedText := request.GetStringArray()[0]

	decodedText := decodeText(encodedText)

	SendProcessedDataMessage(stream, request, fmt.Sprintf("Decoded text: %s", decodedText))
}

func (s *ServerExercises) HandleExercise8(stream pb.RouteGuide_RouteChatServer, request *pb.Request) {
	strings := request.GetStringArray()

	count := 0

	for _, str := range strings {
		num, err := strconv.Atoi(str)
		if err != nil {
			log.Errorf("Error converting string to integer: %v", err)
			continue
		}

		if isPrime(num) {
			count += countDigits(num)
		}
	}

	SendProcessedDataMessage(stream, request, fmt.Sprintf("Total number of digits in prime numbers: %d", count))
}

func (s *ServerExercises) HandleExercise9(stream pb.RouteGuide_RouteChatServer, request *pb.Request) {
	inputStrings := request.GetStringArray()

	count := 0

	for _, str := range inputStrings {
		if hasEvenVowelsOnEvenPositions(str) {
			count++
		}
	}

	SendProcessedDataMessage(stream, request, fmt.Sprintf("Number of words with even vowels on even positions: %d", count))
}

func (s *ServerExercises) HandleExercise10(stream pb.RouteGuide_RouteChatServer, request *pb.Request) {
	inputStrings := request.GetStringArray()

	gcdResult := calculateGCDForAllNumbers(inputStrings)

	SendProcessedDataMessage(stream, request, fmt.Sprintf("GCD for all numbers: %d", gcdResult))
}

func SendErrorProcessingDataMessage(stream pb.RouteGuide_RouteChatServer, request *pb.Request, message string) {
	errSend := stream.Send(&pb.Response{
		ClientInfo: request.ClientInfo,
		Response:   fmt.Sprintf(ERROR_WHILE_PROCESSING_DATA_MESSAGE, message),
	})
	if errSend != nil {
		log.Errorf(ERROR_WHILE_SENDING_ERROR_WHILE_PROCESSING_DATA_MESSAGE, errSend)
	}
}

func SendProcessedDataMessage(stream pb.RouteGuide_RouteChatServer, request *pb.Request, resultString string) {
	errSend := stream.Send(&pb.Response{
		ClientInfo:     request.ClientInfo,
		Response:       fmt.Sprintf(SERVER_PROCESSED_DATA_MESSAGE, resultString),
		ExerciseNumber: request.ExerciseNumber,
		ResultArray:    []string{resultString},
	})
	if errSend != nil {
		log.Errorf(ERROR_WHILE_SERVER_SENDING_PROCESED_DATA, errSend)
	}
}

func convertStringArrayToIntArray(inputStrings []string) ([]int, error) {
	var result []int
	for _, str := range inputStrings {
		num, err := strconv.Atoi(str)
		if err != nil {
			return nil, err
		}
		result = append(result, num)
	}
	return result, nil
}

func calculateSumOfReversedNumbers(numbers []int) int {
	sum := 0
	for _, num := range numbers {
		reversedNum, _ := reverseNumber(num)
		sum += reversedNum
	}
	return sum
}

func calculateAverageOfNumbersInRange(numbers []int, lowerBound, upperBound int) (float64, error) {
	count := 0
	sum := 0

	for _, num := range numbers {
		if isSumOfDigitsInRange(num, lowerBound, upperBound) {
			count++
			sum += num
		}
	}

	if count == 0 {
		return 0, fmt.Errorf("No numbers found in the specified range")
	}

	average := float64(sum) / float64(count)
	return average, nil
}

func reverseNumber(num int) (int, error) {
	strNum := strconv.Itoa(num)
	reversedStrNum, err := reverseString(strNum)
	if err != nil {
		return 0, err
	}

	reversedNum, err := strconv.Atoi(reversedStrNum)
	if err != nil {
		return 0, err
	}

	return reversedNum, nil
}

func reverseString(str string) (string, error) {
	runes := []rune(str)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes), nil
}

func isSumOfDigitsInRange(num, lowerBound, upperBound int) bool {
	sum := sumOfDigits(num)
	return sum >= lowerBound && sum <= upperBound
}

func sumOfDigits(num int) int {
	sum := 0
	for num > 0 {
		sum += num % 10
		num /= 10
	}
	return sum
}

func isBinary(str string) bool {
	_, err := strconv.ParseInt(str, 2, 64)
	return err == nil
}

func caesarCipher(str string, shift int) string {
	var result strings.Builder

	for _, char := range str {
		if 'a' <= char && char <= 'z' {
			result.WriteRune((char-'a'+rune(shift))%26 + 'a')
		} else if 'A' <= char && char <= 'Z' {
			result.WriteRune((char-'A'+rune(shift))%26 + 'A')
		} else {
			result.WriteRune(char)
		}
	}

	return result.String()
}

func decodeText(encodedText string) string {
	var result strings.Builder

	for i := 0; i < len(encodedText); i++ {
		char := rune(encodedText[i])
		if unicode.IsDigit(char) {
			count := int(char - '0')
			i++
			for j := 0; j < count && i < len(encodedText); j++ {
				result.WriteRune(rune(encodedText[i]))
			}
		} else {
			result.WriteRune(char)
		}
	}

	return result.String()
}

func isPrime(num int) bool {
	if num <= 1 {
		return false
	}
	for i := 2; i*i <= num; i++ {
		if num%i == 0 {
			return false
		}
	}
	return true
}

func countDigits(num int) int {
	count := 0
	for num > 0 {
		num /= 10
		count++
	}
	return count
}

func hasEvenVowelsOnEvenPositions(str string) bool {
	vowels := "aeiouAEIOU"
	count := 0

	for i, char := range str {
		if i%2 == 0 && strings.ContainsRune(vowels, char) {
			count++
		}
	}

	return count%2 == 0
}

func calculateGCDForAllNumbers(numbers []string) int {
	intNumbers := make([]int, len(numbers))
	for i, str := range numbers {
		intNum, _ := strconv.Atoi(str)
		intNumbers[i] = intNum
	}

	if len(intNumbers) == 0 {
		return 0
	}

	gcdResult := intNumbers[0]
	for _, num := range intNumbers[1:] {
		gcdResult = gcd(gcdResult, num)
	}

	return gcdResult
}

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}
