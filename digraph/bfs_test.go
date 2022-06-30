package digraph

import (
    "fmt"
    "strings"
    "testing"

    "github.com/stretchr/testify/assert"
)

func TestDigraph_BreadthFirstSearch(t *testing.T) {
    type digraph = Digraph[string, EdgeDontCare, WeightDontCare]
    type vertex  = Vertex [string, EdgeDontCare, WeightDontCare]

    g := &digraph{} // a multidigraph, with cycles and loops

    a := g.AddVertex("a") // vertex 0
    b := g.AddVertex("b") // vertex 1
    c := g.AddVertex("c") // vertex 2
    d := g.AddVertex("d") // vertex 3
    e := g.AddVertex("e") // vertex 4

    x := g.AddVertex("x") // vertex 5
    y := g.AddVertex("y") // vertex 6
    z := g.AddVertex("z") // vertex 7

    g.AddEdge(a, b, nil)
    g.AddEdge(b, c, nil)
    g.AddEdge(b, c, nil) // extra edge
    g.AddEdge(b, e, nil)
    g.AddEdge(c, c, nil) // loop
    g.AddEdge(c, d, nil)
    g.AddEdge(d, e, nil)
    g.AddEdge(e, b, nil) // cycle

    g.AddEdge(x, y, nil) // disconnected (unreachable from abcde)
    g.AddEdge(y, z, nil)
    g.AddEdge(z, x, nil)

    // BFS from start vertex b
    bfs := g.BreadthFirstSearch(nil, b, 0, nil)

    {
        type row struct{
            predecessor *vertex
            distance DistanceT
        }

        const inf = positiveInfiniteEdgeCount
        expected := []row{
            {nil, inf}, // a, unreachable
            {nil,   0}, // b, from start
            {  b,   1}, // c, from b->c
            {  c,   2}, // d, from c->d
            {  b,   1}, // e, from b->e
            {nil, inf}, // x, unreachable
            {nil, inf}, // y, unreachable
            {nil, inf}, // z, unreachable
        }

        for i := 0; i < 8; i++ {
            assert.Equal(t, expected[i].predecessor,  bfs.Predecessor(g.Vertexes[i]), "vertex %d predecessor", i)
            assert.Equal(t, expected[i].distance,     bfs.Distance(   g.Vertexes[i]), "vertex %d distance",    i)
        }
    }

    {
        type row struct{
            name string
            u *vertex
            shortestPath []*vertex
        }

        expected := []row{
            {"bz",  z, []*vertex{}},
            {"b",   b, []*vertex{b}}, // reachable from itself with zero edges
            {"bcd", d, []*vertex{b, c, d}},
            {"be",  e, []*vertex{b, e}},
        }

        for i := 0; i < len(expected); i++ {
            name, u := expected[i].name, expected[i].u
            path := expected[i].shortestPath
            shortest := bfs.ShortestPath(nil, u)

            if len(shortest) == len(path) {
                assert.Equal(t, path, shortest, "shortest path %q (want %s, got %s)",
                    name, sprintPath(path), sprintPath(shortest))
            } else {
                assert.Equal(t, len(path), len(shortest), "shortest path %q length", name)
            }
        }
    }
}

func sprintPath[V any, E any, W Number](v []*Vertex[V, E, W]) string {
    var b strings.Builder
    for i := 0; i < len(v); i++ {
        b.WriteString(fmt.Sprintf("%v", v[i].Value))
    }
    return b.String()
}


func TestDigraph_WeightedSearchGeneral(t *testing.T) {

}
