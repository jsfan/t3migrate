package storage

import (
	"fmt"
	"github.com/jsfan/t3migrate/internal/xml"
	"os"
	"strconv"
	"strings"
)

func (store *Store) FindUsedElements() ([]int64, error) {
	pages, err := store.dal.SelectTVFlexFromPages()
	if err != nil {
		return nil, err
	}
	elementSet := make(map[string]struct{}, 0)
	for _, page := range pages {
		form, err := xml.ParseFlexForm(page.TxTemplaVoilaPlusFlex)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not parse XML: %+v", string(page.TxTemplaVoilaPlusFlex))
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
			return nil, err
		}
		elementList = append(elementList, id)
	}
	return elementList, nil
}
