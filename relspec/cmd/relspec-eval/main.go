package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/puppetlabs/leg/gvalutil/pkg/eval"
	"github.com/puppetlabs/leg/relspec/pkg/evaluate"
	"github.com/puppetlabs/leg/relspec/pkg/pathlang"
	"github.com/puppetlabs/leg/relspec/pkg/query"
	"github.com/puppetlabs/leg/relspec/pkg/ref"
	"github.com/puppetlabs/leg/relspec/pkg/relspec"
)

type Name string

func (n Name) Less(other Name) bool {
	return n < other
}

type References = *ref.Log[Name]

type Data map[string]any

func (d Data) evaluate(name string) *evaluate.Result[References] {
	value, found := d[name]
	if !found {
		return evaluate.ContextualizedResult(evaluate.NewMetadata(
			ref.InitialLog(ref.Errored(Name(name), fmt.Errorf("not found"))),
		))
	}

	r := evaluate.NewResult(
		evaluate.NewMetadata(ref.InitialLog(ref.OK(Name(name)))),
		value,
	)
	r.SetEvaluator(evaluate.DefaultEvaluator[References]())
	return r
}

func (d Data) ResolveMappingType(ctx context.Context, tm map[string]any) (*evaluate.Result[References], error) {
	name, ok := tm["name"].(string)
	if !ok {
		return nil, fmt.Errorf(`missing field "name" for data lookup`)
	}

	return d.evaluate(name), nil
}

var _ eval.Indexable = Data(nil)

func (d Data) Index(ctx context.Context, idx any) (any, error) {
	name, err := eval.StringValue(idx)
	if err != nil {
		return nil, err
	}

	return evaluate.StaticExpandable(d.evaluate(name)), nil
}

func (d Data) Expand(ctx context.Context, depth int) (*evaluate.Result[References], error) {
	if depth == 0 {
		return evaluate.StaticResult[References](d), nil
	}

	r := evaluate.NewResult(evaluate.NewMetadata(ref.NewLog[Name]()), map[string]any(d))
	r.SetEvaluator(evaluate.DefaultEvaluator[References]())
	for name := range d {
		r.References.Set(ref.OK(Name(name)))
	}
	return r, nil
}

func (d Data) String() string {
	return fmt.Sprintf("%v", map[string]any(d))
}

func (d Data) Set(value string) error {
	name, value, _ := strings.Cut(value, "=")
	d[name] = value
	return nil
}

func main() {
	ctx := context.Background()

	data := make(Data)

	quietFlag := flag.Bool("q", false, "If set, just output the result value")
	eagerFlag := flag.Bool("eager", false, "If set, use eager evaluation for functions and operators")
	langFlag := flag.String("l", "path", "The language to use for the expression")
	flag.Var(&data, "d", "Data values to make available for interpolation")

	flag.Parse()

	var tree any
	if err := json.NewDecoder(os.Stdin).Decode(&tree); err != nil {
		fmt.Fprintf(os.Stderr, "document parse error: %s\n", err)
		os.Exit(2)
	}

	ev := relspec.NewEvaluator[References](
		relspec.WithMappingTypeResolvers[References]{"Data": data},
		relspec.WithTemplateEnvironment[References]{"data": data},
		relspec.WithEagerEvaluation[References](*eagerFlag),
	)

	var (
		rv  *evaluate.Result[References]
		err error
	)
	if expr := flag.Arg(0); expr != "" {
		var lang query.Language[References]
		switch *langFlag {
		case "path":
			lang = pathlang.New[References]().Expression
		case "path-template":
			lang = pathlang.New[References]().Template
		case "jsonpath":
			lang = query.JSONPathLanguage[References]
		case "jsonpath-template":
			lang = query.JSONPathTemplateLanguage[References]
		default:
			fmt.Fprintf(os.Stderr, `unknown language %q (expected one of "path", "path-template", "jsonpath", "jsonpath-template")`+"\n", *langFlag)
			os.Exit(2)
		}

		rv, err = query.EvaluateQuery(ctx, ev, lang, tree, expr)
	} else {
		rv, err = evaluate.EvaluateAll(ctx, ev, tree)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "evaluation error: %s\n", err)
		os.Exit(2)
	}

	v, err := json.MarshalIndent(rv.Value, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "data marshal error: %s\n", err)
		os.Exit(2)
	}

	if !*quietFlag {
		fmt.Println("References:")
		rv.References.ForEach(func(ref ref.Reference[Name]) {
			if !ref.Used() {
				fmt.Print("? ")
			} else if ref.Error() != nil {
				fmt.Print("! ")
			} else {
				fmt.Print("* ")
			}
			fmt.Print(ref.ID())
			if !ref.Used() {
				fmt.Printf(" (unused)")
			}
			if ref.Error() != nil {
				fmt.Printf(": %s", ref.Error())
			}
			fmt.Print("\n")
		})

		fmt.Println()
		fmt.Println("Value:")
	}

	fmt.Println(string(v))
}
