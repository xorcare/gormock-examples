//+build !gormock

package main

import (
	"basic/storage"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"io/ioutil"
	"log"
	"path"
	"time"
)

func main() {
	dir, err := ioutil.TempDir("", "")
	ok(err)

	dbPath := path.Join(dir, "gorm.db")
	log.Println(dbPath)
	db, err := gorm.Open("sqlite3", dbPath)
	ok(err)
	defer db.Close()

	ok(db.AutoMigrate(&storage.Model{}).Error)

	s := storage.New(db)
	for i := 0; i < 8; i++ {
		ok(s.Save(&storage.Model{}))
	}

	list(s)

	model, err := s.FindByID(4)
	ok(err)
	model.CreatedAt = time.Now()
	log.Println("Update model", fmt.Sprintf("%+v", model))
	ok(s.Save(model))
	list(s)

	log.Println("DeleteByID", 1)
	ok(s.DeleteByID(1))
	list(s)

	ok(s.DeleteAll())
	list(s)
}

func ok(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func list(s *storage.Storage) {
	if ms, err := s.FindAll(); err != nil {
		log.Fatal(err)
	} else {
		for _, model := range ms {
			log.Println(fmt.Sprintf("%+v", model))
		}
	}
}
