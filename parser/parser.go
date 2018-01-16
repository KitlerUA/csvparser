package parser

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"fmt"

	"time"

	"github.com/KitlerUA/csvparser/config"
	"github.com/KitlerUA/csvparser/csvparser"
	"github.com/KitlerUA/csvparser/policy"
)

//Parse - read and parse file with fileName
//write results to dir
//also returns warning message
func Parse(fileName, dir string) (string, error) {
	var warn string
	if dir == "" {
		var err error
		if dir, err = filepath.Abs(filepath.Dir(os.Args[0])); err != nil {
			return warn, err
		}
		dir += "/"
	} else {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			return warn, err
		}
	}

	if e := config.Init(); e != nil {
		return warn, e
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
		m, b, warn, err = csvparser.XLSX(fileName)
		if err != nil {
			return warn, fmt.Errorf("cannot parse xlsx: %s", err)
		}
	default:
		return warn, fmt.Errorf("format of file isn`t supported")
	}

	for k := range m {
		//channel for getting Policies from parser.Parse
		readerChan := make(chan policy.Policy)
		warnChan := make(chan string)
		quitChan := make(chan struct{})
		go csvparser.Parse(m[k], b[k], readerChan, warnChan, quitChan)
		//if directory already exists we get error, but we need just skip this action, not panic
		dirName := dir + time.Now().Format("2006-01-02_15-04-05") + "_" + k
		if err := os.Mkdir(dirName, os.ModePerm); err != nil && !os.IsExist(err) {
			return warn, fmt.Errorf("cannot create directory for policies: %s", err)
		}
	Listener:
		for {
			select {
			case c := <-readerChan:
				marshaledPolicies, err := json.Marshal(&c)
				if err != nil {
					return warn, fmt.Errorf("cannot marshal policy '%s' : %s", c.Name, err)
				}
				newName := ReplaceRuneWith(c.FileName, ':', '_')
				newName = ReplaceRuneWith(newName, '*', '_')
				if err = ioutil.WriteFile(dirName+"/"+newName+".json", marshaledPolicies, 0666); err != nil {
					return warn, fmt.Errorf("cannot save json file for policy '%s': %s", c.Name, err)
				}
			case w := <-warnChan:
				warn += fmt.Sprintf("<b>%s</b>: %s<br>", k, w)
			case <-quitChan:
				break Listener
			}
		}

	}
	for k := range m {
		delete(m, k)
	}
	for k := range b {
		delete(b, k)
	}
	return warn, nil
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
