package dal

import (
	"fmt"
	"strings"
)

func (dal *DataAccessLayer) Upsert(tableName string, columns []string, values []interface{}) error {
	query := "INSERT INTO `%s` (%s) VALUES (%s) ON DUPLICATE KEY UPDATE %s"
	cols := "`" + strings.Join(columns, "`,`") + "`"
	insertPlaceholders := strings.Join(placeholderList(len(values)), ",")
	updatePlaceholders := strings.Join(updateStmtList(columns), ", ")
	query = fmt.Sprintf(query, tableName, cols, insertPlaceholders, updatePlaceholders)
	valCount := len(values)
	duplicateValues := make([]interface{}, valCount*2)
	for i, v := range values {
		duplicateValues[i] = v
		duplicateValues[valCount+i] = v
	}
	if _, err := dal.db.Exec(query, duplicateValues...); err != nil {
		return err
	}
	return nil
}
