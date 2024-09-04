package main

import (
	"a-library-for-others/csvparser"
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	var s string
	fmt.Scanf("%s", &s)

	file, err := os.Open(s + ".csv")
	if err != nil {
		fmt.Println("Error opening file: ", err)
		os.Exit(1)
	}
	defer file.Close()

	var parser csvparser.CSVParser = csvparser.NewParser()
	for {
		_, err := parser.ReadLine(file)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("Error reading line:", err)
			os.Exit(1)
		}
		var fields []string
		for i := 0; i < parser.GetNumberOfFields(); i++ {
			field, err := parser.GetField(i)
			if err != nil {
				fmt.Println("Error: ", err)
				os.Exit(1)
			}
			// fields = append(fields, fmt.Sprintf("\"%s\"", field))
			fields = append(fields, field)
		}
		fmt.Println("|" + strings.Join(fields, "|") + "|")
	}
}
