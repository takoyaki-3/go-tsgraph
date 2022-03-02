package gotsgraph

import (
	gtfs "github.com/takoyaki-3/go-gtfs"
	et "github.com/takoyaki-3/go-gtfs/edge_timetable"
	. "github.com/takoyaki-3/go-gtfs/pkg"
)

func LoadFromGTFS(dirPath string) (*TSGraph, error) {
	ts := new(TSGraph)
	ts.pointHexes = map[string]int{}
	ts.placeIndex = map[string]int{}
	ts.typeIndex = map[string]int{}
	ts.labelIndex = map[string]int{}

	g, err := gtfs.Load(dirPath, nil)
	if err != nil {
		return ts, err
	}

	// 変換
	edgeTimetable := et.GTFS2TimeTableEdges(g)

	for _, e := range edgeTimetable.Edges {
		if e.PickupType != 1 {
			ts.AddEdge(Edge{
				FromPoint: ts.Point(e.FromStop, HHMMSS2Sec(e.DepartureTime), "PT_Stop"),
				ToPoint:   ts.Point(e.ToStop, HHMMSS2Sec(e.ArrivalTime), "PT_Event"),
				TypeIndex: ts.SetType("ET_Stop2PT"),
			})
		}
		ts.AddEdge(Edge{
			FromPoint:  ts.Point(e.FromStop, HHMMSS2Sec(e.DepartureTime), "PT_Event"),
			ToPoint:    ts.Point(e.ToStop, HHMMSS2Sec(e.ArrivalTime), "PT_Event"),
			TypeIndex:  ts.SetType("ET_PublicTransport"),
			LabelIndex: ts.SetLabel(e.TripId),
		})
		if e.DropOffType != 1 {
			ts.AddEdge(Edge{
				FromPoint: ts.Point(e.FromStop, HHMMSS2Sec(e.DepartureTime), "PT_Event"),
				ToPoint:   ts.Point(e.ToStop, HHMMSS2Sec(e.ArrivalTime), "PT_Stop"),
				TypeIndex: ts.SetType("ET_PT2Stop"),
			})
		}
	}

	return ts, nil
}
