package dynamics

import (
	"errors"
	"time"
)

type Interval string

const (
	IntervalDay   Interval = "day"
	IntervalWeek  Interval = "week"
	IntervalMonth Interval = "month"
)

const (
	lastMonth  = 30
	lastSeason = 3

	season    = 90 * 24 * time.Hour
	twoYears  = 2 * 365 * 24 * time.Hour
	fiveYears = 5 * 365 * 24 * time.Hour
)

var (
	ErrInvalidPeriod   = errors.New("invalid period")
	ErrInvalidInterval = errors.New("invalid interval")
	ErrPeriodTooLarge  = errors.New("period too large")
)

type Period struct {
	From     *time.Time
	To       *time.Time
	Interval Interval
}

func (p *Period) Normalize() {
	now := time.Now().UTC()

	if p.Interval == "" {
		p.Interval = IntervalDay
	}

	if p.To == nil {
		p.To = &now
	}

	if p.From == nil {
		var from time.Time

		switch p.Interval {
		case IntervalDay:
			from = p.To.AddDate(0, 0, -lastMonth)

		case IntervalWeek:
			from = p.To.AddDate(0, -lastSeason, 0)

		case IntervalMonth:
			from = p.To.AddDate(-1, 0, 0)

		default:
			from = p.To.AddDate(0, 0, -lastMonth)
		}

		p.From = &from
	}

	from := normalizeTime(*p.From, p.Interval)
	to := normalizeTime(*p.To, p.Interval)

	p.From = &from
	p.To = &to
}

func (p *Period) Validate() error {
	if p.From == nil || p.To == nil {
		return ErrInvalidPeriod
	}

	if p.From.IsZero() || p.To.IsZero() || p.From.After(*p.To) {
		return ErrInvalidPeriod
	}

	switch p.Interval {
	case IntervalDay:
		if p.To.Sub(*p.From) > season {
			return ErrPeriodTooLarge
		}

	case IntervalWeek:
		if p.To.Sub(*p.From) > twoYears {
			return ErrPeriodTooLarge
		}

	case IntervalMonth:
		if p.To.Sub(*p.From) > fiveYears {
			return ErrPeriodTooLarge
		}

	default:
		return ErrInvalidInterval
	}

	return nil
}

func normalizeTime(t time.Time, interval Interval) time.Time {
	t = t.UTC()

	switch interval {
	case IntervalDay:
		return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)

	case IntervalWeek:
		weekday := int(t.Weekday())

		// sunday -> 7
		if weekday == 0 {
			weekday = 7
		}

		start := t.AddDate(0, 0, -(weekday - 1))

		return time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, time.UTC)

	case IntervalMonth:
		return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.UTC)

	default:
		return t
	}
}
