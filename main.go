package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"text/template"
)

type valueFlags map[string]string

func (i valueFlags) String() string {
	return "valueFlags containing named parameters and their values"
}

func (i valueFlags) Set(val string) error {
	split := strings.Split(val, "=")
	if len(split) != 2 {
		return errors.New("value must be formatted like name=/path/to/file")
	}
	bs, err := ioutil.ReadFile(split[1])
	if err != nil {
		return fmt.Errorf("reading value file %s: %w", split[0], err)
	}
	i[split[0]] = strings.TrimSpace(string(bs))
	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	values := make(valueFlags)
	flag.Var(&values, "val", "named values from the input file pointing to files containing the data to be filled in. Ex: myval=/secrets/myval")
	inFile := flag.String("in", "", "template file to be rendered")
	flag.Parse()

	if *inFile == "" {
		return errors.New("required flag -in")
	}

	if err := renderTemplate(*inFile, values, os.Stdout); err != nil {
		return err
	}
	return nil
}

func renderTemplate(src string, data valueFlags, outFile io.Writer) error {
	tpl, err := template.ParseFiles(src)
	if err != nil {
		return fmt.Errorf("parsing file %s: %w", src, err)
	}

	tpl.Option("missingkey=error")
	if err := tpl.Execute(outFile, data); err != nil {
		return fmt.Errorf("executing template %s: %w", src, err)
	}
	return nil
}
