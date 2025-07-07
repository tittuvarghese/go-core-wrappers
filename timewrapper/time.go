package timewrapper

import (
	"github.com/tittuvarghese/go-core-wrappers/constants"
	"time"
)

type ITime interface {
	GetCurrentTime() time.Time
	GetCurrentTimeStamp() string
	GetTimeDuration(duration int) time.Duration
}

type Time struct {
	ITime
}

func NewTime() ITime {
	return &Time{}

}

func (_time *Time) GetCurrentTime() time.Time {
	return time.Now()
}

func (_time *Time) GetCurrentTimeStamp() string {
	return _time.GetCurrentTime().Format(constants.TimestampFormat)
}

func (_time *Time) GetTimeDuration(duration int) time.Duration {
	return time.Duration(duration)
}
