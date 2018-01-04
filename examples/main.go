package main

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"os"

	"bytes"

	"github.com/KitlerUA/CSVParser/csvparser"
	"github.com/KitlerUA/CSVParser/policy"
)

const defaultCSVFileName = "policies.csv"
const directoryNameForJSONs = "Policies"

func main() {
	parser := csvparser.Parser{}
	//channel for getting Policies from parser.Parse
	readerChan := make(chan policy.Policy, 4)
	//channels for getting error from parser.Parse
	errorChan := make(chan error)
	//if no argument - use default filename
	if len(os.Args) > 1 {
		go parser.Parse(os.Args[1], ',', readerChan, errorChan)
	} else {
		log.Printf("No argument for filename. Use default filename")
		go parser.Parse(defaultCSVFileName, ',', readerChan, errorChan)
	}
	//wait for error
	if err := <-errorChan; err != nil {
		log.Printf("Get parse error: %s", err)
		return
	}

	uniquePolicies := make(map[string]struct{})
	//if error == nil - just start receive Policies
	for c := range readerChan {
		if _, ok := uniquePolicies[c.Name]; ok {
			log.Printf("Find duplicate for policy '%s' on row %d. Skipped", c.Name, c.Row)
			continue
		}
		uniquePolicies[c.Name] = struct{}{}
		marshaledPolicies, err := json.Marshal(&c)
		if err != nil {
			log.Printf("Cannot marshal csv: %s", err)
			return
		}
		newName := ReplaceRuneWith(c.Name, ':', '_')
		os.Mkdir(directoryNameForJSONs, os.ModePerm)
		if err = ioutil.WriteFile(directoryNameForJSONs+"/"+newName+".json", marshaledPolicies, 0666); err != nil {
			log.Printf("Cannot save json file: %s", err)
			return
		}

	}

}

//ReplaceRuneWith - return copy of string with changed rune1 to rune2
func ReplaceRuneWith(str string, char1, char2 rune) string {
	buffer := bytes.Buffer{}
	for _, c := range str {
		if c == char1 {
			buffer.WriteRune(char2)
		} else {
			buffer.WriteRune(c)
		}
	}
	return buffer.String()
}
