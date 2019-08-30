package main

import (
	"log"

	"github.com/astaxie/beego/orm"
	_ "github.com/mattn/go-sqlite3"
)

type DatabasePlayerStore struct {
	o orm.Ormer
}

func (player *Player) TableIndex() [][]string {
	return [][]string{
		[]string{"Name"},
	}
}

func NewDatabasePlayerStore(debug bool) *DatabasePlayerStore {
	orm.Debug = debug
	orm.RegisterModel(new(Player))

	if err := orm.RegisterDataBase("default", "sqlite3", ":memory:"); err != nil {
		log.Printf("error %#v", err)
	}

	if err := orm.RunSyncdb("default", true, debug); err != nil {
		log.Printf("error %#v", err)
	}
	o := orm.NewOrm()
	return &DatabasePlayerStore{o}
}

func (store *DatabasePlayerStore) GetPlayerScore(name string) int {
	player := Player{Name: name}

	if err := store.o.Read(&player, "Name"); err != nil {
		log.Printf("error %#v", err)
	}
	return player.Wins
}

func (store *DatabasePlayerStore) RecordWin(name string) {

	player := Player{Name: name}

	if _, _, err := store.o.ReadOrCreate(&player, "Name"); err != nil {
		log.Printf("error %#v", err)
	}

	player.Wins++

	if _, err := store.o.Update(&player); err != nil {
		log.Printf("error %#v", err)

	}
}

func (store *DatabasePlayerStore) GetLeague() (players []Player) {

	store.o.QueryTable("player").All(&players)
	return
}
