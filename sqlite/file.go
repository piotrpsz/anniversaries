package sqlite

import (
	"os"
)

var databaseFileHeader = []byte{0x53, 0x51, 0x4c, 0x69, 0x74, 0x65, 0x20, 0x66, 0x6f, 0x72, 0x6d, 0x61, 0x74, 0x20, 0x33, 0x00}

func (d *Database) FileExists() error {
	if _, err := os.Stat(d.fpath); err != nil {
		if os.IsNotExist(err) {
			return ErrFileNotExists
		}
		return err
	}
	return nil
}

// Exists
// nil if extsts
func (d *Database) Exists() error {
	if d.db != nil {
		return ErrDatabaseIsInUse
	}

	fh, err := os.Open(d.fpath)
	if err != nil {
		return err
	}
	defer fh.Close()

	nbytes := len(databaseFileHeader)
	header := make([]byte, nbytes)

	n, err := fh.Read(header)
	if err != err {
		return err
	}
	if n != nbytes {
		return ErrWrongDatabaseFormat
	}

	for i, v := range header {
		if v != databaseFileHeader[i] {
			return ErrWrongDatabaseFormat
		}
	}

	return nil
}
