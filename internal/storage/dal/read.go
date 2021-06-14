package dal

import (
	"database/sql"
	"fmt"
	"github.com/jsfan/t3migrate/internal/storage/model"
	"strconv"
	"strings"
)

// DescribeTable describes a table to enable mapping
func (dal *DataAccessLayer) DescribeTable(tableName string) ([]*model.TableDescription, error) {
	query := `DESCRIBE ` + tableName
	rows, err := dal.db.Query(query)
	if err != nil {
		return nil, err
	}
	cols := make([]*model.TableDescription, 0)
	for rows.Next() {
		col := &model.TableDescription{}
		err = rows.Scan(&col.Field, &col.Type, &col.Null, &col.Key, &col.Default, &col.Extra)
		if err != nil {
			return nil, err
		}
		cols = append(cols, col)
	}
	return cols, nil
}

// SelectTVFlexFromPages gets the TemplaVoila Plus flex elements from all pages
func (dal *DataAccessLayer) SelectTVFlexFromPages() ([]*model.TVFlex, error) {
	query := `SELECT uid, tx_templavoilaplus_flex FROM pages WHERE tx_templavoilaplus_flex != ''`
	rows, err := dal.db.Query(query)
	if err != nil {
		return nil, err
	}
	var flexList []*model.TVFlex
	for rows.Next() {
		flex := &model.TVFlex{}
		err := rows.Scan(&flex.Uid, &flex.TxTemplaVoilaPlusFlex)
		if err != nil {
			return nil, err
		}
		flexList = append(flexList, flex)
	}
	return flexList, nil
}

func (dal *DataAccessLayer) SelectAll(tableName string, columns []string, offset, limit *int64, includedIds []int64) ([][]interface{}, error) {
	query := "SELECT %s FROM `%s`"
	where := ""
	if exclusion := makeListCondition(includedIds); exclusion != "" {
		where = exclusion
	}
	if hasColumn(columns, "deleted") {
		if where != "" {
			where += " AND "
		} else {
			where += " WHERE "
		}
		where += "(deleted IS FALSE)"
	}
	if hasColumn(columns, "pid") {
		if where != "" {
			where += " AND "
		} else {
			where += " WHERE "
		}
		where += "(pid != -1)"
	}
	query += where
	if offset != nil {
		query += fmt.Sprintf(" LIMIT %d", *offset)
		if limit != nil {
			query += fmt.Sprintf("%d", limit)
		}
	}
	cols := "`" + strings.Join(columns, "`,`") + "`"
	query = fmt.Sprintf(query, cols, tableName)
	var rows *sql.Rows
	var err error
	if includedIds != nil {
		args := make([]interface{}, len(includedIds))
		for i, id := range includedIds {
			args[i] = &id
		}
		rows, err = dal.db.Query(query, args...)
	} else {
		rows, err = dal.db.Query(query)
	}
	if err != nil {
		return nil, err
	}
	results := make([][]interface{}, 0)
	for rows.Next() {
		values := make([]interface{}, len(columns))
		for i, _ := range columns {
			values[i] = new(interface{})
		}
		if err = rows.Scan(values...); err != nil {
			return nil, err
		}
		results = append(results, values)
	}
	return results, nil
}

func (dal *DataAccessLayer) CountAll(tableName string, columns []string, includedIds []int64) (int64, error) {
	query := "SELECT COUNT(*) c FROM `%s`"
	where := ""
	if exclusion := makeListCondition(includedIds); exclusion != "" {
		where = exclusion
	}
	if hasColumn(columns, "deleted") {
		if where != "" {
			where += " AND "
		} else {
			where += " WHERE "
		}
		where += "(deleted IS FALSE)"
	}
	if hasColumn(columns, "pid") {
		if where != "" {
			where += " AND "
		} else {
			where += " WHERE "
		}
		where += "(pid != -1)"
	}

	query += where
	query = fmt.Sprintf(query, tableName)
	var row *sql.Row
	if includedIds != nil {
		args := make([]interface{}, len(includedIds))
		strIds := make([]string, len(includedIds))
		for i, id := range includedIds {
			args[i] = &id
			strIds[i] = strconv.FormatInt(id, 10)
		}
		row = dal.db.QueryRow(query, args...)
	} else {
		row = dal.db.QueryRow(query)
	}
	count := int64(0)
	if err := row.Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}
