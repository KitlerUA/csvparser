package main

import (
	"log"

	"github.com/KitlerUA/csvparser/parser"
)

const defaultCSVFileName = "List of Actions Test2.xlsx"

func main() {
	var warnings string
	var err error
	if err, warnings = parser.Parse(defaultCSVFileName, ""); err != nil {
		log.Fatalf("%s", err)
	}
	if warnings != "" {
		log.Printf("Parsed with warnings:\n%s", warnings)
	} else {
		log.Printf("Successfully parsed and saved")
	}
}
