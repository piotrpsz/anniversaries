package sqlite

// Ważna uwaga: pomiędzy includami a importem "C"
// nie może być pustej linii.

/*
#include <stdlib.h>
#include <sqlite3.h>
#cgo LDFLAGS: -lsqlite3
*/
import "C"
import (
	"errors"
	"fmt"
	"os"
	"sync"
	"unsafe"
)

type Database struct {
	db    *C.sqlite3
	stmt  *C.sqlite3_stmt
	fpath string
}

var (
	instance *Database
	once     sync.Once
)

func SQLite(fpath []byte) *Database {
	once.Do(func() {
		if fpath == nil {
			panic("no path to database")
		}
		instance = sqlite(string(fpath))
	})
	return instance
}

func sqlite(fpath string) *Database {
	return &Database{fpath: fpath}
}

// Version zwraca numer wersji używanej biblioteki SQLite.
func (d *Database) Version() string {
	return C.GoString(C.sqlite3_libversion())
}

// ErrorCode numer kodu ostatniego błedu
func (d *Database) ErrorCode() C.int {
	return C.sqlite3_errcode(d.db)
}

// ErrorString opis ostatniego błędu
func (d *Database) ErrorString() string {
	return C.GoString(C.sqlite3_errmsg(d.db))
}

// Close zamknięcie bazy danych.
func (d *Database) Close() {
	if d.db == nil {
		return
	}
	if C.sqlite3_close(d.db) == Ok {
		C.sqlite3_shutdown()
		d.db = nil
	}
}

// Remove usunięcie pliku bazy danych z dysku.
// Uwaga: baza danych musi być zamknięta.
func (d *Database) Remove() (err error) {
	if err = d.FileExists(); err != nil {
		if errors.Is(err, ErrDatabaseNotExists) {
			return nil
		}
		return
	}

	err = os.Remove(d.fpath)
	return
}

// Open otwarcie bazy danych do zapisu i odczytu
func (d *Database) Open() (err error) {
	if d.db != nil {
		return ErrDatabaseIsInUse
	}
	if err = d.Exists(); err != nil {
		return
	}

	cfpath := C.CString(d.fpath)
	defer C.free(unsafe.Pointer(cfpath))

	if C.sqlite3_initialize() == Ok {
		if C.sqlite3_open_v2(cfpath, &d.db, C.SQLITE_OPEN_READWRITE, nil) == Ok {
			return
		}
	}

	return fmt.Errorf("%s", d.ErrorString())
}

// Create utworzenie nowej bazy danych.
// Użytkownik sam tworzy tabele.
func (d *Database) Create(doit func(*Database) error) (err error) {
	if d.db != nil {
		return ErrDatabaseIsInUse
	}

	if d.fpath != ":memory:" {
		if err = d.Exists(); !errors.Is(err, os.ErrNotExist) {
			return err
		}
	}

	cfpath := C.CString(d.fpath)
	defer C.free(unsafe.Pointer(cfpath))

	if C.sqlite3_initialize() == Ok {
		if C.sqlite3_open_v2(cfpath, &d.db, C.SQLITE_OPEN_READWRITE|C.SQLITE_OPEN_CREATE, nil) == Ok {
			if doit != nil {
				return doit(d)
			}
			return
		}
	}
	return fmt.Errorf("%s", d.ErrorString())
}

// Exec wykonanie polecenia.
func (d *Database) Exec(query string) (err error) {
	cquery := C.CString(query)
	defer C.free(unsafe.Pointer(cquery))

	if C.sqlite3_exec(d.db, cquery, nil, nil, nil) == Ok {
		return
	}
	return fmt.Errorf("%s", d.ErrorString())
}

func (d *Database) LastInsertedID() int64 {
	return int64(C.sqlite3_last_insert_rowid(d.db))
}
