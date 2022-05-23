package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/PaesslerAG/gval"
	"github.com/puppetlabs/leg/jsonutil/pkg/jsonpath"
	"github.com/puppetlabs/leg/jsonutil/pkg/jsonpath/template"
)

func main() {
	langFlag := flag.String("l", "jsonpath", "The language to use for the query (jsonpath, jsonpath-template)")

	flag.Parse()

	var lang gval.Language
	switch *langFlag {
	case "jsonpath":
		lang = jsonpath.Language(jsonpath.WithMissingKeysAllowed{}, jsonpath.WithPlaceholders{})
	case "jsonpath-template":
		lang = template.TemplateLanguage()
	default:
		fmt.Fprintf(os.Stderr, "unknown language\n")
		os.Exit(2)
	}

	eval, err := lang.NewEvaluable(flag.Arg(0))
	if err != nil {
		fmt.Fprintf(os.Stderr, "query parse error: %+v\n", err)
		os.Exit(2)
	}

	var doc any
	if err := json.NewDecoder(os.Stdin).Decode(&doc); err != nil {
		fmt.Fprintf(os.Stderr, "document parse error: %+v\n", err)
		os.Exit(2)
	}

	selected, err := eval(context.Background(), doc)
	if err != nil {
		fmt.Fprintf(os.Stderr, "query evaluation error: %+v\n", err)
		os.Exit(1)
	}

	switch *langFlag {
	case "jsonpath":
		out, err := json.MarshalIndent(selected, "", "  ")
		if err != nil {
			fmt.Fprintf(os.Stderr, "serialization error: %+v\n", err)
			os.Exit(-1)
		}

		fmt.Println(string(out))
	default:
		fmt.Println(selected)
	}
}
