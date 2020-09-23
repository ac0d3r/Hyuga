package main

import (
	Hyuga "Hyuga/core"
	"Hyuga/database"
)

func main() {
	defer database.Recorder.Close()

	app := Hyuga.Create()
	go Hyuga.DNSDogServe()
	app.Logger.Fatal(app.Start("localhost:5000"))
}
