package main

import "testing"

func TestRecordingWinsAndShowingThem(t *testing.T) {
	t.Run("InMemoryPlayerStore", func(t *testing.T) {
		server := NewPlayerServer(NewInMemoryPlayerStore())

		assertUpdateAndShow(t, server, "Pepper", 3)
		assertUpdateAndShow(t, server, "Candy", 6)
		assertUpdateAndShow(t, server, "Anne", 2)

		players := Players{
			{0, "Candy", 6},
			{0, "Pepper", 3},
			{0, "Anne", 2},
		}
		assertLeague(t, server, players)
		assertLeague(t, server, players)

	})

	t.Run("DatabasePlayerStore", func(t *testing.T) {
		server := NewPlayerServer(NewDatabasePlayerStore(false, true))

		assertUpdateAndShow(t, server, "Pepper", 3)
		assertUpdateAndShow(t, server, "Candy", 6)
		assertUpdateAndShow(t, server, "Anne", 2)

		players := Players{
			{2, "Candy", 6},
			{1, "Pepper", 3},
			{3, "Anne", 2},
		}
		assertLeague(t, server, players)
		assertLeague(t, server, players)
	})

	t.Run("FileSystemPlayerStore", func(t *testing.T) {
		database, cleanup := createTempFile(t, "")
		defer cleanup()

		server := NewPlayerServer(NewFileSystemPlayerStore(database))

		assertUpdateAndShow(t, server, "Pepper", 3)
		assertUpdateAndShow(t, server, "Candy", 6)
		assertUpdateAndShow(t, server, "Anne", 2)

		players := Players{
			{0, "Candy", 6},
			{0, "Pepper", 3},
			{0, "Anne", 2},
		}
		assertLeague(t, server, players)
		assertLeague(t, server, players)
	})

}
