package main

import (
	pb "calendar/internal/proto"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
)

func main() {
	var conn *grpc.ClientConn
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

	defer conn.Close()
	if err != nil {
		log.Fatalf(err.Error())
	}

	client := pb.NewAPIClient(conn)

	event := pb.Event{
		UUID:            1,
		Header:          "1",
		DateTime:        nil,
		Description:     "2",
		Owner:           "3",
		MailingDuration: 0,
		EventDuration:   nil,
	}

	result, err := client.InsertEvent(context.Background(), &event)
	fmt.Print(err)
	if err != nil {
		log.Fatalf("Некая ошибка %v", err.Error())
	}

	log.Print(result)
}
