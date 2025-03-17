package repository

import (
	"fmt"
	"github.com/ankitjha420/pomo/pomodoro"
	"sync"
)

// InMemoryRepo stores all pomodoro intervals ->
type InMemoryRepo struct {
	sync.RWMutex
	intervals []pomodoro.Interval
}

// NewInMemoryRepo returns an in memory repo ->
func NewInMemoryRepo() *InMemoryRepo {
	return &InMemoryRepo{
		intervals: []pomodoro.Interval{},
	}
}

// Repository methods ->

func (r *InMemoryRepo) Create(i pomodoro.Interval) (int64, error) {
	r.Lock()
	defer r.Unlock()

	i.ID = int64(len(r.intervals)) + 1
	r.intervals = append(r.intervals, i)

	return i.ID, nil
}

func (r *InMemoryRepo) Update(i pomodoro.Interval) error {
	r.Lock()
	defer r.Unlock()

	if i.ID == 0 {
		return fmt.Errorf("%w: %d", pomodoro.ErrInvalidID, i.ID)
	}

	r.intervals[i.ID-1] = i
	return nil
}

func (r *InMemoryRepo) ByID(id int64) (pomodoro.Interval, error) {
	r.RLock()
	defer r.RUnlock()

	i := pomodoro.Interval{}
	if id == 0 {
		return i, fmt.Errorf("%w: %d", pomodoro.ErrInvalidID, id)
	}

	i = r.intervals[id-1]
	return i, nil
}

func (r *InMemoryRepo) Last() (pomodoro.Interval, error) {
	r.RLock()
	defer r.RUnlock()

	i := pomodoro.Interval{}
	if len(r.intervals) == 0 {
		return i, pomodoro.ErrNoIntervals
	}
	return r.intervals[len(r.intervals)-1], nil
}

func (r *InMemoryRepo) Breaks(n int8) ([]pomodoro.Interval, error) {
	r.RLock()
	defer r.RUnlock()

	var data []pomodoro.Interval

	for j := len(r.intervals) - 1; j >= 0; j-- {
		if r.intervals[j].Category == pomodoro.CategoryPomodoro {
			continue
		}

		data = append(data, r.intervals[j])
		if len(data) == int(n) {
			return data, nil
		}
	}

	return data, nil
}
