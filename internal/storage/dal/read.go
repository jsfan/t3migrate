package dal

import (
	"fmt"
	"github.com/jsfan/t3migrate/internal/storage/model"
	"strings"
)

// SelectPages gets the TemplaVoila Plus flex elements from all pages
func (dal *DataAccessLayer) SelectPages() ([]*model.TVFlex, error) {
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

// SelectContentElementsNotInList gets all content elements with UIDs not in the passed list
func (dal *DataAccessLayer) SelectContentElementsNotInList(ids []int64) ([]*model.ContentElement, error) {
	query := `SELECT pid, uid, header FROM tt_content WHERE uid NOT IN (%s)`
	placeholders := make([]string, len(ids))
	args := make([]interface{}, len(ids))
	for i := 0; i < len(ids); i++ {
		placeholders[i] = "?"
		args[i] = ids[i]
	}
	query = fmt.Sprintf(query, strings.Join(placeholders, ","))
	rows, err := dal.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	var elementList []*model.ContentElement
	for rows.Next() {
		content := &model.ContentElement{}
		err := rows.Scan(&content.Pid, &content.Uid, &content.Header)
		if err != nil {
			return nil, err
		}
		elementList = append(elementList, content)
	}
	return elementList, nil
}
