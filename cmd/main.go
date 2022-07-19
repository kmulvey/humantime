package main

import (
	"flag"
	"fmt"

	"github.com/kmulvey/humantime"
)

// go run main.go -time "3 days ago"
func main() {
	var tr humantime.TimeRange
	flag.Var(&tr, "time", "time to parse")
	flag.Parse()
	fmt.Println(tr)
}
