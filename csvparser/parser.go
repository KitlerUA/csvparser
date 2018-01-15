package csvparser

import (
	"strings"

	"fmt"

	"github.com/KitlerUA/csvparser/config"
	"github.com/KitlerUA/csvparser/policy"
)

const (
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
	//empty column header error
	headerErrors := make(map[int]struct{})
	//empty page errors
	emptyErrors := make(map[int]struct{})
	//missing sources in config
	misErrors := make(map[string]map[int]struct{})
	//missing binding names
	bindErrors := make(map[string]struct{})
	//first two columns are sources and name of policy
	for j := prefLen; j > 0 && j < len(records[0]); j++ {
		//check header of column, if empty - skipp column
		if records[0][j] == "" {
			headerErrors[j] = struct{}{}
			continue
		}
		srcPol := make(map[string]*policy.Policy)
		//walk down to collect info about actions and form policy for all role-sources pair for current role
		for i := 1; i < len(records); i++ {
			sources := strings.Split(records[i][0], ",")
			for s := range sources {
				src := strings.ToLower(strings.TrimSpace(sources[s]))
				if src == "" {
					emptyErrors[i+1] = struct{}{}
					continue
				}
				//check if source in config list
				if _, ok := config.Get().PagesNames[src]; !ok {
					if _, ok := misErrors[src]; !ok {
						misErrors[src] = make(map[int]struct{})
					}
					misErrors[src][i+1] = struct{}{}
					continue
				}
				//if record for source doesn't exist - create
				if _, ok := srcPol[src]; !ok {
					var name, description, subject, fileName string

					//take info from table, otherwise - send warning and set fields 'missing value'
					if curBind < len(bindings) {
						names := strings.Split(bindings[curBind][1], ":")
						name = fmt.Sprintf("pn:%s:%s:%s", strings.ToLower(bindings[curBind][0]), strings.ToLower(config.Get().PagesNames[src]), strings.ToLower(names[len(names)-1]))
						description = bindings[curBind][2]
						subject = bindings[curBind][1]
						fileName = fmt.Sprintf("%s_%s", strings.ToLower(names[len(names)-1]), config.Get().PagesNames[src])

					} else {
						bindErrors[strings.ToLower(records[0][j])] = struct{}{}
						name = fmt.Sprintf("%s:%s", strings.ToLower(records[0][j]), missing)
						description = missing
						subject = missing
						fileName = fmt.Sprintf("%s_%s", strings.ToLower(records[0][j]), missing)
					}
					srcPol[src] = &policy.Policy{
						Name:        name,
						Description: description,
						Subjects:    []string{subject},
						Effect:      "allow",
						Conditions:  policy.Condition{},
						Resources:   []string{fmt.Sprintf("rn:%s", strings.ToLower(config.Get().PagesNames[src]))},
						Actions:     make([]string, 0),
						FileName:    fileName,
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
	for i := range headerErrors {
		w <- fmt.Sprintf("find empty role-header on %d position", i-prefLen+1)
	}
	for i := range bindErrors {
		w <- fmt.Sprintf("cannot find binding name for '%s'", i)
	}
	for i := range emptyErrors {
		w <- fmt.Sprintf("found empty page-field on row %d", i)
	}
	for s := range misErrors {
		for r := range misErrors[s] {
			w <- fmt.Sprintf("page '%s' (row %d) isn't in config file: skipped", s, r)
		}
	}
}
