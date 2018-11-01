package generator

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/puppetlabs/errawr-gen/doc"
	"github.com/puppetlabs/errawr-gen/golang"
	"github.com/xeipuuv/gojsonschema"
	"gopkg.in/yaml.v2"
)

type Generator interface {
	Generate(pkg string, document *doc.Document, output io.Writer) error
}

type Language string

var (
	LanguageGo Language = "go"
)

type Config struct {
	Package        string
	OutputPath     string
	OutputLanguage Language
	InputPath      string
}

func Generate(conf Config) error {
	if len(conf.Package) == 0 {
		if pkg := os.Getenv("GOPACKAGE"); len(pkg) != 0 {
			conf.Package = pkg
		} else {
			return fmt.Errorf("package name could not be determined; specify one")
		}
	}

	var generator Generator

	switch conf.OutputLanguage {
	case LanguageGo, "":
		generator = golang.NewGenerator()
	default:
		return fmt.Errorf("language %q is not supported", conf.OutputLanguage)
	}

	var input, output *os.File
	var err error

	if len(conf.InputPath) > 0 && conf.InputPath != "-" {
		input, err = os.Open(conf.InputPath)
		if err != nil {
			return fmt.Errorf("could not open input file: %+v", err)
		}
		defer input.Close()
	} else {
		input = os.Stdin
	}

	y, err := ioutil.ReadAll(input)
	if err != nil {
		return fmt.Errorf("could not read file: %+v", err)
	}

	// Pull out the document version.
	var version doc.DocumentVersionFragment
	if err := yaml.Unmarshal(y, &version); err != nil {
		return fmt.Errorf("could not read version from YAML: %+v", err)
	}

	if version.Version != "1" {
		return fmt.Errorf(`unexpected version %q; expected "1"`, version)
	}

	var document doc.Document
	if err := yaml.UnmarshalStrict(y, &document); err != nil {
		return fmt.Errorf("could not parse YAML: %+v", err)
	}

	result, err := doc.Schema.Validate(gojsonschema.NewGoLoader(document))
	if err != nil {
		log.Fatalf("Could not generate YAML validation: %+v", err)
	} else if !result.Valid() {
		errs := make([]string, len(result.Errors()))
		for i, err := range result.Errors() {
			errs[i] = fmt.Sprintf("%s", err)
		}

		return fmt.Errorf("validation errors occurred:\n%+v", strings.Join(errs, "\n"))
	}

	if len(conf.OutputPath) > 0 && conf.OutputPath != "-" {
		output, err = os.Create(conf.OutputPath)
		if err != nil {
			return fmt.Errorf("could not open output file: %+v", err)
		}
		defer output.Close()
	} else {
		output = os.Stdout
	}

	if err := generator.Generate(conf.Package, &document, output); err != nil {
		return fmt.Errorf("failed to generate Go file: %+v", err)
	}

	return nil
}
