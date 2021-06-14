package mapping

import (
	"fmt"
	"github.com/jsfan/t3migrate/internal/storage/model"
	"github.com/spf13/cobra"
	"os"
	"regexp"
)

func MapSrcToDstCols(srcDesc, dstDesc map[string]*model.TableDescription) (ColumnMapping, []string, error) {
	simpleType, err := regexp.Compile(`([a-z]+(int|[a-z]+))\([^)]*\)(\s+(un)?signed)?`)
	if err != nil {
		cobra.CheckErr(fmt.Errorf(`could not compile regex for type comparison: %+v`, err))
	}

	mapping := make(map[string]*string, 0)
	for k, src := range srcDesc {
		if dst, ok := dstDesc[k]; ok {
			srcType := string(simpleType.ReplaceAll([]byte(src.Type), []byte("$2")))
			dstType := string(simpleType.ReplaceAll([]byte(dst.Type), []byte("$2")))
			if dstType == srcType {
				mapping[k] = &dst.Field
			} else {
				fmt.Fprintf(
					os.Stderr,
					`Found column "%s" in both source and destination table but types differ (%s != %s).`+"\n",
					k,
					src.Type,
					dst.Type,
				)
				mapping[k] = nil
			}
		} else {
			mapping[k] = nil
		}
	}
	unused := make([]string, 0)
	for k, v := range dstDesc {
		if src, ok := srcDesc[k]; !ok || src.Field != v.Field {
			unused = append(unused, k)
		}
	}
	return mapping, unused, nil
}
