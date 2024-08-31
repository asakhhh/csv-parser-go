package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"a-library-for-others/csvparser"
)

func main() {
	file, err := os.Open("file.csv")
	if err != nil {
		fmt.Println("Error opening file: ", err)
		os.Exit(1)
	}
	defer file.Close()

	var newParser csvparser.CSVParser = csvparser.NewParser()
	for {
		_, err := newParser.ReadLine(file)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("Error reading line:", err)
			os.Exit(1)
		}
		var fields []string
		for i := 0; i < newParser.GetNumberOfFields(); i++ {
			field, err := newParser.GetField(i)
			if err != nil {
				fmt.Println("Error: ", err)
			}
			fields = append(fields, fmt.Sprintf("%r", field))
		}
		fmt.Println(strings.Join(fields, "|"))
	}
}
