package storage

import (
	"github.com/jsfan/t3migrate/internal/storage/model"
)

func (store *Store) MapTable(table string) (map[string]*model.TableDescription, error) {
	columns, err := store.dal.DescribeTable(table)
	if err != nil {
		return nil, err
	}
	colMap := make(map[string]*model.TableDescription)
	for _, col := range columns {
		colMap[col.Field] = col
	}
	return colMap, nil
}
