package helper

import (
	"database/sql"
	"time"
)

func NewNullString(s string) sql.NullString {
	if len(s) == 0 {
		return sql.NullString{}
	}
	return sql.NullString{
		String: s,
		Valid:  true,
	}
}

func ParseDateStringToTime(date string) (time.Time, error) {
	layout := "2006-01-02"
	result, err := time.Parse(layout, date)

	if err != nil {
		return time.Now(), err
	}

	return result, nil
}

func ParseDateTimeStringToTime(date string) (time.Time, error) {
	layout := "2006-01-02T15:04:05+07:00"
	result, err := time.Parse(layout, date)

	if err != nil {
		return time.Now(), err
	}

	return result, nil
}

func ParseDateTimeToDateString(dateTime time.Time) string {
	return dateTime.Format("2006-01-02")
}

func ReferString(val string) *string {
	return &val
}

func DereferString(s *string) string {
	if s != nil {
		return *s
	}

	return ""
}

func ReferUint(val uint) *uint {
	return &val
}

func DereferUint(s *uint) uint {
	if s != nil {
		return *s
	}

	return 0
}

func ReferUint32(val uint32) *uint32 {
	return &val
}

func DereferUint32(s *uint32) uint32 {
	if s != nil {
		return *s
	}

	return 0
}

func ReferTime(val time.Time) *time.Time {
	return &val
}

func DereferTime(s *time.Time) time.Time {
	if s != nil {
		return *s
	}

	return time.Now()
}
