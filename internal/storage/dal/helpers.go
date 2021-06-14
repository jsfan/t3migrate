package dal

import (
	"fmt"
	"strings"
)

func placeholderList(count int) []string {
	placeholders := make([]string, count)
	for i := 0; i < count; i++ {
		placeholders[i] = "?"
	}
	return placeholders
}

func updateStmtList(columns []string) []string {
	count := len(columns)
	placeholders := make([]string, count)
	for i := 0; i < count; i++ {
		placeholders[i] = columns[i] + " = ?"
	}
	return placeholders
}

func makeListCondition(ids []int64) string {
	if ids != nil {
		placeholders := placeholderList(len(ids))
		return fmt.Sprintf(" WHERE uid NOT IN (%s)", strings.Join(placeholders, ","))
	}
	return ""
}

func hasColumn(cols []string, colName string) bool {
	for _, c := range cols {
		if c == "deleted" {
			return true
		}
	}
	return false
}
