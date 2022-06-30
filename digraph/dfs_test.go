package digraph

import (
    "testing"

    "github.com/stretchr/testify/assert"
)

func TestDigraph_depthFirstSearchRecursive(t *testing.T) {
    type vertex = Vertex[string, EdgeDontCare, WeightDontCare]

    d := &Digraph[string, EdgeDontCare, WeightDontCare]{}

    // constructs the graph given in CLRS "Introduction To Algorithms" p.605
    // and preforms a depth-first-search

    u := d.AddVertex("u") // vertex 0
    v := d.AddVertex("v") // vertex 1
    w := d.AddVertex("w") // vertex 2
    x := d.AddVertex("x") // vertex 3
    y := d.AddVertex("y") // vertex 4
    z := d.AddVertex("z") // vertex 5

    d.AddEdge(u, v, nil)
    d.AddEdge(u, x, nil)
    d.AddEdge(x, v, nil)
    d.AddEdge(y, x, nil)
    d.AddEdge(v, y, nil)
    d.AddEdge(w, y, nil)
    d.AddEdge(w, z, nil)
    d.AddEdge(z, z, nil)

    result := d.depthFirstSearchRecursive(nil)

    type row struct{ vertex *vertex; d, f int }
    times := []row{
        {u,  1,  8},
        {v,  2,  7},
        {w,  9, 12},
        {x,  4,  5},
        {y,  3,  6},
        {z, 10, 11},
    }

    for k, v := range times {
        assert.Equal(t, searchColorFinished, result.vertexes[k].color, "colour for vertex %s", v.vertex.Value)
        assert.Equal(t, v.d, result.vertexes[k].d, "discovery time for vertex %s", v.vertex.Value)
        assert.Equal(t, v.f, result.vertexes[k].f, "finishing time for vertex %s", v.vertex.Value)
    }
}

func TestDFS_edgeType(t *testing.T) {
    // example based on https://courses.cs.duke.edu/fall17/compsci330/lecture12note.pdf

    g := &Digraph[string, string, WeightDontCare]{}

    a := g.AddVertex("a")
    b := g.AddVertex("b")
    c := g.AddVertex("c")
    d := g.AddVertex("d")
    e := g.AddVertex("e")

    g.AddEdge(a, b, "ab")
    g.AddEdge(a, c, "ac")
    g.AddEdge(b, e, "be")
    g.AddEdge(c, c, "cc") // self-loop
    g.AddEdge(d, c, "dc")
    g.AddEdge(e, a, "ea")
    g.AddEdge(e, c, "ec")
    g.AddEdge(e, d, "ed")

    dfs := g.depthFirstSearchRecursive(nil)

    // note: this test assumes a particular DFS transversal order where edges
    // are searched in order they are defined. This doesn't hold true for
    // all possible DFS transversal orders!!! This test is therefore brittle

    assert.Equal(t, searchEdgeTypeTree,    dfs.edgeType(a, b)) // a => b
    assert.Equal(t, searchEdgeTypeTree,    dfs.edgeType(b, e)) // b => e
    assert.Equal(t, searchEdgeTypeBack,    dfs.edgeType(c, c)) // c => c
    assert.Equal(t, searchEdgeTypeTree,    dfs.edgeType(e, c)) // e => c
    assert.Equal(t, searchEdgeTypeForward, dfs.edgeType(a, c)) // a => c
    assert.Equal(t, searchEdgeTypeCross,   dfs.edgeType(d, c)) // d => c
    assert.Equal(t, searchEdgeTypeBack,    dfs.edgeType(e, a)) // e => a
    assert.True(t, g.ContainsCycles(dfs))
}

func TestDigraph_ContainsCycles(t *testing.T) {
    // empty graph
    {
        g := &Digraph[string, string, WeightDontCare]{}
        dfs := g.depthFirstSearchRecursive(nil)
        assert.False(t, g.ContainsCycles(dfs))
    }

    // edgeless graph
    {
        g := &Digraph[string, string, WeightDontCare]{}

        g.AddVertex("a") // vertex 0
        g.AddVertex("b") // vertex 1
        g.AddVertex("c") // vertex 2

        dfs := g.depthFirstSearchRecursive(nil)

        assert.False(t, g.ContainsCycles(dfs))
    }

    // simple graph, no cycles
    {
        g := &Digraph[string, string, WeightDontCare]{}

        a := g.AddVertex("a") // vertex 0
        b := g.AddVertex("b") // vertex 1
        c := g.AddVertex("c") // vertex 2

        g.AddEdge(a, b, "ab")
        g.AddEdge(b, c, "bc")

        dfs := g.depthFirstSearchRecursive(nil)

        assert.False(t, g.ContainsCycles(dfs))
    }

    // simple graph, cycles
    {
        g := &Digraph[string, string, WeightDontCare]{}

        a := g.AddVertex("a") // vertex 0
        b := g.AddVertex("b") // vertex 1
        c := g.AddVertex("c") // vertex 2

        g.AddEdge(a, b, "ab")
        g.AddEdge(b, c, "bc")
        g.AddEdge(c, a, "ca")

        dfs := g.depthFirstSearchRecursive(nil)

        assert.True(t, g.ContainsCycles(dfs))
    }

    // multigraph, no cycles
    {
        g := &Digraph[string, string, WeightDontCare]{}

        a := g.AddVertex("a") // vertex 0
        b := g.AddVertex("b") // vertex 1
        c := g.AddVertex("c") // vertex 2

        g.AddEdge(a, b, "ab1")
        g.AddEdge(a, b, "ab2")
        g.AddEdge(b, c, "bc1")
        g.AddEdge(b, c, "bc2")

        dfs := g.depthFirstSearchRecursive(nil)

        assert.False(t, g.ContainsCycles(dfs))
    }

    // multigraph, cycles
    {
        g := &Digraph[string, string, WeightDontCare]{}

        a := g.AddVertex("a") // vertex 0
        b := g.AddVertex("b") // vertex 1
        c := g.AddVertex("c") // vertex 2

        g.AddEdge(a, b, "ab1")
        g.AddEdge(a, b, "ab2")
        g.AddEdge(b, c, "bc1")
        g.AddEdge(b, c, "bc2")
        g.AddEdge(c, a, "ac1")
        g.AddEdge(c, a, "ac2")

        dfs := g.depthFirstSearchRecursive(nil)

        assert.True(t, g.ContainsCycles(dfs))
    }

    // simple disconnected graph, no cycles
    {
        g := &Digraph[string, string, WeightDontCare]{}

        a := g.AddVertex("a")
        b := g.AddVertex("b")
        c := g.AddVertex("c")

        x := g.AddVertex("x")
        y := g.AddVertex("y")
        z := g.AddVertex("z")

        g.AddEdge(a, b, "ab")
        g.AddEdge(b, c, "bc")

        g.AddEdge(x, y, "xy")
        g.AddEdge(y, z, "yz")

        dfs := g.depthFirstSearchRecursive(nil)

        assert.False(t, g.ContainsCycles(dfs))
    }

    // simple disconnected graph, cycles
    {
        g := &Digraph[string, string, WeightDontCare]{}

        a := g.AddVertex("a")
        b := g.AddVertex("b")
        c := g.AddVertex("c")

        x := g.AddVertex("x")
        y := g.AddVertex("y")
        z := g.AddVertex("z")

        g.AddEdge(a, b, "ab")
        g.AddEdge(b, c, "bc")

        g.AddEdge(x, y, "xy")
        g.AddEdge(y, z, "yz")
        g.AddEdge(z, x, "zx")

        dfs := g.depthFirstSearchRecursive(nil)

        assert.True(t, g.ContainsCycles(dfs))
    }
}


func TestDigraph_TopologicalSort(t *testing.T) {
    // Example from CLRS "Introduction to Algorithms", 3rd ed. page 613
    g := &Digraph[string, EdgeDontCare, WeightDontCare]{}

    shirt       := g.AddVertex("shirt")
    watch       := g.AddVertex("watch")
    undershorts := g.AddVertex("undershorts")
    tie         := g.AddVertex("tie")
    jacket      := g.AddVertex("jacket")
    belt        := g.AddVertex("belt")
    pants       := g.AddVertex("pants")
    shoes       := g.AddVertex("shoes")
    socks       := g.AddVertex("socks")

    g.AddEdge(undershorts, pants,  nil)
    g.AddEdge(undershorts, shoes,  nil)
    g.AddEdge(pants,       shoes,  nil)
    g.AddEdge(socks,       shoes,  nil)
    g.AddEdge(pants,       belt,   nil)
    g.AddEdge(shirt,       tie,    nil)
    g.AddEdge(shirt,       belt,   nil)
    g.AddEdge(tie,         jacket, nil)
    g.AddEdge(belt,        jacket, nil)

    dfs := g.DepthFirstSearch(nil)
    assert.False(t, g.ContainsCycles(dfs))
    topologicalOrdering := dfs.TopologicalSort(nil)

    // there are multiple valid topological sorts. This expectedOrdering is
    // exact only where vertexes and edges are visited in the order they are
    // first added.
    expectedOrdering := []*Vertex[string, EdgeDontCare, WeightDontCare]{
        socks, undershorts, pants, shoes, watch, shirt, belt, tie, jacket,
    }

    assert.Equal(t, expectedOrdering, topologicalOrdering)
}
