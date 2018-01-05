package csvparser

import (
	"encoding/csv"
	"os"
	"strings"

	"github.com/KitlerUA/CSVParser/policy"
)

type Parser struct {
}

//Parse - return slice of policies from given .csv file
func (p *Parser) Parse(fileName string, delimiter rune, c chan policy.Policy, cErr chan error) {
	defer close(c)
	file, err := os.OpenFile(fileName, os.O_RDONLY, 0666)
	if err != nil {
		cErr <- err
	}
	defer file.Close()
	reader := csv.NewReader(file)
	reader.Comma = delimiter
	records, err := reader.ReadAll()
	if err != nil {
		cErr <- err
	}
	cErr <- nil
	var result []policy.Policy
	for i := range records {
		temp := policy.Policy{
			Name:        records[i][0],
			Description: records[i][1],
			Subjects:    []string{},
			Actions:     []string{},
			Conditions:  policy.Condition{},
			Effect:      records[i][4],
			Resources:   []string{},
			Row:         i,
		}
		//parse subjects
		temp.Subjects = strings.Split(records[i][2], string(delimiter))
		for i := range temp.Subjects {
			temp.Subjects[i] = strings.TrimSpace(temp.Subjects[i])
			if temp.Subjects[i] == "" {
				temp.Subjects = append(temp.Subjects[:i], temp.Subjects[i+1:]...)
				continue
			}
		}
		//parse actions
		temp.Actions = strings.Split(records[i][3], string(delimiter))
		//removeEmptyAndSpaces(&temp.Actions)
		for i := len(temp.Actions) - 1; i >= 0; i-- {
			temp.Actions[i] = strings.TrimSpace(temp.Actions[i])
			if temp.Actions[i] == "" {
				temp.Actions = append(temp.Actions[:i], temp.Actions[i+1:]...)
				continue
			}
		}
		//parse resources
		temp.Resources = strings.Split(records[i][6], string(delimiter))
		for i := range temp.Resources {
			temp.Resources[i] = strings.TrimSpace(temp.Resources[i])
			if temp.Resources[i] == "" {
				temp.Resources = append(temp.Resources[:i], temp.Resources[i+1:]...)
				continue
			}
		}
		//append temp Policy to result slice
		result = append(result, temp)
		c <- temp
	}
}
