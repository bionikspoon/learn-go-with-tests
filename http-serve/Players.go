package main

type Players []Player

func (players Players) Find(name string) *Player {
	for i, player := range players {
		if player.Name == name {
			return &players[i]
		}
	}
	return nil
}

type ByWins Players

func (a ByWins) Len() int           { return len(a) }
func (a ByWins) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByWins) Less(i, j int) bool { return a[i].Wins > a[j].Wins }
