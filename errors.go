package gographt

import (
	"errors"
	"fmt"
)

type VertexNotFoundError struct {
	Vertex Vertex
}

func (e *VertexNotFoundError) Error() string {
	return fmt.Sprintf("graph: vertex %q does not exist", e.Vertex)
}

type NotConnectedError struct {
	Source, Target Vertex
}

func (e *NotConnectedError) Error() string {
	return fmt.Sprintf("graph: not connected: %q and %q", e.Source, e.Target)
}

var (
	ErrEdgeAlreadyInGraph = errors.New("graph: edge already present")
	ErrEdgeNotFound       = errors.New("graph: edge does not exist")
	ErrWouldCreateLoop    = errors.New("graph: loop would be created by edge")
)
