package csvparser

import (
	"strings"

	"log"

	"github.com/KitlerUA/csvparser/config"
	"github.com/tealeg/xlsx"
)

//XLSX - read data from .xlsx file and
//return map with key=sheet.Name and matrix of values
//read until first empty row
func XLSX(fileName string) (map[string][][]string, map[string][][]string, error) {
	xlFile, err := xlsx.OpenFile(fileName)
	if err != nil {
		return make(map[string][][]string), make(map[string][][]string), err
	}
	res := make(map[string][][]string)
	bindings := make(map[string][][]string)
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
			//indices of Roles columns
			//roles      []int
			//rolesFound = false
			//indices for roles_start and roles_end
			rolesStart      int
			rolesStartFound = false
			rolesEnd        int
			rolesEndFound   = false
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
			} else if strings.ToLower(cell.String()) == strings.ToLower(config.Get().RolesBegin) {
				rolesStart = j + 1
				rolesStartFound = true
			} else if strings.ToLower(cell.String()) == strings.ToLower(config.Get().RolesEnd) {
				rolesEnd = j - 1
				rolesEndFound = true
			}
		}
		//check if all headers was found
		if !(pageFound && nameFound && rolesStartFound && rolesEndFound) {
			log.Printf("Cannot find %s, %s or bounds fore roles in '%s' sheet", config.Get().Page, config.Get().Name, sheet.Name)
			continue
		}
		//create new record in map after all checks
		res[sheet.Name] = make([][]string, 0)
		//position of the end of table
		var tableEndRow int

		for i, row := range sheet.Rows {
			//first empty row mean the end of the table
			if isRowEmpty(row) {
				tableEndRow = i
				break
			}
			//add new row
			res[sheet.Name] = append(res[sheet.Name], []string{})
			//insert Page
			res[sheet.Name][i] = append(res[sheet.Name][i], row.Cells[page].String())
			//insert Name
			res[sheet.Name][i] = append(res[sheet.Name][i], row.Cells[name].String())
			//insert Roles
			for j := rolesStart; j <= rolesEnd; j++ {
				res[sheet.Name][i] = append(res[sheet.Name][i], row.Cells[j].String())
			}
		}
		//try to find subjects table
		for i := tableEndRow + 1; i < len(sheet.Rows); i++ {
			for j := 0; j < len(sheet.Rows[i].Cells)-2; j++ {
				cell := sheet.Rows[i].Cells[j].String()
				cell1 := sheet.Rows[i].Cells[j+1].String()
				cell2 := sheet.Rows[i].Cells[j+2].String()
				if strings.ToLower(cell) == strings.ToLower(config.Get().Type) &&
					strings.ToLower(cell1) == strings.ToLower(config.Get().TechGroupName) &&
					strings.ToLower(cell2) == strings.ToLower(config.Get().DisplayName) {
					for r := i + 1; r < len(sheet.Rows); r++ {
						if isPartRowEmpty(sheet.Rows[r], j, j+2) {
							break
						}
						bindings[sheet.Name] = append(bindings[sheet.Name], []string{sheet.Rows[r].Cells[j].String(), sheet.Rows[r].Cells[j+1].String(), sheet.Rows[r].Cells[j+2].String()})
					}
				}
			}

		}
	}
	return res, bindings, nil
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

func isPartRowEmpty(row *xlsx.Row, a, b int) bool {
	if len(row.Cells) == 0 {
		return true
	}
	for i := a; i < b && i < len(row.Cells); i++ {
		if row.Cells[i].String() != "" {
			return false
		}
	}
	return true
}
