package main

import (
	"log"
	"time"

	"github.com/astaxie/beego/orm"
	_ "github.com/mattn/go-sqlite3"
)

type DatabasePlayerStore struct {
	o orm.Ormer
}

type DatabasePlayer struct {
	Id      int       `orm:"auto""`
	Name    string    `orm:"unique"`
	Score   int       `orm:"default(0)"`
	Created time.Time `orm:"auto_now_add;type(datetime)"`
	Updated time.Time `orm:"auto_now;type(datetime)"`
}

func (player *DatabasePlayer) TableIndex() [][]string {
	return [][]string{
		[]string{"Name"},
	}
}

func NewDatabasePlayerStore(debug bool) *DatabasePlayerStore {
	orm.Debug = debug
	orm.RegisterModel(new(DatabasePlayer))

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
	player := DatabasePlayer{Name: name}

	if err := store.o.Read(&player, "Name"); err != nil {
		log.Printf("error %#v", err)
	}
	return player.Score
}

func (store *DatabasePlayerStore) RecordWin(name string) {

	player := DatabasePlayer{Name: name}

	if _, _, err := store.o.ReadOrCreate(&player, "Name"); err != nil {
		log.Printf("error %#v", err)
	}

	player.Score++

	if _, err := store.o.Update(&player); err != nil {
		log.Printf("error %#v", err)

	}
}
