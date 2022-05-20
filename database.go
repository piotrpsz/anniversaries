package main

import (
	"errors"
	"os"
	"os/user"
	"path/filepath"

	"calendar/model"
	"calendar/sqlite"
	log "github.com/sirupsen/logrus"
)

const (
	databaseCreated = "database created (%v)"
	databaseOpened  = "database opened (%v)"
)

func dbOpenOrCreate() (err error) {
	fpath, err := dbFilePath()
	if err != nil {
		return
	}

	db := sqlite.SQLite([]byte(fpath))
	// _ = db.Remove()

	err = db.FileExists()
	if err != nil {
		if errors.Is(err, sqlite.ErrFileNotExists) {
			if err = db.Create(createTables); err != nil {
				return
			}
			log.Infof(databaseCreated, fpath)
			return
		}
	}

	if err = db.Open(); err != nil {
		return
	}
	log.Infof(databaseOpened, fpath)
	return
}

func dbFilePath() (fpath string, err error) {
	usr, err := user.Current()
	if err != nil {
		return
	}
	fpath = filepath.Join(usr.HomeDir, ".rocznice")
	if err = os.MkdirAll(fpath, 0744); err != nil {
		return
	}
	fpath = filepath.Join(fpath, "data.sqlite")
	return
}

func createTables(db *sqlite.Database) (err error) {
	if err = model.CreatePersonTable(); err == nil {

	}
	return
}
