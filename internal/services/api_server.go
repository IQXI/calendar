package services

import (
	pb "calendar/internal/proto"
	"context"
	"fmt"
	"github.com/golang/protobuf/ptypes"
	"go.uber.org/zap"
	"time"
)

type Scheduler struct {
	Storage map[int32]pb.Event
	Logger  *zap.Logger
}

func NewScheduler(logger *zap.Logger) *Scheduler {
	sch := Scheduler{
		make(map[int32]pb.Event),
		logger,
	}
	return &sch
}

func (s *Scheduler) InsertEvent(ctx context.Context, event *pb.Event) (*pb.ChangeEventResult, error) {
	if _, ok := s.Storage[event.UUID]; ok {
		s.Logger.Error(fmt.Sprintf("Event with ID %v already in Storage", event.UUID))
		return &pb.ChangeEventResult{Error: fmt.Sprintf("Event with ID %v already in Storage", event.UUID), Result: false}, nil
	}

	s.Storage[event.UUID] = *event
	return &pb.ChangeEventResult{Error: "nil", Result: true}, nil
	//return nil, nil
}

func (s *Scheduler) UpdateEvent(ctx context.Context, req *pb.UpdateRequest) (*pb.ChangeEventResult, error) {
	if _, ok := s.Storage[req.Id]; ok {
		s.Storage[req.Id] = *req.Event
		return &pb.ChangeEventResult{Error: "nil", Result: true}, nil
	}
	s.Logger.Error(fmt.Sprintf("Event with ID %v not in Storage", req.Id))
	return &pb.ChangeEventResult{Error: fmt.Sprintf("Event with ID %v not in Storage", req.Id), Result: false}, nil

}

func (s *Scheduler) RemoveEvent(ctx context.Context, req *pb.RemoveRequest) (*pb.ChangeEventResult, error) {
	if _, ok := s.Storage[req.Id]; ok {
		delete(s.Storage, req.Id)
		return &pb.ChangeEventResult{Error: "nil", Result: true}, nil
	}
	s.Logger.Error(fmt.Sprintf("Event with ID %v not in Storage", req.Id))
	return &pb.ChangeEventResult{Error: fmt.Sprintf("Event with ID %v not in Storage", req.Id), Result: false}, nil
}

func (s *Scheduler) GetDailyEvents(ctx context.Context, req *pb.GetRequest) (*pb.GetResult, error) {

	new_date, _ := ptypes.Timestamp(req.DateTime)
	dateDayStart := time.Date(new_date.Year(), new_date.Month(), new_date.Day(), 0, 0, 0, 0, new_date.Location())
	dateDayEnd := dateDayStart.AddDate(0, 0, 1)

	events := pb.GetResult{Error: "", Events: &pb.EventList{Events: []*pb.Event{}}}
	for _, val := range s.Storage {
		valDateTime, _ := ptypes.Timestamp(val.DateTime)
		if (valDateTime.After(dateDayStart) || valDateTime.Equal(dateDayStart)) && (valDateTime.Before(dateDayEnd) || valDateTime.Equal(dateDayEnd)) {
			events.Events.Events = append(events.Events.Events, &val)
		}
	}
	return &events, nil
}

func (s *Scheduler) GetWeeklyEvents(ctx context.Context, req *pb.GetRequest) (*pb.GetResult, error) {

	new_date, _ := ptypes.Timestamp(req.DateTime)
	dateWeekStart := time.Date(new_date.Year(), new_date.Month(), new_date.Day(), 0, 0, 0, 0, new_date.Location())
	dateWeekEnd := dateWeekStart.AddDate(0, 0, 7)

	events := pb.GetResult{Error: "", Events: &pb.EventList{Events: []*pb.Event{}}}
	for _, val := range s.Storage {
		valDateTime, _ := ptypes.Timestamp(val.DateTime)
		if (valDateTime.After(dateWeekStart) || valDateTime.Equal(dateWeekStart)) && (valDateTime.Before(dateWeekEnd) || valDateTime.Equal(dateWeekEnd)) {
			events.Events.Events = append(events.Events.Events, &val)
		}
	}
	return &events, nil
}

func (s *Scheduler) GetMonthlyEvents(ctx context.Context, req *pb.GetRequest) (*pb.GetResult, error) {

	new_date, _ := ptypes.Timestamp(req.DateTime)
	dateMonthStart := time.Date(new_date.Year(), new_date.Month(), new_date.Day(), 0, 0, 0, 0, new_date.Location())
	dateMonthEnd := dateMonthStart.AddDate(0, 1, 0)

	events := pb.GetResult{Error: "", Events: &pb.EventList{Events: []*pb.Event{}}}
	for _, val := range s.Storage {
		valDateTime, _ := ptypes.Timestamp(val.DateTime)
		if (valDateTime.After(dateMonthStart) || valDateTime.Equal(dateMonthStart)) && (valDateTime.Before(dateMonthEnd) || valDateTime.Equal(dateMonthEnd)) {
			events.Events.Events = append(events.Events.Events, &val)
		}
	}
	return &events, nil
}
