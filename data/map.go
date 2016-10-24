// Portions of this file are derived from GoDS, a data structure library for
// Go.
//
// Copyright (c) 2015, Emir Pasic. All rights reserved.
//
// https://github.com/emirpasic/gods/blob/52d942a0538c185239fa3737047f297d983ac3e0/maps/maps.go

package data

type MapIterationFunc func(key, value interface{}) error

type Map interface {
	Container

	Contains(key interface{}) bool
	Put(key, value interface{}) bool
	Get(key interface{}) (interface{}, bool)
	Remove(key interface{}) bool
	Keys() []interface{}
	ForEach(fn MapIterationFunc) error
}
