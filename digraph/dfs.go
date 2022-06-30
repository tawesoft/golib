package digraph

import (
    "sort"
)

type searchColor int
type searchEdgeType int

const (
    // Vertex state classifications in a search.
    // see CLRS "Introduction to Algorithms", 3rd ed.
    searchColorUndiscovered searchColor = 0 // sometimes called "WHITE"
    searchColorDiscovered   searchColor = 1 // sometimes called "GRAY"
    searchColorFinished     searchColor = 2 // sometimes called "BLACK"

    // Edge type classifications in a search.
    // see CLRS "Introduction to Algorithms", 3rd ed. page 609.
    searchEdgeTypeTree      searchEdgeType = 1
    searchEdgeTypeBack      searchEdgeType = 2
    searchEdgeTypeForward   searchEdgeType = 3
    searchEdgeTypeCross     searchEdgeType = 4
)

// vertexDFS is an annotated visited vertex in a depth-first-search.
type vertexDFS[VertexT any, EdgeT any, WeightT Number] struct {
    vertex      *Vertex[VertexT, EdgeT, WeightT] // matches Digraph.Vertex
    predecessor *vertexDFS[VertexT, EdgeT, WeightT]
    color       searchColor
    d, f        int // discovery and finishing times
}

// DFSResult is the annotated result of a depth first search.
//
// Changing the structure of a graph, such as sorting, adding, or removing
// vertexes and edges, changing edge weights (for a weighted search),or
// changing edge values (for a search constructed using a filter) will
// invalidate the search result.
type DFSResult[VertexT any, EdgeT any, WeightT Number] struct {

    // For any i, where 0 <= i < len(Digraph.Vertexes[i]),
    // Digraph.Vertexes[i] == DFSResult.vertexes[i].vertex
    vertexes []vertexDFS[VertexT, EdgeT, WeightT]
}

// vertexByID returns the annotated vertex for a given vertex ID in a graph.
func (dfs *DFSResult[V, E, W]) vertexByID(id int) *vertexDFS[V, E, W] {
    return &dfs.vertexes[id]
}

// matchingVertex returns the annotated vertex for a given vertex in a graph.
func (dfs *DFSResult[V, E, W]) matchingVertex(v *Vertex[V, E, W]) *vertexDFS[V, E, W] {
    return &dfs.vertexes[v.ID()]
}

// edgeType classifies a DFS edge type between u and v, assuming an edge
// exists. If an edge does not exist, the result is undefined.
func (dfs *DFSResult[V, E, W]) edgeType(u *Vertex[V, E, W], v *Vertex[V, E, W]) searchEdgeType {

    a := dfs.matchingVertex(u)
    b := dfs.matchingVertex(v)

    switch {
        case (a.d < b.d) && (b.d < b.f) && (b.f < a.f) && (b.predecessor == a):
            return searchEdgeTypeTree
        case (b.d < a.d) && (a.d < a.f) && (a.f < b.f): // b.D < a.D < a.F < b.F
            return searchEdgeTypeBack
        case (a.d < b.d) && (b.d < b.f) && (b.f < a.f): // a.D < b.D < b.F < a.F
            return searchEdgeTypeForward
        case (b.d < b.f) && (b.f < a.d) && (a.d < a.f): // b.d < b.F < a.D < a.F
            return searchEdgeTypeCross
        case u == v:
            // a self-loop is also considered a back-edge
            return searchEdgeTypeBack
        default:
            // should never happen
            panic("invalid edge type")
    }
}

// containsCycles is the same as [Digraph.ContainsCycles]
func (dfs *DFSResult[V, E, W]) containsCycles() bool {
    for i := 0; i < len(dfs.vertexes); i++ {
        u := dfs.vertexByID(i)
        edges := u.vertex.Edges
        for j := 0; j < len(edges); j++ {
            v := dfs.matchingVertex(edges[j].Target)
            edgeType := dfs.edgeType(u.vertex, v.vertex)
            if edgeType == searchEdgeTypeBack {
                return true
            }
        }
    }
    return false
}

// TopologicalSort returns the vertexes in topological ordering, using the
// result of a depth-first search as an input. Multiple topological orderings
// may exist, so if you want a specific ordering, see the notes on
// [Digraph.DepthFirstSearch]. Note that only acyclic digraphs (DAGs) have a
// valid topological sort (see [Digraph.ContainsCycles]). This function may
// panic if the input graph contains cycles.
//
// It stores the topological ordering of vertexes in the provided result object,
// resizes the underlying buffer if necessary, and returns that result object
// (or, if nil, creates and returns a new result object).
func (dfs *DFSResult[V, E, W]) TopologicalSort(result []*Vertex[V, E, W]) []*Vertex[V, E, W] {
    if dfs.containsCycles() {
        panic("topological sort is only valid on a DAG")
    }

    if result == nil { result = []*Vertex[V, E, W]{} }
    result = growCap(result, len(dfs.vertexes), len(dfs.vertexes))

    for i := 0; i < len(dfs.vertexes); i++ {
        result[i] = dfs.vertexByID(i).vertex
    }

    sort.Slice(result, func(i int, j int) bool {
        u := dfs.matchingVertex(result[i])
        v := dfs.matchingVertex(result[j])
        return u.f > v.f
    })

    return result
}

// DepthFirstSearch performs a complete "annotated" search of the graph, which
// gives useful properties. It stores this in the provided result object,
// resizes the underlying buffer if necessary, and returns that result object
// (or, if nil, creates and returns a new result object).
//
// If you want a specific ordering of vertexes visited by the search, use
// [Digraph.SortRoots] and [Digraph.SortEdges] first.
func (d *Digraph[V, E, W]) DepthFirstSearch(
    dfs *DFSResult[V, E, W],
    //visitVertex func(*Vertex[V, E, W]),
) *DFSResult[V, E, W] {

    // uses a specific implementation...
    return d.depthFirstSearchRecursive(dfs)
}

// depthFirstSearchRecursive implements a DFS using a recursive method.
func (d *Digraph[V, E, W]) depthFirstSearchRecursive(
    dfs *DFSResult[V, E, W],
) *DFSResult[V, E, W] {
    if dfs == nil {
        dfs = &DFSResult[V, E, W]{}
    }
    z := len(d.Vertexes)
    dfs.vertexes = growCap(dfs.vertexes, z, z)
    clear(dfs.vertexes)

    for i := 0; i < len(dfs.vertexes); i++ {
        u := dfs.vertexByID(i)
        u.vertex = d.Vertexes[i]
    }

    d._dfsRecursiveMain(dfs)

    return dfs
}

// _dfsRecursiveMain corresponds to the DFS function in
// CLRS "Introduction to Algorithms", 3rd ed.
//
// The caller ensures that the dfs argument is non-nil, has an appropriately
// sized buffer, and each element is initialised to its zero value.
//
// Used exclusively by depthFirstSearchRecursive
func (d *Digraph[V, E, W]) _dfsRecursiveMain(dfs *DFSResult[V, E, W]) {
    t := 0

    for i := 0; i < len(d.Vertexes); i++ {
        u := dfs.vertexByID(i)
        if u.color == searchColorUndiscovered {
            t = d._dfsRecursiveVisit(dfs, u, t)
        }
    }
}

// _dfsRecursiveVisit corresponds to the DFS-VISIT function in
// CLRS "Introduction to Algorithms", 3rd ed.
//
// Used exclusively by _dfsRecursiveMain
func (d *Digraph[V, E, W]) _dfsRecursiveVisit(
    dfs *DFSResult[V, E, W],
    u *vertexDFS[V, E, W],
    t int, // time
) int {

    t++
    u.d = t
    u.color = searchColorDiscovered

    edges := u.vertex.Edges
    for j := 0; j < len(edges); j++ {
        v := dfs.matchingVertex(edges[j].Target)
        color := v.color

        if color == searchColorUndiscovered {
            v.predecessor = u
            t = d._dfsRecursiveVisit(dfs, v, t)
        }
    }

    u.color = searchColorFinished
    t++
    u.f = t
    return t
}
