package storage

func (store *Store) CountRecords(tableName string, columns []string, includedIds []int64) (int64, error) {
	return store.dal.CountAll(tableName, columns, includedIds)
}

func (store *Store) GetRecords(tableName string, columns []string, offset, limit *int64, includedIds []int64) ([][]interface{}, error) {
	return store.dal.SelectAll(tableName, columns, offset, limit, includedIds)
}

func (store *Store) PutRecords(tableName string, columns []string, values [][]interface{}) error {
	for i := 0; i < len(values); i++ {
		if err := store.dal.Upsert(tableName, columns, values[i]); err != nil {
			return err
		}
	}
	return nil
}
