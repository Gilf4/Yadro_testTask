package utils

import (
	"time"
)

func ToDuration(strTime string) (time.Duration, error) {
	t, err := time.Parse("15:04:05", strTime)
	if err != nil {
		return 0, err
	}

	d := time.Duration(t.Hour())*time.Hour +
		time.Duration(t.Minute())*time.Minute +
		time.Duration(t.Second())*time.Second
	return d, nil
}
