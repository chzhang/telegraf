package main

import (
	"flag"
	"fmt"
	"log" //nolint:revive
	"os"
	"text/template"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/parser"
)

type PluginData struct {
	SampleConfig string
	Description  string
}

func extractPluginData() (PluginData, error) {
	readMe, err := os.ReadFile("README.md")
	if err != nil {
		return PluginData{}, err
	}
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs
	mdParser := parser.NewWithExtensions(extensions)
	md := markdown.Parse(readMe, mdParser)

	var p PluginData
	var currentSection string

	for _, t := range md.GetChildren() {
		switch tok := t.(type) {
		case *ast.Heading:
			currentSection = tok.HeadingID
		case *ast.CodeBlock:
			if currentSection == "configuration" && string(tok.Info) == "toml" {
				p.SampleConfig = string(tok.Literal)
			}
		}
	}

	return p, nil
}

func generatePluginData(p PluginData) error {
	goPackage := os.Getenv("GOPACKAGE")
	sourceName := fmt.Sprintf("%s.go", goPackage)

	plugin, err := os.ReadFile(sourceName)
	if err != nil {
		return err
	}

	err = os.Rename(sourceName, fmt.Sprintf("%s.tmp", sourceName))
	if err != nil {
		return err
	}

	generatedTemplate := template.Must(template.New("").Parse(string(plugin)))

	f, err := os.Create(sourceName)
	if err != nil {
		return err
	}
	defer f.Close()

	err = generatedTemplate.Execute(f, struct {
		SampleConfig string
	}{
		SampleConfig: p.SampleConfig,
	})
	if err != nil {
		return err
	}

	return nil
}

func main() {
	clean := flag.Bool("clean", false, "Remove generated files")
	flag.Parse()

	if *clean {
		goPackage := os.Getenv("GOPACKAGE")
		sourceName := fmt.Sprintf("%s.go", goPackage)
		err := os.Remove(sourceName)
		if err != nil {
			log.Fatal(err)
		}
		err = os.Rename(fmt.Sprintf("%s.tmp", sourceName), sourceName)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		p, err := extractPluginData()
		if err != nil {
			log.Fatal(err)
		}

		err = generatePluginData(p)
		if err != nil {
			log.Fatal(err)
		}
	}
}
