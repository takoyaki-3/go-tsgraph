package gotsgraph

import (
	"strconv"
)

type TSGraph struct {
	Name        string  `json:"name"`
	Points      []Point `json:"points"`
	Edges       []Edge  `json:"edges"`
	FromEdges   [][]int
	ToEdges     [][]int
	Places      []string `json:"places"`
	Labels      []string `json:"labels"`
	Types       []string `json:"types"`
	PlacePoints [][]int

	placeIndex map[string]int
	typeIndex  map[string]int
	labelIndex map[string]int
	pointHexes map[string]int
}

type Point struct {
	PlaceIndex int `json:"place_index"`
	Time       int `json:"time"`
	TypeIndex  int `json:"type_index"`
	LabelIndex int `json:"label_index"`
}

type Edge struct {
	FromPoint  int       `json:"from_point"`
	ToPoint    int       `json:"to_point"`
	TypeIndex  int       `json:"type_index"`
	LabelIndex int       `json:"label_index"`
	Weights    []float64 `json:"weight"`
}

func (ts *TSGraph) pointHex(p Point) string {
	place := ts.Places[p.PlaceIndex]
	ty := ts.Types[p.TypeIndex]
	return place + ":" + strconv.Itoa(p.Time) + ":" + ty
}

func (ts *TSGraph) Point(place string, t int, ty string) int {
	pid := place + ":" + strconv.Itoa(t) + ":" + ty
	if p, ok := ts.pointHexes[pid]; ok {
		return p
	}
	pi := len(ts.Points)
	ts.pointHexes[pid] = pi
	ts.Points = append(ts.Points, Point{
		PlaceIndex: ts.GetPlaceIndex(place),
		Time:       t,
		TypeIndex:  ts.SetType(ty),
	})
	ts.FromEdges = append(ts.FromEdges, []int{})
	ts.ToEdges = append(ts.ToEdges, []int{})

	return pi
}

func (ts *TSGraph) GetPlaceIndex(place string) int {
	if pi, ok := ts.placeIndex[place]; ok {
		return pi
	}
	pi := len(ts.placeIndex)
	ts.placeIndex[place] = pi
	ts.Places = append(ts.Places, place)
	return pi
}

func (ts *TSGraph) SetType(ty string) int {
	if ti, ok := ts.typeIndex[ty]; ok {
		return ti
	}
	ti := len(ts.typeIndex)
	ts.typeIndex[ty] = ti
	ts.Types = append(ts.Types, ty)
	return ti
}

func (ts *TSGraph) AddEdge(e Edge) {
	ei := len(ts.Edges)
	ts.Edges = append(ts.Edges, e)
	ts.FromEdges[e.FromPoint] = append(ts.FromEdges[e.FromPoint], ei)
	ts.ToEdges[e.ToPoint] = append(ts.ToEdges[e.ToPoint], ei)
}

func (ts *TSGraph) SetLabel(label string) int {
	if li, ok := ts.labelIndex[label]; ok {
		return li
	}
	li := len(ts.labelIndex)
	ts.labelIndex[label] = li
	ts.Labels = append(ts.Labels, label)
	return li
}

func (ts *TSGraph) SetIndexes() {
	ts.placeIndex = map[string]int{}
	ts.typeIndex = map[string]int{}
	ts.labelIndex = map[string]int{}
	ts.pointHexes = map[string]int{}
	for i, v := range ts.Places {
		ts.placeIndex[v] = i
	}
	for i, v := range ts.Types {
		ts.typeIndex[v] = i
	}
	for i, v := range ts.Labels {
		ts.labelIndex[v] = i
	}
	for i, p := range ts.Points {
		ts.pointHexes[ts.pointHex(p)] = i
	}
}
