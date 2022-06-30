package digraph

import (
    "strings"
    "testing"

    "github.com/stretchr/testify/assert"
)

func TestExample_SolarSystem(t *testing.T) {
    type digraph = Digraph[string, string, WeightDontCare]
    type vertex  = Vertex [string, string, WeightDontCare]
    type edge    = Edge   [string, string, WeightDontCare]

    d := &digraph{}

    sun   := d.AddVertex("sun")
    earth := d.AddVertex("earth")
    moon  := d.AddVertex("moon")

    assert.True(t, d.AddUniqueEdge(moon,  earth, "orbits"))
    assert.True(t, d.AddUniqueEdge(earth, sun,   "orbits"))
    assert.True(t, d.AddUniqueEdge(sun,   earth, "heats"))

    // duplicates not allowed
    assert.False(t, d.AddUniqueEdge(sun,   earth, "heats"))

    // matrix for efficient lookup
    filter := func(s *string) bool {
        return *s == "orbits"
    }
    orbitMatrix := d.AdjacencyMatrixFiltered(nil, filter)

    // returns true iff a orbits b
    orbits := func(a *vertex, b *vertex) bool {
        return d.Adjacency(orbitMatrix, a, b) > 0
    }

    assert.True (t, orbits(earth, sun))
    assert.True (t, orbits(moon,  earth))
    assert.False(t, orbits(sun,   earth))
}

func TestExample_TransportRoutes(t *testing.T) {
    type minutes = int
    type digraph = Digraph[string, string, minutes]
    type vertex  = Vertex [string, string, minutes]

    d := &digraph{} // a multidigraph

    swansea := d.AddVertex("swansea")
    neath   := d.AddVertex("neath")
    cardiff := d.AddVertex("cardiff")
    bristol := d.AddVertex("bristol")
    london  := d.AddVertex("london")

    /*
                       /---------(very slow)---------->\
                      |          /<-------------------->\
        swansea <-> neath >-> cardiff >-> bristol >-> london
            \<------------------>/
     */

    d.AddWeightedEdge(swansea, neath,   "bus",   20) // minutes
    d.AddWeightedEdge(swansea, neath,   "rail",  11)
    d.AddWeightedEdge(swansea, cardiff, "rail",  35) // direct line
    d.AddWeightedEdge(neath,   swansea, "bus",   25) // return
    d.AddWeightedEdge(neath,   cardiff, "bus",   52)
    d.AddWeightedEdge(neath,   cardiff, "bus",   60) // a second, slower bus
    d.AddWeightedEdge(neath,   cardiff, "rail",  40)
    d.AddWeightedEdge(neath,   london,  "bus",  360) // VERY slow bus
    d.AddWeightedEdge(cardiff, swansea,  "bus",  50) // return
    d.AddWeightedEdge(cardiff, bristol,  "bus",  62)
    //               (cardiff, bristol,  "rail", 50) // cancelled!
    d.AddWeightedEdge(cardiff, london,  "bus",  182)
    //               (cardiff, london,  "rail", ...) // no direct route!
    d.AddWeightedEdge(bristol, london,  "bus",  147)
    d.AddWeightedEdge(bristol, london,  "rail",  98)
    d.AddWeightedEdge(london, cardiff,  "rail", 130) // return

    // if we want the scenic route, here's a filter for only bus options
    filterOnlyByBus := func(s *string) bool {
        return *s == "bus"
    }

    // is there a bus route between two points with no stops in between?
    directBusRoute := func(source *vertex, dest *vertex) bool {
        return nil != d.FindEdgeFiltered(source, dest, filterOnlyByBus)
    }

    assert.True (t, directBusRoute(swansea, neath))
    assert.False(t, directBusRoute(neath,   bristol))

    // Calculate adjacent points for efficient lookups
    mat := d.AdjacencyMatrix(nil)

    // You can reuse the matrix buffer when calculating a new matrix (e.g. if
    // you modify the graph and invalidate the existing matrix).
    mat = d.AdjacencyMatrix(mat)

    // some properties given by an adjacency matrix
    assert.False(t, d.IsSimple(mat)) // d is a multigraph / multidigraph
    assert.False(t, d.ContainsLoops(mat)) // no edges from a vertex to itself

    // Number of direct connections to a location
    assert.Equal(t, 2, d.Indegree(mat, neath))
    assert.Equal(t, 5, d.Indegree(mat, cardiff))

    // Number of direct connections away from a location
    assert.Equal(t, 5, d.Outdegree(mat, neath))
    assert.Equal(t, 3, d.Outdegree(mat, cardiff))

    // Number of direct connections from a to b
    assert.Equal(t, 3, d.Adjacency(mat, neath, cardiff))
    assert.Equal(t, 1, d.Adjacency(mat, london, cardiff))

    // As above, but only bus routes.
    matBus := d.AdjacencyMatrixFiltered(mat, filterOnlyByBus)

    // Number of direct connections to a location by bus
    assert.Equal(t, 1, d.Indegree(matBus, neath))
    assert.Equal(t, 2, d.Indegree(matBus, cardiff))

    // Number of direct connections away from a location by bus
    assert.Equal(t, 4, d.Outdegree(matBus, neath))
    assert.Equal(t, 3, d.Outdegree(matBus, cardiff))

    // Number of direct connections from a to b by bus only
    assert.Equal(t, 2, d.Adjacency(matBus, neath, cardiff))
    assert.Equal(t, 0, d.Adjacency(matBus, london, cardiff))

    // Calculate weighted adjacent points for efficient lookups...
    // This is a multigraph, so use a minimum function to calculate the best
    // time between two points. This is a "reducer" function in the functional
    // programming sense.

    // Define it manually...
    shortestTime := &EdgeWeightReducer[minutes]{
        Identity: InfWeight[minutes](1),
        Reduce: func(a minutes, b minutes) minutes {
            if a <= b { return a } else { return b }
        },
    }

    // Or alternatively use the builtin...
    shortestTime = NewEdgeWeightReducerMinimum[minutes]()

    // if we were certain d.IsSimple() is true, then we could omit the last
    // two arguments (use zero and nil instead).
    matWeighted    := d.WeightedAdjacencyMatrix        (nil, shortestTime)
    matBusWeighted := d.WeightedAdjacencyMatrixFiltered(nil, shortestTime, filterOnlyByBus)

    // The shortest distance between two adjacent vertexes by any direct method
    // should take this long:
    assert.Equal(t,  40, d.WeightedAdjacency(matWeighted, neath,  cardiff))
    assert.Equal(t, 130, d.WeightedAdjacency(matWeighted, london, cardiff))

    // The shortest distance between two adjacent vertexes by any direct bus
    // should take this long:
    assert.Equal(t, 52, d.WeightedAdjacency(matBusWeighted, neath,  cardiff))

    // no route by bus!
    assert.True(t, d.IsInfiniteWeightedDistance(d.WeightedAdjacency(matBusWeighted, london, cardiff)))

    // properties given by a depth-first search
    dfs := d.DepthFirstSearch(nil)
    assert.True(t, d.ContainsCycles(dfs))

    // shortest distance over a whole path from swansea
    // The general version allows some types of time travelling.
    bfs, ok := d.BreadthFirstSearchWeightedGeneral(nil, swansea)
    assert.True(t, ok) // no negative-weight cycles

    // quickest path from swansea to london
    assert.Equal(t, []*vertex{swansea, cardiff, bristol, london}, bfs.ShortestPath(nil, london))
    assert.Equal(t, 35 + 62 + 98, bfs.WeightedDistance(london))

    // shortest distance over a whole path from london
    bfs, ok = d.BreadthFirstSearchWeightedGeneral(nil, london)
    assert.True(t, ok) // no negative-weight cycles

    // quickest path from london to bristol
    assert.Equal(t, []*vertex{london, cardiff, bristol}, bfs.ShortestPath(nil, bristol))

    // shortest distance over a whole path from bristol
    bfs, ok = d.BreadthFirstSearchWeightedGeneral(nil, bristol)
    assert.True(t, ok) // no negative-weight cycles

    // quickest path from bristol to neath
    assert.Equal(t, []*vertex{bristol, london, cardiff, swansea, neath}, bfs.ShortestPath(nil, neath))
    assert.Equal(t, DistanceT(0), bfs.Distance(bristol))
    assert.Equal(t, DistanceT(1), bfs.Distance(london))
    assert.Equal(t, DistanceT(2), bfs.Distance(cardiff))
    assert.Equal(t, DistanceT(3), bfs.Distance(swansea))
    assert.Equal(t, DistanceT(4), bfs.Distance(neath))
    assert.Equal(t, minutes(98 + 130 + 50 + 11), bfs.WeightedDistance(neath))
    assert.Equal(t, swansea, bfs.Predecessor(neath))
    assert.Equal(t, bristol, bfs.Predecessor(london))
    assert.Nil(t, nil,     bfs.Predecessor(bristol))

    // Check that the guards for the matrix being current are applied
    llanelli  := d.AddVertex("llanelli") // modify the graph
    assert.Panics(t, func() {
        // mat is no longer current as a representation of the graph
        _ = d.Indegree(mat, llanelli) // should panic
    })

    d.AddWeightedEdge(llanelli, swansea, "rail", 35)

    // shortest distance over a whole path from london
    bfs, ok = d.BreadthFirstSearchWeightedGeneral(nil, london)
    assert.True(t, ok) // no negative-weight cycles
    assert.True(t, d.IsInfiniteWeightedDistance(bfs.WeightedDistance(llanelli))) // unreachable
    assert.True(t, d.IsInfiniteDistance(bfs.Distance(llanelli))) // unreachable
    assert.Nil(t, nil, bfs.Predecessor(llanelli))
}

func TestDigraph_SortRoots(t *testing.T) {
    type vertex  = Vertex [string, EdgeDontCare, WeightDontCare]
    g := &Digraph[string, EdgeDontCare, WeightDontCare]{}

    i := g.AddVertex("i") // root
    j := g.AddVertex("j")
    k := g.AddVertex("k")

    z := g.AddVertex("z")
    y := g.AddVertex("y")
    x := g.AddVertex("x") // root

    l := g.AddVertex("l") // a cycle
    m := g.AddVertex("m")
    n := g.AddVertex("n")

    c := g.AddVertex("c")
    b := g.AddVertex("b")
    a := g.AddVertex("a") // root
    d := g.AddVertex("d")

    g.AddEdge(a, b, nil)
    g.AddEdge(a, c, nil)
    g.AddEdge(c, d, nil)

    g.AddEdge(i, j, nil)
    g.AddEdge(j, k, nil)
    g.AddEdge(k, j, nil)
    g.AddEdge(k, k, nil)

    g.AddEdge(x, y, nil)
    g.AddEdge(y, z, nil)

    g.AddEdge(l, m, nil) // a cycle
    g.AddEdge(m, n, nil)
    g.AddEdge(n, l, nil)

    // existing ID order
    assert.Equal(t, VertexID( 0), i.id)
    assert.Equal(t, VertexID( 5), x.id)
    assert.Equal(t, VertexID(11), a.id)

    g.SortRoots(
        g.AdjacencyMatrix(nil),
        func(u *vertex, v *vertex) bool {
            return -1 == strings.Compare(u.Value, v.Value)
        },
    )

    // roots first, in alphabetical order
    assert.Equal(t, "a", g.Vertexes[0].Value)
    assert.Equal(t, "i", g.Vertexes[1].Value)
    assert.Equal(t, "x", g.Vertexes[2].Value)

    // ID reordering
    assert.Equal(t, VertexID(0), a.id)
    assert.Equal(t, VertexID(1), i.id)
    assert.Equal(t, VertexID(2), x.id)
}

func TestDigraph_SortEdges(t *testing.T) {
    type vertex  = Vertex [string, int, WeightDontCare]
    type edge    = Edge   [string, int, WeightDontCare]
    g := &Digraph[string, int, WeightDontCare]{}

    a := g.AddVertex("a")
    b := g.AddVertex("b")
    c := g.AddVertex("c")

    g.AddEdge(a, c, 4)
    g.AddEdge(a, b, 3)
    g.AddEdge(a, b, 1)
    g.AddEdge(a, b, 2)

    g.AddEdge(b, c, 1)
    g.AddEdge(b, c, 2)
    g.AddEdge(b, b, 0)

    lessThan := func(v *vertex, a *edge, b *edge) bool {
        return a.Value < b.Value
    }

    g.SortEdges(lessThan)

    assert.Equal(t, 1, a.Edges[0].Value)
    assert.Equal(t, 2, a.Edges[1].Value)
    assert.Equal(t, 3, a.Edges[2].Value)
    assert.Equal(t, 4, a.Edges[3].Value)
    assert.Equal(t, c, a.Edges[3].Target)

    assert.Equal(t, 0, b.Edges[0].Value)
    assert.Equal(t, 1, b.Edges[1].Value)
    assert.Equal(t, 2, b.Edges[2].Value)

}
