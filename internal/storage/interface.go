package storage

import (
	"github.com/jsfan/t3migrate/internal/config"
	"github.com/jsfan/t3migrate/internal/storage/model"
)

type DAL interface {
	Connect(config *config.MySQLConfig) error

	SelectContentElementsNotInList(ids []int64) ([]*model.ContentElement, error)
	SelectPages() ([]*model.TVFlex, error)

	DeleteContentByIds(ids []int64) error
}
