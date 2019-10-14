package structs

import "time"

type Event struct {
	UUID               string    `db:"uuid" json:"uuid"`                                //ID события
	Header             string    `db:"header" json:"header"`                            //заголовок события
	DateTime           time.Time `db:"datetime" json:"date_time"`                       //дата и время события
	Description        string    `db:"description" json:"description"`                  //описание
	Owner              string    `db:"owner" json:"owner"`                              //владелец события
	MailingDuration    int32     `db:"mailingduration" json:"mailing_duration"`         //за сколько нужно выслать оповещение (в минутах)
	EventDurationStart time.Time `db:"eventduration_start" json:"event_duration_start"` //длительность события начало
	EventDurationStop  time.Time `db:"eventduration_stop" json:"event_duration_stop"`   //длительность события конец
}
