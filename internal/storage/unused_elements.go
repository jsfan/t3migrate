package storage

import (
	"fmt"
	"github.com/golang/glog"
	"github.com/jsfan/t3migrate/internal/xml"
	"strconv"
	"strings"
)

func (store *Store) DeleteUnusedElements() error {
	pages, err := store.dal.SelectPages()
	if err != nil {
		return err
	}
	elementSet := make(map[string]struct{}, 0)
	for _, page := range pages {
		form, err := xml.ParseFlexForm(page.TxTemplaVoilaPlusFlex)
		if err != nil {
			glog.Errorf("Could not parse XML: %+v", string(page.TxTemplaVoilaPlusFlex))
			continue
		}
		for _, f := range form.Data.Sheet.Language.Field {
			if f.Index == "field_content" {
				for _, el := range strings.Split(f.Value, ",") {
					elementSet[el] = struct{}{}
				}
			}
		}
	}
	elementList := make([]int64, 0)
	for idStr, _ := range elementSet {
		if idStr == "" {
			continue
		}
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			return err
		}
		elementList = append(elementList, id)
	}
	unusedElements, err := store.dal.SelectContentElementsNotInList(elementList)
	unusedElementList := make([]int64, len(unusedElements))
	for i, contentElement := range unusedElements {
		glog.Infof("Content element %+v is unused and will be deleted.", contentElement)
		unusedElementList[i] = contentElement.Uid
	}
	if !store.dry {
		glog.Info("Deleting all unreferenced content elements.")
		err := store.dal.DeleteContentByIds(unusedElementList)
		if err != nil {
			return fmt.Errorf("failed to delete unused content elements: %+v", err)
		}
	}
	return nil
}
