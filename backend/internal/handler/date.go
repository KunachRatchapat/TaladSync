package handler

import (
	"time"

	"taladsync/internal/pkg/timeutil"
)

func parseMarketDate(raw string) (time.Time, error) {
	if raw == "" {
		return timeutil.Today(), nil
	}
	t, err := time.Parse("2006-01-02", raw)
	if err != nil {
		return time.Time{}, err
	}
	return timeutil.MarketDate(t), nil
}
