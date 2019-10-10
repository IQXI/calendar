package main

import (
	pb "calendar/internal/proto"
	"context"
	"fmt"
	"github.com/golang/protobuf/ptypes"
	"google.golang.org/grpc"
	"log"
	"time"
)

func main() {
	var conn *grpc.ClientConn
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

	defer conn.Close()
	if err != nil {
		log.Fatalf(err.Error())
	}

	client := pb.NewAPIClient(conn)

	tm, _ := ptypes.TimestampProto(time.Now())

	event := pb.Event{
		UUID:            "1",
		Header:          "1",
		DateTime:        tm,
		Description:     "2",
		Owner:           "3",
		MailingDuration: 0,
		EventDuration: &pb.EventDuration{
			Start: tm,
			Stop:  tm,
		},
	}

	result, err := client.InsertEvent(context.Background(), &event)
	fmt.Print(err)
	if err != nil {
		log.Fatalf("Некая ошибка %v", err.Error())
	}

	log.Print(result)
}
