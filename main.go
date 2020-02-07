package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

func main() {
	valsFile := flag.String("f", "", "values file for rendering template")
	flag.Parse()

	if *valsFile == "" {
		log.Fatal("required flag -f")
	}

	var tplFiles []string
	filepath.Walk(".", func(path string, f os.FileInfo, err error) error {
		if !f.IsDir() && strings.HasSuffix(path, ".tpl") {
			tplFiles = append(tplFiles, path)
		}
		return nil
	})
	log.Infof("found %d template files", len(tplFiles))

	bs, err := ioutil.ReadFile(*valsFile)
	if err != nil {
		log.Fatalf("reading values file: %v", err)
	}
	var values interface{}
	if err := yaml.Unmarshal(bs, &values); err != nil {
		log.Fatalf("unmarshaling yaml values file: %v", err)
	}

	for _, path := range tplFiles {
		log.Infof("rendering %s", path)
		if err := renderTemplate(path, strings.TrimSuffix(path, ".tpl"), values); err != nil {
			log.Error(err)
		}
	}
}

func renderTemplate(src, dst string, vals interface{}) error {
	tpl, err := template.ParseFiles(src)
	if err != nil {
		return fmt.Errorf("parsing file %s: %w", src, err)
	}

	buf := bytes.NewBuffer([]byte{})
	tpl.Option("missingkey=error")
	if err := tpl.Execute(buf, vals); err != nil {
		return fmt.Errorf("executing template %s: %w", src, err)
	}

	f, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("creating file %s: %w", dst, err)
	}

	if _, err := io.Copy(f, buf); err != nil {
		return fmt.Errorf("writing to %s: %w", dst, err)
	}

	return nil
}
