// Portions of this file are derived from Paessler AG's JSONPath implementation.
//
// Copyright (c) 2017 Paessler AG <support@paessler.com>. All rights reserved.
//
// https://github.com/PaesslerAG/jsonpath/blob/c18d0f043db32b5d4295e14c6518fa9160e45d15/parse.go
//
// Licensed under a three-clause BSD-style license. A copy of the full license
// document is included in this distribution in the file `parse.go.LICENSE`.

package jsonpath

import (
	"context"
	"math"
	"text/scanner"

	"github.com/PaesslerAG/gval"
)

type parser struct {
	*gval.Parser
	path *path
}

func parseRootPath(selectorMode selectorMode) func(ctx context.Context, gParser *gval.Parser) (r gval.Evaluable, err error) {
	return func(ctx context.Context, gParser *gval.Parser) (r gval.Evaluable, err error) {
		p := newParser(gParser, rootElement, selectorMode)
		return p.parse(ctx)
	}
}

func parseCurrentPath(selectorMode selectorMode) func(ctx context.Context, gParser *gval.Parser) (r gval.Evaluable, err error) {
	return func(ctx context.Context, gParser *gval.Parser) (r gval.Evaluable, err error) {
		p := newParser(gParser, currentElement, selectorMode)
		return p.parse(ctx)
	}
}

func newParser(gp *gval.Parser, root gval.Evaluable, selectorMode selectorMode) *parser {
	p := &parser{Parser: gp, path: &path{root: root, mode: selectorMode}}
	return p
}

func (p *parser) parse(c context.Context) (r gval.Evaluable, err error) {
	err = p.parsePath(c)

	if err != nil {
		return nil, err
	}
	return p.path.evaluate, nil
}

func (p *parser) parsePath(c context.Context) error {
	scan := func() rune {
		p.SetMode(scanner.GoTokens &^ scanner.Float)
		defer p.SetMode(scanner.GoTokens)
		return p.Scan()
	}()
	switch scan {
	case '.':
		return p.parseSelect(c)
	case '[':
		return p.parseBracket(c, false)
	case '(':
		return p.parseScript(c)
	default:
		p.Camouflage("JSONPath", '.', '[', '(')
		return nil
	}
}

func (p *parser) parseSelect(c context.Context) error {
	scan := p.Scan()
	switch scan {
	case scanner.Ident, scanner.Int:
		p.appendSelector(directSelector(varSelector(p.variableChild(p.Const(p.TokenText())))), selectorKeepErrors)
		return p.parsePath(c)
	case '.':
		p.appendSelector(varSelector(p.variableRecursiveDescent()), selectorDropErrors)
		return p.parseMapper(c)
	case '*':
		p.appendSelector(varSelector(p.variableWildcard()), selectorDropErrors)
		return p.parsePath(c)
	default:
		return p.Expected("JSON select", scanner.Ident, scanner.Int, '.', '*')
	}
}

func (p *parser) parseMapper(c context.Context) error {
	scan := p.Scan()
	switch scan {
	case scanner.Ident:
		p.appendSelector(directSelector(varSelector(p.variableChild(p.Const(p.TokenText())))), selectorKeepErrors)
	case '[':
		return p.parseBracket(c, true)
	case '*':
		p.appendSelector(varSelector(p.variableWildcard()), selectorDropErrors)
	case '(':
		return p.parseScript(c)
	default:
		p.Camouflage("JSON mapper", '[', scanner.Ident, '*', '(')
	}
	return p.parsePath(c)
}

func (p *parser) parseBracket(c context.Context, mapper bool) error {
	keys, seperator, err := p.scanBracket(c)
	if err != nil {
		return err
	}

	switch seperator {
	case ':':
		if mapper {
			return ErrMapperNotCombinable
		} else if len(keys) > 3 {
			return ErrRangeQueryOutOfBounds
		}
		keys = append(keys, []gval.Evaluable{p.Const(0), p.Const(float64(math.MaxInt32)), p.Const(1)}[len(keys):]...)
		p.appendSelector(varSelector(p.variableRange(keys[0], keys[1], keys[2])), selectorDropErrors)
	case '?':
		if len(keys) != 1 {
			return ErrFilterOutOfBounds
		}
		p.appendSelector(filterSelector(varSelector(p.variableWildcard()), keys[0]), selectorDropErrors)
	case '*':
		p.appendSelector(varSelector(p.variableWildcard()), selectorDropErrors)
	case ',':
		selectors := make([]selector, len(keys))
		for i, key := range keys {
			selectors[i] = varSelector(p.variableChild(key))
		}
		p.appendSelector(multiSelector(selectors), selectorDropErrors)
	default:
		if len(keys) != 1 {
			return &UnexpectedSeparatorError{Separator: seperator}
		}
		p.appendSelector(varSelector(p.variableChild(keys[0])), selectorKeepErrors)
	}
	return p.parsePath(c)
}

func (p *parser) scanBracket(c context.Context) (keys []gval.Evaluable, seperator rune, err error) {
	for {
		scan := p.Scan()
		skipScan := false
		switch scan {
		case '?':
			if p.Scan() != '(' {
				return nil, 0, p.Expected("JSONPath filter", '(')
			}
			p.Camouflage("JSONPath filter")
			skipScan = true
		case ':':
			i := float64(0)
			if len(keys) == 1 {
				i = math.MaxInt32
			}
			keys = append(keys, p.Const(i))
			skipScan = true
		case '*':
			if p.Scan() != ']' {
				return nil, 0, p.Expected("JSON bracket star", ']')
			}
			return nil, '*', nil
		case ']':
			if seperator == ':' {
				skipScan = true
				break
			}
			fallthrough
		default:
			p.Camouflage("jsonpath brackets")
			key, err := p.ParseExpression(c)
			if err != nil {
				return nil, 0, err
			}
			keys = append(keys, key)
		}
		if !skipScan {
			scan = p.Scan()
		}
		if seperator == 0 {
			seperator = scan
		}
		switch scan {
		case ':', ',':
		case ']':
			return
		case '?':
			if len(keys) != 0 {
				return nil, 0, p.Expected("JSON filter", ']')
			}
		default:
			return nil, 0, p.Expected("JSON bracket separator", ':', ',')
		}
		if seperator != scan {
			return nil, 0, &MixedSeparatorError{A: seperator, B: scan}
		}
	}
}

func (p *parser) parseScript(c context.Context) error {
	script, err := p.ParseExpression(c)
	if err != nil {
		return err
	}
	if p.Scan() != ')' {
		return p.Expected("JSONPath script", ')')
	}
	p.appendSelector(scriptSelector(script), selectorKeepErrors)
	return p.parsePath(c)
}

func (p *parser) variableWildcard() gval.Evaluable {
	return p.Var(p.Const(variableWildcard{}))
}

func (p *parser) variableRecursiveDescent() gval.Evaluable {
	return p.Var(p.Const(variableRecursiveDescent{}))
}

func (p *parser) variableRange(min, max, step gval.Evaluable) gval.Evaluable {
	return p.Var(func(c context.Context, v any) (any, error) {
		min, err := min.EvalInt(c, v)
		if err != nil {
			return nil, err
		}

		max, err := max.EvalInt(c, v)
		if err != nil {
			return nil, err
		}

		step, err := step.EvalInt(c, v)
		if err != nil {
			return nil, err
		}

		return variableRange{Min: min, Max: max, Step: step}, nil
	})
}

func (p *parser) variableChild(key gval.Evaluable) gval.Evaluable {
	return p.Var(func(c context.Context, v any) (any, error) {
		key, err := key(c, v)
		if err != nil {
			return nil, err
		}

		return variableChild{Key: key}, nil
	})
}

func (p *parser) appendSelector(next selector, mode selectorMode) {
	p.path.appendSelector(next, mode)
}
