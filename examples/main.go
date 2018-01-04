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

const defaultCSVFileName = "examples/qwre.csv"
const directoryNameForJSONs = "Policies"

func main() {
	parser := csvparser.Parser{}
	readerChan := make(chan policy.Policy, 4)
	errorChan := make(chan error)
	if len(os.Args) > 1 {
		go parser.Parse(os.Args[1], ',', readerChan, errorChan)
	} else {
		go parser.Parse(defaultCSVFileName, ',', readerChan, errorChan)
	}
	if err := <-errorChan; err != nil {
		log.Printf("Get parse error: %s", err)
		return
	}

	for c := range readerChan {
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
