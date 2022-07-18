package main

import (
	"flag"
	"fmt"

	"github.com/kmulvey/string2time"
)

// go run main.go -time "3 days ago"
func main() {
	var tr string2time.TimeRange
	flag.Var(&tr, "time", "time to parse")
	flag.Parse()
	fmt.Println(tr)
}
