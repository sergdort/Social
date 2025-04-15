package store

type FakeSqlResult struct {
	InsertID      int64
	InsertError   error
	AffectedRows  int64
	AffectedError error
}

func (f *FakeSqlResult) LastInsertId() (int64, error) {
	return f.InsertID, f.InsertError
}

func (f *FakeSqlResult) RowsAffected() (int64, error) {
	return f.AffectedRows, f.AffectedError
}
