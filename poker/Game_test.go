package poker_test

import (
	"fmt"
	"io"
	"io/ioutil"
	"reflect"
	"testing"
	"time"

	"github.com/bionikspoon/learn-go-with-tests/poker"
)

func TestGame_Start(t *testing.T) {
	store := &poker.StubPlayerStore{
		Players: poker.Players{},
	}

	t.Run("it schedules alerts given 5 players", func(t *testing.T) {
		blindAlerter := &SpyBlindAlerter{}
		game := poker.NewTexasHoldem(blindAlerter, store)
		game.Start(5, ioutil.Discard)

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

		checkSchedulingCases(t, cases, blindAlerter)
	})

	t.Run("it schedules alerts given 7 players", func(t *testing.T) {
		blindAlerter := &SpyBlindAlerter{}
		game := poker.NewTexasHoldem(blindAlerter, store)
		game.Start(7, ioutil.Discard)

		cases := []scheduledAlert{
			{0 * time.Minute, 100},
			{12 * time.Minute, 200},
			{24 * time.Minute, 300},
			{36 * time.Minute, 400},
			{48 * time.Minute, 500},
			{60 * time.Minute, 600},
			{72 * time.Minute, 800},
			{84 * time.Minute, 1000},
			{96 * time.Minute, 2000},
			{108 * time.Minute, 4000},
			{120 * time.Minute, 8000},
		}

		checkSchedulingCases(t, cases, blindAlerter)
	})
}

func TestGame_Finish(t *testing.T) {
	store := &poker.StubPlayerStore{
		Players: poker.Players{},
	}

	t.Run("it records Julia's win", func(t *testing.T) {
		game := poker.NewTexasHoldem(&SpyBlindAlerter{}, store)

		game.Finish("Julia")

		want := []string{"RecordWin Julia"}
		if got := store.CalledWith; !reflect.DeepEqual(got, want) {
			t.Errorf("got %v want %v", got, want)
		}
	})

	t.Run("it records Roger's win", func(t *testing.T) {
		game := poker.NewTexasHoldem(&SpyBlindAlerter{}, store)

		game.Finish("Roger")

		want := []string{"RecordWin Julia", "RecordWin Roger"}
		if got := store.CalledWith; !reflect.DeepEqual(got, want) {
			t.Errorf("got %v want %v", got, want)
		}
	})

}

func checkSchedulingCases(t *testing.T, cases []scheduledAlert, blindAlerter *SpyBlindAlerter) {
	t.Helper()

	for i, want := range cases {
		t.Run(fmt.Sprint(want), func(t *testing.T) {
			if len(blindAlerter.alerts) <= i {
				t.Fatalf("alert %d was not scheduled %v", i, blindAlerter.alerts)
			}
			alert := blindAlerter.alerts[i]
			assertScheduleAlert(t, alert, want)
		})
	}

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

func (alerter *SpyBlindAlerter) ScheduleAlertAt(scheduledAt time.Duration, amount int, _ io.Writer) {
	alerter.alerts = append(alerter.alerts, scheduledAlert{scheduledAt, amount})
}
