package scheduler

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"reflect"
	"testing"
	"time"
)

func TestScheduler_Create(t *testing.T) {
	cfg := zap.Config{
		Encoding:    "json",
		Level:       zap.NewAtomicLevelAt(zapcore.DebugLevel),
		OutputPaths: []string{"main.log"},
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:   "message",
			LevelKey:     "level",
			EncodeLevel:  zapcore.CapitalLevelEncoder,
			TimeKey:      "time",
			EncodeTime:   zapcore.ISO8601TimeEncoder,
			CallerKey:    "caller",
			EncodeCaller: zapcore.ShortCallerEncoder,
		},
	}

	logger, _ := cfg.Build()

	event1 := Event{
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

	event2 := Event{
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

	type fields struct {
		Storage map[int]Event
		Logger  *zap.Logger
	}
	type args struct {
		event Event
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name:   "ОК>",
			fields: fields{Storage: map[int]Event{1: event1}, Logger: logger},
			args:   args{event: event1},
			want:   true,
		},
		{
			name:   "Уже в хранилище",
			fields: fields{Storage: map[int]Event{0: event1, 1: event2}, Logger: logger},
			args:   args{event: event1},
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Scheduler{
				Storage: tt.fields.Storage,
				Logger:  tt.fields.Logger,
			}
			if got := s.Create(tt.args.event); got != tt.want {
				t.Errorf("Create() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestScheduler_GetDailyEvents(t *testing.T) {
	cfg := zap.Config{
		Encoding:    "json",
		Level:       zap.NewAtomicLevelAt(zapcore.DebugLevel),
		OutputPaths: []string{"main.log"},
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:   "message",
			LevelKey:     "level",
			EncodeLevel:  zapcore.CapitalLevelEncoder,
			TimeKey:      "time",
			EncodeTime:   zapcore.ISO8601TimeEncoder,
			CallerKey:    "caller",
			EncodeCaller: zapcore.ShortCallerEncoder,
		},
	}

	logger, _ := cfg.Build()

	event1 := Event{
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

	event2 := Event{
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

	type fields struct {
		Storage map[int]Event
		Logger  *zap.Logger
	}
	type args struct {
		date time.Time
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   []Event
	}{
		{
			name:   "один элемент>",
			fields: fields{Storage: map[int]Event{0: event1, 1: event2}, Logger: logger},
			args:   args{date: time.Now()},
			want:   []Event{event1, event2},
		},
		{
			name:   "нет в хранилище",
			fields: fields{Storage: map[int]Event{}, Logger: logger},
			args:   args{date: time.Now()},
			want:   nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Scheduler{
				Storage: tt.fields.Storage,
				Logger:  tt.fields.Logger,
			}
			if got := s.GetDailyEvents(tt.args.date); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetDailyEvents() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestScheduler_GetMonthlyEvents(t *testing.T) {
	cfg := zap.Config{
		Encoding:    "json",
		Level:       zap.NewAtomicLevelAt(zapcore.DebugLevel),
		OutputPaths: []string{"main.log"},
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:   "message",
			LevelKey:     "level",
			EncodeLevel:  zapcore.CapitalLevelEncoder,
			TimeKey:      "time",
			EncodeTime:   zapcore.ISO8601TimeEncoder,
			CallerKey:    "caller",
			EncodeCaller: zapcore.ShortCallerEncoder,
		},
	}

	logger, _ := cfg.Build()

	event1 := Event{
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

	event2 := Event{
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

	type fields struct {
		Storage map[int]Event
		Logger  *zap.Logger
	}
	type args struct {
		dateMonthStart time.Time
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []Event
	}{
		{
			name:   "два элемента",
			fields: fields{Storage: map[int]Event{0: event1, 1: event2}, Logger: logger},
			args:   args{dateMonthStart: time.Date(time.Now().Year(), time.Now().Month(), 1, 0, 0, 0, 0, time.Now().Location())},
			want:   []Event{event1, event2},
		},
		{
			name:   "ни одного элемента",
			fields: fields{Storage: map[int]Event{}, Logger: logger},
			args:   args{dateMonthStart: time.Date(time.Now().Year(), time.Now().Month()-1, 1, 0, 0, 0, 0, time.Now().Location())},
			want:   nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Scheduler{
				Storage: tt.fields.Storage,
				Logger:  tt.fields.Logger,
			}
			if got := s.GetMonthlyEvents(tt.args.dateMonthStart); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetMonthlyEvents() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestScheduler_GetWeeklyEvents(t *testing.T) {
	cfg := zap.Config{
		Encoding:    "json",
		Level:       zap.NewAtomicLevelAt(zapcore.DebugLevel),
		OutputPaths: []string{"main.log"},
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:   "message",
			LevelKey:     "level",
			EncodeLevel:  zapcore.CapitalLevelEncoder,
			TimeKey:      "time",
			EncodeTime:   zapcore.ISO8601TimeEncoder,
			CallerKey:    "caller",
			EncodeCaller: zapcore.ShortCallerEncoder,
		},
	}

	logger, _ := cfg.Build()

	event1 := Event{
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

	event2 := Event{
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

	type fields struct {
		Storage map[int]Event
		Logger  *zap.Logger
	}
	type args struct {
		dateWeekStart time.Time
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []Event
	}{
		{
			name:   "два элемента",
			fields: fields{Storage: map[int]Event{0: event1, 1: event2}, Logger: logger},
			args:   args{dateWeekStart: time.Date(time.Now().Year(), time.Now().Month(), 26, 0, 0, 0, 0, time.Now().Location())},
			want:   []Event{event1, event2},
		},
		{
			name:   "ни одного элемента",
			fields: fields{Storage: map[int]Event{}, Logger: logger},
			args:   args{dateWeekStart: time.Date(time.Now().Year(), time.Now().Month(), 1, 0, 0, 0, 0, time.Now().Location())},
			want:   nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Scheduler{
				Storage: tt.fields.Storage,
				Logger:  tt.fields.Logger,
			}
			if got := s.GetWeeklyEvents(tt.args.dateWeekStart); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetWeeklyEvents() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestScheduler_Remove(t *testing.T) {
	cfg := zap.Config{
		Encoding:    "json",
		Level:       zap.NewAtomicLevelAt(zapcore.DebugLevel),
		OutputPaths: []string{"main.log"},
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:   "message",
			LevelKey:     "level",
			EncodeLevel:  zapcore.CapitalLevelEncoder,
			TimeKey:      "time",
			EncodeTime:   zapcore.ISO8601TimeEncoder,
			CallerKey:    "caller",
			EncodeCaller: zapcore.ShortCallerEncoder,
		},
	}

	logger, _ := cfg.Build()

	event1 := Event{
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

	event2 := Event{
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

	type fields struct {
		Storage map[int]Event
		Logger  *zap.Logger
	}
	type args struct {
		id int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name:   "Смог удалить",
			fields: fields{Storage: map[int]Event{0: event1, 1: event2}, Logger: logger},
			args:   args{id: 0},
			want:   true,
		},
		{
			name:   "нет такого элемента",
			fields: fields{Storage: map[int]Event{0: event1, 1: event2}, Logger: logger},
			args:   args{id: 25},
			want:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Scheduler{
				Storage: tt.fields.Storage,
				Logger:  tt.fields.Logger,
			}
			if got := s.Remove(tt.args.id); got != tt.want {
				t.Errorf("Remove() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestScheduler_Update(t *testing.T) {
	cfg := zap.Config{
		Encoding:    "json",
		Level:       zap.NewAtomicLevelAt(zapcore.DebugLevel),
		OutputPaths: []string{"main.log"},
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:   "message",
			LevelKey:     "level",
			EncodeLevel:  zapcore.CapitalLevelEncoder,
			TimeKey:      "time",
			EncodeTime:   zapcore.ISO8601TimeEncoder,
			CallerKey:    "caller",
			EncodeCaller: zapcore.ShortCallerEncoder,
		},
	}

	logger, _ := cfg.Build()

	event1 := Event{
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

	event2 := Event{
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

	type fields struct {
		Storage map[int]Event
		Logger  *zap.Logger
	}
	type args struct {
		id    int
		event Event
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name:   "Смог обновить",
			fields: fields{Storage: map[int]Event{0: event1, 1: event2}, Logger: logger},
			args:   args{id: 0, event: event2},
			want:   true,
		},
		{
			name:   "Не смог обновить",
			fields: fields{Storage: map[int]Event{0: event1, 1: event2}, Logger: logger},
			args:   args{id: 25, event: event2},
			want:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Scheduler{
				Storage: tt.fields.Storage,
				Logger:  tt.fields.Logger,
			}
			if got := s.Update(tt.args.id, tt.args.event); got != tt.want {
				t.Errorf("Update() = %v, want %v", got, tt.want)
			}
		})
	}
}
