package main

import (
	Hyuga "Hyuga/core"
	"Hyuga/utils"
)

func main() {
	defer utils.Recorder.Close()

	app := Hyuga.Create()
	go Hyuga.DNSDogServe()
	app.Logger.Fatal(app.Start("localhost:5000"))
}
