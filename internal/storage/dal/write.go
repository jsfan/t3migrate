package dal

import (
	"fmt"
	"strings"
)

// DeleteContentByIds deletes the content elements whose IDs were passed
func (dal *DataAccessLayer) DeleteContentByIds(ids []int64) error {
	query := `DELETE FROM tt_content WHERE uid IN (%s)`
	placeholders := make([]string, len(ids))
	args := make([]interface{}, len(ids))
	for i := 0; i < len(ids); i++ {
		placeholders[i] = "?"
		args[i] = ids[i]
	}
	query = fmt.Sprintf(query, strings.Join(placeholders, ","))
	_, err := dal.db.Exec(query, args...)
	return err
}
