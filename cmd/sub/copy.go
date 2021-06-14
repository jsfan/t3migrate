/*
Copyright © 2021 Christian Lerrahn <github@penpal4u.net>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package sub

import (
	"fmt"
	"github.com/jsfan/t3migrate/internal/mapping"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

var tableList []string
var mapFile *string
var writeMappings *bool

// copyCmd represents the copy command
var copyCmd = &cobra.Command{
	Use:   "copy",
	Short: "Copy data filtering everything that isn't used or deleted.",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if *writeMappings {
			for _, t := range tableList {
				srcMap, err := srcStore.MapTable(t)
				if err != nil {
					cobra.CheckErr(fmt.Errorf(`could not read mapping for table "%s" from source database: %+v`, t, err))
				}
				dstMap, err := dstStore.MapTable(t)
				if err != nil {
					cobra.CheckErr(fmt.Errorf(`could not read mapping for table "%s" from destination database: %+v`, t, err))
				}
				mapped, unused, err := mapping.MapSrcToDstCols(srcMap, dstMap)
				if err != nil {
					cobra.CheckErr(fmt.Errorf(`could not map columns for table "%s": %+v`, t, err))
				}
				fmt.Fprintf(os.Stderr, `The following columns from the source table "%s" could not be mapped to the destination table:
%s
`, t, strings.Join(unused, ","))
				if err := mapping.SaveMappingToFile(*mapFile, t, mapped, *dryRun); err != nil {
					cobra.CheckErr(fmt.Errorf(`could not write to mapping file: %+v`, err))
				}
			}
			return
		}
		mapped, err := mapping.LoadMappings(*mapFile)
		if err != nil {
			cobra.CheckErr(fmt.Errorf(`could not load column mappings rom file "%s": %+v`, *mapFile, err))
		}
		for _, t := range tableList {
			var includeList []int64
			if t == "tt_content" {
				includeList, err = srcStore.FindUsedElements()
			}
			srcCols := make([]string, 0)
			dstCols := make([]string, 0)
			for s, d := range mapped[t] {
				if d != nil {
					srcCols = append(srcCols, s)
					dstCols = append(dstCols, *d)
				}
			}
			count, err := srcStore.CountRecords(t, srcCols, includeList)
			if err != nil {
				cobra.CheckErr(fmt.Errorf(`could not count records for table "%s": %+v`, t, err))
			}
			limit := int64(1000)
			var i int64
			for i = 0; i < (count/limit + 1); i++ {
				offset := i * limit
				res, err := srcStore.GetRecords(t, srcCols, &offset, &limit, includeList)
				if err != nil {
					cobra.CheckErr(fmt.Errorf(`could not retrieve records from source table "%s": %+v`, t, err))
				}
				if err := dstStore.PutRecords(t, dstCols, res); err != nil {
					cobra.CheckErr(fmt.Errorf(`could not write records to destination table "%s": %+v`, t, err))
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(copyCmd)
	var tables string
	copyCmd.PersistentFlags().StringVar(&tables, "tables", "pages,tt_content", "Tables to copy")
	tableList = strings.Split(tables, ",")
	for i, t := range tableList {
		tableList[i] = strings.Trim(t, " \t")
	}
	mapFile = copyCmd.PersistentFlags().String("mappingFile", "./mapping.yaml", "Mapping file")
	writeMappings = copyCmd.PersistentFlags().BoolP("write", "w", false, "Write mapping file")
}
