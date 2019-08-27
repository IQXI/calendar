package main

import (
	cfg "calendar/internal/config"
	lg "calendar/internal/logger"
	"calendar/internal/scheduler"
	"time"
)

func main() {
	logger := lg.GetLogger(cfg.GetConfig())

	logger.Info("Service started!")
	sch := scheduler.Scheduler{
		make(map[int]scheduler.Event),
		logger,
	}

	event1 := scheduler.Event{
		UUID:     0,
		Header:   "Событие 1",
		DateTime: time.Now(),
		EventDuration: struct {
			Start time.Time
			End   time.Time
		}{
			time.Now(),
			time.Now().AddDate(0, 0, 1),
		},
		Description:     "Описание события 1",
		Owner:           "vcherkozyanov",
		MailingDuration: 30,
	}

	event2 := scheduler.Event{
		UUID:     1,
		Header:   "Событие 2",
		DateTime: time.Now(),
		EventDuration: struct {
			Start time.Time
			End   time.Time
		}{
			time.Now(),
			time.Now().AddDate(0, 0, 1),
		},
		Description:     "Описание события 2",
		Owner:           "vcherkozyanov",
		MailingDuration: 30,
	}

	if sch.Create(event1) {
		logger.Info("Event1 is created")
	}

	if sch.Update(0, event2) {
		logger.Info("Event1 is updates")
	}

	if sch.Remove(0) {
		logger.Info("Event1 is removed")
	}

	if sch.Remove(0) {
		logger.Info("Event1 is removed")
	}

	if sch.Create(event2) {
		logger.Info("Event2 is created")
	}

	for _, val := range sch.GetDailyEvents(time.Now()) {
		print("GetDailyEvents", val.UUID, "\n")
	}

	for _, val := range sch.GetWeeklyEvents(time.Date(time.Now().Year(), time.Now().Month(), 1, 0, 0, 0, 0, time.Now().Location())) {
		print("GetWeeklyEvents", val.UUID, "\n")
	}

	for _, val := range sch.GetMonthlyEvents(time.Date(time.Now().Year(), time.Now().Month()+1, 1, 0, 0, 0, 0, time.Now().Location())) {
		print("GetMonthlyEvents", val.UUID, "\n")
	}

}
