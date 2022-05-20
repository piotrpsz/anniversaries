package sqlite

/*
#include <stdlib.h>
#include <sqlite3.h>
#cgo LDFLAGS: -lsqlite3
int bind_text(sqlite3_stmt *stmt, int index, const char* txt) {
	return sqlite3_bind_text(stmt, index, txt, -1, SQLITE_TRANSIENT);
}
const char* column_text(sqlite3_stmt *stmt, int index) {
	return (const char *)sqlite3_column_text(stmt, index);
}
*/
import "C"
import (
	"fmt"
	"unsafe"
)

func (d *Database) step() C.int {
	return C.sqlite3_step(d.stmt)
}

func (d *Database) reset() (err error) {
	if C.sqlite3_reset(d.stmt) == Ok {
		if C.sqlite3_clear_bindings(d.stmt) == Ok {
			return
		}
	}
	return fmt.Errorf(d.ErrorString())
}

func (d *Database) prepare(query string) (err error) {
	cstr := C.CString(query)
	defer C.free(unsafe.Pointer(cstr))

	if C.sqlite3_prepare_v2(d.db, cstr, -1, &d.stmt, nil) == Ok {
		return
	}
	return fmt.Errorf(d.ErrorString())
}

func (d *Database) finalize() (err error) {
	if C.sqlite3_finalize(d.stmt) == Ok {
		return
	}
	return fmt.Errorf(d.ErrorString())
}

func (d *Database) columnCount() int {
	return int(C.sqlite3_column_count(d.stmt))
}

func (d *Database) columnType(idx int) int {
	return int(C.sqlite3_column_type(d.stmt, C.int(idx)))
}

func (d *Database) columnIndex(columnName string) int {
	cstr := C.CString(columnName)
	defer C.free(unsafe.Pointer(cstr))
	return int(C.sqlite3_bind_parameter_index(d.stmt, cstr))
}

func (d *Database) columnName(idx int) string {
	return C.GoString(C.sqlite3_column_name(d.stmt, C.int(idx)))
}

func (d *Database) namedBind(data map[string]interface{}) {
	for k, v := range data {
		d.bindValueAtIndex(d.columnIndex(k), v)
	}
}

// bind bindowanie parametrów do spreparowego zapytania.
// UWAGA: indeksy bindowanych pól zaczynaj się od 1 (nie od 0)
func (d *Database) bind(values ...interface{}) {
	for i, v := range values {
		d.bindValueAtIndex(i+1, v)
	}
}

func (d *Database) bindValueAtIndex(i int, v interface{}) {
	switch v.(type) {
	case int8:
		d.bindInt64(i, int64(v.(int8)))
	case int16:
		d.bindInt64(i, int64(v.(int16)))
	case int32:
		d.bindInt64(i, int64(v.(int32)))
	case int64:
		d.bindInt64(i, v.(int64))
	case int:
		d.bindInt64(i, int64(v.(int)))
	case float32:
		d.bindFloat64(i, float64(v.(float32)))
	case float64:
		d.bindFloat64(i, v.(float64))
	case string:
		d.bindString(i, v.(string))
	case []byte:
		d.bindBlob(i, v.([]byte))
	case nil:
		d.bindNull(i)
	default:
		panic(fmt.Sprintf("unknown value type (%d)", i))
	}
}

func (d *Database) fetchResult() []map[string]interface{} {
	result := make([]map[string]interface{}, 0)

	if n := d.columnCount(); n > 0 {
		for d.step() == StatusRow {
			result = append(result, d.fetchRow(n))
		}
	}

	if len(result) == 0 {
		return nil
	}
	return result
}

func (d *Database) fetchRow(n int) map[string]interface{} {
	row := make(map[string]interface{})

	for i := 0; i < n; i++ {
		name := d.columnName(i)
		switch d.columnType(i) {
		case Int64:
			row[name] = d.fetchInt64(i)
		case Float64:
			row[name] = d.fetchFloat64(i)
		case String:
			row[name] = d.fetchString(i)
		case Blob:
			row[name] = d.fetchBlob(i)
		case Null:
			row[name] = nil
		}
	}

	return row
}

/*                       S E T T E R S                             */

func (d *Database) bindInt64(idx int, v int64) C.int {
	return C.sqlite3_bind_int64(d.stmt, C.int(idx), C.sqlite3_int64(v))
}

func (d *Database) bindFloat64(idx int, v float64) C.int {
	return C.sqlite3_bind_double(d.stmt, C.int(idx), C.double(v))
}

func (d *Database) bindString(idx int, v string) C.int {
	cstr := C.CString(v)
	defer C.free(unsafe.Pointer(cstr))
	retv := C.bind_text(d.stmt, C.int(idx), cstr)
	return retv
}

func (d *Database) bindBlob(idx int, v []byte) C.int {
	return C.sqlite3_bind_blob(d.stmt, C.int(idx), unsafe.Pointer(&v[0]), C.int(len(v)), nil)
}

func (d *Database) bindNull(idx int) C.int {
	return C.sqlite3_bind_null(d.stmt, C.int(idx))
}

/*                       G E T T E R S                             */

func (d *Database) fetchInt64(idx int) int64 {
	return int64(C.sqlite3_column_int64(d.stmt, C.int(idx)))
}

func (d *Database) fetchFloat64(idx int) float64 {
	return float64(C.sqlite3_column_double(d.stmt, C.int(idx)))
}

func (d *Database) fetchString(idx int) string {
	return C.GoString(C.column_text(d.stmt, C.int(idx)))
}

func (d *Database) fetchBlob(idx int) []byte {
	n := C.int(C.sqlite3_column_bytes(d.stmt, C.int(idx)))
	ptr := C.sqlite3_column_blob(d.stmt, C.int(idx))
	return C.GoBytes(ptr, n)
}
