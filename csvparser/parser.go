package csvparser

import (
	"strings"

	"github.com/KitlerUA/csvparser/policy"
)

const (
	sc1 = "SC1"
	sc2 = "SC2"
)

//Parse - read matrix of values and send Policies to the channel
func Parse(records [][]string, c chan policy.Policy) {
	defer close(c)
	//first two columns are sources and name of policy
	for j := 2; j > 0 && j < len(records[0]); j++ {
		tempSc1 := policy.Policy{
			Name:        "pn:fac:pz:" + strings.ToLower(records[0][j]),
			Description: "",
			Subjects:    []string{"gn:fac:" + strings.ToLower(records[0][j]) + " --need to be fixed"},
			Effect:      "allow",
			Conditions:  policy.Condition{},
			Resources:   []string{"rn:pz"},
		}
		tempSc2 := policy.Policy{
			Name:        "pn:fac:pc:" + strings.ToLower(records[0][j]),
			Description: "--need to be fixed",
			Subjects:    []string{"gn:fac:" + strings.ToLower(records[0][j]) + " --need to be fixed"},
			Effect:      "allow",
			Conditions:  policy.Condition{},
			Resources:   []string{"rn:pc"},
		}
		for i := range records {
			if i == 0 {
				continue
			}
			//find sources (may be 1 or 2 sources in one cell)
			temp := strings.Split(records[i][0], ",")
			for k := range temp {
				temp[k] = strings.TrimSpace(temp[k])
				if records[i][j] == "Yes" {
					if temp[k] == sc1 {
						tempSc1.Actions = append(tempSc1.Actions, records[i][1])
					} else if temp[k] == sc2 {
						tempSc2.Actions = append(tempSc2.Actions, records[i][1])
					}
				}
			}
		}
		c <- tempSc1
		c <- tempSc2
	}
}
