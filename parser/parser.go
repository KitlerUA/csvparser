package parser

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path"

	"errors"
	"fmt"

	"github.com/KitlerUA/CSVParser/csvparser"
	"github.com/KitlerUA/CSVParser/policy"
)

func Parse(fileName, dir string) error {
	var m map[string][][]string
	var err error
	var ext string

	ext = path.Ext(fileName)

	if ext == ".csv" {
		m, err = csvparser.CSV(fileName, ';')
		if err != nil {
			return errors.New(fmt.Sprintf("Cannot parse csv: %s", err))
		}
	}
	if ext == ".xlsx" {
		m, err = csvparser.XLSX(fileName)
		if err != nil {
			return errors.New(fmt.Sprintf("Cannot parse xlsx: %s", err))
		}
	}
	for k := range m {
		//channel for getting Policies from parser.Parse
		readerChan := make(chan policy.Policy, 4)
		go csvparser.Parse(m[k], readerChan)
		//if directory already exists we get error, but we need just skip this action, not panic
		if err := os.Mkdir(k, os.ModePerm); err != nil && !os.IsExist(err) {
			return errors.New(fmt.Sprintf("Cannot create directory for policies: %s", err))
		}
		for c := range readerChan {
			marshaledPolicies, err := json.Marshal(&c)
			if err != nil {
				log.Fatalf("Cannot marshal policy '%s' : %s", c.Name, err)
			}
			newName := ReplaceRuneWith(c.Name, ':', '_')
			if err = ioutil.WriteFile(dir+k+"/"+newName+".json", marshaledPolicies, 0666); err != nil {
				log.Fatalf("Cannot save json file for policy '%s': %s", c.Name, err)
			}

		}
	}
	return nil
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
