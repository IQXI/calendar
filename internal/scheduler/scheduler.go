package scheduler

import (
	"fmt"
	"go.uber.org/zap"
	"time"
)

type Event struct {
	UUID          int
	Header        string
	DateTime      time.Time
	EventDuration struct {
		Start time.Time
		End   time.Time
	}
	Description     string
	Owner           string
	MailingDuration int
}

type Scheduler struct {
	Storage map[int]Event
	Logger  *zap.Logger
}

func (s *Scheduler) Create(event Event) bool {
	if _, ok := s.Storage[event.UUID]; ok {
		s.Logger.Error(fmt.Sprintf("Event with ID %v already in Storage", event.UUID))
		return false
	}

	s.Storage[event.UUID] = event
	return true
}

func (s *Scheduler) Update(id int, event Event) bool {
	if _, ok := s.Storage[id]; ok {
		s.Storage[id] = event
		return true
	}
	s.Logger.Error(fmt.Sprintf("Event with ID %v not in Storage", id))
	return false

}

func (s *Scheduler) Remove(id int) bool {
	if _, ok := s.Storage[id]; ok {
		delete(s.Storage, id)
		return true
	}
	s.Logger.Error(fmt.Sprintf("Event with ID %v not in Storage", id))
	return false

}

func (s *Scheduler) GetDailyEvents(date time.Time) []Event {
	dateDayStart := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	dateDayEnd := dateDayStart.AddDate(0, 0, 1)

	var events []Event
	for _, val := range s.Storage {
		if (val.DateTime.After(dateDayStart) || val.DateTime.Equal(dateDayStart)) && (val.DateTime.Before(dateDayEnd) || val.DateTime.Equal(dateDayEnd)) {
			events = append(events, val)
		}
	}
	return events
}

func (s *Scheduler) GetWeeklyEvents(dateWeekStart time.Time) []Event {
	dateWeekStart = time.Date(dateWeekStart.Year(), dateWeekStart.Month(), dateWeekStart.Day(), 0, 0, 0, 0, dateWeekStart.Location())
	dateWeekEnd := dateWeekStart.AddDate(0, 0, 7)

	var events []Event
	for _, val := range s.Storage {
		if (val.DateTime.After(dateWeekStart) || val.DateTime.Equal(dateWeekStart)) && (val.DateTime.Before(dateWeekEnd) || val.DateTime.Equal(dateWeekEnd)) {
			events = append(events, val)
		}
	}
	return events
}

func (s *Scheduler) GetMonthlyEvents(dateMonthStart time.Time) []Event {
	dateMonthStart = time.Date(dateMonthStart.Year(), dateMonthStart.Month(), dateMonthStart.Day(), 0, 0, 0, 0, dateMonthStart.Location())
	dateMonthEnd := dateMonthStart.AddDate(0, 1, 0)

	var events []Event
	for _, val := range s.Storage {
		if (val.DateTime.After(dateMonthStart) || val.DateTime.Equal(dateMonthStart)) && (val.DateTime.Before(dateMonthEnd) || val.DateTime.Equal(dateMonthEnd)) {
			events = append(events, val)
		}
	}
	return events
}
