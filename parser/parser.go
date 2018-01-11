package parser

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"
	"path"

	"fmt"

	"time"

	"github.com/KitlerUA/csvparser/config"
	"github.com/KitlerUA/csvparser/csvparser"
	"github.com/KitlerUA/csvparser/policy"
)

//Parse - read and parse file with fileName
//write results to dir
//also returns warning message
func Parse(fileName, dir string) (error, string) {
	var warn string
	cErr := make(chan error)
	go config.Init(cErr)
	if e := <-cErr; e != nil {
		return e, ""
	}
	var (
		m   map[string][][]string
		b   map[string][][]string
		err error
		ext string
	)

	ext = path.Ext(fileName)

	switch ext {
	case ".xlsx":
		m, b, err = csvparser.XLSX(fileName)
		if err != nil {
			return fmt.Errorf("cannot parse xlsx: %s", err), ""
		}
	default:
		return fmt.Errorf("format of file isn`t supported"), ""
	}

	for k := range m {
		//channel for getting Policies from parser.Parse
		readerChan := make(chan policy.Policy, 4)
		warnChan := make(chan string, 4)
		quitChan := make(chan struct{})
		go csvparser.Parse(m[k], b[k], readerChan, warnChan, quitChan)
		//if directory already exists we get error, but we need just skip this action, not panic
		dirName := dir + time.Now().Format("2006-01-02_15-04-05") + "_" + k
		if err := os.Mkdir(dirName, os.ModePerm); err != nil && !os.IsExist(err) {
			return fmt.Errorf("cannot create directory for policies: %s", err), ""
		}
	Listener:
		for {
			select {
			case c := <-readerChan:
				marshaledPolicies, err := json.Marshal(&c)
				if err != nil {
					return fmt.Errorf("cannot marshal policy '%s' : %s", c.Name, err), ""
				}
				newName := ReplaceRuneWith(c.Name, ':', '_')
				newName = ReplaceRuneWith(newName, '*', '_')
				if err = ioutil.WriteFile(dirName+"/"+newName+".json", marshaledPolicies, 0666); err != nil {
					return fmt.Errorf("cannot save json file for policy '%s': %s", c.Name, err), ""
				}
			case w := <-warnChan:
				warn += k + ": " + w + "\n"
			case <-quitChan:
				break Listener
			}
		}

	}
	return nil, warn
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
