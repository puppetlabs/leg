// Portions of this file are derived from Paessler AG's JSONPath implementation.
//
// Copyright (c) 2017 Paessler AG <support@paessler.com>. All rights reserved.
//
// https://github.com/PaesslerAG/jsonpath/blob/c18d0f043db32b5d4295e14c6518fa9160e45d15/example_test.go
//
// Licensed under a three-clause BSD-style license. A copy of the full license
// document is included in this distribution in the file
// `example_test.go.LICENSE`.

package jsonpath_test

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/PaesslerAG/gval"
	"github.com/puppetlabs/leg/jsonutil/pkg/jsonpath"
)

func ExampleGet() {
	v := any(nil)

	_ = json.Unmarshal([]byte(`{
		"welcome":{
				"message":["Good Morning", "Hello World!"]
			}
		}`), &v)

	welcome, err := jsonpath.Get("$.welcome.message[1]", v)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(welcome)

	// Output:
	// Hello World!
}

func ExampleGet_wildcard() {
	v := any(nil)

	_ = json.Unmarshal([]byte(`{
		"welcome":{
				"message":["Good Morning", "Hello World!"]
			}
		}`), &v)

	welcome, err := jsonpath.Get("$.welcome.message[*]", v)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, value := range welcome.([]any) {
		fmt.Printf("%v\n", value)
	}

	// Output:
	// Good Morning
	// Hello World!
}

func ExampleGet_filter() {
	v := any(nil)

	_ = json.Unmarshal([]byte(`[
		{"key":"a","value" : "I"},
		{"key":"b","value" : "II"},
		{"key":"c","value" : "III"}
		]`), &v)

	values, err := jsonpath.Get(`$[?(@.key=="b")].value`, v)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, value := range values.([]any) {
		fmt.Println(value)
	}

	// Output:
	// II
}

func Example_gval() {
	builder := gval.Full(jsonpath.Language(jsonpath.WithPlaceholders{}))

	path, err := builder.NewEvaluable("{#1: $..[?(@.ping && @.speed > 100)].name}")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	v := any(nil)
	err = json.Unmarshal([]byte(`{
		"device 1":{
			"name": "fancy device",
			"ping": true,
			"speed": 200,
				"subdevice 1":{
					"ping" : true,
					"speed" : 99,
					"name" : "boring subdevice"
				},
				"subdevice 2":{
					"ping" : true,
					"speed" : 150,
					"name" : "fancy subdevice"
				},
				"not an device":{
					"name" : "ping me but I have no speed property",
					"ping" : true
				}
			},
		"fictive device":{
			"ping" : false,
			"speed" : 1000,
			"name" : "dream device"
			}
		}`), &v)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	devices, err := path(context.Background(), v)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for device, name := range devices.(map[string]any) {
		fmt.Printf("%s -> %v\n", device, name)
	}

	// Unordered output:
	// device 1 -> fancy device
	// subdevice 2 -> fancy subdevice
}

func Example_variableSelector() {
	builder := gval.NewLanguage(
		jsonpath.Language(),
		gval.VariableSelector(jsonpath.ChildVariableSelector(func(ctx context.Context, v any, key any, next func(context.Context, jsonpath.PathValue) error) error {
			return jsonpath.DefaultVariableVisitor().VisitChild(ctx, v, key, func(ctx context.Context, pv jsonpath.PathValue) error {
				if s, ok := pv.Value.(string); ok && strings.HasPrefix(s, "base64:") {
					b, err := base64.StdEncoding.DecodeString(s[len("base64:"):])
					if err != nil {
						return fmt.Errorf("could not decode base64 value: %v", err)
					}

					pv.Value = string(b)
				}

				return next(ctx, pv)
			})
		})),
	)

	path, err := builder.NewEvaluable(`$.encoded`)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var v any
	err = json.Unmarshal([]byte(`{
		"encoded": "base64:SGVsbG8sIHdvcmxkIQ=="
	}`), &v)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	decoded, err := path(context.Background(), v)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(decoded)

	// Output:
	// Hello, world!
}
