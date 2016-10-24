// Portions of this file are derived from GoDS, a data structure library for
// Go.
//
// Copyright (c) 2015, Emir Pasic. All rights reserved.
//
// https://github.com/emirpasic/gods/blob/213367f1ca932600ce530ae11c8a8cc444e3a6da/sets/sets.go

package data

type SetIterationFunc func(element interface{}) error

type Set interface {
	Container

	Contains(elements ...interface{}) bool
	Add(elements ...interface{})
	Remove(elements ...interface{})
	ForEach(fn SetIterationFunc) error
}
