package interfaces

import (
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"time"
)

type PSQL struct {
	conn   sqlx.DB
	logger *zap.Logger
	config map[string]interface{}
}

type PSQLEvent struct {
	UUID               string    `db:"uuid"`
	Header             string    `db:"header"`
	DateTime           time.Time `db:"datetime"`
	Description        string    `db:"description"`
	Owner              string    `db:"owner"`
	MailingDuration    int32     `db:"mailingduration"`
	EventDurationStart time.Time `db:"eventduration_start"`
	EventDurationStop  time.Time `db:"eventduration_stop"`
}

type PSQLChangeEvent struct {
	Event PSQLEvent
	UUID  string
}

func NewPSQL(logger *zap.Logger, config map[string]interface{}) PSQL {

	user := config["user"]
	password := config["password"]
	sslmode := config["sslmode"]
	host := config["host"]
	dbname := config["dbname"]

	logger.Info(fmt.Sprintf("user: %v password: %v sslmode: %v host: %v dbname: %v config: %v", user, password, sslmode, host, dbname, config))

	db, err := sqlx.Connect("postgres", fmt.Sprintf("user=%v password=%v host=%v dbname=%v sslmode=%v", user, password, host, dbname, sslmode))
	if err != nil {
		logger.Error(fmt.Sprintf("Error in connect to DB: %v", err))
	}

	ps := PSQL{
		conn:   *db,
		logger: logger,
		config: config,
	}

	return ps
}

func (db *PSQL) InsertEvent(event PSQLEvent) (bool, error) {
	identifier, err := db.GetEventIdByUUID(event.UUID)
	if err != nil {
		return false, err
	}

	if identifier != 0 {

		return false, errors.New(fmt.Sprintf("Event with UUID %v already exyist in DB", event.UUID))
	}

	cursor := db.conn.MustBegin()
	cursor.MustExec("INSERT INTO public.events (uuid, header, datetime, description, owner, eventduration_start, eventduration_stop, mailingduration) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
		event.UUID, event.Header, event.DateTime, event.Description, event.Owner, event.EventDurationStart, event.EventDurationStop, event.MailingDuration)
	err = cursor.Commit()
	if err != nil {
		return false, err
	}
	return true, nil
}

func (db *PSQL) UpdateEvent(req PSQLChangeEvent) (bool, error) {
	identifier, err := db.GetEventIdByUUID(req.UUID)
	if err != nil {
		return false, err
	}

	if identifier != 0 {
		return false, errors.New(fmt.Sprintf("Event with UUID %v not exist in DB", req.UUID))
	}

	cursor := db.conn.MustBegin()
	cursor.MustExec("UPDATE public.events SET uuid=$1, header=$2, datetime=$3, description=$4, owner=$5, eventduration_start=$6, eventduration_stop=$7, mailingduration=$8 where id = $9",
		req.Event.UUID, req.Event.Header, req.Event.DateTime, req.Event.Description, req.Event.Owner, req.Event.EventDurationStart, req.Event.EventDurationStop, req.Event.MailingDuration, identifier)
	err = cursor.Commit()
	if err != nil {
		return false, err
	}
	return true, nil
}

func (db *PSQL) RemoveEvent(req PSQLChangeEvent) (bool, error) {
	identifier, err := db.GetEventIdByUUID(req.UUID)
	if err != nil {
		return false, err
	}

	if identifier != 0 {
		return false, errors.New(fmt.Sprintf("Event with UUID %v not exist in DB", req.UUID))
	}
	cursor := db.conn.MustBegin()
	cursor.MustExec("DELETE FROM public.events WHERE id = $1;", identifier)
	err = cursor.Commit()
	if err != nil {
		return false, err
	}
	return true, nil

}

func (db *PSQL) GetEventIdByUUID(uuid string) (int, error) {
	var identifier []int
	err := db.conn.Select(&identifier, "SELECT id FROM public.events where uuid = $1", uuid)
	if err != nil {
		db.logger.Error(err.Error())
		return 0, err
	}
	if len(identifier) > 0 {
		return identifier[0], nil
	} else {
		return 0, nil
	}

}

func (db *PSQL) GetEvents(start time.Time, stop time.Time) ([]PSQLEvent, error) {
	var selectResult []PSQLEvent
	err := db.conn.Select(&selectResult, "SELECT uuid, header, datetime, description, owner, eventduration_start, eventduration_stop, mailingduration FROM public.events where datetime >= $1 and datetime <= $2",
		start, stop)
	if err != nil {
		db.logger.Error(err.Error())
		return nil, err
	}
	if len(selectResult) > 0 {

		return selectResult, nil
	} else {
		return nil, nil
	}
}
