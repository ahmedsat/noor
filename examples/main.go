package main

import (
	"github.com/ahmedsat/bayaan"
	"github.com/ahmedsat/noor"
)

func main() {

	bayaan.Setup(
		bayaan.WithLevel(bayaan.LoggerLevelDebug),
		bayaan.WithTimeFormat("15:04:05.000"),
	)
	defer bayaan.Close()

	noor.Init(800, 600, "noor example", nil)
	defer noor.Close()

	for !noor.ShouldClose() {
	}

}
