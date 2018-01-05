package csvparser

import (
	"strings"

	"github.com/KitlerUA/CSVParser/policy"
)

const (
	sc1 = "SC1"
	sc2 = "SC2"
)

func Parse(records [][]string, c chan policy.Policy) {
	defer close(c)
	//log.Printf("%d", len(records[0]))
	for j := 2; j > 0 && j < len(records[0])-3; j++ {
		temp_sc1 := policy.Policy{
			Name:        records[0][j] + "_pz",
			Description: "",
			Subjects:    []string{records[0][j]},
			Effect:      "allow",
			Conditions:  policy.Condition{},
			Resources:   []string{"pz"},
		}
		temp_sc2 := policy.Policy{
			Name:        records[0][j] + "_pc",
			Description: "",
			Subjects:    []string{records[0][j]},
			Effect:      "allow",
			Conditions:  policy.Condition{},
			Resources:   []string{records[0][j]},
		}
		for i := range records {
			if i == 0 {
				continue
			}

			temp := strings.Split(records[i][0], ",")
			for k := range temp {
				if temp[k] == sc1 {
					temp_sc1.Actions = append(temp_sc1.Actions, records[i][1])
				} else if temp[k] == sc2 {
					temp_sc2.Actions = append(temp_sc2.Actions, records[i][1])
				}
			}
		}
		c <- temp_sc1
		c <- temp_sc2
	}
}
