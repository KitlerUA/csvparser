package parser

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"
	"path"

	"fmt"

	"github.com/KitlerUA/csvparser/csvparser"
	"github.com/KitlerUA/csvparser/policy"
)

//Parse - read and parse file with fileName
//write results to dir
func Parse(fileName, dir string) error {
	var (
		m   map[string][][]string
		err error
		ext string
	)

	ext = path.Ext(fileName)

	if ext == ".csv" {
		m, err = csvparser.CSV(fileName, ';')
		if err != nil {
			return fmt.Errorf("cannot parse csv: %s", err)
		}
	}
	if ext == ".xlsx" {
		m, err = csvparser.XLSX(fileName)
		if err != nil {
			return fmt.Errorf("cannot parse xlsx: %s", err)
		}
	}
	for k := range m {
		//channel for getting Policies from parser.Parse
		readerChan := make(chan policy.Policy, 4)
		go csvparser.Parse(m[k], readerChan)
		//if directory already exists we get error, but we need just skip this action, not panic
		if err := os.Mkdir(dir+k, os.ModePerm); err != nil && !os.IsExist(err) {
			return fmt.Errorf("cannot create directory for policies: %s", err)
		}
		for c := range readerChan {
			marshaledPolicies, err := json.Marshal(&c)
			if err != nil {
				return fmt.Errorf("cannot marshal policy '%s' : %s", c.Name, err)
			}
			newName := ReplaceRuneWith(c.Name, ':', '_')
			if err = ioutil.WriteFile(dir+k+"/"+newName+".json", marshaledPolicies, 0666); err != nil {
				return fmt.Errorf("cannot save json file for policy '%s': %s", c.Name, err)
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
