package poker

import "io"

type SpyGame struct {
	StartCalled    bool
	StartedWith    int
	FinishedCalled bool
	FinishedWith   string

	BlindAlert []byte
}

func (game *SpyGame) Start(startedWith int, out io.Writer) {
	game.StartedWith = startedWith
	game.StartCalled = true
	_, _ = out.Write(game.BlindAlert)
}
func (game *SpyGame) Finish(finishedWith string) {
	game.FinishedWith = finishedWith
}
