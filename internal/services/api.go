package services

import (
	"calendar/internal/interfaces/postgres"
	pb "calendar/internal/proto"
	"calendar/internal/structs"
	"context"
	"github.com/golang/protobuf/ptypes"
	"go.uber.org/zap"
	"time"
)

type API struct {
	psql   postgres.PSQL
	Logger *zap.Logger
}

func NewAPI(logger *zap.Logger, psql postgres.PSQL) *API {
	sch := API{
		psql,
		logger,
	}
	return &sch
}

//Функции мутации типов

func PBEventToPSQLEvent(event *pb.Event) (structs.Event, error) {

	dt, err := ptypes.Timestamp(event.DateTime)
	if err != nil {
		return structs.Event{
			UUID:               "",
			Header:             "",
			DateTime:           time.Time{},
			Description:        "",
			Owner:              "",
			MailingDuration:    0,
			EventDurationStart: time.Time{},
			EventDurationStop:  time.Time{},
		}, err
	}
	dtStart, err := ptypes.Timestamp(event.EventDuration.Start)
	if err != nil {
		return structs.Event{
			UUID:               "",
			Header:             "",
			DateTime:           time.Time{},
			Description:        "",
			Owner:              "",
			MailingDuration:    0,
			EventDurationStart: time.Time{},
			EventDurationStop:  time.Time{},
		}, err
	}
	dtStop, err := ptypes.Timestamp(event.EventDuration.Stop)
	if err != nil {
		return structs.Event{
			UUID:               "",
			Header:             "",
			DateTime:           time.Time{},
			Description:        "",
			Owner:              "",
			MailingDuration:    0,
			EventDurationStart: time.Time{},
			EventDurationStop:  time.Time{},
		}, err
	}

	psqlEvent := structs.Event{
		UUID:               event.UUID,
		Header:             event.Header,
		DateTime:           time.Date(dt.Year(), dt.Month(), dt.Day(), dt.Hour(), dt.Minute(), dt.Second(), dt.Nanosecond(), dt.Location()),
		Description:        event.Description,
		Owner:              event.Owner,
		MailingDuration:    event.MailingDuration,
		EventDurationStart: time.Date(dtStart.Year(), dtStart.Month(), dtStart.Day(), dtStart.Hour(), dtStart.Minute(), dtStart.Second(), dtStart.Nanosecond(), dtStart.Location()),
		EventDurationStop:  time.Date(dtStop.Year(), dtStop.Month(), dtStop.Day(), dtStop.Hour(), dtStop.Minute(), dtStop.Second(), dtStop.Nanosecond(), dtStop.Location()),
	}

	return psqlEvent, nil
}

func PSQLEventToPBEvent(event structs.Event) (*pb.Event, error) {
	dt, err := ptypes.TimestampProto(event.DateTime)
	if err != nil {
		return nil, err
	}
	dtStart, err := ptypes.TimestampProto(event.EventDurationStart)
	if err != nil {
		return nil, err
	}
	dtStop, err := ptypes.TimestampProto(event.EventDurationStop)
	if err != nil {
		return nil, err
	}

	pbEvent := pb.Event{
		UUID:            event.UUID,
		Header:          event.Header,
		DateTime:        dt,
		Description:     event.Description,
		Owner:           event.Owner,
		MailingDuration: event.MailingDuration,
		EventDuration:   &pb.EventDuration{Start: dtStart, Stop: dtStop},
	}

	return &pbEvent, nil
}

func PBChangeRequestToPSQLChangeRequest(chgRequest *pb.ChangeEventRequest) (postgres.PSQLChangeEvent, error) {

	psqlEvent, err := PBEventToPSQLEvent(chgRequest.Event)
	if err != nil {
		return postgres.PSQLChangeEvent{
				Event: structs.Event{},
				UUID:  "",
			},
			err
	}

	psqlChangeRequest := postgres.PSQLChangeEvent{
		Event: psqlEvent,
		UUID:  chgRequest.Id,
	}

	return psqlChangeRequest, nil
}

func PSQLEventsToPBEventList(events []structs.Event) (*pb.EventList, error) {

	pbEventList := pb.EventList{
		Events: make([]*pb.Event, 0),
	}

	for _, event := range events {

		psqlEvent, err := PSQLEventToPBEvent(event)
		if err != nil {
			return &pb.EventList{}, err
		}

		pbEventList.Events = append(pbEventList.Events, psqlEvent)
	}

	return &pbEventList, nil
}

//INSERT, UPDATE, DELETE

func (s *API) InsertEvent(ctx context.Context, event *pb.Event) (*pb.ChangeEventResult, error) {

	psqlEvent, err := PBEventToPSQLEvent(event)
	if err != nil {
		return &pb.ChangeEventResult{Error: err.Error(), Result: false}, nil
	}

	result, err := s.psql.InsertEvent(psqlEvent)
	if err != nil {
		return &pb.ChangeEventResult{Error: err.Error(), Result: false}, nil
	}

	if result == false {
		return &pb.ChangeEventResult{Error: "Unknown error", Result: false}, nil
	} else {
		return &pb.ChangeEventResult{Error: "nil", Result: true}, nil
	}
}

func (s *API) UpdateEvent(ctx context.Context, req *pb.ChangeEventRequest) (*pb.ChangeEventResult, error) {

	psqlChangeRequest, err := PBChangeRequestToPSQLChangeRequest(req)
	if err != nil {
		return &pb.ChangeEventResult{Error: err.Error(), Result: false}, nil
	}

	result, err := s.psql.UpdateEvent(psqlChangeRequest)
	if err != nil {
		return &pb.ChangeEventResult{Error: err.Error(), Result: false}, nil
	}

	if result == false {
		return &pb.ChangeEventResult{Error: "Unknown error", Result: false}, nil
	} else {
		return &pb.ChangeEventResult{Error: "nil", Result: true}, nil
	}

}

func (s *API) RemoveEvent(ctx context.Context, req *pb.ChangeEventRequest) (*pb.ChangeEventResult, error) {

	psqlChangeRequest, err := PBChangeRequestToPSQLChangeRequest(req)
	if err != nil {
		return &pb.ChangeEventResult{Error: err.Error(), Result: false}, nil
	}

	result, err := s.psql.RemoveEvent(psqlChangeRequest)
	if err != nil {
		return &pb.ChangeEventResult{Error: err.Error(), Result: false}, nil
	}

	if result == false {
		return &pb.ChangeEventResult{Error: "Unknown error", Result: false}, nil
	} else {
		return &pb.ChangeEventResult{Error: "nil", Result: true}, nil
	}
}

//GET methods

func (s *API) GetDailyEvents(ctx context.Context, req *pb.GetRequest) (*pb.GetResult, error) {

	new_date, _ := ptypes.Timestamp(req.DateTime)
	dateDayStart := time.Date(new_date.Year(), new_date.Month(), new_date.Day(), 0, 0, 0, 0, new_date.Location())
	dateDayEnd := dateDayStart.AddDate(0, 0, 1)

	psqlEvents, err := s.psql.GetEvents(dateDayStart, dateDayEnd)
	if err != nil {
		return &pb.GetResult{
			Error:  err.Error(),
			Events: nil,
		}, nil
	}

	pbEventList, err := PSQLEventsToPBEventList(psqlEvents)
	if err != nil {
		return &pb.GetResult{
			Error:  err.Error(),
			Events: nil,
		}, nil
	}

	return &pb.GetResult{Error: "nil", Events: pbEventList}, nil
}

func (s *API) GetWeeklyEvents(ctx context.Context, req *pb.GetRequest) (*pb.GetResult, error) {

	new_date, _ := ptypes.Timestamp(req.DateTime)
	dateWeekStart := time.Date(new_date.Year(), new_date.Month(), new_date.Day(), 0, 0, 0, 0, new_date.Location())
	dateWeekEnd := dateWeekStart.AddDate(0, 0, 7)

	psqlEvents, err := s.psql.GetEvents(dateWeekStart, dateWeekEnd)
	if err != nil {
		return &pb.GetResult{
			Error:  err.Error(),
			Events: nil,
		}, nil
	}

	pbEventList, err := PSQLEventsToPBEventList(psqlEvents)
	if err != nil {
		return &pb.GetResult{
			Error:  err.Error(),
			Events: nil,
		}, nil
	}

	return &pb.GetResult{Error: "nil", Events: pbEventList}, nil
}

func (s *API) GetMonthlyEvents(ctx context.Context, req *pb.GetRequest) (*pb.GetResult, error) {

	new_date, _ := ptypes.Timestamp(req.DateTime)
	dateMonthStart := time.Date(new_date.Year(), new_date.Month(), new_date.Day(), 0, 0, 0, 0, new_date.Location())
	dateMonthEnd := dateMonthStart.AddDate(0, 1, 0)

	psqlEvents, err := s.psql.GetEvents(dateMonthStart, dateMonthEnd)
	if err != nil {
		return &pb.GetResult{
			Error:  err.Error(),
			Events: nil,
		}, nil
	}

	pbEventList, err := PSQLEventsToPBEventList(psqlEvents)
	if err != nil {
		return &pb.GetResult{
			Error:  err.Error(),
			Events: nil,
		}, nil
	}

	return &pb.GetResult{Error: "nil", Events: pbEventList}, nil
}
