package csvparser

import (
	"strings"

	"log"

	"github.com/KitlerUA/csvparser/config"
	"github.com/tealeg/xlsx"
)

func XLSX(fileName string) (map[string][][]string, error) {
	xlFile, err := xlsx.OpenFile(fileName)
	if err != nil {
		return make(map[string][][]string), err
	}
	res := make(map[string][][]string)
	for _, sheet := range xlFile.Sheets {
		//skip empty sheet
		if len(sheet.Rows) < 1 {
			continue
		}
		var (
			//index of Page column
			page      int
			pageFound = false
			//index of Name column
			name      int
			nameFound = false
			////indices of Roles columns
			roles      []int
			rolesFound = false
		)
		//search for "page", "name" and "roles" in first row
		row := sheet.Rows[0]
		for j, cell := range row.Cells {
			if strings.ToLower(cell.String()) == strings.ToLower(config.Get().Page) {
				page = j
				pageFound = true
			} else if strings.ToLower(cell.String()) == strings.ToLower(config.Get().Name) {
				name = j
				nameFound = true
			} else if config.Get().ContainsRole(strings.ToLower(cell.String())) {
				roles = append(roles, j)
				rolesFound = true
			}
		}
		//check if all headers was found
		if !(pageFound && nameFound && rolesFound) {
			log.Printf("Cannot find %s, %s or at least one of the roles in firs row of '%s' sheet", config.Get().Page, config.Get().Name, sheet.Name)
			continue
		}
		//create new record in map after all checks
		res[sheet.Name] = make([][]string, 0)
		//add record for this sheet
		res[sheet.Name] = make([][]string, 0)
		for i, row := range sheet.Rows {
			//first empty row mean the end of the table
			if isRowEmpty(row) {
				break
			}
			//add new row
			res[sheet.Name] = append(res[sheet.Name], []string{})
			//insert Page
			res[sheet.Name][i] = append(res[sheet.Name][i], row.Cells[page].String())
			//insert Name
			res[sheet.Name][i] = append(res[sheet.Name][i], row.Cells[name].String())
			//insert Roles
			for _, r := range roles {
				res[sheet.Name][i] = append(res[sheet.Name][i], row.Cells[r].String())
			}
		}
	}
	return res, nil
}

func isRowEmpty(row *xlsx.Row) bool {
	if len(row.Cells) == 0 {
		return true
	}
	for _, r := range row.Cells {
		if r.String() != "" {
			return false
		}
	}
	return true
}
