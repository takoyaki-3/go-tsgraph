package gotsgraph

import (
	"io/ioutil"

	"github.com/golang/protobuf/proto"
	pb "github.com/takoyaki-3/go-tsgraph/pb"
)

func LoadFromPath(fileName string) (*TSGraph, error) {

	ts := new(TSGraph)

	in, err := ioutil.ReadFile(fileName)
	if err != nil {
		return ts, err
	}
	tsg := &pb.TSGraph{}
	if err := proto.Unmarshal(in, tsg); err != nil {
		return ts, err
	}

	ts.Name = tsg.Name
	ts.Labels = tsg.Labels
	ts.Types = tsg.Types
	ts.Places = tsg.Places

	ts.PlacePoints = make([][]int, len(ts.Places))
	for pi, p := range tsg.Point {
		ts.Points = append(ts.Points, Point{
			PlaceIndex: int(p.Place),
			Time:       int(p.Time),
			LabelIndex: int(p.Label),
			TypeIndex:  int(p.Type),
		})
		ts.PlacePoints[p.Place] = append(ts.PlacePoints[p.Place], pi)
	}

	ts.FromEdges = make([][]int, len(ts.Points))
	ts.ToEdges = make([][]int, len(ts.Points))

	for ei, e := range tsg.Edge {
		ts.Edges = append(ts.Edges, Edge{
			FromPoint:  int(e.FromPoint),
			ToPoint:    int(e.ToPoint),
			TypeIndex:  int(e.Type),
			LabelIndex: int(e.Label),
			Weights:    e.Weight,
		})
		ts.FromEdges[e.FromPoint] = append(ts.FromEdges[e.FromPoint], ei)
		ts.ToEdges[e.ToPoint] = append(ts.ToEdges[e.ToPoint], ei)
	}

	ts.SetIndexes()

	return ts, nil
}

func (ts *TSGraph) DumpToFile(fileName string) error {

	points := []*pb.Point{}
	for _, p := range ts.Points {
		points = append(points, &pb.Point{
			Place: int32(p.PlaceIndex),
			Time:  int32(p.Time),
			Label: int32(p.LabelIndex),
			Type:  int32(p.TypeIndex),
		})
	}

	edges := []*pb.Edge{}
	for _, e := range ts.Edges {
		edges = append(edges, &pb.Edge{
			FromPoint: int32(e.FromPoint),
			ToPoint:   int32(e.ToPoint),
			Type:      int32(e.TypeIndex),
			Label:     int32(e.LabelIndex),
			Weight:    e.Weights,
		})
	}

	tsg := &pb.TSGraph{
		Point:  points,
		Edge:   edges,
		Places: ts.Places,
		Labels: ts.Labels,
		Types:  ts.Types,
	}

	// Write to disk.
	out, err := proto.Marshal(tsg)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(fileName, out, 0644)
}
