// Portions of this file are derived from Paessler AG's JSONPath implementation.
//
// Copyright (c) 2017 Paessler AG <support@paessler.com>. All rights reserved.
//
// https://github.com/PaesslerAG/jsonpath/blob/c18d0f043db32b5d4295e14c6518fa9160e45d15/jsonpath.go
//
// Licensed under a three-clause BSD-style license. A copy of the full license
// document is included in this distribution in the file `jsonpath.go.LICENSE`.

package jsonpath

import (
	"context"
	"fmt"
	"strings"
	"text/scanner"
	"unicode"

	"github.com/PaesslerAG/gval"
	"github.com/generikvault/gvalstrings"
)

// New returns an selector for given JSONPath
func New(path string) (gval.Evaluable, error) {
	return Language().NewEvaluable(path)
}

// Get executes given JSONPath on given value
func Get(path string, value interface{}) (interface{}, error) {
	eval, err := New(path)
	if err != nil {
		return nil, err
	}
	return eval(context.Background(), value)
}

type LanguageOptions struct {
	MissingKeysAllowed bool
	Placeholders       bool
	InitialPath        bool
}

type LanguageOption interface {
	ApplyToLanguageOptions(target *LanguageOptions)
}

func (o *LanguageOptions) ApplyOptions(opts []LanguageOption) {
	for _, opt := range opts {
		opt.ApplyToLanguageOptions(o)
	}
}

// WithMissingKeysAllowed causes the parser not to return an error when a
// selector key is not present in the document.
type WithMissingKeysAllowed struct{}

var _ LanguageOption = WithMissingKeysAllowed{}

func (WithMissingKeysAllowed) ApplyToLanguageOptions(target *LanguageOptions) {
	target.MissingKeysAllowed = true
}

// WithPlaceholders enables the wildcard placeholder feature.
type WithPlaceholders struct{}

var _ LanguageOption = WithPlaceholders{}

func (WithPlaceholders) ApplyToLanguageOptions(target *LanguageOptions) {
	target.Placeholders = true
}

// WithInitialPath allows the path selector characters '.', '[', and '(' to
// begin a selector instead of just '$' or '@' ('$' is implied).
type WithInitialPath struct{}

var _ LanguageOption = WithInitialPath{}

func (WithInitialPath) ApplyToLanguageOptions(target *LanguageOptions) {
	target.InitialPath = true
}

// Language is the JSONPath language.
func Language(opts ...LanguageOption) gval.Language {
	o := &LanguageOptions{}
	o.ApplyOptions(opts)

	selectorMode := selectorKeepErrors
	if o.MissingKeysAllowed {
		selectorMode = selectorDropErrors
	}

	lang := gval.NewLanguage(
		gval.Base(),
		gval.Arithmetic(),
		gval.Text(),
		gval.JSON(),
		gval.PropositionalLogic(),
		gval.PrefixExtension(scanner.Char, func(c context.Context, p *gval.Parser) (gval.Evaluable, error) {
			s, err := gvalstrings.UnquoteSingleQuoted(p.TokenText())
			if err != nil {
				return nil, fmt.Errorf("could not parse string: %v", err)
			}
			return p.Const(s), nil
		}),
		gval.PrefixExtension('$', parseRootPath(selectorMode)),
		gval.PrefixExtension('@', parseCurrentPath(selectorMode)),
		gval.VariableSelector(VariableSelector(VariableVisitorFuncs{})),
		gval.Init(func(ctx context.Context, p *gval.Parser) (gval.Evaluable, error) {
			p.SetIsIdentRuneFunc(func(r rune, pos int) bool {
				return unicode.IsLetter(r) ||
					r == '_' ||
					(pos > 0 && (unicode.IsDigit(r) || r == '/' || r == '-'))
			})

			scan := p.Scan()
			p.Camouflage("JSONPath expression", '$', '@', '.', '[', '(')

			var eval gval.Evaluable
			var err error
			if o.InitialPath && strings.ContainsRune(".[(", scan) {
				eval, err = parseRootPath(selectorMode)(ctx, p)
			} else {
				eval, err = p.ParseExpression(ctx)
			}
			if err != nil {
				return nil, err
			}

			return func(ctx context.Context, parameter interface{}) (interface{}, error) {
				return eval(currentContext(ctx, parameter), parameter)
			}, nil
		}),
	)

	if o.Placeholders {
		lang = gval.NewLanguage(
			lang,
			gval.PrefixExtension('{', parseJSONObject(selectorMode)),
			gval.PrefixExtension('#', parsePlaceholder),
		)
	}

	return lang
}
