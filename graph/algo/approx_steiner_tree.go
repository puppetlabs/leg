// Portions of this file are derived from a Steiner tree approximation
// algorithm written by the University of Southern California.
//
// https://github.com/usc-isi-i2/Web-Karma/blob/cef35dcb1a5042d1e8fabbbd61cb731a78c64454/karma-common/src/main/java/edu/isi/karma/modeling/alignment/SteinerTree.java
//
// Copyright 2012 University of Southern California
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// This code was developed by the Information Integration Group as part
// of the Karma project at the Information Sciences Institute of the
// University of Southern California.  For more information, publications,
// and related projects, please see: http://www.isi.edu/integration

package algo

import (
	"github.com/puppetlabs/horsehead/v2/datastructure"
	"github.com/puppetlabs/horsehead/v2/graph"
)

// ApproximateSteinerTreeSupportedFeatures are the graph features supported by
// the Steiner tree approximation algorithm.
const ApproximateSteinerTreeSupportedFeatures = graph.DeterministicIteration

// An ApproximateSteinerTree is an approximation of the Steiner tree of a graph
// G for a set of vertices Vr. A Steiner tree is a subset of G such that all Vr
// are present in the subset and the weights of the edges connecting Vr are
// minimized.
type ApproximateSteinerTree struct {
	features graph.GraphFeature
	graph    graph.UndirectedGraph
}

// Features returns the graph features being used for this algorithm.
func (ast *ApproximateSteinerTree) Features() graph.GraphFeature {
	return ast.features
}

// Edges returns the set of edges computed as the approximation of the tree.
func (ast *ApproximateSteinerTree) Edges() graph.EdgeSet {
	return ast.graph.Edges()
}

// AsGraph returns a graph representation of the subset of vertices and edges
// represented by this tree.
func (ast *ApproximateSteinerTree) AsGraph() graph.UndirectedGraph {
	return ast.graph
}

type steinerTreeCostMetric struct {
	edges []graph.Edge
}

// ApproximateSteinerTreeOf computes an approximation of the Steiner tree for a
// given graph.
//
// Steiner trees are known to be NP-complete, but an approximation can be found
// by the following algorithm:
//
//   1. Compute the metric closure of the given graph. This forms a complete
//      graph with edge weights corresponding to the distances between vertices;
//      this specialization is computable in polynomial time.
//   2. Remove all vertices not in the desired subset.
//   3. Compute the minimum spanning tree of the metric closure.
//   4. Expand the minimized edges into a graph.
//   5. Replace the edges of the constructed graph with their paths from the
//      given graph.
//   6. Compute the minimum spanning tree of the expanded graph.
//   7. Prepare a graph with only the edges contained in the tree.
//
// If this proves too slow or inaccurate, it can be further optimized. See
// http://dl.acm.org/citation.cfm?doid=1806689.1806769 for more information.
//
// If any of the given vertices do not exist in this graph, an error of type
// VertexNotFoundError is returned. If no path exists between any two of the
// given vertices, an error of type NotConnectedError is returned.
func ApproximateSteinerTreeOf(g graph.UndirectedGraph, required []graph.Vertex) (*ApproximateSteinerTree, error) {
	// Prerequisite: we deduplicate the required vertices.
	vertices := datastructure.NewHashSet()
	for _, vertex := range required {
		vertices.Add(vertex)
	}

	// If we are tasked with deterministic iteration, we need to have the
	// required vertices in a specific order. Unfortunately, given that they're
	// interface{}, the only thing we can depend on for the order is the source
	// graph. So we have to pick them out of its vertices.
	if g.Features()&graph.DeterministicIteration != 0 {
		remaining := vertices
		vertices = datastructure.NewLinkedHashSet()

		g.Vertices().ForEach(func(candidate graph.Vertex) error {
			if remaining.Contains(candidate) {
				vertices.Add(candidate)
				remaining.Remove(candidate)
			}

			if remaining.Empty() {
				return datastructure.ErrStopIteration
			}

			return nil
		})

		if !remaining.Empty() {
			var first graph.Vertex
			remaining.ForEachInto(func(vertex graph.Vertex) error {
				first = vertex
				return datastructure.ErrStopIteration
			})

			return nil, &graph.VertexNotFoundError{Vertex: first}
		}
	}

	// 1 & 2: Compute the metric closure.
	closure := graph.NewSimpleWeightedGraphWithFeatures(g.Features())

	vertices.ForEach(func(vertex interface{}) error {
		closure.AddVertex(vertex)
		return nil
	})

	err := vertices.ForEach(func(v1 interface{}) error {
		paths := BellmanFordShortestPathsOf(g, v1)

		return vertices.ForEach(func(v2 interface{}) error {
			if v1 == v2 || closure.ContainsEdgeBetween(v1, v2) {
				return nil
			}

			// Save the edges so we can recompute them later.
			edges, err := paths.EdgesTo(v2)
			if err != nil {
				return err
			}

			cost, _ := paths.CostTo(v2)

			metric := &steinerTreeCostMetric{edges}
			closure.AddEdgeWithWeight(v1, v2, metric, cost)

			return nil
		})
	})
	if err != nil {
		return nil, err
	}

	// 3: Compute the minimum spanning tree.
	//
	// This can be optimized as well: dense graphs can run Prim's algorithm in
	// linear time.
	mst := PrimMinimumSpanningTreeOf(closure)

	// 4 & 5: Expand the minimum spanning tree into a graph.
	t := graph.NewUndirectedWeightedPseudographWithFeatures(g.Features())

	vertices.ForEach(func(vertex interface{}) error {
		t.AddVertex(vertex)
		return nil
	})

	mst.Edges().ForEach(func(edge graph.Edge) error {
		metric := edge.(*steinerTreeCostMetric)

		// Expand out the edges.
		for _, step := range metric.edges {
			start, _ := g.SourceVertexOf(step)
			end, _ := g.TargetVertexOf(step)
			weight, _ := g.WeightOf(step)

			t.AddVertex(start)
			t.AddVertex(end)
			t.AddEdgeWithWeight(start, end, step, weight)
		}

		return nil
	})

	// 6: Compute the minimum spanning tree of our final expanded graph.
	keep := PrimMinimumSpanningTreeOf(t).Edges()

	// 7: Remove all edges not in the spanning tree.
	for _, edge := range t.Edges().AsSlice() {
		if !keep.Contains(edge) {
			t.RemoveEdge(edge)
		}
	}

	ast := &ApproximateSteinerTree{
		features: g.Features() & ApproximateSteinerTreeSupportedFeatures,
		graph:    t,
	}

	return ast, nil
}
