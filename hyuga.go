package main

import (
	Hyuga "Hyuga/core"
	"Hyuga/database"
)

func main() {
	defer database.Recorder.Close()

	app := Hyuga.Create()
	go Hyuga.DNSDogServe(":53")
	app.Logger.Fatal(app.Start(":5000"))
}
