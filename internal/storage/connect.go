package storage

import (
	"github.com/jsfan/t3migrate/internal/config"
	"github.com/jsfan/t3migrate/internal/storage/dal"
)

type Store struct {
	dal DAL
	dry bool
}

func (store *Store) Connect(config *config.MySQLConfig, dry bool) error {
	store.dal = &dal.DataAccessLayer{}
	store.dry = dry
	return store.dal.Connect(config)
}
