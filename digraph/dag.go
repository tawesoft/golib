package digraph

type vertexDAGSearch[VertexT any, EdgeT any, WeightT Number] struct {
    vertex           *Vertex[VertexT, EdgeT, WeightT] // matches Digraph.Vertex
    predecessor      *vertexDAGSearch[VertexT, EdgeT, WeightT]
    distance         DistanceT // number of edges from source
    weightedDistance WeightT
}

// DAGSearchResult is the annotated result of a directed acyclic graph (a DAG).
//
// Changing the structure of a graph, such as sorting, adding, or removing
// vertexes and edges, changing edge weights (for a weighted search),or
// changing edge values (for a search constructed using a filter) will
// invalidate the search result.
type DAGSearchResult[VertexT any, EdgeT any, WeightT Number] struct {
    start *vertexDAGSearch[VertexT, EdgeT, WeightT]

    // For any i, where 0 <= i < len(Digraph.Vertexes[i]),
    // Digraph.Vertexes[i] == BFSResult.vertexes[i].vertex
    vertexes []vertexDAGSearch[VertexT, EdgeT, WeightT]
}

// vertexByID returns the annotated vertex for a given vertex ID in a graph.
func (dags *DAGSearchResult[V, E, W]) vertexByID(id int) *vertexDAGSearch[V, E, W] {
    return &dags.vertexes[id]
}

// matchingVertex returns the annotated vertex for a given vertex in a graph.
func (dags *DAGSearchResult[V, E, W]) matchingVertex(v *Vertex[V, E, W]) *vertexDAGSearch[V, E, W] {
    return &dags.vertexes[v.ID()]
}

// DAGSearchWeighted performs an "annotated" search of a weighted directed
// acyclic graph (a weighted DAG) from some root node, which gives useful
// properties. It stores this in the provided result object, resizes the
// underlying buffer if necessary, and returns that result object (or, if nil,
// creates and returns a new result object).
//
// The topologicalSort input is an ordering of the vertexes in topological
// order (see [DFSResult.TopologicalSort]).
func (d *Digraph[V, E, W]) DAGSearchWeighted(
    dags *DAGSearchResult[V, E, W],
    start *Vertex[V, E, W],
    topologicalSort []*Vertex[V, E, W],
) *DAGSearchResult[V, E, W] {
    if dags == nil {
        dags = &DAGSearchResult[V, E, W]{}
    }
    z := len(d.Vertexes)
    dags.vertexes = growCap(dags.vertexes, z, z)
    clear(dags.vertexes)

    dags.start = &dags.vertexes[start.ID()]

    for i := 0; i < len(d.Vertexes); i++ {
        v := dags.vertexByID(i)
        v.vertex = d.Vertexes[i]
        v.predecessor = nil
        v.distance = positiveInfiniteEdgeCount
        v.weightedDistance = d.infiniteWeightedDistance(1)
    }

    dags.start.distance = 0
    dags.start.weightedDistance = 0

    inf := infinities[W]{
        positive: d.infiniteWeightedDistance(+1),
        negative: d.infiniteWeightedDistance(-1),
    }

    for i := 0; i < len(topologicalSort) - 1; i++ {
        u := dags.matchingVertex(topologicalSort[i])
        edges := u.vertex.Edges

        for j := 0; j < len(edges); j++ {
            v := dags.matchingVertex(edges[j].Target)
            dags.relax(u, v, edges[j].Weight, inf)
        }
    }

    return dags
}

func (*DAGSearchResult[V, E, W]) relax (
    u *vertexDAGSearch[V, E, W],
    v *vertexDAGSearch[V, E, W],
    w W,
    inf infinities[W],
) {
    sum := inf.sum(u.weightedDistance, w)

    if v.weightedDistance > sum {
        v.weightedDistance = sum
        v.predecessor = u
        v.distance = u.distance + 1
    }
}

// ShortestPath returns a list of the vertexes along the shortest path from the
// source vertex to the given destination vertex, using the result of a
// DAG search from the source vertex as an input. If the search
// was a weighted search, the shortest path is a path with the smallest
// weight. Otherwise, the shortest path is the path with the fewest number of
// edges.
//
// Multiple shortest paths may exist, so if you want a specific ordering, use
// [Digraph.SortRoots] or [Digraph.SortEdges]. Every vertex is reachable by
// itself with a path of zero edges, so if `a == b`, then the shortest path is
// simply a list containing only vertex `a`. Every valid path will begin with
// `a` and end with `b`. If no path exists, the result will be an empty list
// (not nil).
//
// The function stores the vertexes encountered in the provided result object,
// resizes the underlying buffer if necessary, and returns that result object
// (or, if nil, creates and returns a new result object).
func (dags *DAGSearchResult[V, E, W]) ShortestPath(
    result []*Vertex[V, E, W],
    v *Vertex[V, E, W],
) []*Vertex[V, E, W] {
    if result == nil { result = []*Vertex[V, E, W]{} }
    result = growCap(result, 0, len(dags.vertexes))

    // TODO refactor this using interfaces because its the same for all search
    // types.

    s := dags.start
    u := dags.matchingVertex(v)

    for {
        if (s == u) {
            result = append(result, s.vertex)
            break
        } else if u.predecessor == nil {
            result = result[0:0]
            break
        } else {
            result = append(result, v)
            u = u.predecessor
            v = u.vertex
        }
    }

    if len(result) > 0 {  reverse(result) }
    return result
}
