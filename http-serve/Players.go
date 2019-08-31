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
