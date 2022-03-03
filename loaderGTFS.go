package gotsgraph

import (
	"time"
	"sort"

	gtfs "github.com/takoyaki-3/go-gtfs"
	et "github.com/takoyaki-3/go-gtfs/edge_timetable"
	. "github.com/takoyaki-3/go-gtfs/pkg"
	"github.com/takoyaki-3/go-gtfs/tool"
)

func LoadFromGTFSPath(dirPath string,date string) (*TSGraph, error) {
	g, err := gtfs.Load(dirPath, nil)
	if err != nil {
		return &TSGraph{},err
	}
	
	if t,err := time.Parse("20060102",date);err!=nil{
		return &TSGraph{},err
	} else {
		g = tool.ExtractByDate(g,t)
		return LoadFromGTFS(g)
	}
}

func LoadFromGTFS(g *gtfs.GTFS) (*TSGraph,error){
	ts := new(TSGraph)
	ts.pointHexes = map[string]int{}
	ts.placeIndex = map[string]int{}
	ts.typeIndex = map[string]int{}
	ts.labelIndex = map[string]int{}
	ts.placeIndex[""] = 0
	ts.typeIndex[""] = 0
	ts.labelIndex[""] = 0
	ts.pointHexes[""] = 0
	ts.Places = append(ts.Places, "")
	ts.Types = append(ts.Types, "")
	ts.Labels = append(ts.Labels, "")
	ts.PlacePoints = append(ts.PlacePoints, []int{})

	// 変換
	edgeTimetable := et.GTFS2TimeTableEdges(g)
	sort.Slice(edgeTimetable.Edges,func(i,j int)bool{
		iv := edgeTimetable.Edges[i]
		jv := edgeTimetable.Edges[j]
		iStr := iv.FromStop + ":" + iv.ToStop + ":" + iv.DepartureTime + ":" + iv.ArrivalTime + ":" + iv.TripId
		jStr := jv.FromStop + ":" + jv.ToStop + ":" + jv.DepartureTime + ":" + jv.ArrivalTime + ":" + jv.TripId
		return iStr < jStr
	})

	for _, e := range edgeTimetable.Edges {
		if e.PickupType != 1 {
			ts.AddEdge(Edge{
				FromPoint: ts.Point(e.FromStop, HHMMSS2Sec(e.DepartureTime), "PT_Stop"),
				ToPoint:   ts.Point(e.FromStop, HHMMSS2Sec(e.DepartureTime), "PT_Event"),
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
				FromPoint: ts.Point(e.ToStop, HHMMSS2Sec(e.ArrivalTime), "PT_Event"),
				ToPoint:   ts.Point(e.ToStop, HHMMSS2Sec(e.ArrivalTime), "PT_Stop"),
				TypeIndex: ts.SetType("ET_PT2Stop"),
			})
		}
	}

	return ts, nil
}