// Package digraph implements a directed graph (a "digraph") and related
// operations.
//
// Many of the algorithms and definitions are thanks to
// CLRS "Introduction to Algorithms", 3rd edition.
package digraph

import (
    "sort"
)

// EdgeDontCare represents the type for an edge value when you don't care what
// that type is (because edges in your graph do not have values).
type EdgeDontCare   = any

// WeightDontCare represents the type for an edge weight when you don't care
// what that type is (because edges in your graph do not have weights).
type WeightDontCare = int

// Edge represents a directed connection from a source vertex to a target
// vertex. The vertex stores a value of type VertexT.
//
// This edge may optionally have an arbitrary user-supplied value associated
// with it of type EdgeT, and an arbitrary numerical weight (such as a distance
// in metres, or a cost in pounds sterling) of type WeightT. If the edge does
// not store a value, use EdgeT [EdgeDontCare]. If the edge does not store a
// weight, use WeightT [WeightDontCare].
//
// A path, as opposed to a connection, is a list of vertexes and edges from a
// source vertex to a target vertex, including all vertexes and edges passed
// on the way. A path has both a distance in terms of number of edges crossed
// and a weighted distance in terms of the sum of edge weights crossed.
//
// Edge values are stored directly in the edge. Use a pointer type where
// appropriate to avoid copying. There is no useful zero value for an Edge.
type Edge[VertexT any, EdgeT any, WeightT Number] struct {
    Value  EdgeT
    Weight WeightT
    Target *Vertex[VertexT, EdgeT, WeightT]
}

// DistanceT is the type of distance in terms of number of edges (this is
// not a weighted distance, which is instead the graph WeightT type).
type DistanceT int

// Vertex represents an object in a digraph.
//
// This vertex has zero-or-more directed edges to a target vertex (if the graph
// permits loops, this includes directed edges to itself). In the case of a
// multidigraph, there may exist multiple directed edges linking the same two
// vertexes.
//
// The vertex has an arbitrary user-supplied value associated with it of type
// VertexT, If the edges do not store a value, use EdgeT [EdgeDontCare]. If the
// edges do not store a weight, use WeightT [WeightDontCare].
//
// Vertex values are stored directly in the vertex. Use a pointer type where
// appropriate to avoid copying. There is no useful zero value for a Vertex.
type Vertex[VertexT any, EdgeT any, WeightT Number] struct {
    id    VertexID
    Value VertexT
    Edges []Edge[VertexT, EdgeT, WeightT]
}

// VertexID is the type of the index of a vertex in a graph (see [Vertex.ID]).
//
// For any i, where 0 <= i < len(Digraph.Vertexes[i]),
// Digraph.Vertexes[i].ID() == i.
//
// Vertex IDs may change when a graph is sorted (using [Digraph.SortRoots]),
// when new vertexes are added, and when vertexes are removed.
type VertexID int32

// ID returns the current index of a vertex in a graph (see [VertexID]).
func (v *Vertex[V, E, W]) ID() VertexID {
    return v.id
}

// PathVertex is a tuple of a vertex and the edge it was reached from when
// calculating a path.
type PathVertex[VertexT any, EdgeT any, WeightT Number] struct {
    Vertex *Vertex[VertexT, EdgeT, WeightT]
    Via    *Edge[VertexT, EdgeT, WeightT]
}

// Path is an ordered list of zero or more PathVertexes making up a path. A Path can be empty,
// but is never nil.
type Path[VertexT any, EdgeT any, WeightT Number] []PathVertex[VertexT, EdgeT, WeightT]

// Matrix is a two-dimensional array of values, such as a graph's adjacency
// matrix (with values of type VertexID), weighted adjacency matrix (with
// values of type WeightT), or weighted distance matrix (also with values of
// type WeightT).
//
// Changing the structure of a graph, such as sorting, adding, or removing
// vertexes and edges, changing edge weights (for a weighted matrix),or
// changing edge values (for a matrix constructed using a filter) will
// invalidate the matrix.
type Matrix[T Number] struct {
    width  int
    values []T // size is width squared
}

type AdjacencyMatrix = Matrix[DistanceT]
type DistanceMatrix  = Matrix[DistanceT]

func (m Matrix[T]) get(x VertexID, y VertexID) T {
    return m.values[(int(x) * m.width) + int(y)]
}

func (m Matrix[T]) set(x VertexID, y VertexID, w int, value T) {
    m.values[(int(x) * m.width) + int(y)] = value
}

func (m Matrix[T]) add(x VertexID, y VertexID, w int, value T) {
    m.values[(int(x) * m.width) + int(y)] += value
}

// A Digraph is a directed graph of vertexes and (directed) edges.
//
// The graph may be disconnected (a disjoint union of disconnected subgraphs)
// or connected, may be a multigraph/multidigraph (may have more than one edge
// between the same two vertices) or a simple graph, may be weighted (with any
// [Number] type) or unweighted, may contain loops or may not permit
// loops, may contain cycles or be acyclic.
//
// Properties of the graph, such as whether it contains cycles, can be queried
// e.g. to assert that the graph is a valid directed acyclic graph (DAG).
//
// If the edges do not store a value, use EdgeT [EdgeDontCare]. If the graph is
// unweighted, use WeightT [WeightDontCare].
type Digraph[VertexT any, EdgeT any, WeightT Number] struct {
    Vertexes []*Vertex[VertexT, EdgeT, WeightT]

    positiveInfiniteWeight WeightT
    negativeInfiniteWeight WeightT
}

func New[VertexT any, EdgeT any, WeightT Number]() *Digraph[VertexT, EdgeT, WeightT] {
    return &Digraph[VertexT, EdgeT, WeightT]{}
}

// infiniteWeightedDistance returns a value representing infinity for a type
// of weight. This is positive infinity if sign >= 0, negative infinity if
// sign < 0.
func (d *Digraph[V, E, W]) infiniteWeightedDistance(sign int) W {
    if d.positiveInfiniteWeight == 0 { d.positiveInfiniteWeight = InfWeight[W, int]( 1) }
    if d.negativeInfiniteWeight == 0 { d.negativeInfiniteWeight = InfWeight[W, int](-1) }
    if sign >= 0 { return d.positiveInfiniteWeight } else { return d.negativeInfiniteWeight }
}

// IsInfiniteDistance returns true if a distance (in terms of number of edges)
// is infinite, representing an unreachable vertex.
func (d *Digraph[V, E, W]) IsInfiniteDistance(e DistanceT) bool {
    return isInfiniteDistance(e)
}

// IsFiniteDistance returns true if a distance (in terms of number of edges)
// is finite, representing a reachable vertex.
func (d *Digraph[V, E, W]) IsFiniteDistance(e DistanceT) bool {
    return !isInfiniteDistance(e)
}

// IsInfiniteWeightedDistance returns true if a weighted distance (in terms of
// a sum of weights along edges) is infinite, representing an unreachable
// vertex.
func (d *Digraph[V, E, W]) IsInfiniteWeightedDistance(w W) bool {
    return (w == d.infiniteWeightedDistance(1)) || (w == d.infiniteWeightedDistance(-1))
}

// IsFiniteWeightedDistance returns true if a weighted distance (in terms of
// a sum of weights along edges) is finite, representing an reachable vertex.
func (d *Digraph[V, E, W]) IsFiniteWeightedDistance(w W) bool {
    return (w != d.infiniteWeightedDistance(1)) && (w != d.infiniteWeightedDistance(-1))
}

// SortRoots performs an in-place stable sort of the graph's vertexes such that
// root vertexes (if any) appear before non-root vertexes, and the root
// vertexes are sorted by the comparison function lessThan. This can be used
// to achieve a specific ordering for searches or topological sorts.
//
// Note that this can change vertex IDs (see [VertexID]) and therefore
// invalidates any existing calculated matrix, path, search result of the graph
// etc., including the input matrix when this method returns.
//
// The function lessThan defines an ordering by returning true if a vertex
// comes before another. It is only applied to root vertexes.
func (d *Digraph[V, E, W]) SortRoots(
    adjacencyMatrix *AdjacencyMatrix,
    lessThan func(*Vertex[V, E, W], *Vertex[V, E, W]) bool,
) {
    rootsLessThan := func(i, j int) bool {
        u := d.Vertexes[i]
        v := d.Vertexes[j]
        uDegree := d.Indegree(adjacencyMatrix, u)
        vDegree := d.Indegree(adjacencyMatrix, v)
        switch {
            case (uDegree == 0) && (vDegree > 0):
                return true
            case (uDegree > 0) && (vDegree == 0):
                return false
            case (uDegree == 0) && (vDegree == 0):
                return lessThan(u, v)
            default:
                return i < j
        }
    }
    sort.SliceStable(d.Vertexes, rootsLessThan)

    // update IDs to match sorted order
    for i := 0; i < len(d.Vertexes); i++ {
        d.Vertexes[i].id = VertexID(i)
    }
}

// SortEdges performs an in-place stable sort of each vertex's edge. This can
// be used to achieve a specific ordering for searches.
//
// The function lessThan defines an ordering by returning true if an edge
// on the vertex comes before another edge on that vertex.
//
// Sorting edges does not change vertex IDs.
func (d *Digraph[V, E, W]) SortEdges(
    lessThan func(*Vertex[V, E, W], *Edge[V, E, W], *Edge[V, E, W]) bool,
) {
    for i := 0; i < len(d.Vertexes); i++ {
        v := d.Vertexes[i]
        edgeLessThan := func(i, j int) bool {
            a := &v.Edges[i]
            b := &v.Edges[j]
            return lessThan(v, a, b)
        }

        sort.SliceStable(v.Edges, edgeLessThan)
    }
}

// IsWeaklyConnected returns true if there is at least one undirected path
// between every possible vertex pair. A graph with just one vertex is always
// connected.
func (d *Digraph[V, E, W]) IsWeaklyConnected() bool {
    panic("TODO")
}

// IsStronglyConnected returns true if there is at least one directed path
// between every possible vertex pair (in both directions). A graph with
// just one vertex is always connected.
func (d *Digraph[V, E, W]) IsStronglyConnected() bool {
    panic("TODO")
}

// IsComplete returns true if there is exactly one directed edge between
// every possible vertex pair in each direction. For any pair (u,v), u has a
// directed edge to v, and v has a directed edge to u.
func (d *Digraph[V, E, W]) IsComplete() bool {
    panic("TODO")
}

// IsTournament returns true if there is exactly one directed edge between
// every possible vertex pair in either direction. For any pair (u,v), either u
// has a // directed edge to v, or v has a directed edge to u, but not both.
func (d *Digraph[V, E, W]) IsTournament() bool {
    panic("TODO")
}

// IsSimple returns true if for any two vertexes there is at most one edge
// (in the same direction) between them. If false, the digraph is a
// multigraph/multidigraph. A completed adjacency matrix is used as input to do
// this efficiently.
func (d *Digraph[V, E, W]) IsSimple(mat *AdjacencyMatrix) bool {
    for i := 0; i < len(d.Vertexes); i++ {
        u := d.Vertexes[i]
        edges := u.Edges
        for j := 0; j < len(edges); i++ {
            v := edges[j].Target
            if d.Adjacency(mat, u, v) > 1 { return false }
        }
    }
    return true
}

// ContainsCycles returns true if there exists any directed cycle (i.e. any
// path between a vertex and itself). If false, the digraph is a directed
// acyclic graph (DAG). A completed depth-first search of the graph is used as
// input to do this efficiently.
func (d *Digraph[V, E, W]) ContainsCycles(dfs *DFSResult[V, E, W]) bool {
    return dfs.containsCycles()
}

// ContainsLoops returns true if there is any edge that connects a vertex to
// itself. A completed adjacency matrix is used as input to do this
// efficiently.
func (d *Digraph[V, E, W]) ContainsLoops(mat *AdjacencyMatrix) bool {
    for i := 0; i < len(d.Vertexes); i++ {
        v := d.Vertexes[i]
        if d.Adjacency(mat, v, v) > 0 { return true }
    }
    return false
}

// IsTree returns true if for any two vertexes there is exactly one
// undirected path between them (the graph is said to be a "directed tree" or
// polytree). Alternatively, each vertex has exactly one parent, except for the
// root (which has no parent).
func (d *Digraph[V, E, W]) IsTree() bool { return false } // TODO

// IsForest returns true if for any two vertexes there is at most one
// undirected path between them, or (equivalently) the graph is a disjoint
// union of trees (the graph is said to be a "directed forest" or polyforest)
// i.e. every node with no parent is the root node of an individual tree.
func (d *Digraph[V, E, W]) IsForest() bool { return false } // TODO

// Indegree returns the number of edges pointing to the given vertex.
func (d *Digraph[V, E, W]) Indegree(mat *AdjacencyMatrix, v *Vertex[V, E, W]) int {
    sum := 0

    for i := 0; i < len(d.Vertexes); i++ {
        sum += int(mat.get(VertexID(i), v.id))
    }

    return sum
}

// Outdegree returns the number of edges pointing away from the given vertex.
func (d *Digraph[V, E, W]) Outdegree(mat *AdjacencyMatrix, v *Vertex[V, E, W]) int {
    sum := 0

    for i := 0; i < len(d.Vertexes); i++ {
        sum += int(mat.get(v.id, VertexID(i)))
    }

    return sum
}

// Roots returns all vertexes in the graph with no parents. It stores this in
// the provided result object, resizes the underlying buffer if necessary, and
// returns that result object (or, if nil, creates and returns a new result
// object).
func (d *Digraph[V, E, W]) Roots(
    mat *AdjacencyMatrix,
    result []*Vertex[V, E, W],
) []*Vertex[V, E, W] {
    if result == nil { result = make([]*Vertex[V, E, W], len(d.Vertexes)) }
    result = growCap(result, 0, len(d.Vertexes))
    clear(result)

    for i := 0; i < len(d.Vertexes); i++ {
        u := d.Vertexes[i]
        if d.Indegree(mat, u) == 0 {
            result = append(result, u)
        }
    }

    return result
}

// Inputs returns the adjacent vertexes and edges pointing to the given vertex.
// It stores this in the provided result object, resizes the underlying buffer
// if necessary, and returns that result object (or, if nil, creates and
// returns a new result object).
func (d *Digraph[V, E, W]) Inputs(
    mat *AdjacencyMatrix,
    result Path[V, E, W],
    v *Vertex[V, E, W],
) Path[V, E, W] {
    if result == nil { result = Path[V, E, W]{} }
    result = growCap(result, 0, len(d.Vertexes))
    clear(result)

    for i := 0; i < len(d.Vertexes); i++ {
        u := d.Vertexes[i]
        if d.Adjacency(mat, u, v) == 0 { continue }

        for j := 0; j < len(u.Edges); j++ {
            edge := u.Edges[j]
            if edge.Target == v {
                result = append(result, PathVertex[V, E, W]{
                    Vertex: u,
                    Via:    &edge,
                })
            }
        }
    }

    return result
}

// Adjacency returns the number of edges pointing from a to b (where a and b
// are adjacent, or a == b). In a simple graph, this is always zero or one. In
// a multigraph, it could be zero or any positive number.
func (d *Digraph[V, E, W]) Adjacency(
    mat *AdjacencyMatrix,
    a *Vertex[V, E, W],
    b *Vertex[V, E, W],
) int {
    return int(mat.get(a.ID(), b.ID()))
}

// WeightedAdjacency returns the weighted distance from a to b (where a and b
// are adjacent, or a == b), using the weighted adjacency matrix mat. In a
// simple graph, this is always either the weight of exactly one edge, or
// infinity if there is no directed edge from a to b.
//
// In a multigraph, it is always either the weight calculated by the reducer
// function used to generate the weighted adjacency matrix (such as minimum
// weight, or sum of weights), or infinity if there is no directed edge from a
// to b.
func (d *Digraph[V, E, W]) WeightedAdjacency(
    mat *Matrix[W],
    a *Vertex[V, E, W],
    b *Vertex[V, E, W],
) W {
    return mat.get(a.id, b.id)
}

// AdjacencyMatrix calculates the number of directed edges between each
// adjacent vertex, stores this in the provided matrix buffer, resizes the
// underlying buffer if necessary, and returns that matrix (or, if nil, creates
// and returns a new matrix).
func (d *Digraph[V, E, W]) AdjacencyMatrix(mat *AdjacencyMatrix) *AdjacencyMatrix {
    return d.AdjacencyMatrixFiltered(mat, func(*E) bool { return true } )
}

// AdjacencyMatrixFiltered behaves as [AdjacencyMatrix] except that it only
// considers edges where the provided function, which operates on the value of
// the edge, returns true.
func (d *Digraph[V, E, W]) AdjacencyMatrixFiltered(
    mat *AdjacencyMatrix,
    filter func(e *E) bool,
) *AdjacencyMatrix {
    if mat == nil {
        mat = &AdjacencyMatrix{}
    }

    z := len(d.Vertexes)
    z2 := z * z

    mat.width = z
    mat.values = growCap(mat.values, z2, z2)
    clear(mat.values)

    for i := 0; i < len(d.Vertexes); i++ {
        source := d.Vertexes[i]

        for j := 0; j < len(source.Edges); j++ {

            edge := &source.Edges[j]
            if !filter(&edge.Value) { continue }

            dest := edge.Target
            mat.add(source.id, dest.id, z, 1)
        }
    }

    return mat
}

// EdgeWeightReducer defines a method to reduce multiple edge weights into a
// single weight. This can be the case in a multigraph/multidigraph.
//
// The Reduce function is applied left-to-right on a vertex's edges. For a
// given ordering, see [Digraph.SortEdges]. The first argument to the first
// invocation of the reducer is the provided Identity value (i.e. for sum, this
// is zero. For multiply, this is one).
type EdgeWeightReducer[W Number] struct {
    Identity W
    Reduce func(W, W) W
}

func (r *EdgeWeightReducer[W]) start(a W) W {
    return r.run(r.Identity, a)
}

func (r *EdgeWeightReducer[W]) run(a W, b W) W {
    return r.Reduce(a, b)
}

// NewEdgeWeightReducerMinimum returns an edge reducer that reduces a list of
// edge weights of type W to a minimum value.
func NewEdgeWeightReducerMinimum[W Number]() *EdgeWeightReducer[W] {
    inf := InfWeight[W, int](1)
    return &EdgeWeightReducer[W]{
        Identity: inf,
        Reduce: func(a W, b W) W {
            switch
            {
                case (a == inf): return b
                case (b == inf): return a
                case (a <= b):   return a
                default:         return b
            }
        },
    }
}

// NewEdgeWeightReducerMaximum returns an edge reducer that reduces a list of
// edge weights of type W to a maximum value.
func NewEdgeWeightReducerMaximum[W Number]() *EdgeWeightReducer[W] {
    inf := InfWeight[W, int](-1)
    return &EdgeWeightReducer[W]{
        Identity: inf,
        Reduce: func(a W, b W) W {
            switch
            {
                case (a == inf): return a
                case (b == inf): return b
                case (a >= b):   return a
                default:         return b
            }
        },
    }
}

// NewEdgeWeightReducerSum returns an edge reducer that reduces a list of edge
// weights of type W to the sum of those weights.
func NewEdgeWeightReducerSum[W Number]() *EdgeWeightReducer[W] {
    return &EdgeWeightReducer[W]{
        Identity: 0,
        Reduce: func(a W, b W) W {
            return a + b
        },
    }
}

// WeightedAdjacencyMatrix calculates the weights of edges between each
// adjacent vertex, stores this in the provided matrix buffer, resizes the
// underlying buffer if necessary, and returns that matrix (or, if nil, creates
// and returns a new matrix).
//
// In the case of a multigraph (multiple edges between two vertexes), the
// provided reducer is used to reduce multiple weights into a single value. For
// example, sum all weights, or return the minimum weight. If you are certain
// that the graph is not a multigraph (i.e. is simple), then the reducer may be
// nil.
func (d *Digraph[V, E, W]) WeightedAdjacencyMatrix(
    mat *Matrix[W],
    reducer *EdgeWeightReducer[W],
) *Matrix[W] {
    filter := func(*E) bool { return true }
    return d.WeightedAdjacencyMatrixFiltered(mat, reducer, filter)
}

// WeightedAdjacencyMatrixFiltered behaves as [WeightedAdjacencyMatrix] except
// that it only considers edges where the provided filter function, which
// operates on the value of the edge, returns true.
func (d *Digraph[V, E, W]) WeightedAdjacencyMatrixFiltered(
    mat *Matrix[W],
    reducer *EdgeWeightReducer[W],
    filter func(e *E) bool,
) *Matrix[W] {
    if mat == nil {
        mat = &Matrix[W]{}
    }

    z := len(d.Vertexes)
    z2 := z * z

    mat.values = growCap(mat.values, z2, z2)
    mat.width = z
    clear(mat.values)

    inf := d.infiniteWeightedDistance(1)
    for i := 0; i < z2; i++ {
        mat.values[i] = inf
    }

    for i := 0; i < len(d.Vertexes); i++ {
        source := d.Vertexes[i]

        for j := 0; j < len(source.Edges); j++ {

            edge := &source.Edges[j]
            if !filter(&edge.Value) { continue }

            dest := edge.Target

            weight := mat.get(source.id, dest.id)
            var reducedWeight W

            if reducer == nil {
                reducedWeight = edge.Weight
            } else if d.IsInfiniteWeightedDistance(weight) {
                reducedWeight = reducer.start(edge.Weight)
            } else {
                reducedWeight = reducer.run(weight, edge.Weight)
            }

            mat.set(source.id, dest.id, z, reducedWeight)
        }
    }

    return mat
}

// DistanceMatrix returns the distance  TODO
//
// The returned Matrix is no longer current if the graph has been modified by
// adding, changing, or removing vertexes, edges, or their weights or values.
func (d *Digraph[V, E, W]) DistanceMatrix(mat *Matrix[int]) *Matrix[W] {
    return d.DistanceMatrixFiltered(mat)
}

func (d *Digraph[V, E, W]) DistanceMatrixFiltered(mat *Matrix[int]) *Matrix[W] {
    // TODO
    return nil
}

func (d *Digraph[V, E, W]) WeightedDistanceMatrix(mat *Matrix[int]) *Matrix[W] {
    // TODO
    return d.WeightedDistanceMatrixFiltered(mat)
}

func (d *Digraph[V, E, W]) WeightedDistanceMatrixFiltered(mat *Matrix[int]) *Matrix[W] {
    // TODO
    return nil
}

// AddVertex creates a new vertex in the graph, holding the arbitrary
// user-supplied value, and returns a pointer that identifies that vertex in
// the graph.
func (d *Digraph[V, E, W]) AddVertex(value V) *Vertex[V, E, W] {
    v := &Vertex[V, E, W]{
        id:    VertexID(len(d.Vertexes)),
        Value: value,
    }
    d.Vertexes = append(d.Vertexes, v)
    return v
}

// AddEdge creates a new edge between two vertexes in the graph, holding the
// arbitrary user-supplied value.
func (d *Digraph[V, E, W]) AddEdge(from *Vertex[V, E, W], to *Vertex[V, E, W], value E) {
    d.AddWeightedEdge(from, to, value, zero[W]())
}

// AddWeightedEdge creates a new edge between two vertexes in the graph,
// holding the arbitrary user-supplied value, and with a given weight.
func (d *Digraph[V, E, W]) AddWeightedEdge(from *Vertex[V, E, W], to *Vertex[V, E, W], value E, weight W) {
    from.Edges = append(from.Edges, Edge[V, E, W]{
        Value: value,
        Weight: weight,
        Target: to,
    })
}

// AddUniqueEdge creates a new edge between two vertexes in the graph,
// holding the arbitrary user-supplied value, if that edge does not already
// exist, and returns true. Otherwise, if the edge already exists, returns
// false and does nothing.
func (d *Digraph[V, E, W]) AddUniqueEdge(
    from *Vertex[V, E, W],
    to *Vertex[V, E, W],
    value E,
) bool {
    return d.AddUniqueWeightedEdge(from, to, value, zero[W]())
}

// AddUniqueWeightedEdge creates a new edge between two vertexes in the graph,
// holding the arbitrary user-supplied value, and with a given weight, if that
// edge does not already exist, and returns true. Otherwise, if the edge
// already exists, returns false and does nothing.
func (d *Digraph[V, E, W]) AddUniqueWeightedEdge(
    from *Vertex[V, E, W],
    to *Vertex[V, E, W],
    value E,
    weight W,
) bool {
    filterEqual := func(*E, *E) bool { return true }
    return d.AddUniqueWeightedEdgeFiltered(from, to, value, weight, filterEqual)
}

// AddUniqueEdgeFiltered behaves as AddUniqueEdge, however an edge is only
// deemed to be unique if the provided equality function returns false for
// the new edge value and the existing edge value.
func (d *Digraph[V, E, W]) AddUniqueEdgeFiltered(
    from *Vertex[V, E, W],
    to *Vertex[V, E, W],
    value E,
    equal func(a *E, b *E) bool,
) bool {
    return d.AddUniqueWeightedEdgeFiltered(from, to, value, zero[W](), equal)
}

// AddUniqueWeightedEdgeFiltered behaves as AddUniqueWeightedEdge, however an
// edge is only deemed to be unique if the provided equality function returns
// false for the new edge value and the existing edge value. Weights are ignored
// for the purposes of uniqueness.
func (d *Digraph[V, E, W]) AddUniqueWeightedEdgeFiltered(
    from *Vertex[V, E, W],
    to *Vertex[V, E, W],
    value E,
    weight W,
    equal func(a *E, b *E) bool,
) bool {
    filter := func(e *E) bool {
        return equal(&value, e)
    }

    edge := d.FindEdgeFiltered(from, to, filter)
    if edge == nil { // no equal edge found
        d.AddWeightedEdge(from, to, value, weight)
        return true
    } else {
        return false
    }
}

// FindEdge returns an edge, if any, between two adjacent vertexes. Otherwise,
// returns nil. If multiple edges exist between the two vertexes, the specific
// edge returned is arbitrary.
//
// If an adjacency matrix has been constructed, and you only want to know
// if an edge exists and not to find the actual edge, it is more efficient to
// use the [Digraph.Adjacency] method.
func (d *Digraph[V, E, W]) FindEdge(
    from *Vertex[V, E, W],
    to *Vertex[V, E, W],
) *Edge[V, E, W] {
    filterAny := func(*E) bool { return true }
    return d.FindEdgeFiltered(from, to, filterAny)
}

// FindEdgeFiltered returns an edge between two adjacent vertexes if it exists
// and only if the provided function, which operates on the value of the edge,
// returns true. Otherwise, returns nil. If multiple matching edges exist
// between the two vertexes, the specific edge returned is arbitrary.
//
// It can be more efficient to construct and reuse a filtered adjacency matrix
// and use the [Digraph.Adjacency] method instead in cases where you only want
// to know if an edge exists and not to find the actual edge.
func (d *Digraph[V, E, W]) FindEdgeFiltered(
    from *Vertex[V, E, W],
    to *Vertex[V, E, W],
    f func(e *E) bool,
) *Edge[V, E, W] {
    for i := 0; i < len(from.Edges); i++ {
        edge := &from.Edges[i]
        if (edge.Target == to) && f(&edge.Value) {
            return edge
        }
    }

    return nil
}
