package sqlite

import "C"
import (
	"errors"
)

const (
	Ok         C.int = iota // Successful result
	Error                   // SQL error or missing Database
	Internal                // Internal logic error in SQLite
	Perm                    // Access permission denied
	Abort                   // Callback routine requested an abort
	Busy                    // The Database file is locked
	Locked                  // A table in the Database is locked
	NoMem                   // A malloc() failed
	ReadOnly                // Attempt to write a readonly Database
	Interrupt               // Operation terminated by sqlite3_interrupt()
	IoErr                   // Some kind of disk I/O error occurred
	Corrupt                 // The Database disk image is malformed
	NotFound                // NOT USED. Table or record not found
	Full                    // Insertion failed because Database is full
	CantOpen                // Unable to open the Database file
	Protocol                // NOT USED. Database lock protocol error
	Empty                   // Database is empty
	Schema                  // The Database schema changed
	TooBig                  // String or BLOB exceeds size limit
	Constraint              // Abort due to constraint violation
	Mismatch                // Data type mismatch
	Misuse                  // Library used incorrectly
	NoLfs                   // Uses OS features not supported on host
	Auth                    // Authorization denied
	Format                  // Auxiliary Database format error
	Range                   // 2nd parameter to sqlite3_bind out of range
	NotADb                  // File opened that is not a Database file
	StatusRow  = 100        // sqlite3_step() has another row ready
	StatusDone = 101        // sqlite3_step() has finished executing
)

var (
	ErrWrongDatabaseFormat = errors.New("wrong Database format")
	ErrDatabaseNotExists   = errors.New("database file not exists")
	ErrFileNotExists       = errors.New("file not exists")
	ErrDatabaseIsInUse     = errors.New("database is in use (is opened)")
	ErrNoDataForQuery      = errors.New("no data for query")
)

const (
	_ int = iota
	Int64
	Float64
	String
	Blob
	Null
)
