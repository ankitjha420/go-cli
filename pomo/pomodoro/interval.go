package pomodoro

import (
	"context"
	"errors"
	"fmt"
	"time"
)

// constants ->
const (
	CategoryPomodoro   = "Pomodoro"
	CategoryShortBreak = "ShortBreak"
	CategoryLongBreak  = "LongBreak"
)

const (
	StateNotStarted = iota
	StateRunning
	StatePaused
	StateDone
	StateCancelled
)

// Interval type ->
type Interval struct {
	ID              int64
	StartTime       time.Time
	PlannedDuration time.Duration
	ActualDuration  time.Duration
	Category        string
	State           int
}

// IntervalConfig type ->
type IntervalConfig struct {
	repo               *Repository
	PomodoroDuration   time.Duration
	ShortBreakDuration time.Duration
	LongBreakDuration  time.Duration
}

// Repository interface ->
type Repository interface {
	Create(i Interval) (int64, error)
	Update(i Interval) error
	ByID(id int64) (Interval, error)
	Last() (Interval, error)
	Breaks(n int8) ([]Interval, error)
}

// Callback type for when tick functions ends ->
type Callback func(Interval)

// errors ->
var (
	ErrNoIntervals        = errors.New("no intervals")
	ErrIntervalNotRunning = errors.New("interval not running")
	ErrIntervalsCompleted = errors.New("interval completed or cancelled")
	ErrInvalidState       = errors.New("invalid state")
	ErrInvalidID          = errors.New("invalid ID")
)

// NewConfig returns a new config struct ->
func NewConfig(repo *Repository, pomodoro, shortBreak, longBreak time.Duration) *IntervalConfig {
	c := &IntervalConfig{
		repo:               repo,
		PomodoroDuration:   25 * time.Minute,
		ShortBreakDuration: 5 * time.Minute,
		LongBreakDuration:  15 * time.Minute,
	}

	if pomodoro > 0 {
		c.PomodoroDuration = pomodoro
	}
	if shortBreak > 0 {
		c.ShortBreakDuration = shortBreak
	}
	if longBreak > 0 {
		c.LongBreakDuration = longBreak
	}

	return c
}

// nextCategory returns next interval category depending on last interval in repository ->
func nextCategory(r Repository) (string, error) {
	li, err := r.Last()
	if err != nil && !errors.Is(err, ErrNoIntervals) {
		return CategoryPomodoro, nil
	}
	if err != nil {
		return "", err
	}

	if li.Category == CategoryLongBreak || li.Category == CategoryShortBreak {
		return CategoryPomodoro, nil
	}

	lastBreaks, err := r.Breaks(3)
	if err != nil {
		return "", err
	}

	if len(lastBreaks) < 3 {
		return CategoryShortBreak, nil
	}

	for _, i := range lastBreaks {
		if i.Category == CategoryLongBreak {
			return CategoryShortBreak, nil
		}
	}

	return CategoryLongBreak, nil
}

// tick controls current time and state ->
func tick(ctx context.Context, id int64, config *IntervalConfig, start, periodic, end Callback) error {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	repo := *config.repo
	i, err := repo.ByID(id)
	if err != nil {
		return err
	}
	expire := time.After(i.PlannedDuration - i.ActualDuration)

	start(i)

	for {
		select {
		case <-ticker.C:
			i, err := repo.ByID(id)
			if err != nil {
				return err
			}
			if i.State == StatePaused {
				return nil
			}

			i.ActualDuration += time.Second

			if err = repo.Update(i); err != nil {
				return err
			}
			periodic(i)

		case <-expire:
			i, err := repo.ByID(id)
			if err != nil {
				return err
			}

			i.State = StateDone
			end(i)

			return repo.Update(i)

		case <-ctx.Done():
			i, err := repo.ByID(id)
			if err != nil {
				return err
			}

			i.State = StateCancelled
			return repo.Update(i)
		}
	}
}

func newInterval(config *IntervalConfig) (Interval, error) {
	i := Interval{}
	repo := *config.repo
	cat, err := nextCategory(repo)
	if err != nil {
		return i, err
	}

	i.Category = cat
	switch cat {
	case CategoryPomodoro:
		i.PlannedDuration = config.PomodoroDuration
	case CategoryShortBreak:
		i.PlannedDuration = config.ShortBreakDuration
	case CategoryLongBreak:
		i.PlannedDuration = config.LongBreakDuration
	}

	i.ID, err = repo.Create(i)
	if err != nil {
		return i, err
	}
	return i, nil
}

// Following three functions will comprise the API of the Pomodoro application

func GetInterval(config *IntervalConfig) (Interval, error) {
	i := Interval{}
	repo := *config.repo
	var err error

	i, err = repo.Last()
	if err != nil && !errors.Is(err, ErrNoIntervals) {
		return i, err
	}

	if err == nil && i.State != StateCancelled && i.State != StateDone {
		return i, nil
	}

	return newInterval(config)
}

func (i Interval) Start(ctx context.Context, config *IntervalConfig, start, periodic, end Callback) error {
	repo := *config.repo

	switch i.State {
	case StateRunning:
		return nil
	case StateNotStarted:
		i.StartTime = time.Now()
		fallthrough
	case StatePaused:
		i.State = StateRunning
		if err := repo.Update(i); err != nil {
			return err
		}
		return tick(ctx, i.ID, config, start, periodic, end)
	case StateCancelled, StateDone:
		return fmt.Errorf("%w: cannot start", ErrIntervalsCompleted)
	default:
		return fmt.Errorf("%w: %d", ErrInvalidState, i.State)
	}
}

func (i Interval) Pause(config *IntervalConfig) error {
	repo := *config.repo

	if i.State != StateRunning {
		return ErrIntervalNotRunning
	}

	i.State = StatePaused
	return repo.Update(i)
}
