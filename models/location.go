package models

import (
	"database/sql/driver"
	"fmt"
	"strconv"
	"strings"
)

type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

func (l *Location) DecodeText(src string) error {
	if len(src) < 5 {
		return fmt.Errorf("invalid length for location: %v", len(src))
	}

	parts := strings.SplitN(string(src[1:len(src)-1]), ",", 2)
	if len(parts) < 2 {
		return fmt.Errorf("invalid format for location")
	}

	latitude, err := strconv.ParseFloat(parts[0], 64)
	if err != nil {
		return err
	}

	longitude, err := strconv.ParseFloat(parts[1], 64)
	if err != nil {
		return err
	}

	*l = Location{Latitude: latitude, Longitude: longitude}
	return nil
}

// Scan implements the database/sql Scanner interface.
func (l *Location) Scan(src interface{}) error {
	if src == nil {
		return fmt.Errorf("location cannot be nil")
	}

	if l == nil {
		return fmt.Errorf("nil pointer receiver")
	}

	switch src := src.(type) {
	case string:
		return l.DecodeText(src)
	default:
		return fmt.Errorf("cannot scan %T", src)
	}
}

// Value implements the database/sql/driver Valuer interface.
func (l Location) Value() (driver.Value, error) {
	return fmt.Sprintf("POINT(%v,%v)", l.Latitude, l.Longitude), nil
}
