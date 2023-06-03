package common

import (
	"math"
	"math/rand"
	"time"

	"github.com/StepanTita/go-EdgeGPT/common"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"common/convert"
)

func RunEvery(d time.Duration, fs ...func() error) error {
	if err := runFuncs(CurrentTimestamp(), fs...); err != nil {
		return errors.Wrap(err, "failed to run funcs initial")
	}
	for x := range time.Tick(d) {
		if err := runFuncs(x, fs...); err != nil {
			return errors.Wrap(err, "failed to run funcs")
		}
	}
	return nil
}

func RunEveryWithBackoff(d time.Duration, minBackoff, maxBackoff time.Duration, fs ...func() error) {

	backOffs := make([]*funcBackoff, len(fs))
	for i, f := range fs {
		if err := f(); err != nil {
			x := common.CurrentTimestamp()
			logrus.WithError(err).Errorf("failed to run function with backoff: %v -> %v", x, i)
			backOffs[i] = &funcBackoff{lastRun: &x, backOff: minBackoff, trial: 1}
		} else {
			backOffs[i] = nil
		}
	}

	for x := range time.Tick(d) {
		for i, f := range fs {
			if backOffs[i] != nil {
				// if last run of this function + backoff time is after now -> wait more
				// otherwise - retry
				if backOffs[i].lastRun.Add(backOffs[i].backOff).After(x) {
					continue
				}
			}
			if err := f(); err != nil {
				logrus.WithError(err).Errorf("failed to run function with backoff: %v -> %v", x, i)

				// will be just empty struct if nil
				oldBackoff := convert.FromPtr(backOffs[i])
				// min(oldBackoff * 2^(i) + minBackoff * rand.Float[0.5, 1], maxBackoff)
				newBackoffDuration := MinDuration(oldBackoff.backOff*time.Duration(math.Pow(2, float64(oldBackoff.trial)))+minBackoff*time.Duration(rand.Float64()*0.5+0.5), maxBackoff)
				backOffs[i] = &funcBackoff{
					lastRun: MaxTime(oldBackoff.lastRun, convert.ToPtr(x.Add(minBackoff))),
					backOff: newBackoffDuration,
					trial:   oldBackoff.trial + 1,
				}

				logrus.WithField("new-backoff", newBackoffDuration).Warnf("running with new backoff: %d", i)
			} else {
				backOffs[i] = nil
			}
		}
	}
	return
}

func runFuncs(x time.Time, fs ...func() error) error {
	for i, f := range fs {
		if err := f(); err != nil {
			return errors.Wrapf(err, "failed to run function: %v -> %v", x, i)
		}
	}
	return nil
}

// CurrentTimestamp is a utility method to make sure UTC time is used all over the code
func CurrentTimestamp() time.Time {
	return time.Now().UTC()
}

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
