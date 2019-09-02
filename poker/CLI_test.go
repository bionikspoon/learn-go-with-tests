package poker_test

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/bionikspoon/learn-go-with-tests/poker"
)

func TestCLI(t *testing.T) {
	store := &poker.StubPlayerStore{
		Players: poker.Players{},
	}

	t.Run("records Julia's win", func(t *testing.T) {
		in := strings.NewReader("5\nJulia wins\n")
		out := &bytes.Buffer{}
		game := poker.NewGame(&SpyBlindAlerter{}, store)

		cli := poker.NewCLI(in, out, game)
		cli.PlayPoker()

		want := []string{"RecordWin Julia"}
		if got := store.CalledWith; !reflect.DeepEqual(got, want) {
			t.Errorf("got %v want %v", got, want)
		}
	})

	t.Run("records Roger's win", func(t *testing.T) {
		in := strings.NewReader("5\nRoger wins\n")
		out := &bytes.Buffer{}
		game := poker.NewGame(&SpyBlindAlerter{}, store)

		cli := poker.NewCLI(in, out, game)
		cli.PlayPoker()

		want := []string{"RecordWin Julia", "RecordWin Roger"}
		if got := store.CalledWith; !reflect.DeepEqual(got, want) {
			t.Errorf("got %v want %v", got, want)
		}

	})

	t.Run("it schedules printing of blind values", func(t *testing.T) {
		in := strings.NewReader("5\n")
		out := &bytes.Buffer{}
		blindAlerter := &SpyBlindAlerter{}
		game := poker.NewGame(blindAlerter, store)

		cli := poker.NewCLI(in, out, game)
		cli.PlayPoker()

		cases := []scheduledAlert{
			{0 * time.Minute, 100},
			{10 * time.Minute, 200},
			{20 * time.Minute, 300},
			{30 * time.Minute, 400},
			{40 * time.Minute, 500},
			{50 * time.Minute, 600},
			{60 * time.Minute, 800},
			{70 * time.Minute, 1000},
			{80 * time.Minute, 2000},
			{90 * time.Minute, 4000},
			{100 * time.Minute, 8000},
		}

		for i, want := range cases {
			t.Run(fmt.Sprint(want), func(t *testing.T) {
				if len(blindAlerter.alerts) <= i {
					t.Fatalf("alert %d was not scheduled %v", i, blindAlerter.alerts)
				}
				alert := blindAlerter.alerts[i]
				assertScheduleAlert(t, alert, want)
			})
		}
	})

	t.Run("it prompts the user to enter the number of players", func(t *testing.T) {
		in := strings.NewReader("7\n")
		out := &bytes.Buffer{}
		blindAlerter := &SpyBlindAlerter{}
		game := poker.NewGame(blindAlerter, store)

		cli := poker.NewCLI(in, out, game)
		cli.PlayPoker()

		got := out.String()

		if got != poker.PlayerPrompt {
			t.Errorf("got %q want %q", got, poker.PlayerPrompt)
		}

		cases := []scheduledAlert{
			{0 * time.Minute, 100},
			{12 * time.Minute, 200},
			{24 * time.Minute, 300},
			{36 * time.Minute, 400},
		}

		for i, want := range cases {
			t.Run(fmt.Sprint(want), func(t *testing.T) {
				if len(blindAlerter.alerts) <= i {
					t.Fatalf("alert %d was not scheduled %v", i, blindAlerter.alerts)
				}
				alert := blindAlerter.alerts[i]
				assertScheduleAlert(t, alert, want)
			})
		}
	})
}

func assertScheduleAlert(t *testing.T, alert, want scheduledAlert) {
	t.Helper()

	if got := alert.amount; got != want.amount {
		t.Errorf("amount: got %d want %d", got, want.amount)
	}
	if got := alert.at; got != want.at {
		t.Errorf("scheduled at: got %v want %v", got, want.at)
	}
}

type scheduledAlert struct {
	at     time.Duration
	amount int
}

func (alert scheduledAlert) String() string {
	return fmt.Sprintf("%d chips at %v", alert.amount, alert.at)

}

type SpyBlindAlerter struct {
	alerts []scheduledAlert
}

func (alerter *SpyBlindAlerter) ScheduleAlertAt(scheduledAt time.Duration, amount int) {
	alerter.alerts = append(alerter.alerts, scheduledAlert{scheduledAt, amount})
}
