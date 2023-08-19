package math

import "time"

func MinDuration(d1 time.Duration, d2 time.Duration) time.Duration {
	if d1 < d2 {
		return d1
	}
	return d2
}

func MinTime(t1 *time.Time, t2 *time.Time) *time.Time {
	if t1 == nil || t2 == nil {
		return nil
	}
	if t1.Before(*t2) {
		return t1
	}
	return t2
}

func MaxTime(t1 *time.Time, t2 *time.Time) *time.Time {
	if t1 == nil {
		return t2
	} else if t2 == nil {
		return t1
	}

	if t1.After(*t2) {
		return t1
	}
	return t2
}
