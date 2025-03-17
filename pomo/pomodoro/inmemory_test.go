package pomodoro_test

import (
	"github.com/ankitjha420/pomo/pomodoro"
	"github.com/ankitjha420/pomo/pomodoro/repository"
	"testing"
)

func getRepo(t *testing.T) (pomodoro.Repository, func()) {
	t.Helper()
	return repository.NewInMemoryRepo(), func() {}
}