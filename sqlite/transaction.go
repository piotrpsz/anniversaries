package sqlite

// BeginTransaction rozpoczyna transakcję.
func (d *Database) BeginTransaction() error {
	return d.Exec("BEGIN IMMEDIATE TRANSACTION")
}

// CommitTransaction zatwierdza transakcję.
func (d *Database) CommitTransaction() error {
	return d.Exec("COMMIT TRANSACTION")
}

// RollbackTransaction przywraca stan sprzed rozpoczęcia transakcji.
func (d *Database) RollbackTransaction() error {
	return d.Exec("ROLLBACK TRANSACTION")
}

// FinishTransaction kończy transakcję w sposób zależny
// od przysłanej flagi true := Commit, false := Rollback.
func (d *Database) FinishTransaction(success bool) error {
	if success {
		return d.CommitTransaction()
	}
	return d.RollbackTransaction()
}
