package dal

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang/glog"
	"github.com/jsfan/t3migrate/internal/config"
)

type DataAccessLayer struct {
	db *sql.DB
}

// Connect connects a DAL to a database using the passed config
func (dal *DataAccessLayer) Connect(config *config.MySQLConfig) error {
	if config.Port == 0 {
		config.Port = 3306
		glog.Infof("Setting default port %d", config.Port)
	}
	connStr := fmt.Sprintf(
		`%s:%s@tcp(%s:%d)/%s`,
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.Database,
	)
	var err error
	dal.db, err = sql.Open("mysql", connStr)
	if err != nil {
		return err
	}
	return nil
}
