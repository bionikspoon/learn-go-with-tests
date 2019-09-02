package poker

import (
	"bufio"
	"io"
	"strings"
	"time"
)

type CLI struct {
	store   PlayerStore
	in      *bufio.Scanner
	alerter BlindAlerter
}

func NewCLI(store PlayerStore, in io.Reader, alerter BlindAlerter) *CLI {
	return &CLI{store, bufio.NewScanner(in), alerter}
}

func (cli *CLI) PlayPoker() {
	cli.scheduleBlindAlerts()
	userInput := cli.readline()
	cli.store.RecordWin(extractWinner(userInput))
}

func (cli *CLI) scheduleBlindAlerts() {
	blinds := []int{100, 200, 300, 400, 500, 600, 800, 1000, 2000, 4000, 8000}
	blindTime := 0 * time.Minute
	for _, blind := range blinds {
		cli.alerter.ScheduleAlertAt(blindTime, blind)
		blindTime += 10 * time.Minute
	}
}

func (cli *CLI) readline() string {
	cli.in.Scan()

	return cli.in.Text()
}

func extractWinner(userInput string) string {
	return strings.Replace(userInput, " wins", "", 1)
}
