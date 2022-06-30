package digraph

import (
    "testing"

    "github.com/stretchr/testify/assert"
)

func TestDigraph_DAGSearchWeighted(t *testing.T) {
    g := &Digraph[string, EdgeDontCare, int]{}
    type vertex = Vertex[string, EdgeDontCare, int]

    // CLRS "Introduction to Algorithms", 3rd ed. p. 656

    // jumble the order for topological sort
    x := g.AddVertex("x")
    c := g.AddVertex("c")
    a := g.AddVertex("a")
    z := g.AddVertex("z")
    b := g.AddVertex("b")
    y := g.AddVertex("y")

    g.AddWeightedEdge(a, b, nil,  5)
    g.AddWeightedEdge(a, c, nil,  3)
    g.AddWeightedEdge(b, c, nil,  2)
    g.AddWeightedEdge(b, x, nil,  6)
    g.AddWeightedEdge(c, x, nil,  7)
    g.AddWeightedEdge(c, y, nil,  4)
    g.AddWeightedEdge(c, z, nil,  2)
    g.AddWeightedEdge(x, y, nil, -1)
    g.AddWeightedEdge(y, z, nil, -2)

    dfs := g.DepthFirstSearch(nil)
    assert.False(t, g.ContainsCycles(dfs))

    topological := dfs.TopologicalSort(nil)
    assert.Equal(t, []*vertex{a, b, c, x, y, z}, topological)

    dags := g.DAGSearchWeighted(nil, b, topological)

    // shortest path from b to z
    assert.Equal(t, []*vertex{b, x, y, z}, dags.ShortestPath(nil, z))

    // unreachable
    w := g.AddVertex("w")
    dfs = g.DepthFirstSearch(dfs)
    topological = dfs.TopologicalSort(topological)
    dags = g.DAGSearchWeighted(dags, a, topological)
    assert.Equal(t, []*vertex{}, dags.ShortestPath(nil, w))

}
