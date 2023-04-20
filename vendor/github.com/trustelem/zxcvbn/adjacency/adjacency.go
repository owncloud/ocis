package adjacency

type Graph struct {
	Graph         map[string][]string
	Name          string
	AverageDegree float64
}

// Graphs exports the adjacency graphs by name
var Graphs = make(map[string]*Graph)

func init() {
	initGraph("qwerty", adjacencyGraphQwerty)
	initGraph("dvorak", adjacencyGraphDvorak)
	initGraph("keypad", adjacencyGraphKeypad)
	initGraph("mac_keypad", adjacencyGraphMacKeypad)
}

func initGraph(name string, data map[string][]string) {
	g := &Graph{
		Name:          name,
		Graph:         data,
		AverageDegree: calculateAvgDegree(data),
	}
	Graphs[name] = g
}

func calculateAvgDegree(g map[string][]string) float64 {
	var avg float64
	for _, value := range g {
		for _, chars := range value {
			if chars != "" {
				avg += float64(1)
			}
		}

	}

	return avg / float64(len(g))
}
