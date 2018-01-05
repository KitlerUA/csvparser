package csvparser

import (
	"github.com/tealeg/xlsx"
)

func XLSX(fileName string) (map[string][][]string, error) {
	xlFile, err := xlsx.OpenFile(fileName)
	if err != nil {
		return make(map[string][][]string), err
	}
	res := make(map[string][][]string)
	for _, sheet := range xlFile.Sheets {
		res[sheet.Name] = make([][]string, 0)
		for i, row := range sheet.Rows {
			res[sheet.Name] = append(res[sheet.Name], []string{})
			for _, cell := range row.Cells {
				res[sheet.Name][i] = append(res[sheet.Name][i], cell.String())
			}
		}
	}
	return res, nil
}
