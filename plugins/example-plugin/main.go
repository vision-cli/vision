package main

import (
	"fmt"
	"io"
	"os"
	"text/template"
)

type Dog struct {
	Name     string
	Age      int
	Breed    string
	NoOfLegs int
}

func main() {
	f, err := os.OpenFile("example.txt.tmpl", os.O_RDWR, 0666)
	if err != nil {
		fmt.Errorf("opening file: %w", err)
	}

	b, err := io.ReadAll(f)
	if err != nil {
		fmt.Errorf("reading file: %w", err)
	}

	fileText := string(b)

	tmplEx, err := template.New("templateFile").Parse(fileText)
	bernie := Dog{
		Name:     "Bernie",
		Age:      7,
		Breed:    "Jack Russell",
		NoOfLegs: 1241,
	}

	err = f.Truncate(0)
	if err != nil {
		fmt.Errorf("truncating: %w", err)
	}

	_, err = f.Seek(0, 0)
	if err != nil {
		fmt.Errorf("seeking: %w", err)
	}

	err = tmplEx.Execute(f, bernie)
	if err != nil {
		panic(err)
	}
	err = f.Close()
	if err != nil {
		panic(err)
	}

}
