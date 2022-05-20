package sqlite

import (
	"fmt"
	"strings"
)

type Row = map[string]any
type Where = map[string]any

// ExecQuery wykonanie zapytania, które nie zwraca danych.
func (d *Database) ExecQuery(query string, values ...any) error {
	if d.prepare(query) == nil {
		defer d.finalize()
		if len(values) > 0 {
			d.bind(values...)
		}
		retv := d.step()
		if retv == Ok || retv == StatusDone {
			return nil
		}
	}
	return fmt.Errorf(d.ErrorString())
}

// UpdateQuery
// example UPDATE table SET f1 = v1, f2 = v2, .... WHERE id=?
func (d *Database) UpdateQuery(query string, values ...any) error {
	return d.ExecQuery(query, values...)
}

// Update wykonanie zapytania UPDATE, skonstruowanego z przysłanych komponentów.
func (d *Database) Update(table string, values Row, where Where) error {
	n := len(values)
	if n == 0 {
		return ErrNoDataForQuery
	}

	binds := make([]string, 0, n)
	params := make([]any, 0, n)
	for k, v := range values {
		binds = append(binds, k+"=?")
		params = append(params, v)
	}
	query := fmt.Sprintf("UPDATE %s SET %s", table, strings.Join(binds, ","))

	if len(where) > 0 {
		binds = binds[:0]
		for k, v := range where {
			binds = append(binds, k+"=?")
			params = append(params, v)
		}
		query += fmt.Sprintf(" WHERE %s", strings.Join(binds, " AND "))
	}

	return d.UpdateQuery(query, params...)
}

// Insert wykonanie zapytania INSERT, skonstruowanego z przysłanych komponentów.
// Jeśli wszystko poszło dobrze zwraca wartość ID dodanego wiersza.
func (d *Database) Insert(table string, values Row) (int64, error) {
	delete(values, "id")

	n := len(values)
	if n == 0 {
		return -1, ErrNoDataForQuery
	}

	names := make([]string, 0, n)
	binds := make([]string, 0, n)
	params := make([]any, 0, n)

	for k, v := range values {
		names = append(names, k)
		binds = append(binds, "?")
		params = append(params, v)
	}

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
		table,
		strings.Join(names, ","),
		strings.Join(binds, ","))

	return d.InsertQuery(query, params...)
}

// InsertQuery wykonanie gotowego zapytania INSERT.
// Jeśli wszystko poszło dobrze zwraca wartość ID dodanego wiersza.
func (d *Database) InsertQuery(query string, values ...any) (int64, error) {
	if d.ExecQuery(query, values...) == nil {
		return d.LastInsertedID(), nil
	}
	return -1, fmt.Errorf(d.ErrorString())
}

// SelectQuery - wykonanie gotowego zapytania SELECT.
// Np. SELECT * WHERE id=?
func (d *Database) SelectQuery(query string, values ...any) ([]Row, error) {
	if d.prepare(query) == nil {
		defer d.finalize()

		if len(values) > 0 {
			d.bind(values...)
		}
		if result := d.fetchResult(); result != nil {
			return result, nil
		}
	}
	return nil, fmt.Errorf(d.ErrorString())
}

// Select wykonanie zapytania SELECT, skonstruowanego z przysłanych komponentów.
func (d *Database) Select(table string, where Where, names ...string) ([]Row, error) {
	fields := "*"
	if len(names) > 0 {
		fields = strings.Join(names, ",")
	}
	query := fmt.Sprintf("SELECT %s FROM %s", fields, table)

	n := len(where)
	var params []any
	if n > 0 {
		whereNames := make([]string, 0, n)
		for k, v := range where {
			whereNames = append(whereNames, fmt.Sprintf("%s=?", k))
			params = append(params, v)
		}
		query += fmt.Sprintf(" WHERE %s", strings.Join(whereNames, " AND "))
	}

	return d.SelectQuery(query, params...)
}

// SelectJoin zapytanie wykorzystujące JOIN'y.
// JOIN'ów jest tyle ile jest wskazanych tablic (minus 1).
// Tablica z indeksem 0 uznawana jest za tablicę 'master'.
func (d *Database) SelectJoin(tablesAndFields []Row, where Where) ([]Row, error) {
	if len(tablesAndFields) == 0 {
		return nil, ErrNoDataForQuery
	}

	var masterTableName string
	var tables []string
	var fields []string

	// Pobranie nazwy tablicy 'master' oraz zebranie pól.
	for tableName, tableFields := range tablesAndFields[0] {
		masterTableName = tableName
		if names, ok := tableFields.([]string); ok {
			for _, name := range names {
				fields = append(fields, "\t"+nameWithAlias(tableName, name))
			}
		}
	}
	// Zbieranie nazw tablic i pól z pozostałych tablic.
	for _, data := range tablesAndFields[1:] {
		for tableName, tableFields := range data {
			if names, ok := tableFields.([]string); ok {
				tables = append(tables, tableName)
				for _, name := range names {
					fields = append(fields, "\t"+nameWithAlias(tableName, name))
				}
			}
		}
	}
	if len(tables) == 0 || len(fields) == 0 {
		return nil, ErrNoDataForQuery
	}

	// Początkowa postać zapytania.
	query := fmt.Sprintf("SELECT\n%s\nFROM\n\t%s", strings.Join(fields, ",\n"), masterTableName)

	// JOIN
	var joins []string
	for _, joinedTableName := range tables {
		joins = append(joins, joinCmd(masterTableName, joinedTableName))
	}
	// Zapytanie z JOIN'ami.
	query += "\n" + strings.Join(joins, "\n")

	// WHERE jeśli jest
	var params []interface{}
	if len(where) > 0 {
		var binds []string
		for k, v := range where {
			binds = append(binds, fmt.Sprintf("%s.%s=?", masterTableName, k))
			params = append(params, v)
		}
		// Zapytanie z WHERE
		query += fmt.Sprintf("\nWHERE %s\n", strings.Join(binds, " AND "))
	}

	// Wykonujemy zapytanie do bazy danych.
	result, err := d.SelectQuery(query, params...)
	if err != nil {
		return nil, err
	}

	result = onlyNotEmpty(result)

	// Grupowanie wyników i zwrot
	retval := make([]Row, 0, len(result))
	for _, row := range result {
		retval = append(retval, groupResult(masterTableName, dispathTables(row)))
	}
	return retval, nil
}

func onlyNotEmpty(data []Row) []Row {
	n := len(data)
	if n == 0 {
		return nil
	}

	result := make([]Row, 0, len(data))
	for _, row := range data {
		vdata := make(Row)
		for k, v := range row {
			if v != nil {
				vdata[k] = v
			}
		}
		if len(vdata) > 0 {
			result = append(result, vdata)
		}
	}

	if len(result) == 0 {
		return nil
	}
	return result
}

// nameWithAlias w zapytaniach z joinami pola muszą nazywać się <table>.<name>
// ale muszą mieć dodane aliasy aby można rozpoznać je w odpowiedzi.
// Np. pole title w tabeli book: book.title AS book_title
// W zwróconym wyniku zapytania będziemy widzieli aliasy.
func nameWithAlias(tableName, fieldName string) string {
	return fmt.Sprintf("%s.%s AS %s_%s", tableName, fieldName, tableName, fieldName)
}

// displatchTables Pogrupowanie pól (nazwa/wartość) zwróconych
// z zapytania w/g tabel. Każde zwrócone pole ma w nazwie
// prefiks '<table>_' (alias).
func dispathTables(row Row) map[string]Row {
	tables := make(map[string]map[string]interface{})

	for name, value := range row {
		idx := strings.Index(name, "_")
		table := name[:idx]
		field := name[idx+1:]
		if _, ok := tables[table]; !ok {
			tables[table] = make(map[string]interface{})
		}
		tables[table][field] = value
	}

	return tables
}

// groupResult rekursywne grupowanie wyników zapytania.
// Wartość pola o nazwie <table>_id jest zastępowana przez zbiór pól tablicy 'table'.
// Operacja wykonywana jest rekursywnie póki się da.
func groupResult(tableName string, tables map[string]Row) Row {
	master := tables[tableName]
	for name := range master {
		if strings.HasSuffix(name, "_id") {
			nameWithoutSuffix := name[:len(name)-3]
			if _, ok := tables[nameWithoutSuffix]; ok {
				delete(master, name)
				master[nameWithoutSuffix] = groupResult(nameWithoutSuffix, tables)
			} else {
				// nie zwracamy pól, które mają wartość nil
				// master[nameWithoutSuffix] = nil
			}
		}
	}
	return master
}

// joinCmd kontruowanie polecenia JOIN dla podanych tablic.
func joinCmd(masterTableName, joindeTableName string) string {
	masterField := fmt.Sprintf("%s.%s_id", masterTableName, joindeTableName)
	joinedField := fmt.Sprintf("%s.id", joindeTableName)
	return fmt.Sprintf("JOIN %s ON %s=%s", joindeTableName, joinedField, masterField)
}
