// Copyright (c) 2014 Dataence, LLC. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package xtime is a time parser that parses the time without knowning the
// exact format.
package xtime

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"
)

// TimeFormats is a list of commonly seen time formats from log messages
var TimeFormats = []string{
	"Mon Jan _2 15:04:05 2006",
	"Mon Jan 02 15:04:05 -0700 2006",
	"Mon Jan 02 15:04:05 -07:00 2006",
	"02 Jan 06 15:04 -0700",
	"02 Jan 06 15:04 -07:00",
	"Monday, 02-Jan-06 15:04:05",
	"Monday, 02-Jan-06 15:04:05 -0700",
	"Monday, 02-Jan-06 15:04:05 -07:00",
	"Mon, 02 Jan 2006 15:04:05",
	"Mon, 02 Jan 2006 15:04:05 -0700",
	"Mon, 02 Jan 2006 15:04:05 -07:00",
	"2006-01-02T15:04:05",
	"2006-01-02T15:04:05Z0700",
	"2006-01-02T15:04:05Z07:00",
	"2006-01-02T15:04:05.999999999Z07:00",
	"2006-01-02T15:04:05.999Z07:00",
	"Jan _2 15:04:05",
	"Jan _2 15:04:05.000",
	"Jan _2 15:04:05.000000",
	"Jan _2 15:04:05.000000000",
	"_2/Jan/2006:15:04:05 -0700",
	"_2/Jan/2006:15:04:05 -07:00",
	"Jan 2, 2006 3:04:05 PM",
	"Jan 2 2006 15:04:05",
	"Jan 2 15:04:05 2006",
	"Jan 2 15:04:05 -0700",
	"2006-01-02 15:04:05,000 -0700",
	"2006-01-02 15:04:05,000 -07:00",
	"2006-01-02 15:04:05 -0700",
	"2006-01-02 15:04:05 -07:00",
	"2006-01-02 15:04:05-0700",
	"2006-01-02 15:04:05-07:00",
	"2006-01-02 15:04:05,000",
	"2006-01-02 15:04:05",
	"2006/01/02 15:04:05",
	"2006/01/02",
	"2006-01-02",
	"01/02/2006",
	"01/02/2006 15:04:05",
	"06-01-02 15:04:05,000 -0700",
	"06-01-02 15:04:05,000 -07:00",
	"06-01-02 15:04:05,000",
	"06-01-02 15:04:05",
	"06/01/02 15:04:05",
	"15:04:05,000",
	"1/2/2006 3:04:05 PM",
	"1/2/06 3:04:05.000 PM",
	"1/2/2006 15:04",
	"1/2/2006",
	"2006/1/2",
}

type TimeTree struct {
	formats []string
	root    timeNode
}

func (tt *TimeTree) ParseInLocation(t string, loc *time.Location) (time.Time, error) {
	if loc == nil {
		return tt.Parse(t)
	}

	tnv := newTimeNodeVisitor()
	tnv.Visit(tt.root, t)

	candidates := tnv.Matches()

	for _, candidate := range candidates {
		t, err := time.ParseInLocation(candidate, t, loc)
		if _, ok := err.(*time.ParseError); ok {
			continue
		}

		return t, err
	}

	return time.Time{}, ErrInvalidTime
}

func (tt *TimeTree) Parse(t string) (time.Time, error) {
	tnv := newTimeNodeVisitor()
	tnv.Visit(tt.root, t)

	candidates := tnv.Matches()

	for _, candidate := range candidates {
		t, err := time.Parse(candidate, t)
		if _, ok := err.(*time.ParseError); ok {
			continue
		}

		return t, err
	}

	return time.Time{}, ErrInvalidTime
}

func (tt *TimeTree) IsTime(t string) bool {
	_, err := tt.Parse(t)
	return err == nil
}

func Root() *TimeTree {
	return timeTreeRoot
}

func Compile(formats []string) *TimeTree {
	return &TimeTree{
		formats: formats,
		root:    buildTimeTree(formats),
	}
}

func ParseInLocation(t string, loc *time.Location) (time.Time, error) {
	return timeTreeRoot.ParseInLocation(t, loc)
}

func Parse(t string) (time.Time, error) {
	return timeTreeRoot.Parse(t)
}

func IsTime(t string) bool {
	return timeTreeRoot.IsTime(t)
}

var (
	timeTreeRoot *TimeTree
)

func init() {
	timeTreeRoot = Compile(TimeFormats)
}

type timeNodeVisitor struct {
	matches map[string]struct{}
}

func (tnv *timeNodeVisitor) Matches() []string {
	matches := make([]string, len(tnv.matches))

	i := 0
	for m := range tnv.matches {
		matches[i] = m
		i++
	}

	sort.Slice(matches, func(i, j int) bool {
		return len(matches[i]) > len(matches[j])
	})
	return matches
}

func (tnv *timeNodeVisitor) Visit(tn timeNode, t string) {
	values := tn.Values()
	if len(values) > 0 && len(t) == 0 {
		for _, value := range values {
			tnv.matches[value] = struct{}{}
		}
	} else {
		for _, c := range tn.Children() {
			c.Visit(tnv, t)
		}
	}
}

func newTimeNodeVisitor() *timeNodeVisitor {
	return &timeNodeVisitor{
		matches: make(map[string]struct{}),
	}
}

type timeNode interface {
	Values() []string
	Children() []timeNode
	Visit(tnv *timeNodeVisitor, t string)
}

type baseTimeNode struct {
	values   []string
	children []timeNode
}

func (tn *baseTimeNode) Values() []string {
	return tn.values
}

func (tn *baseTimeNode) Children() []timeNode {
	return tn.children
}

type rootTimeNode struct {
	*baseTimeNode
}

func (rtn *rootTimeNode) Visit(tnv *timeNodeVisitor, t string) {
	tnv.Visit(rtn, t)
}

func (rtn *rootTimeNode) String() string {
	return fmt.Sprintf(`{RootTimeNode children=%s}`, rtn.children)
}

type digitsTimeNode struct {
	*baseTimeNode

	min, max int
	padded   bool
}

func (dtn *digitsTimeNode) Visit(tnv *timeNodeVisitor, t string) {
	var padded bool
	sz := 0

	for len(t) > 0 && sz < dtn.max {
		r, l := utf8.DecodeRuneInString(t)
		if !padded && dtn.padded && unicode.IsSpace(r) {
			sz++
			t = t[l:]

			continue
		} else if !unicode.IsDigit(r) {
			return
		}

		padded = true

		sz++
		t = t[l:]

		if sz >= dtn.min {
			tnv.Visit(dtn, t)
		}
	}
}

func (dtn *digitsTimeNode) String() string {
	return fmt.Sprintf(
		`{DigitsTimeNode min=%d max=%d padded=%v children=%s values=%s}`,
		dtn.min, dtn.max, dtn.padded, dtn.children, dtn.values,
	)
}

type spaceTimeNode struct {
	*baseTimeNode
}

func (stn *spaceTimeNode) Visit(tnv *timeNodeVisitor, t string) {
	var ok bool

	for len(t) > 0 {
		r, l := utf8.DecodeRuneInString(t)
		if !unicode.IsSpace(r) {
			break
		}

		ok = true
		t = t[l:]
	}

	if !ok {
		return
	}

	tnv.Visit(stn, t)
}

func (stn *spaceTimeNode) String() string {
	return fmt.Sprintf(`{SpaceTimeNode children=%s values=%s}`, stn.children, stn.values)
}

type zoneTimeNode struct {
	*baseTimeNode

	len int
}

func (ztn *zoneTimeNode) Visit(tnv *timeNodeVisitor, t string) {
	sz := 0

	// Could be "Z" or +/-HH:MM.
	r, l := utf8.DecodeRuneInString(t)

	sz++
	t = t[l:]

	if r == 'Z' {
		tnv.Visit(ztn, t)
		return
	} else if r != '+' && r != '-' {
		// Not a valid zone offset.
		return
	}

	// Now we need some digits.
	var ok bool

	for len(t) > 0 && sz < ztn.len {
		r, l := utf8.DecodeRuneInString(t)
		if !unicode.IsNumber(r) {
			break
		}

		ok = true

		sz++
		t = t[l:]
	}

	if !ok {
		return
	}

	// Now maybe we have a ':'?
	if sz < ztn.len {
		if len(t) == 0 || t[0] != ':' {
			return
		}

		sz++
		t = t[1:]

		// And now we need the rest of the digits.
		ok = false

		for len(t) > 0 && sz < ztn.len {
			r, l := utf8.DecodeRuneInString(t)
			if !unicode.IsNumber(r) {
				break
			}

			sz++
			t = t[l:]
		}

		if sz < ztn.len {
			return
		}
	}

	tnv.Visit(ztn, t)
}

func (ztn *zoneTimeNode) String() string {
	return fmt.Sprintf(`{ZoneTimeNode len=%d children=%s values=%s}`, ztn.len, ztn.children, ztn.values)
}

type exactTimeNode struct {
	*baseTimeNode

	expected string
}

func (etn *exactTimeNode) Visit(tnv *timeNodeVisitor, t string) {
	if !strings.HasPrefix(t, etn.expected) {
		return
	}

	t = t[len(etn.expected):]
	tnv.Visit(etn, t)
}

func (etn *exactTimeNode) String() string {
	return fmt.Sprintf(`{ExactTimeNode expected=%s children=%s values=%s}`, etn.expected, etn.children, etn.values)
}

type otherTimeNode struct {
	*baseTimeNode
}

func (otn *otherTimeNode) Values() []string {
	return otn.values
}

func (otn *otherTimeNode) Children() []timeNode {
	return otn.children
}

func (otn *otherTimeNode) Visit(tnv *timeNodeVisitor, t string) {
	var ok bool

	for len(t) > 0 {
		r, l := utf8.DecodeRuneInString(t)
		if unicode.IsSpace(r) || unicode.IsDigit(r) {
			break
		}

		ok = true
		t = t[l:]
	}

	if !ok {
		return
	}

	tnv.Visit(otn, t)
}

func (otn *otherTimeNode) String() string {
	return fmt.Sprintf(`{OtherTimeNode children=%s values=%s}`, otn.children, otn.values)
}

func (tn *baseTimeNode) withZoneNode(len int) *baseTimeNode {
	for _, c := range tn.children {
		ztn, ok := c.(*zoneTimeNode)
		if !ok {
			continue
		}

		if ztn.len != len {
			continue
		}

		return ztn.baseTimeNode
	}

	base := &baseTimeNode{}
	c := &zoneTimeNode{baseTimeNode: base, len: len}

	tn.children = append(tn.children, c)
	return base
}

func (tn *baseTimeNode) withOtherNode() *baseTimeNode {
	for _, c := range tn.children {
		otn, ok := c.(*otherTimeNode)
		if !ok {
			continue
		}

		return otn.baseTimeNode
	}

	base := &baseTimeNode{}
	c := &otherTimeNode{baseTimeNode: base}

	tn.children = append(tn.children, c)
	return base
}

func (tn *baseTimeNode) withSpaceNode() *baseTimeNode {
	for _, c := range tn.children {
		otn, ok := c.(*spaceTimeNode)
		if !ok {
			continue
		}

		return otn.baseTimeNode
	}

	base := &baseTimeNode{}
	c := &spaceTimeNode{baseTimeNode: base}

	tn.children = append(tn.children, c)
	return base
}

func (tn *baseTimeNode) withDigitsNode(min, max int, padded bool) *baseTimeNode {
	for _, c := range tn.children {
		dtn, ok := c.(*digitsTimeNode)
		if !ok {
			continue
		}

		if dtn.min != min || dtn.max != max || dtn.padded != padded {
			continue
		}

		return dtn.baseTimeNode
	}

	base := &baseTimeNode{}
	c := &digitsTimeNode{
		baseTimeNode: base,
		min:          min,
		max:          max,
		padded:       padded,
	}

	tn.children = append(tn.children, c)
	return base
}

func (tn *baseTimeNode) withFracNode(len int) *baseTimeNode {
	return tn.withExpected(".").withDigitsNode(1, len, false)
}

func (tn *baseTimeNode) withExpected(s string) *baseTimeNode {
	for _, c := range tn.children {
		etn, ok := c.(*exactTimeNode)
		if !ok {
			continue
		}

		if etn.expected != s {
			continue
		}

		return etn.baseTimeNode
	}

	base := &baseTimeNode{}
	c := &exactTimeNode{
		baseTimeNode: base,
		expected:     s,
	}

	tn.children = append(tn.children, c)
	return base
}

var (
	timeNodeMatcher = regexp.MustCompile(`(?P<zone>[Z-]07(?::?00)?)|(?P<space>\s+)|(?P<year>(?:20)?06)|(?P<hour>(?:15|[0_]?3))|(?P<month>[0_]?1)|(?P<day>[0_]?2)|(?P<minute>[0_]?4)|(?P<second>[0_]?5)|(?P<frac>\.[90]+)|(?P<digit>\d+)`)
)

func buildTimeTree(formats []string) timeNode {
	base := &baseTimeNode{}

	for _, f := range formats {
		matches := timeNodeMatcher.FindAllStringSubmatchIndex(f, -1)

		// Cursor into the time tree.
		cur := base

		i := 0
		for len(matches) > 0 {
			if matches[0][0] > i {
				cur = cur.withOtherNode()
			}

			for n, name := range timeNodeMatcher.SubexpNames()[1:] {
				start := matches[0][(n+1)*2]
				if start == -1 {
					continue
				}

				end := matches[0][(n+1)*2+1]

				switch name {
				case "zone":
					cur = cur.withZoneNode(end - start)
				case "space":
					cur = cur.withSpaceNode()
				case "year", "digit":
					cur = cur.withDigitsNode(end-start, end-start, false)
				case "hour", "month", "day", "minute", "second":
					cur = cur.withDigitsNode(1, 2, f[start] == '_')
				case "frac":
					cur = cur.withFracNode(end - start - 1)
				default:
					panic(fmt.Errorf("unknown subexpression name %s", name))
				}

				break
			}

			i = matches[0][1]
			matches = matches[1:]
		}

		if i < len(f) {
			cur = cur.withOtherNode()
		}

		cur.values = append(cur.values, f)
	}

	return &rootTimeNode{base}
}
