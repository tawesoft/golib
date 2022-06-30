# TODO

* Use new Path / PathVertex types
* Consider how to implement undirected graphs
  (maybe as a facade over a  digraph?)
  + Orientation to turn a graph into a digraph and a complete graph into a 
    tournament
* Removing vertexes and edges
* Serialisation
  + Digraph.NamedVertexes map and Digraph.AddNamedVertex?
  + (De)serialisation function pointers for vertex and edge values?
* More DAG stuff
  + d.CriticalPathLength(DAGS)
* Tree stuff
* Functions:
  + Digraph.New
  + Digraph.Clear
  + Digraph.Copy
  + Digraph.Transpose (in place) (or Digraph.Reverse?)
  + Digraph.ReduceEdges (for a multigraph)
  + Graph.FromDigraph (remove directions)
  + Graph.ToDigraph / Graph.Orientate (make a Digraph)
  + Digraph.Reorientate? (change directions)
  + <Search>.Walk
  + Digraph.NewFromSubgraph(walker) // uses <Search>.Walk as a walker
* Refactoring:
  + Digraph.VertexByID() return vertex pointer by ID
  + Common search vertex type (which can share the RELAX method) (see also 
    Perf note)
  + A common search interface type? (probably not)
  + Filter variants of searches
* Perf:
  + using IDs instead of pointers in annotated search vertexes would reduce 
    the work the garbage collector needs to do.
* Testing:
  + Tests for Digraph.IsSomething / Digraph.Contains functions
* Features:
  + BreadthFirstSearchWeighted (Dijkstra)
* Nice to haves:
  + Digraph.StronglyConnectedComponents, etc.

## External TODO

* Loader based on Digraph
* Task runner based on Digraph
