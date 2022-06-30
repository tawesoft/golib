package digraph

// vertexBFS is an annotated visited vertex in a breadth-first-search
type vertexBFS[VertexT any, EdgeT any, WeightT Number] struct {
    vertex           *Vertex[VertexT, EdgeT, WeightT] // matches Digraph.Vertex
    predecessor      *vertexBFS[VertexT, EdgeT, WeightT]
    color            searchColor
    distance         DistanceT // edges from source
    weightedDistance WeightT // sum of edge weights from source
}

// BFSResult is the annotated result of a breadth-first search
//
// Changing the structure of a graph, such as sorting, adding, or removing
// vertexes and edges, changing edge weights (for a weighted search),or
// changing edge values (for a search constructed using a filter) will
// invalidate the search result.
type BFSResult[VertexT any, EdgeT any, WeightT Number] struct {

    // For any i, where 0 <= i < len(Digraph.Vertexes[i]),
    // Digraph.Vertexes[i] == BFSResult.vertexes[i].vertex
    vertexes []vertexBFS[VertexT, EdgeT, WeightT]

    start    *vertexBFS[VertexT, EdgeT, WeightT]
    queue    []*vertexBFS[VertexT, EdgeT, WeightT]
}

// vertexByID returns the annotated vertex for a given vertex ID in a graph.
func (bfs *BFSResult[V, E, W]) vertexByID(id int) *vertexBFS[V, E, W] {
    return &bfs.vertexes[id]
}

// matchingVertex returns the annotated vertex for a given vertex in a graph.
func (bfs *BFSResult[V, E, W]) matchingVertex(v *Vertex[V, E, W]) *vertexBFS[V, E, W] {
    return &bfs.vertexes[v.ID()]
}

// BreadthFirstSearch performs a complete search of the reachable graph from a
// given start vertex, taking the shortest number of edges. The resulting
// search tree gives useful properties.
//
// The search stores a result in the provided result object, resizing its
// underlying buffer if necessary, and returns that result object (or, if nil,
// creates and returns a new result object). Use the same result object across
// multiple searches to minimise memory allocations.
//
// If you want a specific ordering of vertexes visited by the search, use
// [Digraph.SortRoots] and [Digraph.SortEdges] first.
//
// The visitVertex callback function is called when each vertex is first
// visited in the search. If nil, it is ignored.
//
// Computed vertex distances from source, in terms of the number of edges
// crossed, are of type int32. If maxDepth is greater than zero, the search
// excludes any vertex with a distance greater than maxDepth.
//
// If the graph is modified, the resulting search is no longer valid and must
// not be queried.
func (d *Digraph[V, E, W]) BreadthFirstSearch(
    bfs *BFSResult[V, E, W],
    start *Vertex[V, E, W],
    maxDepth DistanceT,
    visitVertex func(*Vertex[V, E, W]),
) *BFSResult[V, E, W] {
    if maxDepth <= 0 { maxDepth = positiveInfiniteEdgeCount }
    if bfs == nil {
        bfs = &BFSResult[V, E, W]{}
    }
    z := len(d.Vertexes)
    bfs.vertexes = growCap(bfs.vertexes, z, z)
    clear(bfs.vertexes)

    bfs.queue = growCap(bfs.queue, 0, z)

    bfs.start = bfs.matchingVertex(start)

    for i := 0; i < len(d.Vertexes); i++ {
        v := &bfs.vertexes[i]
        v.vertex = d.Vertexes[i]
        v.predecessor = nil
        v.distance = positiveInfiniteEdgeCount
        v.weightedDistance = d.infiniteWeightedDistance(1)
    }

    bfs.start.distance = 0
    bfs.start.weightedDistance = 0
    bfs.start.color = searchColorDiscovered
    // if (visitVertex != nil) && (!visitVertex(s)) { return bfs }
    if (visitVertex != nil) { visitVertex(start) }

    bfs.queue = append(bfs.queue, bfs.start)

    for len(bfs.queue) != 0 {
        // dequeue
        u := bfs.queue[len(bfs.queue)-1]
        bfs.queue = bfs.queue[:len(bfs.queue)-1]

        edges := u.vertex.Edges
        for i := 0; i < len(edges); i++ {
            v := bfs.matchingVertex(edges[i].Target)
            if v.color != searchColorUndiscovered { continue }

            v.color       = searchColorDiscovered
            v.distance    = u.distance + 1
            v.predecessor = u

            if visitVertex != nil { visitVertex(v.vertex) }
            if v.distance >= maxDepth { continue }

            // queue
            bfs.queue = append(bfs.queue, v)
        }
    }

    return bfs
}

// Distance returns the distance of the shortest path from the source vertex
// to the given destination vertex, using the result of a breadth-first search
// from the source vertex as an input, in terms of the number of edges crossed
// (i.e. NOT the weighted distance). If the search was weighted this may not be
// the shortest path in terms of edges, but is instead the number of edges on
// the shortest path by distance. If no path exists, the distance is infinite
// (represented by [maths.MaxInt32]) - see [Digraph.IsInfiniteDistance].
func (bfs *BFSResult[V, E, W]) Distance(v *Vertex[V, E, W]) DistanceT {
    return bfs.matchingVertex(v).distance
}

// WeightedDistance returns the distance of the shortest path from the source
// vertex to the given destination vertex, in terms of the minimum weight,
// using the result of a breadth-first search from the source vertex as an
// input. If no path exists, or if the search was not weighted, the distance is
// infinite ([Digraph.IsInfiniteWeight] returns true).
func (bfs *BFSResult[V, E, W]) WeightedDistance(v *Vertex[V, E, W]) W {
    return bfs.matchingVertex(v).weightedDistance
}

// Predecessor returns, from a search starting at the source vertex (forming a
// breadth-first search tree), the parent of a given vertex. If the given
// vertex is not reachable from search, returns nil. If the given vertex is
// the source vertex (and has no predecessor), returns nil.
func (bfs *BFSResult[V, E, W]) Predecessor(v *Vertex[V, E, W]) *Vertex[V, E, W] {
    u := bfs.matchingVertex(v)
    if isInfiniteDistance(u.distance) { return nil }
    predecessor := u.predecessor

    if predecessor == nil {
        if u.distance != 0 {
            panic("reachable, non-root vertex has no predecessor")
        }
        return nil
    }

    return predecessor.vertex
}

// ShortestPath returns a list of the vertexes along the shortest path from the
// source vertex to the given destination vertex, using the result of a
// breadth-first search from the source vertex as an input. If the search
// was a weighted search, the shortest path is a path with the smallest
// weight. Otherwise, the shortest path is the path with the fewest number of
// edges.
//
// Multiple shortest paths may exist, so if you want a specific  ordering, see
// the notes on [Digraph.BreadthFirstSearch]. Every vertex is reachable by
// itself with a path of zero edges, so if `a == b`, then the shortest path is
// simply a list containing only vertex `a`. Every valid path will begin with
// `a` and end with `b`. If no path exists, the result will be an empty list
// (not nil).
//
// The function stores the vertexes encountered in the provided result object,
// resizes the underlying buffer if necessary, and returns that result object
// (or, if nil, creates and returns a new result object).
func (bfs *BFSResult[V, E, W]) ShortestPath(
    result []*Vertex[V, E, W],
    v *Vertex[V, E, W],
) []*Vertex[V, E, W] {
    if result == nil { result = []*Vertex[V, E, W]{} }
    result = growCap(result, 0, len(bfs.vertexes))

    s := bfs.start
    u := bfs.matchingVertex(v)

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

// BreadthFirstSearchWeightedGeneral performs a complete search of the
// reachable graph from a given start vertex, calculating the shortest path by
// weight (distance). The resulting search tree gives useful properties. It
// stores this in the provided result object, resizes the underlying buffer if
// necessary, and returns that result object (or, if nil, creates and returns a
// new result object).
//
// This search is performed in the "general case" where edges may have negative
// weights, but not negative-weight cycles. When the boolean return value is
// false, the search has detected negative-weight cycles and cannot proceed.
//
// If you are certain that the graph does not contain edges with negative
// weights, then the alternative method [Digraph.WeightedSearch] is likely
// to be faster.
//
// If you want a specific ordering of vertexes visited by the search, use
// [Digraph.SortRoots] and [Digraph.SortEdges] first.
func (d *Digraph[V, E, W]) BreadthFirstSearchWeightedGeneral(
    bfs *BFSResult[V, E, W],
    start *Vertex[V, E, W],
) (*BFSResult[V, E, W], bool) {

    // Bellman-Ford algorithm

    if bfs == nil {
        bfs = &BFSResult[V, E, W]{}
    }
    z := len(d.Vertexes)
    bfs.vertexes = growCap(bfs.vertexes, z, z)
    clear(bfs.vertexes)

    bfs.start = bfs.matchingVertex(start)

    for i := 0; i < len(d.Vertexes); i++ {
        v := bfs.vertexByID(i)
        v.vertex = d.Vertexes[i]
        v.predecessor = nil
        v.distance = positiveInfiniteEdgeCount
        v.weightedDistance = d.infiniteWeightedDistance(1)
    }

    bfs.start.distance = 0
    bfs.start.weightedDistance = 0

    inf := infinities[W]{
        positive: d.infiniteWeightedDistance(+1),
        negative: d.infiniteWeightedDistance(-1),
    }

    for n := 0; n < len(d.Vertexes) - 1; n++ {
        for i := 0; i < len(d.Vertexes); i++ {
            u := bfs.vertexByID(i)
            edges := u.vertex.Edges
            for j := 0; j < len(edges); j++ {
                v := bfs.matchingVertex(edges[j].Target)
                w := edges[j].Weight
                bfs.relax(u, v, w, inf)
            }
        }
    }

    // check for negative cycles
    for i := 0; i < len(d.Vertexes); i++ {
        u := bfs.vertexByID(i)
        edges := u.vertex.Edges
        for j := 0; j < len(edges); j++ {
            v := bfs.matchingVertex(edges[j].Target)
            w := edges[j].Weight
            sum := inf.sum(u.weightedDistance, w)
            if v.weightedDistance > sum {
                return bfs, false
            }
        }
    }

    return bfs, true
}

func (bfs *BFSResult[V, E, W]) relax (
    u *vertexBFS[V, E, W],
    v *vertexBFS[V, E, W],
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
