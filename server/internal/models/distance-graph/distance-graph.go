package distancegraph

import (
	"fmt"

	"github.com/starwander/goraph"
)

type Distancegraph struct {
	*goraph.Graph
	minDistances map[string]map[string]float64
}

type path struct {
	from   string
	to     string
	length int
}

func New() (*Distancegraph, error) {
	op := "models.distancegraph.New"
	// TODO: transfer distances and vertices to config or database

	distances := []path{
		{"A", "A", 0}, {"B", "B", 0}, {"C", "C", 0},
		{"A", "B", 2}, {"B", "A", 2},
		{"B", "C", 5}, {"C", "B", 5},
	}
	vertices := []string{"A", "B", "C"}

	graph := Distancegraph{goraph.NewGraph(), nil}
	for _, vertex := range vertices {
		err := graph.AddVertex(vertex, nil)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
	}

	for _, path := range distances {
		err := graph.AddEdge(path.from, path.to, float64(path.length), nil)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
	}

	minDistances := make(map[string]map[string]float64)
	for _, from := range vertices {
		dist, _, err := graph.Dijkstra(from)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		vertexDistances := make(map[string]float64)
		for _, to := range vertices {
			vertexDistances[to] = dist[to]
		}
		minDistances[from] = vertexDistances
	}
	graph.minDistances = minDistances

	return &graph, nil
}

func (g *Distancegraph) MinDistance(from string, to string) float64 {
	return g.minDistances[from][to]
}
