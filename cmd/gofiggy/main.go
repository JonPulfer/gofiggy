package main

import (
	"github.com/JonPulfer/gofiggy/pkg/controller"
	"github.com/JonPulfer/gofiggy/pkg/events"
)

func main() {
	var eventHandler = events.MockHandler{}
	controller.Start("default", eventHandler)
}
