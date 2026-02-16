package builder

import (
	"errors"
	"fmt"

	"github.com/mna/pigeon/ast"
)

var (
	// ErrNoLeader is no leader error.
	ErrNoLeader = errors.New(
		"SCC has no leadership candidate (no element is included in all cycles)")
	// ErrHaveLeftRecursion is recursion error.
	ErrHaveLeftRecursion = errors.New("grammar contains left recursion")
)

// PrepareGrammar evaluates parameters associated with left recursion.
func PrepareGrammar(grammar *ast.Grammar) (bool, error) {
	mapRules := make(map[string]*ast.Rule, len(grammar.Rules))
	for _, rule := range grammar.Rules {
		mapRules[rule.Name.Val] = rule
	}
	ComputeNullables(mapRules)
	haveLeftRecursion, err := ComputeLeftRecursives(mapRules)
	if err != nil {
		return false, fmt.Errorf("error compute left recursive: %w", err)
	}
	return haveLeftRecursion, nil
}

// ComputeNullables evaluates nullable nodes.
func ComputeNullables(rules map[string]*ast.Rule) {
	// Compute which rules in a grammar are nullable
	for _, rule := range rules {
		rule.NullableVisit(rules)
	}
}

func findLeader(
	graph map[string]map[string]struct{}, scc map[string]struct{},
) (string, error) {
	// Try to find a leader such that all cycles go through it.
	leaders := make(map[string]struct{}, len(scc))
	for k := range scc {
		leaders[k] = struct{}{}
	}
	for start := range scc {
		cycles, err := FindCyclesInSCC(graph, scc, start)
		if err != nil {
			return "", fmt.Errorf("error find cycles: %w", err)
		}
		for _, cycle := range cycles {
			mapCycle := make(map[string]struct{}, len(cycle))
			for _, k := range cycle {
				mapCycle[k] = struct{}{}
			}
			for k := range scc {
				if _, okCycle := mapCycle[k]; !okCycle {
					delete(leaders, k)
				}
			}
			if len(leaders) == 0 {
				return "", ErrNoLeader
			}
		}
	}
	// Pick an arbitrary leader from the candidates.
	var leader string
	for k := range leaders {
		leader = k // The only element.
		break
	}
	return leader, nil
}

// ComputeLeftRecursives evaluates left recursion.
func ComputeLeftRecursives(rules map[string]*ast.Rule) (bool, error) {
	graph := MakeFirstGraph(rules)
	vertices := make([]string, 0, len(graph))
	haveLeftRecursion := false
	for k := range graph {
		vertices = append(vertices, k)
	}
	sccs := StronglyConnectedComponents(vertices, graph)
	for _, scc := range sccs {
		if len(scc) > 1 {
			for name := range scc {
				rules[name].LeftRecursive = true
				haveLeftRecursion = true
			}
			leader, err := findLeader(graph, scc)
			if err != nil {
				return false, fmt.Errorf("error find leader %v: %w", scc, err)
			}
			rules[leader].Leader = true
		} else {
			var name string
			for k := range scc {
				name = k // The only element.
				break
			}
			if _, ok := graph[name][name]; ok {
				rules[name].LeftRecursive = true
				rules[name].Leader = true
				haveLeftRecursion = true
			}
		}
	}
	return haveLeftRecursion, nil
}

// MakeFirstGraph compute the graph of left-invocations.
// There's an edge from A to B if A may invoke B at its initial position.
// Note that this requires the nullable flags to have been computed.
func MakeFirstGraph(rules map[string]*ast.Rule) map[string]map[string]struct{} {
	graph := make(map[string]map[string]struct{})
	vertices := make(map[string]struct{})
	for rulename, rule := range rules {
		names := rule.InitialNames()
		graph[rulename] = names
		for name := range names {
			vertices[name] = struct{}{}
		}
	}
	for vertex := range vertices {
		if _, ok := graph[vertex]; !ok {
			graph[vertex] = make(map[string]struct{})
		}
	}
	return graph
}
