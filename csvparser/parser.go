package csvparser

import (
	"strings"

	"fmt"

	"github.com/KitlerUA/csvparser/policy"
)

const (
	sc1 = "SC1"
	sc2 = "SC2"
)

//Parse - read matrix of values and send Policies to the channel
func Parse(records [][]string, bindings [][]string, c chan policy.Policy, w chan string, q chan struct{}) {
	defer func() {
		q <- struct{}{}
	}()
	//first two columns are sources and name of policy
	for j := 2; j > 0 && j < len(records[0]); j++ {
		tempSc1 := policy.Policy{
			Name:        "pn:fac:pz:" + strings.ToLower(records[0][j]),
			Description: "",
			Subjects:    []string{"gn:fac:" + strings.ToLower(records[0][j])},
			Effect:      "allow",
			Conditions:  policy.Condition{},
			Resources:   []string{"rn:pz"},
			Actions:     make([]string, 0),
		}
		tempSc2 := policy.Policy{
			Name:        "pn:fac:pc:" + strings.ToLower(records[0][j]),
			Description: "",
			Subjects:    []string{"gn:fac:" + strings.ToLower(records[0][j])},
			Effect:      "allow",
			Conditions:  policy.Condition{},
			Resources:   []string{"rn:pc"},
			Actions:     make([]string, 0),
		}
		if j-2 >= len(bindings) {
			w <- fmt.Sprintf("cannot find binding name for '%s', default name used", strings.ToLower(records[0][j]))
		} else {
			tempSc1.Subjects[0] = bindings[j-2][1]
			tempSc1.Description = bindings[j-2][2]
			names := strings.Split(bindings[j-2][1], ":")
			tempSc1.Name = "pn:" + strings.ToLower(bindings[j-2][0]) + ":pz:" + strings.ToLower(names[len(names)-1])
			tempSc2.Subjects[0] = bindings[j-2][1]
			tempSc2.Description = bindings[j-2][2]
			tempSc2.Name = "pn:" + strings.ToLower(bindings[j-2][0]) + ":pc:" + strings.ToLower(names[len(names)-1])
		}
		for i := range records {
			if i == 0 {
				continue
			}
			//find sources (may be 1 or 2 sources in one cell)
			temp := strings.Split(records[i][0], ",")
			for k := range temp {
				temp[k] = strings.TrimSpace(temp[k])
				if strings.ToLower(records[i][j]) == strings.ToLower("Yes") {
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
