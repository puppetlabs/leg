// Portions of this file are derived from GoDS, a data structure library for
// Go.
//
// Copyright (c) 2015, Emir Pasic. All rights reserved.
//
// https://github.com/emirpasic/gods/blob/213367f1ca932600ce530ae11c8a8cc444e3a6da/containers/containers.go

package data

type Container interface {
	Empty() bool
	Size() int
	Clear()
	Values() []interface{}
}
