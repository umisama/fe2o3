package main

import (
	"github.com/zenazn/goji"
	"github.com/zenazn/goji/web/middleware"
	"log"
	"time"
)

func main() {
	initDatabase()
	go Scheduled()
	goji.Get("/", index)
	goji.Use(middleware.NoCache)
	goji.Serve()
}

func Scheduled() {
	fetch(time.Now())
	for t := range time.Tick(time.Hour) {
		err := fetch(t)
		if err != nil {
			log.Println(err)
		}
	}
}

func fetch(t time.Time) error {
	log.Println("fetch start on ", t)
	hatena_list, err := GetHatenaFeedListFromInternet()
	if err != nil {
		return err
	}
	qiita_list, err := GetQiitaFeedListFromInternet()
	if err != nil {
		return err
	}

	tx, _ := db.Begin()
	for _, entry := range hatena_list {
		err = entry.Save(tx)
		if err != nil {
			log.Println(err)
			continue
		}
	}
	for _, entry := range qiita_list {
		err = entry.Save(tx)
		if err != nil {
			log.Println(err)
			continue
		}
	}
	err = tx.Commit()
	if err != nil {
		return err
	}

	log.Println("complete in success")
	return nil
}
