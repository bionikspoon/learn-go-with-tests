package poker

import (
	"io"
	"time"
)

type TexasHoldemGame struct {
	alerter BlindAlerter
	store   PlayerStore
}

func NewTexasHoldem(alerter BlindAlerter, store PlayerStore) *TexasHoldemGame {
	return &TexasHoldemGame{alerter: alerter, store: store}
}

func (game *TexasHoldemGame) Start(numberOfPlayers int, w io.Writer) {
	blindIncrement := time.Duration(5+numberOfPlayers) * time.Minute
	blinds := []int{100, 200, 300, 400, 500, 600, 800, 1000, 2000, 4000, 8000}
	blindTime := 0 * time.Minute
	for _, blind := range blinds {
		game.alerter.ScheduleAlertAt(blindTime, blind, w)
		blindTime += blindIncrement
	}
}

func (game TexasHoldemGame) Finish(winner string) {
	game.store.RecordWin(winner)
}
