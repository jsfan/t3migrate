package storage

import (
	"github.com/jsfan/t3migrate/internal/config"
	"github.com/jsfan/t3migrate/internal/storage/model"
)

type DAL interface {
	Connect(config *config.MySQLConfig) error

	CountAll(tableName string, columns []string, includedIds []int64) (int64, error)
	DescribeTable(tableName string) ([]*model.TableDescription, error)
	SelectAll(tableName string, columns []string, offset, limit *int64, includedIds []int64) ([][]interface{}, error)
	SelectTVFlexFromPages() ([]*model.TVFlex, error)

	Upsert(tableName string, columns []string, values []interface{}) error
}
