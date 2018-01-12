package csvparser

import (
	"strings"

	"fmt"

	"github.com/KitlerUA/csvparser/config"
	"github.com/KitlerUA/csvparser/policy"
)

const (
	sc1     = "SC1"
	sc2     = "SC2"
	missing = "MISSING"
)

//Parse - read matrix of values and send Policies to the channel
func Parse(records [][]string, bindings [][]string, c chan policy.Policy, w chan string, q chan struct{}) {
	defer func() {
		q <- struct{}{}
	}()
	if len(records) == 1 {
		w <- fmt.Sprintf("action table is empty")
	}
	//count of prefix parameters like source and name of policy
	prefLen := 2
	//current binding position
	curBind := 0
	//first two columns are sources and name of policy
	for j := prefLen; j > 0 && j < len(records[0]); j++ {
		srcPol := make(map[string]*policy.Policy)
		//walk down to collect info about actions and form policy for all role-sources pair for current role
		for i := 1; i < len(records); i++ {
			sources := strings.Split(records[i][0], ",")
			for s := range sources {
				src := strings.ToLower(strings.TrimSpace(sources[s]))
				if src == "" {
					w <- fmt.Sprintf("found empty page on row %d", i+1)
					continue
				}
				//check if source in config list
				if _, ok := config.Get().PagesNames[src]; !ok {
					w <- fmt.Sprintf("page '%s' (row %d) isn't in config file: skipped", src, i+1)
					continue
				}
				//if record for source doesn't exist - create
				if _, ok := srcPol[src]; !ok {
					var name, description, subject string

					//take info from table, otherwise - send warning and set fields 'missing value'
					if curBind < len(bindings) {
						names := strings.Split(bindings[curBind][1], ":")
						name = fmt.Sprintf("pn:%s:%s:%s", strings.ToLower(bindings[curBind][0]), strings.ToLower(config.Get().PagesNames[src]), strings.ToLower(names[len(names)-1]))
						description = bindings[curBind][2]
						subject = bindings[curBind][1]
					} else {
						w <- fmt.Sprintf("cannot find binding name for '%s'", strings.ToLower(records[0][j]))
						name = missing
						description = missing
						subject = missing
					}
					srcPol[src] = &policy.Policy{
						Name:        name,
						Description: description,
						Subjects:    []string{subject},
						Effect:      "allow",
						Conditions:  policy.Condition{},
						Resources:   []string{fmt.Sprintf("rn:%s", strings.ToLower(config.Get().PagesNames[src]))},
						Actions:     make([]string, 0),
					}
				}
				if strings.ToLower(records[i][j]) == strings.ToLower("Yes") {
					srcPol[src].Actions = append(srcPol[src].Actions, records[i][1])
				}
			}
		}
		for _, v := range srcPol {
			c <- *v
		}
		curBind++
	}
}
