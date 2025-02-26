package builder

import (
	"errors"
	"fmt"
)

// ErrInvalidParameters is parameters error.
var ErrInvalidParameters = errors.New("invalid parameters passed to function")

func min(a1 int, a2 int) int {
	if a1 <= a2 {
		return a1
	}
	return a2
}

// StronglyConnectedComponents compute strongly сonnected сomponents of a graph.
// Tarjan's strongly connected components algorithm.
func StronglyConnectedComponents(
	vertices []string, edges map[string]map[string]struct{},
) []map[string]struct{} {
	// Tarjan's strongly connected components algorithm
	var (
		identified = map[string]struct{}{}
		stack      = []string{}
		index      = map[string]int{}
		lowlink    = map[string]int{}
		dfs        func(v string) []map[string]struct{}
	)

	dfs = func(vertex string) []map[string]struct{} {
		index[vertex] = len(stack)
		stack = append(stack, vertex)
		lowlink[vertex] = index[vertex]

		sccs := []map[string]struct{}{}
		for w := range edges[vertex] {
			if _, ok := index[w]; !ok {
				sccs = append(sccs, dfs(w)...)
				lowlink[vertex] = min(lowlink[vertex], lowlink[w])
			} else if _, ok := identified[w]; !ok {
				lowlink[vertex] = min(lowlink[vertex], lowlink[w])
			}
		}

		if lowlink[vertex] == index[vertex] {
			scc := map[string]struct{}{}
			for _, v := range stack[index[vertex]:] {
				scc[v] = struct{}{}
			}
			stack = stack[:index[vertex]]
			for v := range scc {
				identified[v] = struct{}{}
			}
			sccs = append(sccs, scc)
		}
		return sccs
	}

	sccs := []map[string]struct{}{}
	for _, v := range vertices {
		if _, ok := index[v]; !ok {
			sccs = append(sccs, dfs(v)...)
		}
	}
	return sccs
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func reduceGraph(
	graph map[string]map[string]struct{}, scc map[string]struct{},
) map[string]map[string]struct{} {
	reduceGraph := map[string]map[string]struct{}{}
	for src, dsts := range graph {
		if _, ok := scc[src]; !ok {
			continue
		}
		reduceGraph[src] = map[string]struct{}{}
		for dst := range dsts {
			if _, ok := scc[dst]; !ok {
				continue
			}
			reduceGraph[src][dst] = struct{}{}
		}
	}
	return reduceGraph
}

// FindCyclesInSCC find cycles in SCC emanating from start.
// Yields lists of the form ['A', 'B', 'C', 'A'], which means there's
// a path from A -> B -> C -> A.  The first item is always the start
// argument, but the last item may be another element, e.g.  ['A',
// 'B', 'C', 'B'] means there's a path from A to B and there's a
// cycle from B to C and back.
func FindCyclesInSCC(
	graph map[string]map[string]struct{}, scc map[string]struct{}, start string,
) ([][]string, error) {
	// Basic input checks.
	if _, ok := scc[start]; !ok {
		return nil, fmt.Errorf(
			"%w: scc %v does not contain %q", ErrInvalidParameters, scc, start)
	}
	extravertices := []string{}
	for k := range scc {
		if _, ok := graph[k]; !ok {
			extravertices = append(extravertices, k)
		}
	}
	if len(extravertices) != 0 {
		return nil, fmt.Errorf(
			"%w: graph does not contain scc. %v",
			ErrInvalidParameters, extravertices)
	}

	// Reduce the graph to nodes in the SCC.
	graph = reduceGraph(graph, scc)
	if _, ok := graph[start]; !ok {
		return nil, fmt.Errorf(
			"%w: graph %v does not contain %q",
			ErrInvalidParameters, graph, start)
	}

	// Recursive helper that yields cycles.
	var dfs func(node string, path []string) [][]string
	dfs = func(node string, path []string) [][]string {
		ret := [][]string{}
		if contains(path, node) {
			t := make([]string, 0, len(path)+1)
			t = append(t, path...)
			t = append(t, node)
			ret = append(ret, t)
			return ret
		}
		path = append(path, node) // TODO: Make this not quadratic.
		for child := range graph[node] {
			ret = append(ret, dfs(child, path)...)
		}
		return ret
	}

	return dfs(start, []string{}), nil
}
