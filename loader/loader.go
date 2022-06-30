// SPDX-License-Identifier: MIT
// x-doc-short-desc: concurrent dependency graph solver
// x-doc-stable: candidate
// x-doc-copyright: 2022 Ben Golightly <ben@tawesoft.co.uk>
// x-doc-copyright: 2022 Tawesoft Ltd <open-source@tawesoft.co.uk>
// x-doc-copyright: 2022 CONTRIBUTORS

// Package loader implements the ability to define a graph of tasks and
// dependencies that produce results, and solve the graph incrementally or
// totally, including concurrently.
//
// For example, this could be used to implement a loading screen for a computer
// game with a progress bar that updates in real time, with images being
// decoded concurrently with files being loaded from disk, and synchronised
// with the main thread for safe OpenGL operations such as creating texture
// objects on the GPU.
package loader

import (
    "fmt"

    "github.com/tawesoft/golib/digraph"
)

type Result struct {
    Value any
    Err error
}

type Inputs struct {
    Named map[string]Result
    Direct []Result
}

// GetNamed gets a named input from the input set. It panics if an input
// with that name does not exist
func (n *Inputs) GetNamed(name string) Result {
    if r, ok := n.Named[name]; ok {
        return r
    } else {
        panic(fmt.Errorf("missing %s (required named input)", name))
    }
}

// AssertNone panics unless the inputs are empty
func (n *Inputs) AssertNone() {
    if (len(n.Direct) > 0) || (len(n.Named) > 0) {
        panic("unexpected inputs (expected none)")
    }
}

// AssertNoneDirect panics unless the direct inputs are empty
func (n *Inputs) AssertNoneDirect() {
    if len(n.Direct) > 0 {
        panic("unexpected directed inputs (expected none)")
    }
}

// AssertNoneNamed panics unless the named inputs are empty
func (n *Inputs) AssertNoneNamed() {
    if len(n.Named) > 0 {
        panic("unexpected named inputs (expected none)")
    }
}


// AssertNoErrors panics unless all the inputs (if any) have nil error values.
func (n *Inputs) AssertNoErrors() {
    for i, v := range n.Direct {
        err := v.Err
        if err != nil {
            panic(fmt.Errorf("unexpected input error for direct input %d: %w", i, err))
        }
    }

    for k, v := range n.Named {
        err := v.Err
        if err != nil {
            panic(fmt.Errorf("unexpected input error for named input %q: %w", k, err))
        }
    }
}

type taskLifetime int
const (
    taskLifetimeTemporary taskLifetime = 0
    taskLifetimeKeep      taskLifetime = 1
    taskLifetimeRequired  taskLifetime = 2
)

type taskID string
type Task struct {
    lifetime taskLifetime  // if required, an error stops the entire load process

    // Description is a present-tense sentence fragment describing what the
    // task is doing e.g. "reading file" or "decoding image".
    Description string

    Load func(inputs Inputs) (any, error)
    Free func(input Result) error

    // TODO there are mutable, so move them out
    childrenCompleted int
    parentsCompleted  int
}

func (t *Task) load(inputs Inputs) Result {
    guarded := func(f func(inputs Inputs) (any, error)) (value any, err error) {
        defer func() {
            if r := recover(); r != nil {
                if rAsErr, ok := r.(error); ok {
                    err = fmt.Errorf("panic: %w", rAsErr)
                } else {
                    err = fmt.Errorf("panic: %s", r)
                }
            }
        }()
        return f(inputs)
    }

    value, err := guarded(t.Load)

    if err != nil {
        if len(t.Description) > 0 {
            err = fmt.Errorf("error %s: %w", t.Description, err)
        } else {
            err = fmt.Errorf("error in task: %w", err)
        }
    }

    return Result{value, err}
}

// Vertex represents a task in the loader dependency graph.
type Vertex struct {
    vertex *digraph.Vertex[Task, string, digraph.WeightDontCare]
}

// Loader is a dependency graph (represented as a DAG) modelling tasks that
// perform a unit of work, either alone or on the results of some previous
// tasks that must complete first.
type Loader struct {
    graph *digraph.Digraph[Task, string, digraph.WeightDontCare]
    results []Result // by vertex ID
}

// New creates a new [Loader].
func New() *Loader {
    return &Loader{
        graph: digraph.New[Task, string, digraph.WeightDontCare](),
    }
}

// Add creates a new vertex in a loader dependency graph. Inputs are provided
// as arguments to the [Task.Load] method.
func (l *Loader) Add(t Task, namedInputs map[string]Vertex, directInputs []Vertex) Vertex {
    u := Vertex{vertex: l.graph.AddVertex(t)}
    i := 0

    for k, v := range namedInputs {
        l.graph.AddEdge(v.vertex, u.vertex, k)
        i++
    }

    for _, v := range directInputs {
        l.graph.AddEdge(v.vertex, u.vertex, "")
        i++
    }

    return u
}

func (v Vertex) Keep() {
    t := &v.vertex.Value
    t.lifetime = taskLifetimeKeep
}

func (v Vertex) Require() {
    t := &v.vertex.Value
    t.lifetime = taskLifetimeRequired
}

type Manager interface {
    Queue(WorkerTask)
    Start(chan WorkerResult)
    Stop()
}

type ManagerAsync struct {
    MaxWorkers int
    jobChan chan WorkerTask
}

type ManagerSync struct {
    resultChan chan <- WorkerResult
}

func worker(jobs <-chan WorkerTask, results chan <- WorkerResult) {
    for {
        job, ok := <- jobs
        if !ok { return }

        results <- WorkerResult{
            ID:     job.ID,
            Result: job.Work(job.Inputs),
        }
    }
}

func (m *ManagerAsync) Queue(t WorkerTask) {
    m.jobChan <- t
}

func (m *ManagerAsync) Start(resultChan chan <- WorkerResult, bufSize int) {
    m.jobChan = make(chan WorkerTask, bufSize)

    for i := 0; i < 3; i++ {
        go worker(m.jobChan, resultChan)
    }
}

func (m *ManagerAsync) Stop() {
    close(m.jobChan)
}

type WorkerTask struct {
    ID int
    Inputs Inputs
    Work func(Inputs) Result
}

type WorkerResult struct {
    ID int
    Result Result
}


func (m *ManagerSync) Queue(t WorkerTask) {
    m.resultChan <- WorkerResult{
        ID:     t.ID,
        Result: t.Work(t.Inputs),
    }
}

func (m *ManagerSync) Start(resultChan chan <- WorkerResult, bufSize int) {
    m.resultChan = resultChan
}

func (m *ManagerSync) Stop() {}



// Walk performs each task in the dependency graph, where possible using
// concurrency. It blocks until complete.
func (l *Loader) Walk() {

    mat := l.graph.AdjacencyMatrix(nil)
    roots := l.graph.Roots(mat, nil)
    l.results = make([]Result, len(l.graph.Vertexes), len(l.graph.Vertexes))

    // active is a set of all vertexes that are not blocked from loading
    active := make(map[digraph.VertexID]any)

    // buffered result channel, returning a vertex ID and result
    resultChan := make(chan WorkerResult, len(l.graph.Vertexes))

    //mgr := ManagerAsync{MaxWorkers: 3}
    mgr := ManagerSync{}
    mgr.Start(resultChan, len(l.graph.Vertexes))

    // initialise queue
    for i := 0; i < len(roots); i++ {
        v := roots[i]
        active[v.ID()] = nil
        task := &v.Value
        inputs := makeInputs(l.graph.Inputs(mat, nil, v), l.results)
        mgr.Queue(WorkerTask{
            ID:     int(v.ID()),
            Inputs: inputs,
            Work:   task.load,
        })
    }

    for len(active) > 0 {
        // get a result
        result := <-resultChan

        // remove from active vertexes
        delete(active, digraph.VertexID(result.ID))

        // store the result
        l.results[result.ID] = result.Result

        // lookup in graph as vertex u
        u := l.graph.Vertexes[result.ID]

        // for each parent of u
        parents := l.graph.Inputs(mat, nil, u) // TODO reuse buffer
        for i := 0; i < len(parents); i++ {
            v := parents[i].Vertex
            v.Value.childrenCompleted++
            if v.Value.childrenCompleted == l.graph.Outdegree(mat, v) {
                if v.Value.lifetime == taskLifetimeTemporary {
                    fmt.Printf("can free %q task result\n", v.Value.Description)
                }
            }
        }

        // for each child of u
        edges := u.Edges
        for i := 0; i < len(edges); i++ {
            v := edges[i].Target
            v.Value.parentsCompleted++
            if v.Value.parentsCompleted == l.graph.Indegree(mat, v) {
                inputs := makeInputs(l.graph.Inputs(mat, nil, v), l.results)
                active[v.ID()] = nil
                mgr.Queue(WorkerTask{
                    ID:     int(v.ID()),
                    Inputs: inputs,
                    Work:   v.Value.load,
                })
            }
        }
    }

    mgr.Stop()

    // process queue
    /*
    for {
        complete u <- read from channel
            remove u from queue
            for each of u's parent p:
                increase childrenCompleted[p]
                if childrenCompleted[p] == outdegree(p):
                    if lifetime is Temporary:
                        free o parent
            for each edge u -> v:
                increase parentsCompleted[v]
                if parentsCompleted[v] == indegree[v]:
                    add v to queue
                    send v to loader with input results
    }

    TODO what about freeing???
     */

}

// Step performs (or continues to perform) each task in the dependency graph,
// where possible using concurrency. It does not block, and must be repeatedly
// called until done, when it returns false. If busy-waiting, it is a good
// idea to have some sort of small delay between each iteration.
func (l *Loader) Step() bool {
    // TODO
    return false
}

// makeInputs returns a loader [Inputs] struct from named and unnamed
// ("direct") dependency graph input vertexes.
func makeInputs(
    inputs digraph.Path[Task, string, digraph.WeightDontCare],
    results []Result,
) Inputs {
    directInputs := make([]Result, 0)
    namedInputs  := make(map[string]Result)

    for j := 0; j < len(inputs); j++ {
        n := inputs[j]
        name := n.Via.Value
        if len(name) == 0 {
            directInputs = append(directInputs, results[n.Vertex.ID()])
        } else {
            namedInputs[name] = results[n.Vertex.ID()]
        }
    }

    return Inputs{
        Direct: directInputs,
        Named:  namedInputs,
    }
}

func (l *Loader) Result(v Vertex) Result {
    return l.results[v.vertex.ID()]
}

// WalkSynchronous performs each task in the dependency graph, in topological
// sorted order, one-by-one in sequence until complete. The visit function is a
// callback, called for each task, where i is the progress out of the total
// number of tasks.
func (l *Loader) WalkSynchronous(
    visit func(t Task, i int, total int),
) {
    dfs := l.graph.DepthFirstSearch(nil)
    topological := dfs.TopologicalSort(nil)
    mat := l.graph.AdjacencyMatrix(nil)
    l.results = make([]Result, len(l.graph.Vertexes), len(l.graph.Vertexes))

    for i := 0; i < len(topological); i++ {
        v := topological[i]
        task := &v.Value

        visit(*task, i, len(topological))

        result := task.load(makeInputs(l.graph.Inputs(mat, nil, v), l.results))
        l.results[v.ID()] = result
    }
}
