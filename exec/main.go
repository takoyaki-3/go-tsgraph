package main

import (
	"fmt"
	"log"

	tsg "github.com/takoyaki-3/go-tsgraph"
)

func main() {
	// ts,err := tsg.LoadFromGTFS("./GTFS")
	ts, err := tsg.LoadFromPath("./ts.tsg.pbf")
	if err != nil {
		log.Fatalln(err)
	}

	ts.DumpToFile("./ts.tsg.pbf")
	if false {
		fmt.Println(ts)
	}
}
