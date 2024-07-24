package utils

import (
	"database/sql/driver"
	"fmt"
	"strconv"
	"time"
)

type Timestamp struct {
	time.Time
}

func NowTimestamp() *Timestamp {
	return &Timestamp{time.Now()}
}

func (t Timestamp) MarshalJSON() ([]byte, error) {
	return ([]byte)(strconv.FormatInt(t.Time.Unix(), 10)), nil
}

func (t *Timestamp) UnmarshalJSON(data []byte) error {
	if stamp, err := strconv.ParseInt(string(data), 10, 64); err != nil {
		return err
	} else {
		t.Time = time.Unix(stamp, 0)
		return nil
	}
}

// Value timestamp into sql need this function
func (t Timestamp) Value() (driver.Value, error) {
	var zeroTime time.Time
	if t.Time.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return t.Time, nil
}

// Scan value of time.Time implements the mysql.Scanner interface.
func (t *Timestamp) Scan(v interface{}) error {
	if value, ok := v.(time.Time); ok {
		t.Time = value
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)
}
